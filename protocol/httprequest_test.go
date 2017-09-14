package protocol_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/config"
	"github.com/gkewl/pulsecheck/constant"
	"github.com/gkewl/pulsecheck/dbhandler"
	"github.com/gkewl/pulsecheck/protocol"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"

	"github.com/gkewl/pulsecheck/utilities"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestStruct struct {
	Name string `json:"name"`
	Vals []int  `json:"vals"`
}

var _ = Describe("HTTPRequest", func() {
	var router *mux.Router
	var server *httptest.Server
	var ctx = common.AppContext{}
	var err error
	var response *http.Response
	config.LoadConfigurations()

	ctx.Db, _ = dbhandler.CreateConnection()

	BeforeEach(func() {
		router = mux.NewRouter().StrictSlash(true)
		server = httptest.NewServer(router)
	})

	AfterEach(func() {
		Expect(err).To(BeNil())
		Expect(response.StatusCode).To(Equal(401))
		server.Close()
	})

	It("makes url parameters available", func() {
		var testfunc = func(w http.ResponseWriter, r *http.Request) {
			defer GinkgoRecover()
			rc, _ := protocol.NewHTTPRequestContext(&ctx, r)
			Expect(rc.Value("name", "")).To(Equal("foo"))
			Expect(rc.Value("notthere", "bar")).To(Equal("bar"))
			Expect(rc.IntValue("id", 0)).To(Equal(int64(42)))
			Expect(rc.IntValue("notthere", 24)).To(Equal(int64(24)))
			Expect(rc.IntValue("name", 24)).To(Equal(int64(24)))
			w.WriteHeader(401)
		}
		router.HandleFunc("/test/{name}/sub/{id}", testfunc).Methods("GET")
		request, _ := http.NewRequest("GET", fmt.Sprintf("%s/test/foo/sub/42", server.URL), nil)
		response, err = http.DefaultClient.Do(request)
	})

	It("makes form url parameters available", func() {
		var testfunc = func(w http.ResponseWriter, r *http.Request) {
			defer GinkgoRecover()
			rc, _ := protocol.NewHTTPRequestContext(&ctx, r)
			Expect(rc.Value("name", "")).To(Equal("foo"))
			Expect(rc.IntValue("id", 0)).To(Equal(int64(42)))
			w.WriteHeader(401)
		}
		router.HandleFunc("/test", testfunc).Methods("GET")
		request, _ := http.NewRequest("GET", fmt.Sprintf("%s/test?name=foo&id=42", server.URL), nil)
		response, err = http.DefaultClient.Do(request)
	})

	It("makes posted data parameters available to xxxValue() calls", func() {
		var testfunc = func(w http.ResponseWriter, r *http.Request) {
			defer GinkgoRecover()
			rc, _ := protocol.NewHTTPRequestContext(&ctx, r)
			Expect(rc.Value("name", "")).To(Equal("foo"))
			Expect(rc.IntValue("num", 0)).To(Equal(int64(34)))
			Expect(rc.IntValue("largenum", 0)).To(Equal(int64(1234567890)))
			Expect(rc.FloatValue("floater", 0)).To(Equal(float64(1234567890.123456)))
			Expect(rc.BoolValue("bool", false)).To(Equal(true))
			Expect(rc.Value("notthere", "def")).To(Equal("def"))
			w.WriteHeader(401)
		}
		router.HandleFunc("/test", testfunc).Methods("POST")
		samplePost := `{"name":"foo", "num":34, "largenum":1234567890, "floater":1234567890.123456, "bool":true}`
		request, _ := http.NewRequest("POST", fmt.Sprintf("%s/test", server.URL), bytes.NewBufferString(samplePost))
		response, err = http.DefaultClient.Do(request)
	})

	It("scans posted content", func() {
		var testfunc = func(w http.ResponseWriter, r *http.Request) {
			dest := TestStruct{}
			defer GinkgoRecover()
			rc, _ := protocol.NewHTTPRequestContext(&ctx, r)
			err = rc.Scan("ignored", &dest)
			Expect(dest.Name).To(Equal("foo"))
			Expect(dest.Vals).To(Equal([]int{1, 2, 3}))
			w.WriteHeader(401)
		}
		router.HandleFunc("/test", testfunc).Methods("POST")
		request, _ := http.NewRequest("POST", fmt.Sprintf("%s/test", server.URL), bytes.NewBufferString(`{"name":"foo","vals":[1,2,3]}`))
		request.Header.Add("Content-type", "application/json")
		response, err = http.DefaultClient.Do(request)
	})

	It("file posted formfile", func() {
		var testfunc = func(w http.ResponseWriter, r *http.Request) {
			defer GinkgoRecover()
			rc, err := protocol.NewHTTPRequestContext(&ctx, r)
			if err != nil {
				Expect(err).To(BeNil())
			}
			dest := rc.RequestUploadFiles()
			Expect(len(dest)).To(Equal(1))
			//Contentcheck
			if b, err := ioutil.ReadAll(dest[0].File); err == nil {
				Expect(string(b)).To(Equal("MOS S3 Upload demo"))
			}
			w.WriteHeader(401)
		}

		path := "/tmp/temp.txt"
		if runtime.GOOS == "windows" {
			path = "c:/temp/temp.txt"
		}

		// detect if file exists
		var _, err = os.Stat(path)

		// create file if not exists
		if os.IsNotExist(err) {
			var file, err = os.Create(path)
			if err != nil {
				log.Fatal(err)
				os.Exit(0)
			}

			defer file.Close()
		}
		// Adding Content to file

		wf, err := os.OpenFile(path, os.O_RDWR, 0644)
		n, err := io.WriteString(wf, "MOS S3 Upload demo")
		if err != nil {
			log.Fatal(err, n)
			return
		}
		// save changes
		err = wf.Sync()
		if err != nil {
			return
		}
		wf.Close()
		f, err := os.Open(path)
		defer f.Close()
		// returns file info
		fi, err := f.Stat()

		//Prepare a form that you will submit
		buf := new(bytes.Buffer)

		extraParams := map[string]string{
			"title":       "Test Document",
			"author":      "Rama",
			"description": "Multipart form Test",
		}
		writer := multipart.NewWriter(buf)

		part, err := writer.CreateFormFile("file", fi.Name())

		if err != nil {
			log.Fatal(err)
		}
		if _, err = io.Copy(part, f); err != nil {
			log.Fatal(err)
		}
		for key, val := range extraParams {
			_ = writer.WriteField(key, val)
		}

		if err = writer.Close(); err != nil {
			log.Fatal(err)
		}
		//Now that you have form, you can submit it
		router.HandleFunc("/test", testfunc).Methods("POST")
		request, err := http.NewRequest("POST", fmt.Sprintf("%s/test", server.URL), buf)
		if err != nil {
			return
		}
		// Don't forget to set the content type, this will contain the boundary.
		request.Header.Set("Content-Type", writer.FormDataContentType())
		response, err = http.DefaultClient.Do(request)
	})

	It("captures xid from header and puts it in response", func() {
		var testfunc = func(w http.ResponseWriter, r *http.Request) {
			defer GinkgoRecover()
			r = r.WithContext(context.WithValue(r.Context(), constant.Xid, "foo"))
			rc, _ := protocol.NewHTTPRequestContext(&ctx, r)
			Expect(rc.Xid()).To(Equal("foo"))
			utilities.WriteJSON(w, 401, common.StructuredResponse{Xid: "foobar", StatusCode: 401, Response: "ok"})
		}
		router.HandleFunc("/test", testfunc).Methods("GET")
		request, _ := http.NewRequest("GET", fmt.Sprintf("%s/test?name=foo&id=42", server.URL), nil)
		response, err = http.DefaultClient.Do(request)
		Expect(err).To(BeNil())
		Expect(response.Header.Get("X-Xid")).To(Equal("foobar"))
	})
})

package protocol_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/samv/sse"

	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/config"
	"github.com/gkewl/pulsecheck/constant"
	"github.com/gkewl/pulsecheck/dbhandler"

	"github.com/gkewl/pulsecheck/protocol"
)

func decodeResponseBasic(resp *http.Response) (map[string]interface{}, error) {
	var decoded map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, err
	}
	return decoded, nil
}

var _ = Describe("HTTPProtocol", func() {

	var wRec *httptest.ResponseRecorder
	var appCtx = &common.AppContext{}
	var dummyHandler *protocol.NewHTTPHandler
	var dummyRequest *http.Request

	config.LoadConfigurations()
	var dbErr error
	appCtx.Db, dbErr = dbhandler.CreateConnection()
	if dbErr != nil {
		panic(dbErr)
	}

	BeforeEach(func() {
		wRec = httptest.NewRecorder()
		dummyRequest = httptest.NewRequest("GET", "/dummy", nil)
		dummyHandler = &protocol.NewHTTPHandler{
			AppContext: appCtx,
			Route: &common.Route{
				Name:    "GetDummy",
				Method:  "GET",
				Pattern: "/dummy",
				ControllerFunc: func(common.RequestContext) (interface{}, error) {
					return nil, fmt.Errorf("Test missing Controller func")
				},
				AuthRequired:   constant.Guest,
				NormalHttpCode: http.StatusOK,
			},
		}
	})

	It("Can handle basic requests", func() {
		dummyHandler.Route.ControllerFunc = func(common.RequestContext) (interface{}, error) {
			return "hello", nil
		}
		dummyHandler.ServeHTTP(wRec, dummyRequest)
		resp := wRec.Result()
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		Expect(resp.Header.Get("Content-Type")).To(HavePrefix("application/json"))
		decoded, err := decodeResponseBasic(resp)
		Expect(err).To(BeNil())

		Expect(decoded).To(HaveKey("xid"))
		Expect(decoded["xid"]).To(MatchRegexp(`^[0-9a-zA-Z_\-]{10,}$`))

		Expect(decoded).To(HaveKey("statuscode"))
		Expect(decoded["statuscode"]).To(BeNumerically("==", http.StatusOK))

		Expect(decoded).To(HaveKey("response"))
		Expect(decoded["response"]).To(Equal("hello"))
	})

	It("can handle stream responses", func() {
		// for the sake of this test's brevity (so that .Equal() can be used at the end),
		// these are all types returned by json.Unmarshal into an interface{}
		testData := []interface{}{
			map[string]interface{}{"hello": "realtime"},
			"bob",
			[]interface{}{"a", "quick", "brown", "fox", "yada", "yada"},
		}

		dummyHandler.Route.ControllerFunc = func(common.RequestContext) (interface{}, error) {
			return sse.NewJSONEncoderFeed(&feedMeSeymour{testData}), nil
		}

		By("Starting a test HTTP server")
		server := httptest.NewServer(dummyHandler)
		client := sse.NewSSEClient()
		By("Connecting to that server with an SSE client")
		err := client.GetStream(server.URL)
		Expect(err).To(BeNil())

		timer := time.NewTimer(2 * time.Second)
		var readMessages []interface{}
		By("Reading messages for up to 2s...")
	readLoop:
		for {
			select {
			case ev, ok := <-client.Messages():
				if ok {
					By(fmt.Sprintf("Reading an event: %s", string(ev.Data)))
					var decoded interface{}
					err = json.Unmarshal(ev.Data, &decoded)
					Expect(err).To(BeNil())
					readMessages = append(readMessages, decoded)
					if len(readMessages) == len(testData) {
						timer.Reset(100 * time.Millisecond)
					}
				} else {
					By("Reading eof (?)")
					break readLoop
				}
			case <-timer.C:
				By("giving up after 2s")
				break readLoop
			}
		}
		// use goroutines because who cares if they clean up, amirite?
		go client.Close()
		go server.Close()

		Expect(readMessages).To(Equal(testData))
	})
})

type feedMeSeymour struct {
	messages []interface{}
}

// GetEventChan is an example SSE server which just returns the
// events, one by one, with a brief delay in between.
func (fms *feedMeSeymour) GetEventChan(ccc <-chan struct{}) <-chan interface{} {
	sinkChan := make(chan interface{})
	go func() {
		var next = fms.messages[0]
		for {
			time.Sleep(50 * time.Millisecond)
			select {
			case sinkChan <- next:
				fms.messages = fms.messages[1:]
				if len(fms.messages) == 0 {
					sinkChan = nil
				} else {
					next = fms.messages[0]
				}
			case _, ok := <-ccc:
				if !ok {
					return
				}
			}
		}
	}()
	return sinkChan
}

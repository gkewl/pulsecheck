package protocol

import (
	"context"
	"encoding/json"
	"fmt"

	auth "github.com/gkewl/pulsecheck/authentication"
	"github.com/gkewl/pulsecheck/common"
	//	"github.com/gkewl/pulsecheck/config"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gkewl/pulsecheck/constant"
	"github.com/gkewl/pulsecheck/utilities"

	"github.com/gorilla/mux"
)

// HTTPRequestContext implements the RequestContext interface using information from an
// inbound HTTP request
type HTTPRequestContext struct {
	common.RequestContextBase
	r           *http.Request
	requestBody []byte
	files       []common.Upload
	urlVars     map[string]string
	bodyVars    interface{}
}

// NewHTTPRequestContext returns an initialized RequestContext ready for request processing
// or an error if the HTTP request body could not be read or there was a database error
func NewHTTPRequestContext(appCtx *common.AppContext, r *http.Request) (hrc HTTPRequestContext, err error) {
	// todo: restore this check when front end is sending correct content-type
	// if err = utilities.CheckForJSON(r); err != nil {
	// 	return
	// }
	//To Check Memory available 24MB
	const MEMORY24MB = (1 << 20) * 24
	err = r.ParseMultipartForm(MEMORY24MB)
	//Max memory Argument
	if r.MultipartForm != nil && r.MultipartForm.File != nil {
		for _, fheaders := range r.MultipartForm.File {
			for _, hdr := range fheaders {
				// open uploaded
				var infile multipart.File
				if infile, err = hdr.Open(); err != nil {
					continue
				}
				// Only the first 512 bytes are used to sniff the content type.
				var size int64 = 512

				buffer := make([]byte, size)

				// read file content to buffer
				if err != nil && err != io.EOF {
					log.Fatal(err)
				}

				defer infile.Close()

				// filetype gets the Content Type of the File
				filetype := http.DetectContentType(buffer)
				// Reset the read pointer if necessary.
				infile.Seek(0, 0)

				file := common.Upload{infile, hdr.Filename, filetype}

				hrc.files = append(hrc.files, file)
			}
		}
	}

	if r.Body != nil && r.ContentLength > 0 {
		defer func() { _ = r.Body.Close() }()
		if hrc.requestBody, err = ioutil.ReadAll(r.Body); err != nil {
			return
		}
	}

	// this timeout mostly affects MOS internal clients, eg database handles
	timeout, _ := time.ParseDuration("60s")

	hrc.Context = context.TODO() //context.Background()

	hrc.Context, hrc.CancelFunc = context.WithTimeout(hrc.Context, timeout)

	if hrc.Txn, err = appCtx.Db.BeginTxx(hrc.Context, nil); err != nil {
		return
	}

	hrc.MaxDeadlockRetries = constant.MaxDeadlockRetries
	hrc.Xnid, _ = r.Context().Value(constant.Xid).(string)
	hrc.Username, hrc.Userid, hrc.Userrole, hrc.Companyid = auth.GetUserInfoFromContext(r.Context())
	hrc.AppCtx = appCtx
	hrc.r = r
	hrc.urlVars = mux.Vars(r)
	return
}

type CancelFunc func()

// Token returns the token presented with this request
func (hrc *HTTPRequestContext) Token() string {
	return strings.Replace(hrc.r.Header.Get("Authorization"), "Bearer ", "", 1)
}

// Done calls the cancel func to clean up resources linked to the context
func (hrc *HTTPRequestContext) Done() {
	hrc.CancelFunc()
}

// Value returns the named input variable as a string
func (hrc *HTTPRequestContext) Value(name string, defValue string) (val string) {
	var present bool
	if val, present = hrc.urlVars[name]; present == false {
		val = hrc.r.FormValue(name)
		if val == "" {
			return hrc.BodyValue(name, defValue)
		}
	}
	return val
}

// BodyValue tries to get named parameter from json parse of body
func (hrc *HTTPRequestContext) BodyValue(name, defValue string) string {
	if len(hrc.requestBody) > 0 {
		if hrc.bodyVars == nil {
			_ = json.Unmarshal(hrc.requestBody, &hrc.bodyVars)
		}
		if hrc.bodyVars != nil {
			if varMap, ok := hrc.bodyVars.(map[string]interface{}); ok {
				if val, present := varMap[name]; present {
					// avoid default exponent for numerics
					switch value := val.(type) {
					case int, int16, int32, int64:
						return fmt.Sprintf("%d", value)
					case float32, float64:
						return fmt.Sprintf("%f", value)
					}
					return fmt.Sprintf("%v", val)
				}
			}
		}
	}
	return defValue
}

// IntValue returns the named input variable as an int64
func (hrc *HTTPRequestContext) IntValue(name string, defValue int64) int64 {
	val := hrc.Value(name, "default")
	if val == "default" {
		return defValue
	}
	if fVal, err := strconv.ParseFloat(val, 64); err == nil {
		return int64(fVal)
	}
	return defValue
}

// IntValue32 returns the named input variable as an int64
func (hrc *HTTPRequestContext) IntValue32(name string, defValue int) int {
	val := hrc.Value(name, "default")
	if val == "default" {
		return defValue
	}
	if iVal, err := strconv.Atoi(val); err == nil {
		return iVal
	}
	return defValue
}

// BoolValue returns the named input variable as an bool
func (hrc *HTTPRequestContext) BoolValue(name string, defValue bool) bool {
	val := hrc.Value(name, "default")
	if val == "default" {
		return defValue
	}
	if iVal, err := strconv.ParseBool(val); err == nil {
		return iVal
	}
	return defValue
}

// FloatValue returns the named input variable as a float64
func (hrc *HTTPRequestContext) FloatValue(name string, defValue float64) float64 {
	val := hrc.Value(name, "default")
	if val == "default" {
		return defValue
	}
	if fVal, err := strconv.ParseFloat(val, 64); err == nil {
		return fVal
	}
	return defValue
}

// Time Value returns the named input variable as a time.Time
func (hrc *HTTPRequestContext) TimeValue(name string, defValue time.Time) time.Time {
	val := hrc.Value(name, "default")
	if val == "default" {
		return defValue
	}
	if parsedVal, err := utilities.ParseStringtoTime(val); err == nil {
		return parsedVal
	}
	return defValue
}

// Scan will unmarshal json context from the request body into the destination
// Callers should pass the address of their destination struct
// If there is no body no action is taken. Returns error for JSON parse problem
func (hrc *HTTPRequestContext) Scan(name string, dest interface{}) error {
	if len(hrc.requestBody) > 0 {
		return json.Unmarshal(hrc.requestBody, dest)
	}
	return nil
}

func (hrc *HTTPRequestContext) RequestBody() []byte {
	return hrc.requestBody
}

func (hrc *HTTPRequestContext) RequestUploadFiles() []common.Upload {
	return hrc.files
}

// ResetForRetry clears out deferred functions, log fields and rolls back
// the transaction
func (req *HTTPRequestContext) ResetForRetry() (err error) {
	req.Tx().Rollback()
	req.Txn, err = req.AppCtx.Db.Beginx()
	//req.RequestContextBase.ClearDeferredRequests()
	//req.LogValues = map[string]interface{}{}
	return
}

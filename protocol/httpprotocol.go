package protocol

import (
	"context"
	"net/http"
	"time"

	"github.com/samv/sse"

	auth "github.com/gkewl/pulsecheck/authentication"
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/constant"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/logger"
	"github.com/gkewl/pulsecheck/utilities"
	"github.com/gkewl/pulsecheck/xid"
)

// NewHTTPHandler is created by the router for new-style controller management
type NewHTTPHandler struct {
	AppContext *common.AppContext
	Route      *common.Route
}

// ServeHTTP standard handler for an HTTP request. Delegates to the controller
// function for this path
func (ah NewHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// verify user and authorized to perform
	if ah.Route.AuthRequired > constant.Guest {
		r = auth.AnnotateRequestWithAuthorizedUser(w, r, ah.Route.AuthRequired)
	}

	if r == nil {
		return // error response already written
	}

	var requestEnd, processingEnd, dbEnd time.Time

	start := time.Now()
	xidvalue := r.Header.Get("X-Xid")
	if xidvalue == "" {
		xidvalue = xid.UniqueIdGenerator()
	}
	goCtx := context.WithValue(r.Context(), constant.Xid, xidvalue)
	r = r.WithContext(goCtx)

	reqCtx, err := NewHTTPRequestContext(ah.AppContext, r)
	defer reqCtx.Done()
	requestEnd = time.Now()

	code := http.StatusOK
	var evSource sse.EventFeed
	if err == nil {
		var result interface{}
		// Handle deadlock retry
		var retry = true
		for retry {
			result, err = ah.Route.ControllerFunc(&reqCtx)
			retry, err = ProcessDeadlock(&reqCtx, err)
		}
		processingEnd = time.Now()
		if err == nil {
			err = reqCtx.Tx().Commit()
			if err == nil {
				if ah.Route.NormalHttpCode != code && ah.Route.NormalHttpCode != 0 {
					code = ah.Route.NormalHttpCode
				}
				if reqCtx.IsRawResponse {
					err = utilities.WriteRawResponse(w, code, reqCtx.RawContentType, result)
				} else if source, ok := result.(sse.EventFeed); ok {
					evSource = source
				} else {
					response := common.StructuredResponse{
						Xid:        xidvalue,
						StatusCode: code,
						Response:   result,
					}
					err = utilities.WriteJSON(w, code, response)
				}
			}
		} else {
			_ = reqCtx.Tx().Rollback()
		}
		dbEnd = time.Now()
		RunDeferredRequests(&reqCtx, (err == nil))
	}
	fields := map[string]interface{}{
		"protocol":   "http",
		"reading":    utilities.DurationTruncated(start, requestEnd, 4),
		"processing": utilities.DurationTruncated(requestEnd, processingEnd, 4),
		"commit":     utilities.DurationTruncated(processingEnd, dbEnd, 4),
		"retries":    reqCtx.DeadlockRetryCount(),
		"method":     r.Method,
		"path":       r.URL.Path,
		"status":     code,
		"action":     ah.Route.Name,
		"user":       reqCtx.UserName(),
	}
	for k, v := range *reqCtx.LogValues() {
		fields[k] = v
	}
	lm := logger.LogModel{
		Level:    logger.INFO,
		Msg:      "success",
		Xid:      xidvalue,
		Fields:   fields,
		Duration: utilities.DurationTruncated(start, time.Now(), 4),
	}
	if err == nil {
		logger.Log(lm)
		if evSource != nil {
			// disable nginx buffering
			w.Header().Set("X-Accel-Buffering", "no")
			err := sse.SinkEvents(w, code, evSource)
			if err != nil {
				http.Error(w, "failed to sink JSON events", http.StatusInternalServerError)
			}
		}
	} else {
		ah.respondWithError(w, r, err, xidvalue)
		if !ah.Route.SecureBody {
			lm.Input = string(reqCtx.RequestBody())
		}
		lm.Level = logger.ERROR
		lm.Fields["status"] = http.StatusInternalServerError
		if ee, ok := err.(eh.Error); ok {
			lm.Fields["status"] = ee.NamedError.HTTPStatus()
			lm.Err = err
			lm.Msg = ee.DetailStack("; ")
			lm.Caller = ee.LocationStack("; ")
		}
		logger.ErrorLog(lm)
	}
}

type clientError struct {
	ErrorInfo string
	Message   string
	Xid       string
}

// respondWithError writes a response with error info back to client
func (ah NewHTTPHandler) respondWithError(w http.ResponseWriter, r *http.Request, err error, xid string) {
	statuscode, content := StructuredError(err, xid)
	utilities.WriteJSON(w, statuscode, content)
}

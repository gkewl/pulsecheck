package common

import (
	"context"
	//"encoding/json"
	"fmt"

	"github.com/gkewl/pulsecheck/constant"
	"github.com/gkewl/pulsecheck/logger"

	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gkewl/pulsecheck/xid"
	"github.com/jmoiron/sqlx"
)

//this interface is used for every class which can be
type Routes []Route

type AppHandlerFunc func(*AppContext, http.ResponseWriter, *http.Request) (int, error)

type ControllerFunc func(RequestContext) (interface{}, error)

type Route struct {
	Name               string         // unique action name
	Method             string         // HTTP only method name
	Pattern            string         // HTTP only routing pattern
	HandlerFunc        AppHandlerFunc // old HTTP handler
	ControllerFunc     ControllerFunc // new controller func w/RequestContext
	AuthRequired       int            // describes level of auth needed for this action
	NormalHttpCode     int            // HTTP only on success what HTTP status (e.g. 201 created)
	StructuredResponse bool           // New-style structure to the client response
	SecureBody         bool           // Do not log body when it is set true

}

type Router interface {
	RouteForAction(action string) *Route
}

type APIRoutes []Routes

type job struct {
	name     string
	duration time.Duration
}

// MQTTPub interface for MQTTPublisher publish capability
type MQTTPub interface {
	Publish(PublishRequest) error
	Stop()
}

// AppContext contains our local context; our database pool, session store, template
// registry and anything else our handlers need to access. We'll create an instance of it
// in our main() function and then explicitly pass a reference to it for our handlers to access.
type AppContext struct {
	Context       context.Context
	CancelFunc    context.CancelFunc
	Db            *sqlx.DB
	MQTTPublisher MQTTPub
	Version       string
	BuildTime     string
	GitHash       string
	UseMock       bool
}

// PublishRequest represents a payload to be published to a message broker
type PublishRequest struct {
	Xid           string
	Topic         string
	Key           interface{}
	Payload       interface{}
	QOS           int
	ResponseTopic string
	CommandType   uint32 // 0 = default
}

// RequestContext is an interface returning information about the current request
type RequestContext interface {
	AppContext() *AppContext
	GetContext() context.Context
	UserID() int64
	UserName() string
	UserRole() string
	ActorID() int64
	ActorName() string
	SetActor(int64, string)
	Tx() *sqlx.Tx
	Xid() string
	Value(name string, defValue string) string           // string request parameter
	IntValue(name string, defValue int64) int64          // numeric request parameter
	IntValue32(name string, defValue int) int            // numeric request parameter
	BoolValue(name string, defValue bool) bool           // boolean request parameter
	FloatValue(name string, defValue float64) float64    // float request parameter
	TimeValue(name string, defValue time.Time) time.Time // time request parameter
	Scan(name string, dest interface{}) error            // stores entity in request into dest
	AddLogValue(key string, value interface{})
	LogValues() *map[string]interface{}
	RequestBody() []byte
	RequestUploadFiles() []Upload
	ResetForRetry() error
	SetMaxDeadlockRetries(int)
	BumpDeadlocks() bool     // increment deadlock count
	DeadlockRetryCount() int // how many deadlocks occurred
	Token() string
	SetUserId(int64)
	SetIsRawResponse(bool, string)
}

type RequestContextBase struct {
	Context            context.Context
	CancelFunc         context.CancelFunc
	AppCtx             *AppContext
	Txn                *sqlx.Tx
	Userid             int64
	Username           string
	Userrole           string
	Actorid            int64
	Actorname          string
	Xnid               string
	logValues          map[string]interface{}
	MaxDeadlockRetries int
	Deadlockretrycount int
	IsRawResponse      bool
	RawContentType     string
}

type Upload struct {
	File        io.Reader
	Filename    string
	ContentType string
}

// AppContext returns the application context for this request
func (rc *RequestContextBase) AppContext() *AppContext {
	return rc.AppCtx
}

// GetContext returns the stdlib context.Context for this request Context
func (rc *RequestContextBase) GetContext() context.Context {
	return rc.Context
}

// Tx returns the database transaction wrapping this request
func (rc *RequestContextBase) Tx() *sqlx.Tx {
	return rc.Txn
}

// UserId returns the authenticated user db id
func (rc *RequestContextBase) UserID() int64 {
	return rc.Userid
}

// UserName returns the authenticated user name
func (rc *RequestContextBase) UserName() string {
	return rc.Username
}

// UserRole returns the authenticated user's role
func (rc *RequestContextBase) UserRole() string {
	return rc.Userrole
}

// ActorID returns the actor ID for this request
func (rc *RequestContextBase) ActorID() int64 {
	return rc.Actorid
}

// ActorName returns the actor name for this request
func (rc *RequestContextBase) ActorName() string {
	return rc.Actorname
}

// SetActor sets the actor info for this request
func (rc *RequestContextBase) SetActor(id int64, name string) {
	rc.Actorid = id
	rc.Actorname = name
}

// SetUserId sets the UserId info for this request
func (rc *RequestContextBase) SetUserId(id int64) {
	rc.Userid = id
}

// Xid returns the transaction id for this request
func (rc *RequestContextBase) Xid() string {
	return rc.Xnid
}

// AddLogValue adds a key/value pair to be logged
func (rc *RequestContextBase) AddLogValue(name string, value interface{}) {
	(*rc.LogValues())[name] = value
}

// LogValues returns a map of key/values to be logged
func (rc *RequestContextBase) LogValues() *map[string]interface{} {
	if rc.logValues == nil {
		rc.logValues = map[string]interface{}{}
	}
	return &rc.logValues
}

// SetMaxDeadlockRetries configures how many retries will be performed
func (rc *RequestContextBase) SetMaxDeadlockRetries(max int) {
	rc.MaxDeadlockRetries = max
}

// BumpDeadlocks increments the deadlock count by 1 and returns
// whether retry still available
func (rc *RequestContextBase) BumpDeadlocks() (retry bool) {
	rc.Deadlockretrycount += 1
	return rc.Deadlockretrycount <= rc.MaxDeadlockRetries
}

// DeadlockRetryCount returns how many deadlocks occurred during this request
func (rc *RequestContextBase) DeadlockRetryCount() int {
	return rc.Deadlockretrycount
}

// SetRawResponse sets the IsRawResponse
func (rc *RequestContextBase) SetIsRawResponse(isRawResponse bool, contentType string) {
	rc.IsRawResponse = isRawResponse
	rc.RawContentType = contentType
}

// AppHandler - wrap http.handler embed field *AppContext
// avoid globals and improve middleware chaining and error handling
type AppHandler struct {
	*AppContext
	H     AppHandlerFunc
	Route *Route
}

//Recovery Handler handles panic in our application and returns a 500
func RecoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				debug.PrintStack()
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

type ErrorResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Details   string `json:"details"`
	Locations string `json:"locations"`
}

type StructuredResponse struct {
	Xid        string         `json:"xid"`
	StatusCode int            `json:"statuscode"`
	Response   interface{}    `json:"response,omitempty"`
	Error      *ErrorResponse `json:"error,omitempty"`
}

// ServeHTTP - http handler for AppHandler type
func (AH AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	xidvalue := xid.UniqueIdGenerator()
	goCtx := context.WithValue(r.Context(), constant.Xid, xidvalue)
	r = r.WithContext(goCtx)
	defer r.Context().Done()
	_, err := AH.H(AH.AppContext, w, r) // errors already logged
	if err == nil {
		fields := map[string]interface{}{
			"protocol": "http",
			"method":   r.Method,
			"path":     r.URL.Path,
			"action":   AH.Route.Name,
		}
		logger.Log(logger.LogModel{
			Level:    logger.INFO,
			Msg:      "success",
			Xid:      xidvalue,
			Fields:   fields,
			Duration: time.Since(start).Seconds(),
		})
	} else {

		lm := logger.LogModel{}
		if !AH.Route.SecureBody {
			respbody, _ := ioutil.ReadAll(r.Body)
			lm.Input = string(respbody)
		}
		lm.Level = logger.ERROR
		lm.Err = err
		logger.ErrorLog(lm)
	}

}

func CaselessMatcher(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.ToLower(r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

var varPattern = regexp.MustCompile(`\{\w+(?::([^}]+))?\}`)

// CaptureRegexp returns a capturing, anchored regular expression matching the route's HTTP URL
func (route *Route) CaptureRegexp() (*regexp.Regexp, error) {
	re := varPattern.ReplaceAllStringFunc(route.Pattern, func(varString string) string {
		def := varString[1 : len(varString)-1]
		split := strings.SplitN(def, ":", 2)
		varName := split[0]
		var pattern string
		if len(split) == 2 {
			pattern = split[1]
		} else {
			pattern = "[^/]*"
		}
		return fmt.Sprintf("(?P<%s>%s)", varName, pattern)
	})
	return regexp.Compile("^" + re + "$")
}

package protocol

import (
	"context"
	"encoding/json"
	"regexp"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/constant"
	"github.com/gkewl/pulsecheck/dbhandler"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	//	"github.com/gkewl/pulsecheck/model"
	"github.com/gkewl/pulsecheck/utilities"
	"github.com/gkewl/pulsecheck/xid"
)

// TestCaller allows tests to make generic calls without knowing the protocol
// involved
type TestCaller interface {
	MakeTestCall(*TestConfig) string
	Protocol() string
}

// TestConfig holds information necessary to make a request in a test
// context. It is mostly protocol independent but has some fields that
// are protocol specific, like expected HTTP status
type TestConfig struct {
	Action             string
	AuthLevel          int
	ExpectedHTTPStatus int
	ExpectedError      string
	Params             map[string]string
	ContentName        string
	Content            interface{}
}

// MakeTestConfig constructs a test request configuration with the specified
// action and authLevel and expected HTTP status
func MakeTestConfig(action string, authLevel int, expectedStatus int, expectedError *eh.NamedError,
	params map[string]string, contentName string, content interface{}) *TestConfig {
	tc := TestConfig{
		Action:             action,
		ExpectedHTTPStatus: expectedStatus,
		AuthLevel:          authLevel,
		Params:             params,
		ContentName:        contentName,
		Content:            content,
	}
	if expectedError != nil {
		tc.ExpectedError = expectedError.Description()
	}
	return &tc
}

// TestRequestContext implements the RequestContext interface and allows
// stubbing of inbound request values
type TestRequestContext struct {
	common.RequestContextBase

	Userid          int
	Username        string
	Userrole        string
	Actorid         int64
	Actorname       string
	ScanSource      string
	Values          map[string]string
	appContext      common.AppContext
	tx              *sqlx.Tx
	valuesRequested []string
	logValues       map[string]interface{}

	MaxDeadlockRetries int
	Deadlockretrycount int
	IsRawResponse      bool
	RawContentType     string
	xid                string
}

// NewTestRequestContext allows wrapping existing tx for legacy code
func NewTestRequestContext(tx *sqlx.Tx) *TestRequestContext {
	return &TestRequestContext{
		tx: tx,
	}
}

// Complete cleans up the test request committing or rolling back the Tx
func (req *TestRequestContext) Complete(commit bool) {
	if req.tx != nil {
		if commit {
			_ = req.tx.Commit()
		} else {
			_ = req.tx.Rollback()
		}
		req.tx = nil
	}
}

// AppContext returns the context
func (req *TestRequestContext) AppContext() *common.AppContext {
	if req.appContext.Db == nil {
		req.appContext.Db, _ = dbhandler.CreateConnection()
	}
	return &req.appContext
}

// GetContext returns a context object.
func (req *TestRequestContext) GetContext() context.Context {
	claims := make(map[string]interface{})
	claims["iat"] = time.Now().Unix()
	claims["sub"] = req.Username
	claims["userid"] = float64(req.Userid)
	claims["scope"] = req.Userrole
	return context.WithValue(context.TODO(), constant.Claims, claims)
}

// ServiceName returns the service name
func (req *TestRequestContext) ServiceName() string {
	return "dummytestservice"
}

// UserID returns user id
func (req *TestRequestContext) UserID() int {
	return req.Userid
}

// UserName returns user name
func (req *TestRequestContext) UserName() string {
	return req.Username
}

// UserRole returns user role
func (req *TestRequestContext) UserRole() string {
	return req.Userrole
}

// ActorID returns the actor ID for this request
func (req *TestRequestContext) ActorID() int64 {
	return req.Actorid
}

// ActorName returns the actor name for this request
func (req *TestRequestContext) ActorName() string {
	return req.Actorname
}

// Token returns the token for this request
func (req *TestRequestContext) Token() string {
	return ""
}

// SetUserId sets the User id for this request
func (req *TestRequestContext) SetUserId(id int) {
	req.Userid = id
}

// Xid returns request transaction id
func (req *TestRequestContext) Xid() string {
	if req.xid == "" {
		req.xid = xid.UniqueIdGenerator()
	}
	return req.xid
}

// Tx creates and returns an open database transaction
func (req *TestRequestContext) Tx() *sqlx.Tx {
	var err error

	if req.tx == nil {
		if req.appContext.Db == nil {
			req.appContext.Db, err = dbhandler.CreateConnection()
			if err != nil {
				println("ERROR CREATING DB CONNECTION", err.Error())
			}
		}
		req.tx, err = dbhandler.CreateTx(&req.appContext)
		if err != nil {
			println("CREATING DB TX fail  ", err.Error())

		}
	}
	return req.tx
}

// VerifyParams ensures the code asked for all parameters in the route pattern
func (req *TestRequestContext) VerifyParams(action string, routes []common.Route) bool {
	paramRegex := regexp.MustCompile(`\{[^/]*\}`)
	nameRegex := regexp.MustCompile(`[a-z]+`)
	for _, r := range routes {
		if r.Name == action {
			matches := paramRegex.FindAll([]byte(r.Pattern), -1)
			for _, m := range matches {
				name := nameRegex.Find(m)
				if name != nil && req.WasValueRequested(string(name)) == false {
					return false
				}
			}
		}
	}
	return true
}

// WasValueRequested returns whether the named value was requested from the request object
func (req *TestRequestContext) WasValueRequested(name string) bool {
	for _, v := range req.valuesRequested {
		if v == name {
			return true
		}
	}
	return false
}

// Value returns string value from parameters
func (req *TestRequestContext) Value(name string, defValue string) string {
	req.valuesRequested = append(req.valuesRequested, name)
	if val, present := req.Values[name]; present {
		return val
	}
	return defValue
}

// BoolValue returns bool value from parameters
func (req *TestRequestContext) BoolValue(name string, defValue bool) bool {
	val := req.Value(name, "default")
	if val == "default" {
		return defValue
	}
	if iVal, err := strconv.ParseBool(val); err == nil {
		return iVal
	}
	return defValue
}

// IntValue returns int64 value from parameters
func (req *TestRequestContext) IntValue(name string, defValue int64) int64 {
	val := req.Value(name, "default")
	if val == "default" {
		return defValue
	}
	if iVal, err := strconv.ParseInt(val, 10, 64); err == nil {
		return iVal
	}
	return defValue
}

// IntValue32 returns the named input variable as an int64
func (hrc *TestRequestContext) IntValue32(name string, defValue int) int {
	val := hrc.Value(name, "default")
	if val == "default" {
		return defValue
	}
	if iVal, err := strconv.Atoi(val); err == nil {
		return iVal
	}
	return defValue
}

// FloatValue returns the named input variable as a float64
func (req *TestRequestContext) FloatValue(name string, defValue float64) float64 {
	val := req.Value(name, "default")
	if val == "default" {
		return defValue
	}
	if fVal, err := strconv.ParseFloat(val, 64); err == nil {
		return fVal
	}
	return defValue
}

// TimeValue returns the named input variable as a time.Time
func (req *TestRequestContext) TimeValue(name string, defValue time.Time) time.Time {
	val := req.Value(name, "default")
	if val == "default" {
		return defValue
	}
	if fVal, err := utilities.ParseStringtoTime(val); err == nil {
		return fVal
	}
	return defValue
}

// Scan parses the scan source into the destination
func (req *TestRequestContext) Scan(name string, dest interface{}) error {
	return json.Unmarshal([]byte(req.ScanSource), dest)
}

// AddLogValue adds a key/value pair to be logged
func (req *TestRequestContext) AddLogValue(name string, value interface{}) {
	(*req.LogValues())[name] = value
}

// LogValues returns a map of key/values to be logged
func (req *TestRequestContext) LogValues() *map[string]interface{} {
	if req.logValues == nil {
		req.logValues = map[string]interface{}{}
	}
	return &req.logValues
}

// RequestBody returns request body
func (req *TestRequestContext) RequestBody() []byte {
	return []byte{}
}

// GetPedigree implements the RequestContext method; this is copied
// from the 'superclass' because it's not possible for an 'inherited'
// method to call a 'subclass' method (go is not polymorphic)
// func (req *TestRequestContext) GetPedigree(getpf common.GetAuthInfoFunc) *model.Pedigree {
// 	if req.DefaultPedigree == nil {
// 		pedigree := getpf(req.GetContext())
// 		pedigree.ActorID = req.ActorID()
// 		pedigree.ActorName = req.ActorName()
// 		pedigree.Xid = req.Xid()
// 		pedigree.Service = req.ServiceName()
// 		req.DefaultPedigree = pedigree
// 	}
// 	return req.DefaultPedigree
// }

// RequestUploadFiles return request uploaded files
func (req *TestRequestContext) RequestUploadFiles() []common.Upload {
	return []common.Upload{}
}

// SetMaxDeadlockRetries configures how many retries will be performed
func (req *TestRequestContext) SetMaxDeadlockRetries(max int) {
	req.MaxDeadlockRetries = max
}

// BumpDeadlocks increments the deadlock count by 1 and returns
// whether retry still available
func (req *TestRequestContext) BumpDeadlocks() (retry bool) {
	req.Deadlockretrycount++
	return req.Deadlockretrycount <= req.MaxDeadlockRetries
}

// DeadlockRetryCount returns how many deadlocks occurred during this request
func (req *TestRequestContext) DeadlockRetryCount() int {
	return req.Deadlockretrycount
}

// RecordEvent records an event like RequestContextBase but also saves it for easy
// retrieval by tests
// func (req *TestRequestContext) RecordEvent(event common.SourceEvent) {
// 	req.RequestContextBase.RecordEvent(event)
// 	req.events = append(req.events, event)
// }

// GetRecordedEvents returns all the source events recorded since it was last called, if any.
// func (req *TestRequestContext) GetRecordedEvents() []common.SourceEvent {
// 	events := req.RequestContextBase.GetRecordedEvents()
// 	req.ClearRecordedEvents()
// 	return events
// }

// ResetForRetry clears out deferred functions, log fields and rolls back
// the transaction
func (req *TestRequestContext) ResetForRetry() (err error) {
	req.Tx().Rollback()
	req.tx, err = req.appContext.Db.Beginx()
	//req.RequestContextBase.ClearDeferredRequests()
	req.logValues = map[string]interface{}{}
	return
}

// SetIsRawResponse sets the IsRawResponse
func (req *TestRequestContext) SetIsRawResponse(isRawResponse bool, contentType string) {
	req.IsRawResponse = isRawResponse
	req.RawContentType = contentType
}

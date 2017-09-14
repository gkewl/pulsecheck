package errorhandler

import (
	"encoding/json"
	"fmt"
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/logger"
	"github.com/gkewl/pulsecheck/utilities"
	"github.com/gkewl/pulsecheck/xid"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/juju/errors"
)

const (
	area_validation = 1000
	area_processing = 2000
	area_database   = 3000
	area_json       = 4000
	area_http       = 5000

	fatal = 1

	VALIDATION_ERROR       = area_validation + 10
	PROCESSING_ERROR       = area_processing + 10
	DB_CONNECTFAILURE      = area_database + 10 + fatal
	DB_QUERYERROR          = area_database + 20
	DB_DATANOTFOUND        = area_database + 30
	JSON_PARSINGERROR      = area_json + 10 + fatal
	JSON_ENCODINGERROR     = area_json + 20 + fatal
	HTTP_MISSINGPARAMETERS = area_http + 10
	HTTP_INVALIDPARAMETERS = area_http + 20
)

var namedErrors = map[int]*NamedError{}

// NamedError holds the error code, short description, and HTTP status
type NamedError struct {
	codeval int
	desc    string
	status  int
}

// Code returns the code
func (e *NamedError) Code() int {
	return e.codeval
}

// Description returns the description
func (e *NamedError) Description() string {
	return e.desc
}

// String returns the code and description
func (e *NamedError) String() string {
	return fmt.Sprintf("%d - %s", e.codeval, e.desc)
}

// HTTPStatus returns the http status for this error
func (e *NamedError) HTTPStatus() int {
	return e.status
}

// newNamedError creates a new named error
func newNamedError(code int, description string, httpStatus int) *NamedError {
	if ne, exists := namedErrors[code]; exists {
		fmt.Printf("FATAL ERROR: duplicate named error code %d/%s: %v\n", code, description, ne)
		os.Exit(1)
	}
	ne := &NamedError{codeval: code, desc: description, status: httpStatus}
	namedErrors[code] = ne
	return ne
}

// Error wraps a NamedError with details and location
type Error struct {
	NamedError *NamedError
	Details    []string
	Locations  []string
}

// Error implements the standard error interface and just returns the
// named error NNNN - description
func (e Error) Error() string {
	return e.NamedError.String()
}

// DetailStack returns all details joined by the separator
func (e Error) DetailStack(sep string) string {
	return strings.Join(e.Details, sep)
}

// LocationStack returns all detail locations joined by the separator
func (e Error) LocationStack(sep string) string {
	return strings.Join(e.Locations, sep)
}

// NewError returns an Error based on a NamedError
func NewError(namedErr *NamedError, format string, args ...interface{}) Error {
	return Error{
		NamedError: namedErr,
		Details:    []string{fmt.Sprintf("%d: ", namedErr.Code()) + fmt.Sprintf(format, args...)},
		Locations:  []string{CallingFunction(1)},
	}
}

// NewErrorFromError copies details from other error with new named error
func NewErrorFromError(namedErr *NamedError, err Error) Error {
	return Error{
		NamedError: namedErr,
		Details:    err.Details,
		Locations:  err.Locations,
	}
}

// WrapError embeds details from other error into a new error
func WrapError(namedErr *NamedError, err error, format string, args ...interface{}) Error {
	e := Error{
		NamedError: namedErr,
		Details:    []string{fmt.Sprintf("%d: ", namedErr.Code()) + fmt.Sprintf(format, args...)},
		Locations:  []string{CallingFunction(1)},
	}
	if eherr, ok := err.(Error); ok {
		e.Details = append(e.Details, eherr.Details...)
		e.Locations = append(e.Locations, eherr.Locations...)
	} else if err != nil {
		e.Details = append(e.Details, err.Error())
		e.Locations = append(e.Locations, "unknown")
	}
	return e
}

// NewErrorNotFound can be called when data was not returned OR there was a DB error
func NewErrorNotFound(namedErr *NamedError, dberr error, format string, args ...interface{}) Error {
	detail := fmt.Sprintf(format, args...)
	if dberr != nil && NotNoRowsError(dberr) {
		detail = detail + " DB Error: " + dberr.Error()
	}
	return Error{
		NamedError: namedErr,
		Details:    []string{fmt.Sprintf("%d: ", namedErr.Code()) + detail},
		Locations:  []string{CallingFunction(1)},
	}
}

// HasNoRowsError returns true if the error or error details contain "sql no rows returned"
func HasNoRowsError(err error) bool {
	txt := err.Error()
	if eherr, ok := err.(Error); ok {
		txt = txt + " " + eherr.DetailStack(";")
	}
	return strings.Contains(txt, ErrDBNoRows.Error())
}

// NotNoRowsError returns true if the error does NOT contain "sql no rows returned"
func NotNoRowsError(err error) bool {
	return !HasNoRowsError(err)
}

// ContainsError returns true if this error or any of its details matches
func ContainsError(err error, namedErr *NamedError) bool {
	if eherr, ok := err.(Error); ok {
		if eherr.NamedError.Code() == namedErr.Code() {
			return true
		}
		check := fmt.Sprintf("%d: ", namedErr.Code())
		for _, d := range eherr.Details {
			if strings.HasPrefix(d, check) {
				return true
			}
		}
	}
	return false
}

// ContainsErrorText checks the details stack for the passed text
func ContainsErrorText(err error, text string) bool {
	if err == nil {
		return false
	}
	if eherr, ok := err.(Error); ok {
		return strings.Contains(eherr.DetailStack(";"), text)
	}
	return strings.Contains(err.Error(), text)
}

// AddDetail adds detail to the err and returns the enhanced error
func AddDetail(err error, format string, args ...interface{}) error {
	if eherr, ok := err.(Error); ok {
		eherr.Details = append(eherr.Details, fmt.Sprintf(format, args...))
		eherr.Locations = append(eherr.Locations, CallingFunction(1))
		return eherr
	}
	return errors.Annotatef(err, format, args...)
}

type clientError struct {
	ErrorInfo string
	Message   string
	Xid       string
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	err := common.ErrorResponse{Code: http.StatusNotFound, Message: "Page not found"}
	result := common.StructuredResponse{Xid: xid.UniqueIdGenerator(), StatusCode: http.StatusNotFound, Error: &err}
	_ = utilities.WriteJSON(w, http.StatusNotFound, result)
}

// RespWithError is the legacy routine called by controllers to send error back
// to HTTP client and log it
func RespWithError(w http.ResponseWriter, r *http.Request, input interface{}, err error) error {
	xid := utilities.GetXid(r)
	var statuscode int = http.StatusInternalServerError
	var content interface{}

	var locations string = ""
	var message string = ""

	eherr, ok := err.(Error)
	if ok {
		statuscode = eherr.NamedError.HTTPStatus()
		locations = eherr.LocationStack("; ")
		message = eherr.NamedError.Description()
		content = common.StructuredResponse{
			Xid:        xid,
			StatusCode: eherr.NamedError.HTTPStatus(),
			Error: &common.ErrorResponse{
				Code:      eherr.NamedError.Code(),
				Message:   eherr.NamedError.Description(),
				Details:   eherr.DetailStack("; "),
				Locations: eherr.LocationStack("; "),
			},
		}
	} else {
		content = common.StructuredResponse{
			Xid: xid,
			Error: &common.ErrorResponse{
				Message: err.Error(),
			},
		}
	}

	utilities.WriteJSON(w, statuscode, content)
	var inputStr string
	if !strings.Contains(r.URL.Path, "token-auth") {
		b, errj := json.Marshal(input)
		if errj == nil {
			inputStr = string(b)
		}
	}

	fields := map[string]interface{}{
		"protocol": "http",
		"method":   r.Method,
		"path":     r.URL.Path,
		"status":   statuscode,
	}

	logger.ErrorLog(logger.LogModel{
		Level:  logger.ERROR,
		Caller: locations,
		Msg:    message,
		Err:    err,
		Xid:    xid,
		Input:  inputStr,
		Fields: fields,
	})
	return err
}

var prefixPath = regexp.MustCompile(".*/src/github.com/gkewl/pulsecheck")

// CallingFunction returns a description of the calling function skip+1 back in the stack
func CallingFunction(skip int) string {
	pc, file, line, _ := runtime.Caller(skip + 1)
	return fmt.Sprintf("%s [%s:%d]", runtime.FuncForPC(pc).Name(), prefixPath.ReplaceAllString(file, ""), line)
}

//FullErrorText takes an error and returns its full string representation including
//all details from the stack if it is an Error struct
func FullErrorText(err error) string {
	errText := err.Error()
	if ehErr, ok := err.(Error); ok {
		errText += ". Details = [" + ehErr.DetailStack("; ") + "]"
	}
	return errText
}

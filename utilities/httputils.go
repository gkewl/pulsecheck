package utilities

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"

	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/constant"
)

// CheckForJSON makes sure that the request's Content-Type is application/ json.
func CheckForJSON(r *http.Request) error {
	ct := r.Header.Get("Content-Type")

	// No Content-Type header is ok as long as there's no Body
	if ct == "" && (r.Body == nil || r.ContentLength == 0) {
		return nil
	}

	// Otherwise it better be json
	if MatchesContentType(ct, "application/json") {
		return nil
	}
	return fmt.Errorf("Content-Type specified (%s) must be 'application/json'", ct)
}

// ParseForm ensures the request form is parsed even with invalid content types.
// If we don't do this, POST method without Content-type (even with empty body) will fail.
func ParseForm(r *http.Request) error {
	if r == nil {
		return nil
	}
	if err := r.ParseForm(); err != nil && !strings.HasPrefix(err.Error(), "mime:") {
		return err
	}
	return nil
}

// WriteJSON writes the value v to the http response stream as json with standard json encoding.
func WriteJSON(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Connection", "close")
	w.Header().Set("Content-Type", "application/json")
	if sresp, ok := v.(common.StructuredResponse); ok {
		w.Header().Set("X-XID", sresp.Xid)
	}
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(code)
	w.Write(data)

	return nil
}

// WriteJSONStructuredResponse writes the value v to the http response stream as json with standard json encoding.
func WriteJSONStructuredResponse(r *http.Request, w http.ResponseWriter, code int, v interface{}) error {
	xid := GetXid(r)
	result := common.StructuredResponse{Xid: xid, StatusCode: code, Response: v}

	w.Header().Set("Connection", "close")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-XID", xid)
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(result)
}

// GetBytes -
func GetBytes(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// WriteRawResponse writes the value v to the http response stream
func WriteRawResponse(w http.ResponseWriter, code int, contentType string, v interface{}) error {
	w.Header().Set("Connection", "close")
	w.Header().Set("Content-Type", contentType)

	var data []byte
	var err error
	switch dataType := v.(type) {
	case []byte:
		data = v.([]byte)
	default:
		data, err = GetBytes(v)
		if err != nil {
			logrus.Errorf("Error converting datatype %v to byte array for WriteRawResponse error: %v", dataType, err)
		}
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(code)
	w.Write(data)
	return err
}

// MatchesContentType validates the content type against the expected one
func MatchesContentType(contentType, expectedType string) bool {
	mimetype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		logrus.Errorf("Error parsing media type: %s error: %v", contentType, err)
	}
	return err == nil && mimetype == expectedType
}

// SendToHTTP -
func SendToHTTP(url string, method string, data interface{}, token string) (*http.Response, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	reader := strings.NewReader(string(b))

	request, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	res, err := http.DefaultClient.Do(request)
	return res, err
}

// GetXid -
func GetXid(r *http.Request) string {
	xid := r.Context().Value(constant.Xid)
	if xid != nil {
		return xid.(string)
	}
	return ""
}

package utilities

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

// HTTPClientPerformer provides a rudimentary interface to wrap external http calls
type HTTPClientPerformer interface {
	Post(xid string, url string, content interface{}) (response []byte, httpStatusCode int, err error)
}

// JSONRestHTTPClient performs http operations that expect and return JSON
type JSONRestHTTPClient struct {
}

// Post performs an http post to the given URL
func (c JSONRestHTTPClient) Post(xid string, url string, content interface{}) (response []byte, httpStatusCode int, err error) {
	var output []byte
	var bodyStr = ""

	if content != nil {
		body, err := json.Marshal(content)
		if err != nil {
			return nil, -1, errors.New("JSON Marshal error")
		}
		bodyStr = string(body)
	}

	reader := strings.NewReader(bodyStr)
	req, err := http.NewRequest("post", url, reader)
	req.Header.Set("Content-type", "application/json")

	tr := &http.Transport{TLSClientConfig: &tls.Config{}}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)

	if err != nil {
		return nil, -1, errors.New("Endpoint is not available, Error: " + err.Error())
	}

	defer func() { _ = resp.Body.Close() }()

	output, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, resp.StatusCode, errors.New("Unable to read response body, Error: " + err.Error())
	}

	return output, resp.StatusCode, nil
}

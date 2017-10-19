package protocol

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"runtime"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	auth "github.com/gkewl/pulsecheck/authentication"
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/model"
)

// HTTPTestCaller implements the TestCaller interface for making HTTP requests
type HTTPTestCaller struct {
	BaseURL string
	Router  common.Router
}

// Protocol reports the protocol for this test caller
func (c *HTTPTestCaller) Protocol() string {
	return "http"
}

// MakeTestCall makes an HTTP REST call to its configured URL, confirms the response
// is successful or meets configured errors
func (c *HTTPTestCaller) MakeTestCall(config *TestConfig) (responseBody string) {
	// note caller info in case of failed expectation in this method
	_, fn, line, _ := runtime.Caller(1)
	By(fmt.Sprintf("Test HTTP call from %s:%d", fn, line))

	// find route information
	route := c.Router.RouteForAction(config.Action)
	Expect(route).ToNot(BeNil())

	// set up body reader
	reader := strings.NewReader("")
	if config.Content != nil {
		if str, ok := config.Content.(string); ok {
			reader = strings.NewReader(str)
		} else {
			j, err := json.Marshal(config.Content)
			Expect(err).To(BeNil())
			reader = strings.NewReader(string(j))
		}
	}

	// construct request to path
	path, remaining := makePath(config, route)
	urlParams := []string{}
	for _, p := range remaining {
		urlParams = append(urlParams, fmt.Sprintf("%s=%s", p, url.QueryEscape(config.Params[p])))
	}
	request, err := http.NewRequest(route.Method, fmt.Sprintf("%s%s?%s", c.BaseURL, path, strings.Join(urlParams, "&")), reader)
	Expect(err).To(BeNil())
	request.Header.Set("Content-type", "application/json")

	// set up authentication

	authBackend := auth.InitJWTAuthenticationBackend()
	token, _ := authBackend.GenerateToken("rgunari@gmail.com", model.UserCompany{UserID: 1, CompanyID: 1}, "USER")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token.Token))

	response, err := http.DefaultClient.Do(request)
	Expect(err).To(BeNil())
	expectedStatus := config.ExpectedHTTPStatus
	if expectedStatus == 0 {
		expectedStatus = 200
	}
	Expect(response.StatusCode).To(Equal(expectedStatus))
	if response.ContentLength >= -1 {
		defer func() { _ = response.Body.Close() }()
		resBody, err := ioutil.ReadAll(response.Body)
		Expect(err).To(BeNil())
		responseBody = string(resBody)
		if config.ExpectedError != "" {
			Expect(responseBody).To(ContainSubstring(config.ExpectedError))
		}
	}
	return
}

// makePath returns the route path with substitution for {xxx} identifiers
func makePath(config *TestConfig, route *common.Route) (path string, remainingParams []string) {
	path = route.Pattern
	usedParams := []string{}
	regex := regexp.MustCompile(`\{[^/]*\}`)
	remove := regexp.MustCompile(`\{|\}|:.*`)
	matches := regex.FindAll([]byte(path), -1)
	if matches != nil {
		for _, m := range matches {
			p := remove.ReplaceAllString(string(m), "")
			sub, found := config.Params[p]
			Expect(found).To(BeTrue(), "parameter %s not found in test config", p)
			path = strings.Replace(path, string(m), sub, -1)
			usedParams = append(usedParams, p)
		}
	}
	usedParamsString := strings.Join(usedParams, "|")
	for p, _ := range config.Params {
		if !strings.Contains(usedParamsString, p) {
			remainingParams = append(remainingParams, p)
		}
	}
	return
}

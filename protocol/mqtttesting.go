package protocol

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/gkewl/pulsecheck/constant"
	"github.com/gkewl/pulsecheck/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	auth "github.com/gkewl/pulsecheck/authentication"
	"github.com/gkewl/pulsecheck/common"
)

// MQTTTestCaller implements the TestCaller interface for making MQTT requests
type MQTTTestCaller struct {
	Client *MQTTClient
}

// Protocol reports the protocol for this test caller
func (c *MQTTTestCaller) Protocol() string {
	return "mqtt"
}

// MakeTestCall makes an MQTT call to its configured client, confirms the response
// is successful or meets configured errors, and returns the response
func (c *MQTTTestCaller) MakeTestCall(config *TestConfig) string {
	// note caller info in case of failed expectation in this method
	_, fn, line, _ := runtime.Caller(1)
	By(fmt.Sprintf("Test MQTT call from %s:%d", fn, line))

	// find route information
	route := c.Client.Router.RouteForAction(config.Action)
	Expect(route).ToNot(BeNil())

	request := map[string]interface{}{}
	// set up authentication
	meta := map[string]interface{}{}
	if config.AuthLevel != constant.Guest {
		authBackend := auth.InitJWTAuthenticationBackend()
		token, _ := authBackend.GenerateToken("rgunari@gmail.com", model.UserCompany{UserID: 1, CompanyID: 1}, "USER")
		meta["token"] = token.Token
	}
	meta["action"] = config.Action
	request["_meta_"] = meta
	for k, v := range config.Params {
		request[k] = v
	}
	if config.Content != nil {
		if _, ok := config.Content.(string); ok {
			var content interface{}
			_ = json.Unmarshal([]byte(config.Content.(string)), &content)
			request[config.ContentName] = content
		} else {
			request[config.ContentName] = config.Content
		}
	}
	payload, err := json.Marshal(request)
	Expect(err).To(BeNil())

	_, replyPayload := c.Client.Receive("/actor/somemachine/out", string(payload))
	var reply common.StructuredResponse
	err = json.Unmarshal([]byte(replyPayload), &reply)
	Expect(err).To(BeNil())
	if config.ExpectedError != "" {
		Expect(reply.Error.Message).To(ContainSubstring(config.ExpectedError))
	}
	return replyPayload
}

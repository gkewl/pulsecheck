package actor_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gkewl/pulsecheck/apis/actor"
	"github.com/gkewl/pulsecheck/constant"
	"github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/protocol"
)

var httpCaller protocol.HTTPTestCaller
var mqttCaller protocol.MQTTTestCaller

func TestSuite(t *testing.T) {
	httpCaller = protocol.HTTPTestCaller{BaseURL: baseUrl, Router: router}
	//mqttCaller = protocol.MQTTTestCaller{Client: mqttClient}

	RegisterFailHandler(Fail)
	RunSpecs(t, "Actor Test Suite")
}

var _ = Describe("Actor API protocol tests", func() {
	var mockBL actor.MockBL
	var noParams = map[string]string{}
	var nameParams = map[string]string{"name": "foo"}
	var idParams = map[string]string{"id": "42"}
	var termParams = map[string]string{"type": "USER", "term": "foo"}
	var usertypeParams = map[string]string{"type": "USER"}
	var fakeActor = `{"name": "foo"}`

	BeforeEach(func() {
		mockBL = actor.MockBL{}
		actor.TestingBizLogic = &mockBL
	})

	AfterEach(func() {
		actor.TestingBizLogic = nil
	})

	var allTests = func(caller protocol.TestCaller) {
		It("routes get", func() {
			config := protocol.MakeTestConfig("GetActor", constant.Guest, 200, nil, nameParams, "actor", nil)
			caller.MakeTestCall(config)
			Expect(mockBL.Name).To(Equal("foo"))
		})

		It("routes get all", func() {
			config := protocol.MakeTestConfig("GetActors", constant.Guest, 200, nil, usertypeParams, "actor", nil)
			caller.MakeTestCall(config)
		})

		It("routes create", func() {
			config := protocol.MakeTestConfig("CreateActor", constant.Superuser, 201, nil, noParams, "actor", fakeActor)
			body := caller.MakeTestCall(config)
			Expect(body).To(ContainSubstring("foo"))
		})

		It("detects insufficent auth to create", func() {
			config := protocol.MakeTestConfig("CreateActor", constant.Operator, 401, errorhandler.ErrUnAuthorizedUserForAPI, noParams, "actor", nil)
			caller.MakeTestCall(config)
		})

		It("routes update", func() {
			config := protocol.MakeTestConfig("UpdateActor", constant.Superuser, 200, nil, idParams, "actor", fakeActor)
			caller.MakeTestCall(config)
			Expect(mockBL.ID).To(Equal(int64(42)))
		})

		It("detects insufficent auth to update", func() {
			config := protocol.MakeTestConfig("UpdateActor", constant.Operator, 401, errorhandler.ErrUnAuthorizedUserForAPI, idParams, "actor", nil)
			caller.MakeTestCall(config)
		})

		It("routes delete", func() {
			config := protocol.MakeTestConfig("DeleteActor", constant.Superuser, 200, nil, idParams, "actor", nil)
			caller.MakeTestCall(config)
			Expect(mockBL.ID).To(Equal(int64(42)))
		})

		It("detects insufficent auth to delete", func() {
			config := protocol.MakeTestConfig("DeleteActor", constant.Operator, 401, errorhandler.ErrUnAuthorizedUserForAPI, idParams, "actor", nil)
			caller.MakeTestCall(config)
		})

		It("routes search", func() {
			config := protocol.MakeTestConfig("SearchActor", constant.Guest, 200, nil, termParams, "actor", nil)
			caller.MakeTestCall(config)
			Expect(mockBL.Term).To(Equal("foo"))
			Expect(mockBL.Type).To(Equal("USER"))
		})

		It("detects bad json input (http only)", func() {
			if caller.Protocol() == "http" {
				config := protocol.MakeTestConfig("UpdateActor", constant.Superuser, 400, errorhandler.ErrJsonDecodeFail, idParams, "actor", `{invalid:json}`)
				caller.MakeTestCall(config)
			}
		})

	}
	allTests(&httpCaller)
	//allTests(&mqttCaller)
})

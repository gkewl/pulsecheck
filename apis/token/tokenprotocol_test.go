package token_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gkewl/pulsecheck/apis/token"

	"github.com/gkewl/pulsecheck/constant"
	"github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/protocol"
)

var _ = Describe("Token API protocol tests", func() {
	var mockBL token.MockBLToken
	var getParams = map[string]string{"id": "42"}
	var getNameParams = map[string]string{"name": "foo"}
	var limitParams = map[string]string{"limit": "42"}
	var fakeToken = `{"actorid": 42}`

	BeforeEach(func() {
		mockBL = token.MockBLToken{}
		token.TestingBizLogic = &mockBL
	})

	AfterEach(func() {
		token.TestingBizLogic = nil
	})

	var allTests = func(caller protocol.TestCaller) {
		It("routes get", func() {
			config := protocol.MakeTestConfig("GetToken", constant.Superuser, 200, nil, getParams, "tk", nil)
			caller.MakeTestCall(config)
			Expect(mockBL.ID).To(Equal(int64(42)))
		})
		It("routes get by token string", func() {
			config := protocol.MakeTestConfig("GetTokenbyString", constant.Superuser, 200, nil, getNameParams, "tk", nil)
			caller.MakeTestCall(config)
			Expect(mockBL.TokenString).To(Equal("foo"))
		})

		It("routes get all", func() {
			config := protocol.MakeTestConfig("GetTokens", constant.Superuser, 200, nil, limitParams, "tk", nil)
			caller.MakeTestCall(config)
			Expect(mockBL.Limit).To(Equal(int64(42)))
		})

		It("routes update", func() {
			config := protocol.MakeTestConfig("UpdateToken", constant.Superuser, 200, nil, getParams, "tk", fakeToken)
			caller.MakeTestCall(config)
			Expect(mockBL.ID).To(Equal(int64(42)))
		})

		It("detects insufficent auth to update", func() {
			config := protocol.MakeTestConfig("UpdateToken", constant.Guest, 401, errorhandler.ErrUnAuthorizedUserForAPI, getParams, "tk", nil)
			caller.MakeTestCall(config)
		})

		It("detects bad json input (http only)", func() {
			if caller.Protocol() == "http" {
				config := protocol.MakeTestConfig("UpdateToken", constant.Superuser, 400, errorhandler.ErrJsonDecodeFail, getParams, "tk", `{invalid:json`)
				caller.MakeTestCall(config)
			}
		})
	}
	allTests(&httpCaller)
	allTests(&mqttCaller)
})

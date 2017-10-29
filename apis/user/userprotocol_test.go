package user_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gkewl/pulsecheck/apis/user"

	"github.com/gkewl/pulsecheck/constant"
	"github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/protocol"
)

var _ = Describe("User API protocol tests", func() {
	var mockBL user.MockBLUser
	var noParams = map[string]string{}
	var getParams = map[string]string{"id": "42"}
	var limitParams = map[string]string{"limit": "42"}
	var fakeUser = `{"email": "test@email.com"}`

	BeforeEach(func() {
		mockBL = user.MockBLUser{}
		user.TestingBizLogic = &mockBL
	})

	AfterEach(func() {
		user.TestingBizLogic = nil
	})

	var allTests = func(caller protocol.TestCaller) {
		It("routes get", func() {
			config := protocol.MakeTestConfig("GetUser", constant.User, 200, nil, getParams, "usr", nil)
			caller.MakeTestCall(config)
			Expect(mockBL.ID).To(Equal(int(42)))
		})

		// codegen todo: reference and unique key gets

		It("routes get all", func() {
			config := protocol.MakeTestConfig("GetUsers", constant.User, 200, nil, limitParams, "usr", nil)
			caller.MakeTestCall(config)
			Expect(mockBL.Limit).To(Equal(int64(42)))
		})

		It("routes create", func() {
			config := protocol.MakeTestConfig("CreateUser", constant.Superuser, 201, nil, noParams, "usr", fakeUser)
			body := caller.MakeTestCall(config)
			Expect(body).To(ContainSubstring("test@email.com"))
		})

		It("detects insufficent auth to create", func() {
			config := protocol.MakeTestConfig("CreateUser", constant.Guest, 401, errorhandler.ErrUnAuthorizedUserForAPI, noParams, "usr", nil)
			caller.MakeTestCall(config)
		})

		It("routes update", func() {
			config := protocol.MakeTestConfig("UpdateUser", constant.User, 200, nil, getParams, "usr", fakeUser)
			caller.MakeTestCall(config)
			Expect(mockBL.ID).To(Equal(int(42)))
		})

		It("detects insufficent auth to update", func() {
			config := protocol.MakeTestConfig("UpdateUser", constant.Guest, 401, errorhandler.ErrUnAuthorizedUserForAPI, getParams, "usr", nil)
			caller.MakeTestCall(config)
		})

		It("routes delete", func() {
			config := protocol.MakeTestConfig("DeleteUser", constant.Admin, 200, nil, getParams, "usr", nil)
			caller.MakeTestCall(config)
			Expect(mockBL.ID).To(Equal(int(42)))
		})

		It("detects insufficent auth to delete", func() {
			config := protocol.MakeTestConfig("DeleteUser", constant.Guest, 401, errorhandler.ErrUnAuthorizedUserForAPI, getParams, "usr", nil)
			caller.MakeTestCall(config)
		})

		It("detects bad json input (http only)", func() {
			if caller.Protocol() == "http" {
				config := protocol.MakeTestConfig("UpdateUser", constant.User, 400, errorhandler.ErrJsonDecodeFail, getParams, "usr", `{invalid:json}`)
				caller.MakeTestCall(config)
			}
		})
	}
	allTests(&httpCaller)
	//	allTests(&mqttCaller)
})

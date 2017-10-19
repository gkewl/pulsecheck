package company_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gkewl/pulsecheck/apis/company"

	"github.com/gkewl/pulsecheck/constant"
	"github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/protocol"
)

var _ = Describe("Company API protocol tests", func() {
	var mockBL company.MockBLCompany
	var noParams = map[string]string{}
	var getParams = map[string]string{"id": "42"}
	var limitParams = map[string]string{"limit": "42"}
	var fakeCompany = `{"name": "foo"}`

	BeforeEach(func() {
		mockBL = company.MockBLCompany{}
		company.TestingBizLogic = &mockBL
	})

	AfterEach(func() {
		company.TestingBizLogic = nil
	})

	var allTests = func(caller protocol.TestCaller) {
		It("routes get", func() {
			config := protocol.MakeTestConfig("GetCompany", constant.Guest, 200, nil, getParams, "comp", nil)
			caller.MakeTestCall(config)
			Expect(mockBL.ID).To(Equal(int(42)))
		})

		// codegen todo: reference and unique key gets

		It("routes get all", func() {
			config := protocol.MakeTestConfig("GetCompanys", constant.Guest, 200, nil, limitParams, "comp", nil)
			caller.MakeTestCall(config)
			Expect(mockBL.Limit).To(Equal(int(42)))
		})

		It("routes create", func() {
			config := protocol.MakeTestConfig("CreateCompany", constant.Guest, 201, nil, noParams, "comp", fakeCompany)
			body := caller.MakeTestCall(config)
			Expect(body).To(ContainSubstring("foo"))
		})

		// It("detects insufficent auth to create", func() {
		// 	config := protocol.MakeTestConfig("CreateCompany", constant.Guest, 401, errorhandler.ErrUnAuthorizedUserForAPI, noParams, "comp", nil)
		// 	caller.MakeTestCall(config)
		// })

		It("routes update", func() {
			config := protocol.MakeTestConfig("UpdateCompany", constant.Guest, 200, nil, getParams, "comp", fakeCompany)
			caller.MakeTestCall(config)
			Expect(mockBL.ID).To(Equal(int(42)))
		})

		// It("detects insufficent auth to update", func() {
		// 	config := protocol.MakeTestConfig("UpdateCompany", constant.Guest, 401, errorhandler.ErrUnAuthorizedUserForAPI, getParams, "comp", nil)
		// 	caller.MakeTestCall(config)
		// })

		It("routes delete", func() {
			config := protocol.MakeTestConfig("DeleteCompany", constant.Guest, 200, nil, getParams, "comp", nil)
			caller.MakeTestCall(config)
			Expect(mockBL.ID).To(Equal(int(42)))
		})

		// It("detects insufficent auth to delete", func() {
		// 	config := protocol.MakeTestConfig("DeleteCompany", constant.Guest, 401, errorhandler.ErrUnAuthorizedUserForAPI, getParams, "comp", nil)
		// 	caller.MakeTestCall(config)
		// })

		It("detects bad json input (http only)", func() {
			if caller.Protocol() == "http" {
				config := protocol.MakeTestConfig("UpdateCompany", constant.Guest, 400, errorhandler.ErrJsonDecodeFail, getParams, "comp", `{invalid:json}`)
				caller.MakeTestCall(config)
			}
		})
	}
	allTests(&httpCaller)
	allTests(&mqttCaller)
})

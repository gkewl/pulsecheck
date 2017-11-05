package employee_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gkewl/pulsecheck/apis/employee"

	"github.com/gkewl/pulsecheck/constant"
	"github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/protocol"
)

var _ = Describe("Employee API protocol tests", func() {
	var mockBL employee.MockBLEmployee
	var noParams = map[string]string{}
	var getParams = map[string]string{"id": "42"}
	var limitParams = map[string]string{"limit": "42"}
	var fakeEmployee = `{"firstname": "foo","dateofbirth":"2017-10-28"}`
	var fakeEmployees = `[{"firstname": "foo1"},{"firstname": "foo2"}]`

	BeforeEach(func() {
		mockBL = employee.MockBLEmployee{}
		employee.TestingBizLogic = &mockBL
	})

	AfterEach(func() {
		employee.TestingBizLogic = nil
	})

	var allTests = func(caller protocol.TestCaller) {
		It("routes get", func() {
			config := protocol.MakeTestConfig("GetEmployee", constant.User, 200, nil, getParams, "emp", nil)
			caller.MakeTestCall(config)
			Expect(mockBL.ID).To(Equal(int64(42)))
		})

		It("routes get all", func() {
			config := protocol.MakeTestConfig("GetEmployees", constant.User, 200, nil, limitParams, "emp", nil)
			caller.MakeTestCall(config)
			Expect(mockBL.Limit).To(Equal(int64(42)))
		})

		It("routes create", func() {
			config := protocol.MakeTestConfig("CreateEmployee", constant.User, 201, nil, noParams, "emp", fakeEmployee)
			body := caller.MakeTestCall(config)
			Expect(body).To(ContainSubstring("foo"))
			Expect(body).To(ContainSubstring("2017-10-28"))
			Expect(mockBL.DOB).To(Equal("2017-10-28"))

		})

		It("detects insufficent auth to create", func() {
			config := protocol.MakeTestConfig("CreateEmployee", constant.Guest, 401, errorhandler.ErrUnAuthorizedUserForAPI, noParams, "emp", nil)
			caller.MakeTestCall(config)
		})

		It("routes update", func() {
			config := protocol.MakeTestConfig("UpdateEmployee", constant.User, 200, nil, getParams, "emp", fakeEmployee)
			caller.MakeTestCall(config)
			Expect(mockBL.ID).To(Equal(int64(42)))
		})

		It("detects insufficent auth to update", func() {
			config := protocol.MakeTestConfig("UpdateEmployee", constant.Guest, 401, errorhandler.ErrUnAuthorizedUserForAPI, getParams, "emp", nil)
			caller.MakeTestCall(config)
		})

		It("routes delete", func() {
			config := protocol.MakeTestConfig("DeleteEmployee", constant.User, 200, nil, getParams, "emp", nil)
			caller.MakeTestCall(config)
			Expect(mockBL.ID).To(Equal(int64(42)))
		})

		It("detects insufficent auth to delete", func() {
			config := protocol.MakeTestConfig("DeleteEmployee", constant.Guest, 401, errorhandler.ErrUnAuthorizedUserForAPI, getParams, "emp", nil)
			caller.MakeTestCall(config)
		})

		It("detects bad json input (http only)", func() {
			if caller.Protocol() == "http" {
				config := protocol.MakeTestConfig("UpdateEmployee", constant.User, 400, errorhandler.ErrJsonDecodeFail, getParams, "emp", `{invalid:json}`)
				caller.MakeTestCall(config)
			}
		})
		It("routes Upload employees", func() {
			config := protocol.MakeTestConfig("UploadEmployees", constant.User, 200, nil, noParams, "emp", fakeEmployees)
			body := caller.MakeTestCall(config)
			Expect(body).To(ContainSubstring("foo1"))
			Expect(body).To(ContainSubstring("foo2"))
		})
		It("routes searchstatus", func() {
			config := protocol.MakeTestConfig("SearchStatus", constant.User, 200, nil, noParams, "emp", fakeEmployee)
			body := caller.MakeTestCall(config)
			Expect(body).To(ContainSubstring("foo"))
			Expect(body).To(ContainSubstring("2017-10-28"))
			Expect(mockBL.DOB).To(Equal("2017-10-28"))

		})

	}
	allTests(&httpCaller)
	//	allTests(&mqttCaller)
})

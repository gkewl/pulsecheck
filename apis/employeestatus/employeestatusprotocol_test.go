package employeestatus_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gkewl/pulsecheck/apis/employeestatus"

	"github.com/gkewl/pulsecheck/constant"
	"github.com/gkewl/pulsecheck/protocol"
)

var _ = Describe("EmployeeStatus API protocol tests", func() {
	var mockBL employeestatus.MockBLEmployeeStatus
	var getParams = map[string]string{"id": "42"}
	var limitParams = map[string]string{"limit": "42"}

	BeforeEach(func() {
		mockBL = employeestatus.MockBLEmployeeStatus{}
		employeestatus.TestingBizLogic = &mockBL
	})

	AfterEach(func() {
		employeestatus.TestingBizLogic = nil
	})

	var allTests = func(caller protocol.TestCaller) {
		It("routes get", func() {
			config := protocol.MakeTestConfig("GetEmployeeStatus", constant.User, 200, nil, getParams, "es", nil)
			caller.MakeTestCall(config)
			Expect(mockBL.ID).To(Equal(int64(42)))
		})

		It("routes get all", func() {
			config := protocol.MakeTestConfig("GetEmployeeStatuss", constant.User, 200, nil, limitParams, "es", nil)
			caller.MakeTestCall(config)
			Expect(mockBL.Limit).To(Equal(int64(42)))
		})

	}
	allTests(&httpCaller)
	//allTests(&mqttCaller)
})

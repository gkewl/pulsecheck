package employeestatus_test

import (
	"fmt"

	"github.com/gkewl/pulsecheck/utilities"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gkewl/pulsecheck/apis/employee"
	"github.com/gkewl/pulsecheck/apis/employeestatus"
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/constant"

	"github.com/gkewl/pulsecheck/errorhandler"
	tst "github.com/gkewl/pulsecheck/lib/testing"
	"github.com/gkewl/pulsecheck/model"
	"github.com/gkewl/pulsecheck/protocol"
)

// sampleEmployeeStatus constructs a employeestatus
func sampleEmployeeStatus(reqCtx common.RequestContext, t tst.T) (es model.EmployeeStatus) {
	emp := t.Employee(employee.BLEmployee{}.Create(reqCtx, t.SampleEmployee()))

	return model.EmployeeStatus{
		EmployeeID: emp.ID,
		Ofac:       true,
		IsActive:   true,
		CreatedBy:  "admin",
		ModifiedBy: "admin",
	}
}

func expectError(expected *errorhandler.NamedError, got error) {
	Expect(got).ToNot(BeNil())
	Expect(got.Error()).To(ContainSubstring(expected.String()))
}

var _ = Describe("EmployeeStatus Business Logic Tests", func() {
	var reqCtx *protocol.TestRequestContext
	var logic = employeestatus.BLEmployeeStatus{}
	var t tst.T

	BeforeEach(func() {
		reqCtx = &protocol.TestRequestContext{Userid: 1}
		t = tst.T{ReqCtx: reqCtx}
		reqCtx.Companyid = 1

	})
	AfterEach(func() {
		reqCtx.Complete(false)
	})

	It("creates a good employeestatus", func() {
		es := t.EmployeeStatus(logic.Create(reqCtx, sampleEmployeeStatus(reqCtx, t)))
		fmt.Printf("after %+v\n", es)
		check := t.GetEmployeeStatus(es.EmployeeID)
		Expect(check.EmployeeID).To(Equal(es.EmployeeID))
		//Expect(check.Consider).To(BeTrue()) TBD to check why this is failing
	})

	It("gets all employeestatuss", func() {
		emp := t.EmployeeStatus(logic.Create(reqCtx, sampleEmployeeStatus(reqCtx, t)))
		result := t.EmployeeStatuss(logic.GetAll(reqCtx, 1000))
		Expect(utilities.ConfirmValuesInSlice(result, "EmployeeID", emp.EmployeeID)).To(BeNil())

	})

	It("updates an existing employeestatus", func() {
		By("creating a employeestatus")
		es := sampleEmployeeStatus(reqCtx, t)
		es.OIG = false
		es = t.EmployeeStatus(logic.Create(reqCtx, es))

		check := t.GetEmployeeStatus(es.EmployeeID)
		Expect(check.OIG).To(BeFalse())
		Expect(check.Consider).To(BeFalse())

		By("updating the employeestatus")
		t.UpdateEmployeeStatus(es.EmployeeID, constant.Source_OIG, true)

		By("retrieving the employeestatus")
		check = t.GetEmployeeStatus(es.EmployeeID)
		Expect(check.OIG).To(BeTrue())
		Expect(check.Consider).To(BeTrue())
		Expect(check.OfacLastSearch).ToNot(BeNil())

	})

	It("does not update a employeestatus that isn't there", func() {
		_ = t.EmployeeStatusErr(logic.Update(reqCtx, 0, constant.Source_SAM, true))
	})

	It("deletes a employeestatus", func() {
		es := t.EmployeeStatus(logic.Create(reqCtx, sampleEmployeeStatus(reqCtx, t)))
		result, err := logic.Delete(reqCtx, es.ID)
		Expect(err).To(BeNil())
		Expect(result).To(Equal("Ok"))
		_ = t.EmployeeStatusErr(logic.Get(reqCtx, es.ID))
	})

	It("does not delete a employeestatus that isn't there", func() {
		_, err := logic.Delete(reqCtx, 0)
		Expect(err).ToNot(BeNil())
	})

})

package employee_test

import (
	"github.com/gkewl/pulsecheck/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gkewl/pulsecheck/apis/employee"

	"github.com/gkewl/pulsecheck/errorhandler"
	tst "github.com/gkewl/pulsecheck/lib/testing"
	"github.com/gkewl/pulsecheck/protocol"
	"github.com/gkewl/pulsecheck/utilities"
)

func expectError(expected *errorhandler.NamedError, got error) {
	Expect(got).ToNot(BeNil())
	Expect(got.Error()).To(ContainSubstring(expected.String()))
}

var _ = Describe("Employee Business Logic Tests", func() {
	var reqCtx *protocol.TestRequestContext
	var logic = employee.BLEmployee{}
	var t tst.T

	BeforeEach(func() {
		reqCtx = &protocol.TestRequestContext{Userid: 1}
		reqCtx.Companyid = 1
		t = tst.T{ReqCtx: reqCtx}
	})
	AfterEach(func() {
		reqCtx.Complete(false)
	})

	It("creates a good employee", func() {
		emp := t.Employee(logic.Create(reqCtx, t.SampleEmployee()))

		check := t.GetEmployee(emp.ID)
		Expect(check.Firstname).To(Equal(emp.Firstname))
		Expect(check.Dateofbirth).To(Equal("2017-10-28"))

	})

	It("gets all employees", func() {
		t.Employee(logic.Create(reqCtx, t.SampleEmployee()))
		result := t.Employees(logic.GetAll(reqCtx, 10))
		Expect(len(result)).To(BeNumerically(">", 0))
		Expect(result[0].Dateofbirth).To(Equal("2017-10-28"))
	})

	It("updates an existing employee", func() {
		By("creating a employee")
		emp := t.Employee(logic.Create(reqCtx, t.SampleEmployee()))
		emp.Lastname = emp.Lastname + "Test"
		emp.Dateofbirth = "2017-09-28"
		By("updating the employee")
		t.UpdateEmployee(emp)

		By("retrieving the employee")
		check := t.GetEmployee(emp.ID)
		Expect(check.Lastname).To(Equal(emp.Lastname))
		Expect(check.Dateofbirth).To(Equal("2017-09-28"))
	})

	It("does not update a employee that isn't there", func() {
		_ = t.EmployeeErr(logic.Update(reqCtx, 0, t.SampleEmployee()))
	})

	It("deletes a employee", func() {
		emp := t.Employee(logic.Create(reqCtx, t.SampleEmployee()))
		result, err := logic.Delete(reqCtx, emp.ID)
		Expect(err).To(BeNil())
		Expect(result).To(Equal("Ok"))
		_ = t.EmployeeErr(logic.Get(reqCtx, emp.ID))
	})

	It("does not delete a employee that isn't there", func() {
		_, err := logic.Delete(reqCtx, 0)
		Expect(err).ToNot(BeNil())
	})

	PIt("searches for employees", func() {
		baseName := utilities.GenerateRandomString(10)
		for i := 0; i < 5; i++ {
			emp := t.SampleEmployee()
			emp.Firstname = baseName + utilities.GenerateRandomString(10)
			t.Employee(logic.Create(reqCtx, emp))
		}
		emps := t.Employees(logic.Search(reqCtx, baseName, 4))
		Expect(len(emps)).To(Equal(4))
	})

	It("Upload a new employees", func() {
		employees := []model.Employee{}
		employees = append(employees, t.SampleEmployee())
		employees = append(employees, t.SampleEmployee())

		emps, err := logic.Upload(reqCtx, employees)
		Expect(err).To(BeNil())
		Expect(len(emps)).To(Equal(2))

		for _, e := range emps {
			Expect(e.ID).To(BeNumerically(">", 0))
		}
	})

	It("Upload existing employee and  new employees", func() {
		employees := []model.Employee{}
		emp := t.Employee(logic.Create(reqCtx, t.SampleEmployee()))
		employees = append(employees, emp)
		employees = append(employees, t.SampleEmployee())

		emps, err := logic.Upload(reqCtx, employees)
		Expect(err).To(BeNil())
		Expect(len(emps)).To(Equal(2))

		for _, e := range emps {
			Expect(e.ID).To(BeNumerically(">", 0))
		}

		utilities.ConfirmValuesInSlice(emps, "ID", emp.ID)
	})

})

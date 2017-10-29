package testing

import (
	//. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	null "gopkg.in/guregu/null.v3"

	"github.com/gkewl/pulsecheck/apis/employee"
	"github.com/gkewl/pulsecheck/model"
	"github.com/gkewl/pulsecheck/utilities"
)

// SampleEmployee -  constructs a employee
func (t *T) SampleEmployee() (emp model.Employee) {

	return model.Employee{
		CompanyID:   1,
		Firstname:   "UT_Emp_" + utilities.GenerateRandomString(8),
		Middlename:  null.StringFrom("foo"),
		Lastname:    "foo",
		Dateofbirth: "2017-10-28",
		Type:        1,
		IsActive:    true,
		CreatedBy:   "admin",
		ModifiedBy:  "admin",
	}
}

// Employee expects a employee and an error and verifies the error is nil
// and returns the employee
func (t *T) Employee(emp model.Employee, e error) model.Employee {
	Expect(e).To(BeNil(), callers())
	return emp
}

// GetEmployee fetches a employee using an ID and verifies no error
func (t *T) GetEmployee(ID int64) model.Employee {
	Expect(t.ReqCtx).ToNot(BeNil(), callers())
	emp, e := employee.BLEmployee{}.Get(t.ReqCtx, ID)
	Expect(e).To(BeNil(), callers())
	return emp
}

// ReGetEmployee expects a employee and an error and verifies the error is nil
// and re-gets the employee and returns it
func (t *T) ReGetEmployee(emp model.Employee, e error) model.Employee {
	Expect(e).To(BeNil(), callers())
	Expect(t.ReqCtx).ToNot(BeNil(), callers())
	emp, e = employee.BLEmployee{}.Get(t.ReqCtx, emp.ID)
	Expect(e).To(BeNil(), callers())

	return emp
}

// UpdateEmployee expects a employee struct and updates it, verifying no error
func (t *T) UpdateEmployee(emp model.Employee) model.Employee {
	Expect(t.ReqCtx).ToNot(BeNil(), callers())
	emp, e := employee.BLEmployee{}.Update(t.ReqCtx, emp.ID, emp)
	Expect(e).To(BeNil(), callers())
	return emp
}

// Employees expects a employee slice and an error and verifies the error is nil
// and returns the employees
func (t *T) Employees(emps []model.Employee, e error) []model.Employee {
	Expect(e).To(BeNil(), callers())
	return emps
}

// EmployeeErr expects a employee and an error and verifies the error is
// not nil and returns the error
func (t *T) EmployeeErr(emp model.Employee, e error) error {
	Expect(e).ToNot(BeNil(), callers())
	return e
}

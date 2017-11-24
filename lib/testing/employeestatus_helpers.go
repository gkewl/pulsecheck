package testing

import (
	. "github.com/onsi/gomega"
	//null "gopkg.in/guregu/null.v3"

	"github.com/gkewl/pulsecheck/apis/employeestatus"
	"github.com/gkewl/pulsecheck/model"
	//	"github.com/gkewl/pulsecheck/utilities"
)

// EmployeeStatus expects a employeestatus and an error and verifies the error is nil
// and returns the employeestatus
func (t *T) EmployeeStatus(es model.EmployeeStatus, e error) model.EmployeeStatus {
	Expect(e).To(BeNil(), callers())
	return es
}

// GetEmployeeStatus fetches a employeestatus using an ID and verifies no error
func (t *T) GetEmployeeStatus(ID int64) model.EmployeeStatus {
	Expect(t.ReqCtx).ToNot(BeNil(), callers())
	es, e := employeestatus.BLEmployeeStatus{}.Get(t.ReqCtx, ID)
	Expect(e).To(BeNil(), callers())
	return es
}

// ReGetEmployeeStatus expects a employeestatus and an error and verifies the error is nil
// and re-gets the employeestatus and returns it
func (t *T) ReGetEmployeeStatus(es model.EmployeeStatus, e error) model.EmployeeStatus {
	Expect(e).To(BeNil(), callers())
	Expect(t.ReqCtx).ToNot(BeNil(), callers())
	es, e = employeestatus.BLEmployeeStatus{}.Get(t.ReqCtx, es.ID)
	Expect(e).To(BeNil(), callers())
	return es
}

// UpdateEmployeeStatus expects a employeestatus struct and updates it, verifying no error
func (t *T) UpdateEmployeeStatus(employeeID int64, source string, value bool) model.EmployeeStatus {

	Expect(t.ReqCtx).ToNot(BeNil(), callers())
	es, e := employeestatus.BLEmployeeStatus{}.Update(t.ReqCtx, employeeID, source, value)
	Expect(e).To(BeNil(), callers())
	return es
}

// EmployeeStatuss expects a employeestatus slice and an error and verifies the error is nil
// and returns the employeestatuss
func (t *T) EmployeeStatuss(ess []model.EmployeeStatus, e error) []model.EmployeeStatus {
	Expect(e).To(BeNil(), callers())
	return ess
}

// EmployeeStatusErr expects a employeestatus and an error and verifies the error is
// not nil and returns the error
func (t *T) EmployeeStatusErr(es model.EmployeeStatus, e error) error {
	Expect(e).ToNot(BeNil(), callers())
	return e
}

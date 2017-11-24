package employeestatus

import (
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/model"
)

// MockBLEmployeeStatus is a struct for mocking employeestatus business logic methods
type MockBLEmployeeStatus struct {
	ID    int64
	Term  string
	Limit int64
	Err   error
}

// Create mocks create
func (bl *MockBLEmployeeStatus) Create(reqCtx common.RequestContext, es model.EmployeeStatus) (model.EmployeeStatus, error) {
	return es, bl.Err
}

// Get mocks get
func (bl *MockBLEmployeeStatus) Get(reqCtx common.RequestContext, employeeID int64) (es model.EmployeeStatus, err error) {
	bl.ID = employeeID
	return es, bl.Err
}

// GetAll mocks return all employeestatuss
func (bl *MockBLEmployeeStatus) GetAll(reqCtx common.RequestContext, limit int64) (ess []model.EmployeeStatus, err error) {
	bl.Limit = limit
	return ess, bl.Err
}

// Update mocks update of a single employeestatus
func (bl *MockBLEmployeeStatus) Update(reqCtx common.RequestContext, id int64, source string, value bool) (model.EmployeeStatus, error) {
	bl.ID = id
	return model.EmployeeStatus{}, bl.Err
}

// Delete mocks delete
func (bl *MockBLEmployeeStatus) Delete(reqCtx common.RequestContext, id int64) (string, error) {
	bl.ID = id
	return "ok", bl.Err
}

// // Search finds employeestatuss matching the term
// func (bl *MockBLEmployeeStatus) Search(reqCtx common.RequestContext, term string, limit int64) (ess []model.EmployeeStatus, err error) {
// 	bl.Term = term
// 	bl.Limit = limit
// 	return ess, bl.Err
// }

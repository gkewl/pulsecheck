package employee

import (
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/model"
)

// MockBLEmployee is a struct for mocking employee business logic methods
type MockBLEmployee struct {
	ID int64

	Term  string
	Limit int64
	Err   error
	DOB   string
}

// Create mocks create
func (bl *MockBLEmployee) Create(reqCtx common.RequestContext, emp model.Employee) (model.Employee, error) {
	bl.DOB = emp.Dateofbirth
	return emp, bl.Err
}

// Get mocks get
func (bl *MockBLEmployee) Get(reqCtx common.RequestContext, id int64) (emp model.Employee, err error) {
	bl.ID = id
	return emp, bl.Err
}

// GetAll mocks return all employees
func (bl *MockBLEmployee) GetAll(reqCtx common.RequestContext, limit int64) (emps []model.Employee, err error) {
	bl.Limit = limit
	return emps, bl.Err
}

// Update mocks update of a single employee
func (bl *MockBLEmployee) Update(reqCtx common.RequestContext, id int64, emp model.Employee) (model.Employee, error) {
	bl.ID = id
	return emp, bl.Err
}

// Delete mocks delete
func (bl *MockBLEmployee) Delete(reqCtx common.RequestContext, id int64) (string, error) {
	bl.ID = id
	return "ok", bl.Err
}

// Search finds employees matching the term
func (bl *MockBLEmployee) Search(reqCtx common.RequestContext, term string, limit int64) (emps []model.Employee, err error) {
	bl.Term = term
	bl.Limit = limit
	return emps, bl.Err
}

func (bl *MockBLEmployee) Upload(reqCtx common.RequestContext, employees []model.Employee) ([]model.Employee, error) {
	return employees, bl.Err
}

func (bl *MockBLEmployee) SearchStatus(reqCtx common.RequestContext, emp model.Employee) (model.Employee, error) {
	bl.DOB = emp.Dateofbirth
	return emp, bl.Err
}

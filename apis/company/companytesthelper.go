package company

import (
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/model"
)

// MockBLCompany is a struct for mocking company business logic methods
type MockBLCompany struct {
	ID int

	Term  string
	Limit int
	Err   error
}

// Create mocks create
func (bl *MockBLCompany) Create(reqCtx common.RequestContext, comp model.Company) (model.Company, error) {
	return comp, bl.Err
}

// Get mocks get
func (bl *MockBLCompany) Get(reqCtx common.RequestContext, id int) (comp model.Company, err error) {
	bl.ID = id
	return comp, bl.Err
}

// GetAll mocks return all companys
func (bl *MockBLCompany) GetAll(reqCtx common.RequestContext, limit int) (comps []model.Company, err error) {
	bl.Limit = limit
	return comps, bl.Err
}

// Update mocks update of a single company
func (bl *MockBLCompany) Update(reqCtx common.RequestContext, id int, comp model.Company) (model.Company, error) {
	bl.ID = id
	return comp, bl.Err
}

// Delete mocks delete
func (bl *MockBLCompany) Delete(reqCtx common.RequestContext, id int) (string, error) {
	bl.ID = id
	return "ok", bl.Err
}

// Search finds companys matching the term
func (bl *MockBLCompany) Search(reqCtx common.RequestContext, term string, limit int) (comps []model.Company, err error) {
	bl.Term = term
	bl.Limit = limit
	return comps, bl.Err
}

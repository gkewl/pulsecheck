package token

import (
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/model"
)

// MockBLToken is a struct for mocking token business logic methods
type MockBLToken struct {
	ID int64

	Term        string
	Limit       int64
	Err         error
	TokenString string
}

// Create mocks create
func (bl *MockBLToken) Create(reqCtx common.RequestContext, tk model.Token) (model.Token, error) {
	return tk, bl.Err
}

// Get mocks get
func (bl *MockBLToken) Get(reqCtx common.RequestContext, id int64) (tk model.Token, err error) {
	bl.ID = id
	return tk, bl.Err
}

// Get mocks get
func (bl *MockBLToken) GetByToken(reqCtx common.RequestContext, tokenStr string) (tk model.Token, err error) {
	bl.TokenString = tokenStr
	return tk, bl.Err
}

// GetAll mocks return all tokens
func (bl *MockBLToken) GetAll(reqCtx common.RequestContext, limit int64) (tks []model.Token, err error) {
	bl.Limit = limit
	return tks, bl.Err
}

// Update mocks update of a single token
func (bl *MockBLToken) Update(reqCtx common.RequestContext, id int64, tk model.Token) (model.Token, error) {
	bl.ID = id
	return tk, bl.Err
}

// Delete mocks delete
func (bl *MockBLToken) Delete(reqCtx common.RequestContext, id int64) (string, error) {
	bl.ID = id
	return "ok", bl.Err
}

// Search finds tokens matching the term
func (bl *MockBLToken) Search(reqCtx common.RequestContext, term string, limit int64) (tks []model.Token, err error) {
	bl.Term = term
	bl.Limit = limit
	return tks, bl.Err
}

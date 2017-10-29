package user

import (
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/model"
)

// MockBLUser is a struct for mocking user business logic methods
type MockBLUser struct {
	ID int

	Term  string
	Limit int64
	Err   error
}

// Create mocks create
func (bl *MockBLUser) Create(reqCtx common.RequestContext, usr model.RegisterUser) (model.User, error) {
	u := model.User{Email: usr.Email}
	return u, bl.Err
}

// Get mocks get
func (bl *MockBLUser) Get(reqCtx common.RequestContext, id int) (usr model.User, err error) {
	bl.ID = id
	return usr, bl.Err
}

// GetAll mocks return all users
func (bl *MockBLUser) GetAll(reqCtx common.RequestContext, limit int64) (usrs []model.User, err error) {
	bl.Limit = limit
	return usrs, bl.Err
}

// Update mocks update of a single user
func (bl *MockBLUser) Update(reqCtx common.RequestContext, id int, usr model.User) (model.User, error) {
	bl.ID = id
	return usr, bl.Err
}

// Delete mocks delete
func (bl *MockBLUser) Delete(reqCtx common.RequestContext, id int) (string, error) {
	bl.ID = id
	return "ok", bl.Err
}

// Search finds users matching the term
func (bl *MockBLUser) Search(reqCtx common.RequestContext, term string, limit int64) (usrs []model.User, err error) {
	bl.Term = term
	bl.Limit = limit
	return usrs, bl.Err
}

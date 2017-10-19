package authuser

import (
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/model"
)

// MockBLCompany is a struct for mocking company business logic methods
type MockBLAuthUser struct {
	ID int

	Term  string
	Limit int
	Err   error
}

func (bl *MockBLAuthUser) CreateOrUpdate(common.RequestContext, model.AuthUser) (model.AuthUser, error) {
	return model.AuthUser{}, bl.Err
}
func (bl *MockBLAuthUser) Get(common.RequestContext, int) (model.AuthUser, error) {
	return model.AuthUser{}, bl.Err
}
func (bl *MockBLAuthUser) Delete(common.RequestContext, int) (string, error) {
	return "ok", bl.Err
}
func (bl *MockBLAuthUser) Authenticate(common.RequestContext, int, string) (bool, error) {
	return true, bl.Err
}

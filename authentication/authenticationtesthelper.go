package authentication

import (
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/model"
)

// MockBLlooppartinfo is a struct for mocking looppartinfo business logic methods
type MockBLauthentication struct {
	name  string
	Token string
	Exp   int64
	Err   error
}

// Mock LoginUser
func (bl *MockBLauthentication) LoginUser(reqCtx common.RequestContext, input model.User) (model.TokenInfo, error) {
	bl.Token = "42"
	return model.TokenInfo{Token: "42"}, bl.Err
}

// Mock MachineToken
func (bl *MockBLauthentication) MachineToken(reqCtx common.RequestContext, name string) (model.TokenInfo, error) {
	bl.name = name
	bl.Token = "42"
	return model.TokenInfo{Token: "42"}, bl.Err
}

// Mock MachineToken
func (bl *MockBLauthentication) ValidateToken(reqCtx common.RequestContext, token string) (int64, error) {
	return 42, bl.Err
}

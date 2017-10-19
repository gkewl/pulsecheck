package authentication

import (
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/model"
)

// MockBLauthentication is a struct for mocking looppartinfo business logic methods
type MockBLauthentication struct {
	name  string
	Token string
	Exp   int64
	Err   error
}

// LoginUser - Mock LoginUser
func (bl *MockBLauthentication) LoginUser(reqCtx common.RequestContext, input model.AuthUser) (model.TokenInfo, error) {
	bl.Token = "42"
	return model.TokenInfo{Token: "42"}, bl.Err
}

// ValidateToken mocks ValidateToken
func (bl *MockBLauthentication) ValidateToken(reqCtx common.RequestContext, token string) (int64, error) {
	return 42, bl.Err
}

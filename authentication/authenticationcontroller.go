package authentication

import (
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/constant"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
)

// TestingBizLogic will be set in testing to a MockBL
var TestingBizLogic BizLogic

func authenticationLogic() BizLogic {
	if TestingBizLogic != nil {
		return TestingBizLogic
	}
	return BLAuthentication{}
}

type authenticationControllerFunc func(reqCtx common.RequestContext, authenticationBL BizLogic, authenticationInput model.AuthenticateUser) (interface{}, error)

// ControlWrapper extracts information from the request and calls the wrapped
// controller function
func ControlWrapper(handler authenticationControllerFunc) func(common.RequestContext) (interface{}, error) {
	return func(reqCtx common.RequestContext) (interface{}, error) {
		var authenticationInterface BizLogic
		authenticationInterface = BLAuthentication{}
		if TestingBizLogic != nil {
			authenticationInterface = TestingBizLogic
		}
		authentication := model.AuthenticateUser{}
		err := reqCtx.Scan("", &authentication)
		if err != nil {
			return nil, eh.NewError(eh.ErrJsonDecodeFail, "JSON Error: "+err.Error())
		}
		return handler(reqCtx, authenticationInterface, authentication)
	}
}

// GetRoutes returns all Authentication-related routes
func GetRoutes() common.Routes {
	return common.Routes{
		common.Route{
			Name:           "LoginUser",
			Method:         "POST",
			Pattern:        "/auth/token-auth",
			ControllerFunc: ControlWrapper(LoginUser),
			SecureBody:     true,
		},

		common.Route{
			Name:           "ValidateToken",
			Method:         "GET",
			Pattern:        "/auth/validatetoken",
			ControllerFunc: ValidateToken,
		},
	}
}

// LoginUser authenticates User
func LoginUser(reqCtx common.RequestContext, authenticationBL BizLogic, input model.AuthenticateUser) (interface{}, error) {
	reqCtx.SetUserId(constant.DefaultAdmin)
	return authenticationBL.LoginUser(reqCtx, input)
}

// ValidateToken Validate Token New authentication Process
func ValidateToken(reqCtx common.RequestContext) (interface{}, error) {
	token := reqCtx.Token()
	var authenticationInterface BizLogic
	authenticationInterface = BLAuthentication{}
	return authenticationInterface.ValidateToken(reqCtx, token)
}

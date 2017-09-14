package token

import (
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/constant"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
)

// TestingBizLogic will be set in testing to a MockBL
var TestingBizLogic BizLogic

type tokenControllerFunc func(reqCtx common.RequestContext, tokenBL BizLogic, tokenInput model.Token) (interface{}, error)

// ControlWrapper extracts information from the request and calls the wrapped
// controller function
func ControlWrapper(handler tokenControllerFunc) func(common.RequestContext) (interface{}, error) {
	return func(reqCtx common.RequestContext) (interface{}, error) {
		var tokenInterface BizLogic
		tokenInterface = BLToken{}
		if TestingBizLogic != nil {
			tokenInterface = TestingBizLogic
		}
		tk := model.Token{}
		err := reqCtx.Scan("tk", &tk)
		if err != nil {
			return nil, eh.NewError(eh.ErrJsonDecodeFail, "JSON Error: "+err.Error())
		}
		return handler(reqCtx, tokenInterface, tk)
	}
}

// GetRoutes returns all token-related routes
func GetRoutes() common.Routes {

	return common.Routes{

		common.Route{
			Name:           "SearchToken",
			Method:         "GET",
			Pattern:        "/token/search",
			ControllerFunc: ControlWrapper(SearchToken),
			AuthRequired:   constant.Superuser,
		},

		common.Route{
			Name:           "GetToken",
			Method:         "GET",
			Pattern:        "/token/{id:[0-9]+}",
			ControllerFunc: ControlWrapper(GetToken),
			AuthRequired:   constant.Superuser,
		},
		common.Route{
			Name:           "GetTokenbyString",
			Method:         "GET",
			Pattern:        "/token/{name}",
			ControllerFunc: ControlWrapper(GetByToken),
			AuthRequired:   constant.Superuser,
		},
		common.Route{
			Name:           "GetTokens",
			Method:         "GET",
			Pattern:        "/token",
			ControllerFunc: ControlWrapper(GetAllTokens),
			AuthRequired:   constant.Superuser,
		},

		common.Route{
			Name:           "UpdateToken",
			Method:         "PUT",
			Pattern:        "/token/{id}",
			ControllerFunc: ControlWrapper(UpdateToken),
			AuthRequired:   constant.Superuser,
		},
	}
}

// GetToken gets a token by primary key
func GetToken(reqCtx common.RequestContext, tokenBL BizLogic, tk model.Token) (interface{}, error) {
	id := reqCtx.IntValue("id", 0)
	return tokenBL.Get(reqCtx, id)
}

// GetByToken gets a token by token string
func GetByToken(reqCtx common.RequestContext, tokenBL BizLogic, tk model.Token) (interface{}, error) {
	tokenString := reqCtx.Value("name", "")
	return tokenBL.GetByToken(reqCtx, tokenString)
}

// GetAllTokens gets all tokens
func GetAllTokens(reqCtx common.RequestContext, tokenBL BizLogic, tk model.Token) (interface{}, error) {
	limit := reqCtx.IntValue("limit", 50)
	return tokenBL.GetAll(reqCtx, limit)
}

// SearchToken finds tokens that match a term
func SearchToken(reqCtx common.RequestContext, tokenBL BizLogic, tk model.Token) (interface{}, error) {
	return tokenBL.Search(reqCtx, reqCtx.Value("term", ""), reqCtx.IntValue("limit", 1000))
}

// UpdateToken updates a single token
func UpdateToken(reqCtx common.RequestContext, tokenBL BizLogic, tk model.Token) (interface{}, error) {
	id := reqCtx.IntValue("id", 0)
	return tokenBL.Update(reqCtx, id, tk)
}

package user

import (
	"net/http"

	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/constant"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
)

// TestingBizLogic will be set in testing to a MockBL
var TestingBizLogic BizLogic

type userControllerFunc func(reqCtx common.RequestContext, userBL BizLogic, userInput model.User) (interface{}, error)

func getInterface() BizLogic {
	userInterface := BizLogic(BLUser{})
	if TestingBizLogic != nil {
		userInterface = TestingBizLogic
	}
	return userInterface
}

// ControlWrapper extracts information from the request and calls the wrapped
// controller function
func ControlWrapper(handler userControllerFunc) func(common.RequestContext) (interface{}, error) {
	return func(reqCtx common.RequestContext) (interface{}, error) {
		userInterface := getInterface()
		usr := model.User{}
		err := reqCtx.Scan("usr", &usr)
		if err != nil {
			return nil, eh.NewError(eh.ErrJsonDecodeFail, "JSON Error: "+err.Error())
		}
		return handler(reqCtx, userInterface, usr)
	}
}

// GetRoutes returns all user-related routes
func GetRoutes() common.Routes {

	return common.Routes{

		common.Route{
			Name:           "SearchUser",
			Method:         "GET",
			Pattern:        "/user/search",
			ControllerFunc: ControlWrapper(SearchUser),
			AuthRequired:   constant.User,
		},

		common.Route{
			Name:           "GetUsers",
			Method:         "GET",
			Pattern:        "/user/all",
			ControllerFunc: ControlWrapper(GetAllUsers),
			AuthRequired:   constant.User,
		},

		common.Route{
			Name:           "GetUser",
			Method:         "GET",
			Pattern:        "/user/{id}",
			ControllerFunc: ControlWrapper(GetUser),
			AuthRequired:   constant.User,
		},

		common.Route{
			Name:           "CreateUser",
			Method:         "POST",
			Pattern:        "/user",
			ControllerFunc: CreateUser,
			AuthRequired:   constant.Superuser,
			NormalHttpCode: http.StatusCreated,
		},
		common.Route{
			Name:           "UpdateUser",
			Method:         "PUT",
			Pattern:        "/user/{id}",
			ControllerFunc: ControlWrapper(UpdateUser),
			AuthRequired:   constant.User,
		},
		common.Route{
			Name:           "DeleteUser",
			Method:         "DELETE",
			Pattern:        "/user/{id}",
			ControllerFunc: ControlWrapper(DeleteUser),
			AuthRequired:   constant.Admin,
		},
	}
}

// GetUser gets a user by primary key
func GetUser(reqCtx common.RequestContext, userBL BizLogic, usr model.User) (interface{}, error) {
	id := reqCtx.IntValue32("id", reqCtx.UserID())
	return userBL.Get(reqCtx, id)
}

// GetAllUsers gets all users
func GetAllUsers(reqCtx common.RequestContext, userBL BizLogic, usr model.User) (interface{}, error) {
	limit := reqCtx.IntValue("limit", 50)
	return userBL.GetAll(reqCtx, limit)
}

// // CreateUser creates a user
// func CreateUser(reqCtx common.RequestContext, userBL BizLogic, usr model.User) (interface{}, error) {
// 	return userBL.Create(reqCtx, usr)
// }

// CreateUser creates a user
func CreateUser(reqCtx common.RequestContext) (interface{}, error) {

	userInterface := getInterface()
	usr := model.RegisterUser{}
	err := reqCtx.Scan("usr", &usr)
	if err != nil {
		return nil, eh.NewError(eh.ErrJsonDecodeFail, "JSON Error: "+err.Error())
	}

	return userInterface.Create(reqCtx, usr)
}

// UpdateUser updates a single user
func UpdateUser(reqCtx common.RequestContext, userBL BizLogic, usr model.User) (interface{}, error) {
	id := reqCtx.IntValue32("id", 0)
	return userBL.Update(reqCtx, id, usr)
}

// DeleteUser deletes a single user
func DeleteUser(reqCtx common.RequestContext, userBL BizLogic, usr model.User) (interface{}, error) {
	id := reqCtx.IntValue32("id", 0)
	return userBL.Delete(reqCtx, id)
}

// SearchUser finds users that match a term
func SearchUser(reqCtx common.RequestContext, userBL BizLogic, usr model.User) (interface{}, error) {
	return userBL.Search(reqCtx, reqCtx.Value("term", ""), reqCtx.IntValue("limit", 1000))
}

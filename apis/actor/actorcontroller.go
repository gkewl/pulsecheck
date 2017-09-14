package actor

import (
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/constant"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
	"net/http"
)

//TestingBizLogic will be set in testing to a MockBL
var TestingBizLogic BizLogic

type actorControllerFunc func(reqCtx common.RequestContext, actorBL BizLogic, actorInput model.Actor) (interface{}, error)

// ControlWrapper extracts information from the request and calls the wrapped
// controller function
func ControlWrapper(handler actorControllerFunc) func(common.RequestContext) (interface{}, error) {
	return func(reqCtx common.RequestContext) (interface{}, error) {
		var actorInterface BizLogic
		actorInterface = BLActor{}
		if TestingBizLogic != nil {
			actorInterface = TestingBizLogic
		}
		actor := model.Actor{}
		err := reqCtx.Scan("actor", &actor)
		if err != nil {
			return nil, eh.NewError(eh.ErrJsonDecodeFail, "JSON Error: "+err.Error())
		}
		return handler(reqCtx, actorInterface, actor)
	}
}

// GetRoutes returns all step-related routes
func GetRoutes() common.Routes {

	return common.Routes{
		common.Route{
			Name:           "CreateActor",
			Method:         "POST",
			Pattern:        "/actor",
			ControllerFunc: ControlWrapper(CreateActor),
			AuthRequired:   constant.Superuser,
			NormalHttpCode: http.StatusCreated,
		},
		common.Route{
			Name:           "SearchActor",
			Method:         "GET",
			Pattern:        "/actor/search", //URl is /Actor/Search?Type="User"&Term="anything here"
			ControllerFunc: ControlWrapper(Search),
		},
		common.Route{
			Name:           "GetActor",
			Method:         "GET",
			Pattern:        "/actor/{name}",
			ControllerFunc: ControlWrapper(GetActor),
		},
		common.Route{
			Name:           "UpdateActor",
			Method:         "PUT",
			Pattern:        "/actor/{id}",
			ControllerFunc: ControlWrapper(UpdateActor),
			AuthRequired:   constant.Superuser,
		},
		common.Route{
			Name:           "DeleteActor",
			Method:         "DELETE",
			Pattern:        "/actor/{id}",
			ControllerFunc: ControlWrapper(DeleteActor),
			AuthRequired:   constant.Superuser,
		},

		common.Route{
			Name:           "GetActors",
			Method:         "GET",
			Pattern:        "/actor",
			ControllerFunc: ControlWrapper(GetActorAll),
		},
	}
}

// CreateActor creates an actor
func CreateActor(reqCtx common.RequestContext, actorBL BizLogic, actor model.Actor) (interface{}, error) {
	return actorBL.Create(reqCtx, actor)
}

// GetActor returns an actor
func GetActor(reqCtx common.RequestContext, actorBL BizLogic, actor model.Actor) (interface{}, error) {
	return actorBL.Get(reqCtx, reqCtx.Value("name", ""))
}

// GetActorAll returns all the actors
func GetActorAll(reqCtx common.RequestContext, actorBL BizLogic, actor model.Actor) (interface{}, error) {
	return actorBL.GetAll(reqCtx, reqCtx.Value("type", ""))
}

// UpdateActor updates a single actor
func UpdateActor(reqCtx common.RequestContext, actorBL BizLogic, actor model.Actor) (interface{}, error) {
	id := reqCtx.IntValue("id", 0)
	return actorBL.Update(reqCtx, id, actor)
}

// DeleteActor deletes a single actor
func DeleteActor(reqCtx common.RequestContext, actorBL BizLogic, actor model.Actor) (interface{}, error) {
	id := reqCtx.IntValue("id", 0)
	return actorBL.Delete(reqCtx, id)
}

// Search finds actors that match a term and a single type
func Search(reqCtx common.RequestContext, actorBL BizLogic, actor model.Actor) (interface{}, error) {
	userType := reqCtx.Value("type", "")
	term := reqCtx.Value("term", "")
	return actorBL.Search(reqCtx, userType, term)
}

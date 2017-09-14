package actor

import (
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/constant"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
	"gopkg.in/guregu/null.v3"
)

// BizLogic is the interface for all actor business logic methods
type BizLogic interface {
	Create(common.RequestContext, model.Actor) (model.Actor, error)
	Get(common.RequestContext, string) (model.Actor, error)
	GetAll(common.RequestContext, string) ([]model.Actor, error)
	Update(common.RequestContext, int64, model.Actor) (model.Actor, error)
	Delete(common.RequestContext, int64) (string, error)
	Search(common.RequestContext, string, string) ([]model.ActorSearchResponse, error)
	UpdateLastLoginTime(common.RequestContext, int64) error
}

// BLActor implements the actor.BizLogic interface
type BLActor struct {
}

//ValidateActor does the preprocesing checks the ensure that preconditions are met before creating or updating an actor.
func ValidateActor(reqCtx common.RequestContext, actor model.Actor, insert int) error {
	newErr := eh.ErrActorInsert
	if insert == 0 {
		newErr := eh.ErrActorUpdate
		if actor.Id <= 0 {
			return eh.NewError(newErr, "Id cannot be 0 or less")
		}
	}

	if len(actor.Name) == 0 {
		return eh.NewError(newErr, "Actor name cannot be empty")
	}

	if actor.Type != constant.ActorType_System || actor.Type != constant.ActorType_User {
		return eh.NewError(newErr, "Actor type needs to be SYSTEM or USER")
	}

	if len(actor.Role) == 0 {
		return eh.NewError(newErr, "Role cannot be empty")
	}

	if actor.Role != constant.Role_Guest || actor.Type != constant.Role_Admin || actor.Type != constant.Role_SuperUser {
		return eh.NewError(newErr, "Check Role Name")
	}

	return nil
}

// Create will insert a new Actor into the db
func (bl BLActor) Create(reqCtx common.RequestContext, actor model.Actor) (model.Actor, error) {

	err := ValidateActor(reqCtx, actor, 1)
	if err != nil {
		return model.Actor{}, err
	}

	if actor.Manager.Name.Valid == true {
		parent, err := bl.Get(reqCtx, actor.Manager.Name.String)
		if err != nil {
			return model.Actor{}, eh.NewError(eh.ErrActorParentNotFound, "DB Error: "+err.Error())
		}
		actor.Manager = model.NullableNameDescription{Id: null.NewInt(parent.Id, true)}
	}

	actor, err = DLCreate(reqCtx, actor)

	if err != nil {
		return model.Actor{}, eh.NewError(eh.ErrActorInsert, "DB Error: "+err.Error())
	}

	return actor, err
}

// Get returns a single actor by name
func (bl BLActor) Get(reqCtx common.RequestContext, name string) (actor model.Actor, err error) {
	actor, err = DLGet(reqCtx, name)

	if err != nil || len(actor.Name) == 0 {
		return model.Actor{}, eh.NewErrorNotFound(eh.ErrActorDataNotFound, err, "actor %s not found", name)
	}
	return
}

// GetAll returns all actors
func (bl BLActor) GetAll(reqCtx common.RequestContext, userTypes string) (actors []model.Actor, err error) {

	actors, err = DLGetAll(reqCtx, userTypes)

	if err != nil {
		return []model.Actor{}, eh.NewErrorNotFound(eh.ErrActorDataNotFound, err, "actor %s not found")
	}
	return
}

//Update updates the actor by id if its valid
func (bl BLActor) Update(reqCtx common.RequestContext, id int64, actor model.Actor) (model.Actor, error) {

	err := ValidateActor(reqCtx, actor, 0)
	if err != nil {
		return model.Actor{}, err
	}

	err = DLUpdate(reqCtx, id, actor)

	if err != nil {
		return model.Actor{}, eh.NewError(eh.ErrActorUpdate, "DB Error: "+err.Error())
	}

	return actor, err
}

//UpdateLastLoginTime updates last login time for actor
func (bl BLActor) UpdateLastLoginTime(reqCtx common.RequestContext, id int64) error {
	err := DLUpdateLastlogin(reqCtx, id)
	if err != nil {
		return eh.NewError(eh.ErrActorUpdate, "Updating last login failed. DB Error: "+err.Error())
	}

	return nil
}

// Delete marks an Actor as inactive
func (bl BLActor) Delete(reqCtx common.RequestContext, ID int64) (string, error) {

	affect, err := DLMarkInactive(reqCtx, ID)

	if err != nil {
		return "", eh.NewError(eh.ErrActorDelete, "DB Error: "+err.Error())
	}

	if err == nil && affect != 1 {
		err = eh.NewError(eh.ErrActorDataNotFound, "actor id %d not found", ID)
	}

	return "ok", err
}

//Search returns the actor based on the search term used as the filter on the where clause.
func (bl BLActor) Search(reqCtx common.RequestContext, userType string, term string) ([]model.ActorSearchResponse, error) {

	actors, err := DLGetActors(reqCtx, userType, term, 50)
	if err != nil {
		return []model.ActorSearchResponse{}, eh.NewError(eh.ErrActorDataNotFound, "DB Error: "+err.Error())
	}
	return actors, err

}

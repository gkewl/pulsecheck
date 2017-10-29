package user

import (
	"github.com/gkewl/pulsecheck/apis/authuser"
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/constant"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
)

// BizLogic is the interface for all user business logic methods
type BizLogic interface {
	Create(common.RequestContext, model.RegisterUser) (model.User, error)
	Get(common.RequestContext, int) (model.User, error)
	GetAll(common.RequestContext, int64) ([]model.User, error)
	Update(common.RequestContext, int, model.User) (model.User, error)
	Delete(common.RequestContext, int) (string, error)
	Search(common.RequestContext, string, int64) ([]model.User, error)
}

// BLUser implements the user.BizLogic interface
type BLUser struct {
}

// Create will insert a new user into the db
func (bl BLUser) Create(reqCtx common.RequestContext, usr model.RegisterUser) (model.User, error) {

	newUser := model.User{
		Email:      usr.Email,
		FirstName:  usr.FirstName,
		MiddleName: usr.MiddleName,
		LastName:   usr.LastName,
		CompanyID:  usr.CompanyID,
	}

	var err error
	newUser, err = dlCreate(reqCtx, newUser)
	if err != nil {
		return model.User{}, eh.NewError(eh.ErrUserInsert, "DB Error: "+err.Error())
	}

	if newUser.ID > 0 {
		a := model.AuthUser{UserID: newUser.ID, Password: usr.Password}
		a, err = authuser.BLAuthUser{}.CreateOrUpdate(reqCtx, a)
		if err != nil {
			return model.User{}, eh.WrapError(eh.ErrUserInsert, err, "DB Error: "+err.Error())
		}
	}
	return newUser, err
}

// Get returns a single user by primary key
func (bl BLUser) Get(reqCtx common.RequestContext, id int) (usr model.User, err error) {
	usr, err = dlGet(reqCtx, id, reqCtx.CompanyID())
	if err != nil || usr.ID == 0 {
		return model.User{}, eh.NewErrorNotFound(eh.ErrUserDataNotFound, err, `User not found: id %d`, id)
	}
	return
}

// GetAll will return all users
func (bl BLUser) GetAll(reqCtx common.RequestContext, limit int64) (usrs []model.User, err error) {
	usrs, err = dlGetAll(reqCtx, limit, reqCtx.CompanyID())
	if err != nil {
		return []model.User{}, eh.NewError(eh.ErrUserDataNotFound, "DB Error: "+err.Error())
	}
	return
}

// Update updates a single user
func (bl BLUser) Update(reqCtx common.RequestContext, id int, usr model.User) (model.User, error) {
	// todo: add validation here
	if reqCtx.CompanyID() != usr.CompanyID {
		return model.User{}, eh.NewError(eh.ErrUserUpdate, "User %s does not belong to your company ", usr.Email)

	}
	result, err := dlUpdate(reqCtx, id, usr)
	if err != nil {
		return model.User{}, eh.NewError(eh.ErrUserUpdate, "DB Error: "+err.Error())
	}
	return result, err
}

// Delete marks a single user inactive
func (bl BLUser) Delete(reqCtx common.RequestContext, id int) (string, error) {
	// todo: add validation here

	err := dlMarkInactive(reqCtx, id)
	if err != nil {
		return constant.ResultFail, eh.NewError(eh.ErrUserDelete, "DB Error: "+err.Error())
	}
	return constant.ResultOk, err
}

// HardDelete physically deletes a user, usually for testing
func (bl BLUser) HardDelete(reqCtx common.RequestContext, id int) (string, error) {
	err := dlDelete(reqCtx, id)
	if err != nil {
		return constant.ResultFail, eh.NewError(eh.ErrUserDelete, "DB Error: "+err.Error())
	}
	return constant.ResultOk, err
}

// Search finds users matching the term
func (bl BLUser) Search(reqCtx common.RequestContext, term string, limit int64) (usrs []model.User, err error) {
	usrs, err = dlSearch(reqCtx, term, limit)
	if err != nil {
		return []model.User{}, eh.NewError(eh.ErrUserDataNotFound, "DB Error: "+err.Error())
	}
	return
}

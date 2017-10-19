package authuser

import (
	"crypto/md5"
	"fmt"
	"io"

	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/constant"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
)

// BizLogic is the interface for all AuthUser business logic methods
type BizLogic interface {
	CreateOrUpdate(common.RequestContext, model.AuthUser) (model.AuthUser, error)
	Get(common.RequestContext, int) (model.AuthUser, error)
	Delete(common.RequestContext, int) (string, error)
	Authenticate(common.RequestContext, int, string) (bool, error)
}

var TestingBizLogic BizLogic

// BLAuthUser implements the AuthUser.BizLogic interface
type BLAuthUser struct {
}

func GetInterface() BizLogic {
	if TestingBizLogic != nil {
		return TestingBizLogic
	}
	return BLAuthUser{}
}

// Create will insert a new AuthUser into the db
func (bl BLAuthUser) CreateOrUpdate(reqCtx common.RequestContext, au model.AuthUser) (model.AuthUser, error) {

	usr, err := dlGet(reqCtx, au.UserID)
	if err != nil && !eh.HasNoRowsError(err) {
		return model.AuthUser{}, err
	}

	au.Password = bl.getHash(au.Password)
	if usr.ID > 0 {
		//update
		usr, err = dlUpdate(reqCtx, au)
	} else {
		//create
		usr, err = dlCreate(reqCtx, au)
	}
	if err != nil {
		return model.AuthUser{}, eh.NewError(eh.ErrAuthUserInsert, "DB Error: "+err.Error())
	}
	return usr, err
}

// Get returns a single AuthUser by primary key
func (bl BLAuthUser) Get(reqCtx common.RequestContext, userID int) (au model.AuthUser, err error) {
	au, err = dlGet(reqCtx, userID)

	if err != nil || au.ID == 0 {
		return model.AuthUser{}, eh.NewErrorNotFound(eh.ErrAuthUserDataNotFound, err, `AuthUser not found: userID %d`, userID)
	}
	return
}

// Update updates a single AuthUser
func (bl BLAuthUser) Authenticate(reqCtx common.RequestContext, userID int, pwd string) (bool, error) {

	usr, err := bl.Get(reqCtx, userID)
	if err != nil {
		return false, eh.NewError(eh.ErrAuthUserUpdate, "DB Error: "+err.Error())
	}

	hash := bl.getHash(pwd)
	if hash == usr.Password {
		return true, nil
	}
	return false, nil
}

func (bl BLAuthUser) getHash(value string) string {
	h := md5.New()
	io.WriteString(h, value)
	str := fmt.Sprintf("%x", h.Sum(nil))
	return str
}

// Delete marks a single AuthUser inactive
func (bl BLAuthUser) Delete(reqCtx common.RequestContext, userID int) (string, error) {
	// todo: add validation here
	err := dlMarkInactive(reqCtx, userID)
	if err != nil {
		return constant.ResultFail, eh.NewError(eh.ErrAuthUserDelete, "DB Error: "+err.Error())
	}
	return constant.ResultOk, err
}

package authuser

import (
	"github.com/gkewl/pulsecheck/common"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
)

const (
	getQuery = `select au.id, au.userid, au.password , au.isactive from authuser au`
)

// dlGet retrieves the specified AuthUser
func dlGet(reqCtx common.RequestContext, userID int) (au model.AuthUser, err error) {
	query := getQuery + ` where au.userid = ?`
	err = reqCtx.Tx().Get(&au, query, userID)
	return
}

// dlCreate creates a AuthUser
func dlCreate(reqCtx common.RequestContext, au model.AuthUser) (model.AuthUser, error) {
	params := map[string]interface{}{
		"userid":     au.UserID,
		"password":   au.Password,
		"isactive":   1,
		"createdby":  reqCtx.UserName(),
		"modifiedby": reqCtx.UserName(),
	}

	result, err := reqCtx.Tx().NamedExec(
		`insert into authuser (userid, password, isactive, createdby, modifiedby)
	 		values (:userid, :password, :isactive, :createdby, :modifiedby)`,
		params)

	if err == nil {
		au.ID, _ = result.LastInsertId()
	}

	return au, err
}

// dlUpdate updates fields on a AuthUser and returns full updated object
func dlUpdate(reqCtx common.RequestContext, au model.AuthUser) (model.AuthUser, error) {
	params := map[string]interface{}{
		"userid":     au.UserID,
		"password":   au.Password,
		"modifiedby": reqCtx.UserName(),
	}
	_, err := reqCtx.Tx().NamedExec(
		`update authuser
			set  password=:password,
			modifiedby=:modifiedby
			where userid=:userid`, params)
	if err == nil {
		return dlGet(reqCtx, au.UserID)
	}
	return model.AuthUser{}, err
}

// dlMarkInactive set the isactive flag to zero for this AuthUser
func dlMarkInactive(reqCtx common.RequestContext, userID int) error {
	result, err := reqCtx.Tx().Exec(
		`update authuser au
		set isactive = 0, modifiedby = ?
		where au.userid = ? and au.isactive = 1`,
		reqCtx.UserName(), userID)
	if err == nil {
		if affect, _ := result.RowsAffected(); affect != 1 {
			err = eh.ErrDBNoRows
		}
	}
	return err
}

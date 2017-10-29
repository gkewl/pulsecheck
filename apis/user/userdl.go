package user

import (
	"fmt"

	//	"github.com/gkewl/pulsecheck/apis/search"

	"github.com/gkewl/pulsecheck/common"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
)

const (
	getQuery = `select user.id as id, user.email as email, user.firstname as firstname,
user.middlename as middlename, user.lastname as lastname,
user.companyid as companyid, co.name as companyname, user.isactive as isactive, user.createdby
as createdby, user.created as created, user.modifiedby as modifiedby,
user.modified as modified from user user 
join company co on co.id = user.companyid and co.isactive = 1`
)

// dlGet retrieves the specified user
func dlGet(reqCtx common.RequestContext, id int, companyid int) (usr model.User, err error) {
	query := getQuery + ` where user.id = ? and user.companyid=? and  user.isactive = 1`
	err = reqCtx.Tx().Get(&usr, query, id, companyid)
	return
}

// dlGetAll retrieves all users
func dlGetAll(reqCtx common.RequestContext, limit int64, companyid int) (usrs []model.User, err error) {
	usrs = []model.User{}
	query := getQuery + ` where user.companyid=? and user.isactive=1  `

	if limit != 0 {
		query = query + fmt.Sprintf(" limit %d", limit)
	}
	err = reqCtx.Tx().Select(&usrs, query, companyid)
	return
}

// dlCreate creates a user
func dlCreate(reqCtx common.RequestContext, usr model.User) (model.User, error) {
	params := map[string]interface{}{
		"email":      usr.Email,
		"firstname":  usr.FirstName,
		"middlename": usr.MiddleName,
		"lastname":   usr.LastName,
		"companyid":  usr.CompanyID,
		"isactive":   1,
		"createdby":  reqCtx.UserName(),
		"modifiedby": reqCtx.UserName(),
	}

	result, err := reqCtx.Tx().NamedExec(
		`insert into user (email, firstname, middlename, lastname, companyid, isactive,
createdby, modifiedby)
	 		values (:email, :firstname, :middlename, :lastname, :companyid, :isactive,
:createdby, :modifiedby)`,
		params)

	if err == nil {
		id, _ := result.LastInsertId()
		usr.ID = int(id)
	}

	return usr, err
}

// dlUpdate updates fields on a user and returns full updated object
func dlUpdate(reqCtx common.RequestContext, id int, usr model.User) (model.User, error) {
	params := map[string]interface{}{
		"id":         id,
		"email":      usr.Email,
		"firstname":  usr.FirstName,
		"middlename": usr.MiddleName,
		"lastname":   usr.LastName,
		"companyid":  usr.CompanyID,
		"modifiedby": reqCtx.UserName(),
	}
	_, err := reqCtx.Tx().NamedExec(
		`update user
			set email = :email, firstname = :firstname, middlename
= :middlename, lastname = :lastname, companyid =
:companyid, modifiedby = :modifiedby
			where id = :id`, params)
	if err == nil {
		return dlGet(reqCtx, id, usr.CompanyID)
	}
	return model.User{}, err
}

// dlMarkInactive set the isactive flag to zero for this user
func dlMarkInactive(reqCtx common.RequestContext, id int) error {
	result, err := reqCtx.Tx().Exec(
		`update user user
		set isactive = 0, modifiedby = ?
		where user.id = ? and user.isactive = 1`,
		reqCtx.UserName(), id)
	if err == nil {
		if affect, _ := result.RowsAffected(); affect != 1 {
			err = eh.ErrDBNoRows
		}
	}
	return err
}

// dlDelete deletes this user
func dlDelete(reqCtx common.RequestContext, id int) error {
	result, err := reqCtx.Tx().Exec(`delete from user where id = ?`, id)
	if err == nil {
		if affect, _ := result.RowsAffected(); affect != 1 {
			err = eh.ErrDBNoRows
		}
	}
	return err
}

// dlSearch searches for users
func dlSearch(reqCtx common.RequestContext, term string, limit int64) ([]model.User, error) {
	data := []model.User{}

	// searchCols := []string{"user.email", "user.firstname", "user.middlename", "user.lastname"}
	// searchTerm := search.GetSearchCondition(searchCols, term)
	// isActiveWhere := ` user.isactive=1 `
	// qry := getQuery + ` where  ` + isActiveWhere
	// if len(searchTerm) > 0 && len(isActiveWhere) > 0 {
	// 	qry += ` and `
	// }
	// qry += searchTerm

	// if limit > 0 {
	// 	qry = qry + fmt.Sprintf(` limit %d`, limit)
	// }

	// stmt, err := reqCtx.Tx().Preparex(qry)
	// if err != nil {
	// 	return data, err
	// }

	// err = stmt.Select(&data)
	return data, nil
}

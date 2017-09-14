package token

import (
	"fmt"

	"github.com/gkewl/pulsecheck/apis/search"

	"github.com/gkewl/pulsecheck/common"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
)

const (
	getQuery = `select tk.id as id, tk.actorid as actorid, 
					tk.token as token, tk.expireson as expireson, tk.blocked as blocked,
					tk.created as created, tk.modified as modified, 
					tk.createdby as "createdby.id", a1.name as "createdby.name",  a1.description as "createdby.description",  
					tk.modifiedby as "modifiedby.id", a1.name as "modifiedby.name",  a1.description as "modifiedby.description",  
					tk.rowversion as rowversion
				from tokenaudit tk 
				left join actor a1 on a1.id = tk.createdby
				left join actor a2 on a2.id = tk.modifiedby
				
				`
)

// dlGet retrieves the specified token
func dlGet(reqCtx common.RequestContext, id int64) (tk model.Token, err error) {
	query := getQuery + ` where tk.id = ?`
	err = reqCtx.Tx().Get(&tk, query, id)
	return
}

// dlGetByToken retrieves the specified token
func dlGetByToken(reqCtx common.RequestContext, tokenString string) (tk model.Token, err error) {
	query := getQuery + ` where tk.token = ?`
	err = reqCtx.Tx().Get(&tk, query, tokenString)
	return
}

// dlGetAll retrieves all tokens
func dlGetAll(reqCtx common.RequestContext, limit int64) (tks []model.Token, err error) {
	tks = []model.Token{}
	query := getQuery
	query = query + ``
	if limit != 0 {
		query = query + fmt.Sprintf(" limit %d", limit)
	}
	err = reqCtx.Tx().Select(&tks, query)
	return
}

// dlCreate creates a token
func dlCreate(reqCtx common.RequestContext, tk model.Token) (model.Token, error) {
	params := map[string]interface{}{
		"actorid":    tk.ActorID,
		"token":      tk.Token,
		"expireson":  tk.ExpiresOn,
		"blocked":    tk.Blocked,
		"createdby":  reqCtx.UserID(),
		"modifiedby": reqCtx.UserID(),
	}

	result, err := reqCtx.Tx().NamedExec(
		`insert into tokenaudit (actorid,  token, expireson, blocked, createdby, modifiedby)
	 		values (:actorid, :adusername, :token, :expireson, :blocked, :createdby,:modifiedby)`,
		params)

	if err == nil {
		tk.ID, _ = result.LastInsertId()
	}

	return tk, err
}

// dlUpdate updates fields on a token and returns full updated object
func dlUpdate(reqCtx common.RequestContext, id int64, tk model.Token) (model.Token, error) {
	params := map[string]interface{}{
		"id":      id,
		"blocked": tk.Blocked,
	}
	_, err := reqCtx.Tx().NamedExec(
		`update tokenaudit
			set blocked=:blocked
			where id = :id`, params)
	if err == nil {
		return dlGet(reqCtx, id)
	}
	return model.Token{}, err
}

// dlDelete deletes this token
func dlDelete(reqCtx common.RequestContext, id int64) error {
	result, err := reqCtx.Tx().Exec(`delete from tokenaudit where id = ?`, id)
	if err == nil {
		if affect, _ := result.RowsAffected(); affect != 1 {
			err = eh.ErrDBNoRows
		}
	}
	return err
}

// dlSearch searches for tokens
func dlSearch(reqCtx common.RequestContext, term string, limit int64) ([]model.Token, error) {
	data := []model.Token{}

	searchCols := []string{"tk.adusername"}
	searchTerm := search.GetSearchCondition(searchCols, term)

	qry := getQuery
	if len(searchTerm) > 0 {
		qry += ` where  `
	}
	qry += searchTerm

	if limit > 0 {
		qry = qry + fmt.Sprintf(` limit %d`, limit)
	}

	stmt, err := reqCtx.Tx().Preparex(qry)
	if err != nil {
		return data, err
	}

	err = stmt.Select(&data)
	return data, err
}

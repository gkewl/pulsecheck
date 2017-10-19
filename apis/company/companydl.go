package company

import (
	"fmt"

	"github.com/gkewl/pulsecheck/common"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
)

const (
	getQuery = `select  id, name, isactive,createdby, created,
				modifiedby, modified from company `
)

// dlGet retrieves the specified company
func dlGet(reqCtx common.RequestContext, id int) (comp model.Company, err error) {
	query := getQuery + ` where company.id = ? and company.isactive = 1`
	err = reqCtx.Tx().Get(&comp, query, id)
	return
}

// dlGetAll retrieves all companys
func dlGetAll(reqCtx common.RequestContext, limit int) (comps []model.Company, err error) {
	comps = []model.Company{}
	query := getQuery + ` where company.isactive=1  `

	if limit != 0 {
		query = query + fmt.Sprintf(" limit %d", limit)
	}
	err = reqCtx.Tx().Select(&comps, query)
	return
}

// dlCreate creates a company
func dlCreate(reqCtx common.RequestContext, comp model.Company) (model.Company, error) {
	params := map[string]interface{}{
		"name":       comp.Name,
		"isactive":   1,
		"createdby":  reqCtx.UserName(),
		"modifiedby": reqCtx.UserName(),
	}

	result, err := reqCtx.Tx().NamedExec(
		`insert into company (name, isactive, createdby, modifiedby)
	 		values (:name, :isactive, :createdby, :modifiedby)`,
		params)

	if err == nil {
		id, _ := result.LastInsertId()
		comp.ID = int(id)
	}

	return comp, err
}

// dlUpdate updates fields on a company and returns full updated object
func dlUpdate(reqCtx common.RequestContext, id int, comp model.Company) (model.Company, error) {
	params := map[string]interface{}{
		"id":         id,
		"name":       comp.Name,
		"modifiedby": reqCtx.UserName(),
	}
	_, err := reqCtx.Tx().NamedExec(
		`update company
			set name = :name, modifiedby = :modifiedby
			where id = :id`, params)
	if err == nil {
		return dlGet(reqCtx, id)
	}
	return model.Company{}, err
}

// dlMarkInactive set the isactive flag to zero for this company
func dlMarkInactive(reqCtx common.RequestContext, id int) error {
	result, err := reqCtx.Tx().Exec(
		`update company 
		set isactive = 0, modifiedby = ?
		where id = ? and isactive = 1`,
		reqCtx.UserName(), id)
	if err == nil {
		if affect, _ := result.RowsAffected(); affect != 1 {
			err = eh.ErrDBNoRows
		}
	}
	return err
}

// dlDelete deletes this company
func dlDelete(reqCtx common.RequestContext, id int) error {
	result, err := reqCtx.Tx().Exec(`delete from company where id = ?`, id)
	if err == nil {
		if affect, _ := result.RowsAffected(); affect != 1 {
			err = eh.ErrDBNoRows
		}
	}
	return err
}

// dlSearch searches for companys
func dlSearch(reqCtx common.RequestContext, term string, limit int64) ([]model.Company, error) {
	data := []model.Company{}

	// searchCols := []string{"c.name"}
	// searchTerm := search.GetSearchCondition(searchCols, term)
	// isActiveWhere := ` c.isactive=1 `
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

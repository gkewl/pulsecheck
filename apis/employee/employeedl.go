package employee

import (
	"fmt"

	//	"github.com/gkewl/pulsecheck/apis/search"
	"github.com/gkewl/pulsecheck/common"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
)

const (
	getQuery = `select emp.id as id, emp.companyid as companyid, emp.firstname as firstname,
emp.middlename as middlename, emp.lastname as lastname,
emp.dateofbirth as dateofbirtht, emp.type as type, emp.isactive as
isactive, emp.createdby as createdby, emp.created as created,
emp.modifiedby as modifiedby, emp.modified as modified , es.consider consider
from employee emp  
join employeestatus es on es.employeeid = emp.id
`
)

// dlGet retrieves the specified employee
func dlGet(reqCtx common.RequestContext, id int64) (emp model.Employee, err error) {
	query := getQuery + ` where emp.id = ? and emp.isactive = 1`
	err = reqCtx.Tx().Get(&emp, query, id)
	return
}

// dlGetAll retrieves all employees
func dlGetAll(reqCtx common.RequestContext, limit int64, companyID int) (emps []model.Employee, err error) {
	emps = []model.Employee{}
	query := getQuery + ` where emp.companyid=? and emp.isactive=1  `

	if limit != 0 {
		query = query + fmt.Sprintf(" limit %d", limit)
	}
	err = reqCtx.Tx().Select(&emps, query, companyID)
	return
}

// dlCreate creates a employee
func dlCreate(reqCtx common.RequestContext, emp model.Employee) (model.Employee, error) {
	params := map[string]interface{}{
		"companyid":   emp.CompanyID,
		"firstname":   emp.Firstname,
		"middlename":  emp.Middlename,
		"lastname":    emp.Lastname,
		"dateofbirth": emp.DateofbirthT,
		"type":        emp.Type,
		"isactive":    1,
		"createdby":   reqCtx.UserName(),
		"modifiedby":  reqCtx.UserName(),
	}
	result, err := reqCtx.Tx().NamedExec(
		`insert into employee (companyid, firstname, middlename, lastname, dateofbirth, type,
isactive, createdby, modifiedby)
	 		values (:companyid, :firstname, :middlename, :lastname, :dateofbirth, :type,
:isactive, :createdby, :modifiedby)`,
		params)

	if err == nil {
		emp.ID, _ = result.LastInsertId()
	}

	return emp, err
}

// dlUpdate updates fields on a employee and returns full updated object
func dlUpdate(reqCtx common.RequestContext, id int64, emp model.Employee) (model.Employee, error) {
	params := map[string]interface{}{
		"id":          id,
		"companyid":   emp.CompanyID,
		"firstname":   emp.Firstname,
		"middlename":  emp.Middlename,
		"lastname":    emp.Lastname,
		"dateofbirth": emp.DateofbirthT,
		"type":        emp.Type,
		"modifiedby":  reqCtx.UserName(),
	}
	_, err := reqCtx.Tx().NamedExec(
		`update employee
			set companyid = :companyid, firstname = :firstname,
middlename = :middlename, lastname = :lastname,
dateofbirth = :dateofbirth, type = :type,
modifiedby = :modifiedby
			where id = :id`, params)
	if err == nil {
		return dlGet(reqCtx, id)
	}
	return model.Employee{}, err
}

// dlMarkInactive set the isactive flag to zero for this employee
func dlMarkInactive(reqCtx common.RequestContext, id int64) error {
	result, err := reqCtx.Tx().Exec(
		`update employee emp
		set isactive = 0, modifiedby = ?
		where emp.id = ? and emp.isactive = 1`,
		reqCtx.UserName(), id)
	if err == nil {
		if affect, _ := result.RowsAffected(); affect != 1 {
			err = eh.ErrDBNoRows
		}
	}
	return err
}

// dlDelete deletes this employee
func dlDelete(reqCtx common.RequestContext, id int64) error {
	result, err := reqCtx.Tx().Exec(`delete from employee where id = ?`, id)
	if err == nil {
		if affect, _ := result.RowsAffected(); affect != 1 {
			err = eh.ErrDBNoRows
		}
	}
	return err
}

// dlSearch searches for employees
func dlSearch(reqCtx common.RequestContext, term string, limit int64) ([]model.Employee, error) {
	data := []model.Employee{}

	// searchCols := []string{ "emp.firstname", "emp.middlename", "emp.lastname" }
	// searchTerm := search.GetSearchCondition(searchCols, term)
	// isActiveWhere := ` emp.isactive=1 `
	// qry := getQuery + ` where  ` +  isActiveWhere
	// if len(searchTerm) > 0 && len(isActiveWhere) > 0 {
	// 	qry +=  ` and `
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

// dlExists retrieves the specified employee
func dlExists(reqCtx common.RequestContext, emp model.Employee) (empOutput model.Employee, err error) {
	query := getQuery + ` where emp.firstname = ? and emp.lastname = ?
						  and  emp.dateofbirth=? and emp.companyid=?	
							`
	err = reqCtx.Tx().Get(&empOutput, query, emp.Firstname, emp.Lastname,
		emp.DateofbirthT, emp.CompanyID)
	return
}

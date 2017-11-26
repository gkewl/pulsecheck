package employeestatus

import (
	"fmt"

	//"stash.teslamotors.com/mos/factory/services/api/design/search"

	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/utilities"
	null "gopkg.in/guregu/null.v3"

	"github.com/gkewl/pulsecheck/constant"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
)

const (
	getQuery = `select es.id as id, es.employeeid as employeeid, 
				 es.oig , es.oiglastsearch,es.oigreference,
				 es.sam , es.samlastsearch, es.samreference, 
				 es.ofac , es.ofaclastsearch,es.ofacreference,  es.consider,
				 es.isactive as isactive, es.createdby as createdby, 
				 es.created as created,
			es.modifiedby as modifiedby, es.modified as modified from
			employeestatus es 
			join employee e on e.id = es.employeeid `
)

// dlGet retrieves the specified employeestatus
func dlGet(reqCtx common.RequestContext, employeeID int64) (es model.EmployeeStatus, err error) {
	query := getQuery + ` where es.employeeid = ? and es.isactive = 1`
	err = reqCtx.Tx().Get(&es, query, employeeID)
	return
}

// dlGetAll retrieves all employeestatuss
func dlGetAll(reqCtx common.RequestContext, limit int64) (ess []model.EmployeeStatus, err error) {
	ess = []model.EmployeeStatus{}
	query := getQuery + ` where es.isactive=1 and e.companyid = ?`

	if limit != 0 {
		query = query + fmt.Sprintf(" limit %d", limit)
	}
	err = reqCtx.Tx().Select(&ess, query, reqCtx.CompanyID())
	return
}

// dlCreate creates a employeestatus
func dlCreate(reqCtx common.RequestContext, es model.EmployeeStatus) (model.EmployeeStatus, error) {

	updateTime(&es)

	params := map[string]interface{}{
		"employeeid":     es.EmployeeID,
		"oig":            es.OIG,
		"oiglastsearch":  es.OIGLastSearch,
		"oigreference":   es.OIGReference,
		"sam":            es.Sam,
		"samlastsearch":  es.SamLastSearch,
		"samreference":   es.SamReference,
		"ofac":           es.Ofac,
		"ofaclastsearch": es.OfacLastSearch,
		"ofacreference":  es.OfacReference,
		"isactive":       1,
		"createdby":      reqCtx.UserName(),
		"modifiedby":     reqCtx.UserName(),
	}

	result, err := reqCtx.Tx().NamedExec(
		`insert into employeestatus (employeeid, oig, oiglastsearch, sam,samlastsearch,   ofac,ofaclastsearch,isactive, createdby,
										modifiedby)
	 	values (:employeeid, :oig, :oiglastsearch, :sam, :samlastsearch,:ofac, ofaclastsearch, :isactive, :createdby,
					:modifiedby)`,
		params)

	if err == nil {
		es.ID, _ = result.LastInsertId()
	}
	return dlGet(reqCtx, es.EmployeeID)
}

// dlUpdate updates fields on a employeestatus and returns full updated object
func dlUpdate(reqCtx common.RequestContext, employeeID int64, source string, value bool, referene string) (model.EmployeeStatus, error) {
	var setString string
	switch source {
	case constant.Source_OIG:
		setString = `set oig=?, oiglastsearch=? , oigreference=?`
	case constant.Source_SAM:
		setString = `set sam=?, samlastsearch=?  , samreference=?`
	case constant.Source_OFAC:
		setString = `set ofac=?, ofaclastsearch=? , ofacreference=?`
	default:
		return model.EmployeeStatus{}, eh.NewError(eh.ErrEmployeeStatusUpdate, "Source %s not implemented")
	}
	now := utilities.Now()
	query := `update employeestatus ` + setString + `, modifiedby =? where employeeid=?`
	_, err := reqCtx.Tx().Exec(query, value, now, referene, reqCtx.UserName(), employeeID)
	if err == nil {
		return dlGet(reqCtx, employeeID)
	}
	return model.EmployeeStatus{}, err
}

// dlMarkInactive set the isactive flag to zero for this employeestatus
func dlMarkInactive(reqCtx common.RequestContext, id int64) error {
	result, err := reqCtx.Tx().Exec(
		`update employeestatus es
		set isactive = 0, modifiedby = ?
		where es.id = ? and es.isactive = 1`,
		reqCtx.UserName(), id)
	if err == nil {
		if affect, _ := result.RowsAffected(); affect != 1 {
			err = eh.ErrDBNoRows
		}
	}
	return err
}

// dlDelete deletes this employeestatus
func dlDelete(reqCtx common.RequestContext, id int64) error {
	result, err := reqCtx.Tx().Exec(`delete from employeestatus where id = ?`, id)
	if err == nil {
		if affect, _ := result.RowsAffected(); affect != 1 {
			err = eh.ErrDBNoRows
		}
	}
	return err
}

// dlSearch searches for employeestatuss
func dlSearch(reqCtx common.RequestContext, term string, limit int64) ([]model.EmployeeStatus, error) {
	data := []model.EmployeeStatus{}

	// searchCols := []string{"es.source", "es.status"}
	// searchTerm := search.GetSearchCondition(searchCols, term)
	// isActiveWhere := ` es.isactive=1 `
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

func updateTime(es *model.EmployeeStatus) {
	if es.OIG {
		es.OIGLastSearch = null.TimeFrom(utilities.Now())
	}

	if es.Sam {
		es.SamLastSearch = null.TimeFrom(utilities.Now())
	}

	if es.Ofac {
		es.OfacLastSearch = null.TimeFrom(utilities.Now())
	}
}

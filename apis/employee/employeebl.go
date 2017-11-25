package employee

import (
	"fmt"

	"github.com/gkewl/pulsecheck/logger"

	"github.com/gkewl/pulsecheck/apis/elasticsearch"
	"github.com/gkewl/pulsecheck/apis/employeestatus"

	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/constant"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
	"github.com/gkewl/pulsecheck/utilities"
)

// BizLogic is the interface for all employee business logic methods
type BizLogic interface {
	Create(common.RequestContext, model.Employee) (model.Employee, error)
	Get(common.RequestContext, int64) (model.Employee, error)
	GetAll(common.RequestContext, int64) ([]model.Employee, error)
	Update(common.RequestContext, int64, model.Employee) (model.Employee, error)
	Delete(common.RequestContext, int64) (string, error)
	//upload multiple employees
	Upload(common.RequestContext, []model.Employee) ([]model.Employee, error)
	//search for one employee
	SearchStatus(common.RequestContext, model.Employee) (model.Employee, error)
	Search(common.RequestContext, string, int64) ([]model.Employee, error)
}

// BLEmployee implements the employee.BizLogic interface
type BLEmployee struct {
}

// Create will insert a new employee into the db
func (bl BLEmployee) Create(reqCtx common.RequestContext, emp model.Employee) (model.Employee, error) {
	var err error
	emp.CompanyID = reqCtx.CompanyID()
	emp.DateofbirthT, err = utilities.ParseStringToDate(emp.Dateofbirth)
	if err != nil {
		return model.Employee{}, eh.WrapError(eh.ErrEmployeeInsert, err, "Date parsing error ")
	}
	e, err := dlExists(reqCtx, emp)
	if err != nil || e.ID == 0 {
		e, err = dlCreate(reqCtx, emp)
		if err != nil {
			return model.Employee{}, eh.NewError(eh.ErrEmployeeInsert, "DB Error: "+err.Error())
		}
		// create employee status also here
		es := model.EmployeeStatus{
			EmployeeID: e.ID,
		}
		_, err = employeestatus.BLEmployeeStatus{}.Create(reqCtx, es)
		if err != nil {
			return model.Employee{}, eh.WrapError(eh.ErrEmployeeInsert, err, "Unable to create employee status")
		}
	}
	e.Dateofbirth = utilities.ParseDateToString(e.DateofbirthT)
	return e, nil

}

// Get returns a single employee by primary key
func (bl BLEmployee) Get(reqCtx common.RequestContext, id int64) (emp model.Employee, err error) {
	emp, err = dlGet(reqCtx, id)

	if err != nil || emp.ID == 0 {
		return model.Employee{}, eh.NewErrorNotFound(eh.ErrEmployeeDataNotFound, err, `Employee not found: id %d`, id)
	}

	emp.Dateofbirth = utilities.ParseDateToString(emp.DateofbirthT)
	if emp.CompanyID != reqCtx.CompanyID() {
		return model.Employee{}, eh.NewError(eh.ErrEmployeeDelete, fmt.Sprintf("Employee id %d  not available in company %d", id, emp.CompanyID))

	}
	return
}

// GetAll will return all employees
func (bl BLEmployee) GetAll(reqCtx common.RequestContext, limit int64) (emps []model.Employee, err error) {

	emps, err = dlGetAll(reqCtx, limit, reqCtx.CompanyID())
	if err != nil {
		return []model.Employee{}, eh.NewError(eh.ErrEmployeeDataNotFound, "DB Error: "+err.Error())
	}

	for idx, emp := range emps {
		emps[idx].Dateofbirth = utilities.ParseDateToString(emp.DateofbirthT)

	}
	return
}

// Update updates a single employee
func (bl BLEmployee) Update(reqCtx common.RequestContext, id int64, emp model.Employee) (model.Employee, error) {
	var err error
	emp.CompanyID = reqCtx.CompanyID()
	emp.DateofbirthT, err = utilities.ParseStringToDate(emp.Dateofbirth)
	if err != nil {
		return model.Employee{}, eh.WrapError(eh.ErrEmployeeUpdate, err, "Date parsing error ")
	}

	result, err := dlUpdate(reqCtx, id, emp)
	if err != nil {
		return model.Employee{}, eh.NewError(eh.ErrEmployeeUpdate, "DB Error: "+err.Error())
	}
	return result, err
}

// Delete marks a single employee inactive
func (bl BLEmployee) Delete(reqCtx common.RequestContext, id int64) (string, error) {
	e, err := bl.Get(reqCtx, id)
	if err != nil {
		return constant.ResultFail, eh.NewError(eh.ErrEmployeeDelete, "DB Error: "+err.Error())
	}

	if e.CompanyID != reqCtx.CompanyID() {
		return constant.ResultFail, eh.NewError(eh.ErrEmployeeDelete, fmt.Sprintf("Employee id %d  not available in company %d", id, e.CompanyID))
	}

	err = dlMarkInactive(reqCtx, id)
	if err != nil {
		return constant.ResultFail, eh.NewError(eh.ErrEmployeeDelete, "DB Error: "+err.Error())
	}
	return constant.ResultOk, err
}

// HardDelete physically deletes a employee, usually for testing
func (bl BLEmployee) HardDelete(reqCtx common.RequestContext, id int64) (string, error) {
	err := dlDelete(reqCtx, id)
	if err != nil {
		return constant.ResultFail, eh.NewError(eh.ErrEmployeeDelete, "DB Error: "+err.Error())
	}
	return constant.ResultOk, err
}

// Search finds employees matching the term
func (bl BLEmployee) Search(reqCtx common.RequestContext, term string, limit int64) (emps []model.Employee, err error) {
	emps, err = dlSearch(reqCtx, term, limit)
	if err != nil {
		return []model.Employee{}, eh.NewError(eh.ErrEmployeeDataNotFound, "DB Error: "+err.Error())
	}
	for idx, emp := range emps {
		emps[idx].Dateofbirth = utilities.ParseDateToString(emp.DateofbirthT)
	}

	return
}

// Exists verifies whether employee exists in the company
func (bl BLEmployee) exists(reqCtx common.RequestContext, emp model.Employee) (ret bool) {
	var err error
	emp.DateofbirthT, err = utilities.ParseStringToDate(emp.Dateofbirth)
	output, err := dlExists(reqCtx, emp)
	if err != nil {
		return
	}
	if err != nil {
		return

	} else if output.ID > 0 {
		ret = true
	}
	return
}

//Upload - creates not existing employees
func (bl BLEmployee) Upload(reqCtx common.RequestContext, employees []model.Employee) ([]model.Employee, error) {

	output := []model.Employee{}

	for _, emp := range employees {
		// e1, err := bl.Create(reqCtx, emp)
		// if err != nil {
		// 	return []model.Employee{}, eh.WrapError(eh.ErrEmployeeUpload, err, "error in uploading employee. %s", e1.ToString())
		// }

		e2, err := bl.SearchStatus(reqCtx, emp)
		if err != nil {
			return []model.Employee{}, eh.WrapError(eh.ErrEmployeeUpload, err, "error in uploading employee. %s", emp.ToString())
		}
		output = append(output, e2)
	}

	return output, nil
}

//SearchStatus - search employee in employee table or creates employee
//and get the result from elastic
func (bl BLEmployee) SearchStatus(reqCtx common.RequestContext, emp model.Employee) (model.Employee, error) {
	e1, err := bl.Create(reqCtx, emp)
	if err != nil {
		return model.Employee{}, eh.WrapError(eh.ErrEmployeeSearch, err, "error in searching employee. %s", e1.ToString())
	}
	//get the sources for this empployee to be searched
	// search in elastic
	sources := []string{constant.Source_OIG}
	var retErr error
	for _, src := range sources {
		err = bl.searchESAndUpdate(reqCtx, e1, src)
		if err != nil {
			logger.LogError(fmt.Sprintf("Error in searching employee id %d for source %s", e1.ID, src), reqCtx.Xid())
			retErr = err
		}
	}
	if retErr == nil {
		return bl.Get(reqCtx, e1.ID)
	}

	return model.Employee{}, eh.WrapError(eh.ErrEmployeeSearch, retErr, "error in searching employee. %s", e1.ToString())
}

// searchESAndUpdate - search elastic search for source and update for employee status
func (bl BLEmployee) searchESAndUpdate(reqCtx common.RequestContext, emp model.Employee, source string) error {
	
	
	result, err := elasticsearch.BLElasticSearch{}.SearchOne(reqCtx, emp, source)
	if err == nil {
		if len(result) == 0 {
			//update the result as not found
			employeestatus.BLEmployeeStatus{}.Update(reqCtx, emp.ID, source, false)
		}
		for idx, res := range result {
			if idx > 0 {
				//log to check for alternate way
				logger.LogInfo(fmt.Sprintf("Getting more result for the employee id %d   , result %+v", emp.ID, res), reqCtx.Xid())
				continue
			}

			//update employee status
			_, err = employeestatus.BLEmployeeStatus{}.Update(reqCtx, emp.ID, source, true)
			if err != nil {
				return err
			}

			//insert into employee transaction

		}
	}

	return err
}

package employeestatus

import (
	//	"github.com/gkewl/pulsecheck/apis/companysource"
	"github.com/gkewl/pulsecheck/apis/companysource"
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/constant"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
)

// BizLogic is the interface for all employeestatus business logic methods
type BizLogic interface {
	Create(common.RequestContext, model.EmployeeStatus) (model.EmployeeStatus, error)
	Get(common.RequestContext, int64) (model.EmployeeStatus, error)

	GetAll(common.RequestContext, int64) ([]model.EmployeeStatus, error)
	Update(common.RequestContext, int64, string, bool, string) (model.EmployeeStatus, error)
	Delete(common.RequestContext, int64) (string, error)

	//Search(common.RequestContext, string, int64) ([]model.EmployeeStatus, error)
}

// BLEmployeeStatus implements the employeestatus.BizLogic interface
type BLEmployeeStatus struct {
}

// Create will insert a new employeestatus into the db
func (bl BLEmployeeStatus) Create(reqCtx common.RequestContext, es model.EmployeeStatus) (model.EmployeeStatus, error) {
	es, err := dlCreate(reqCtx, es)
	if err != nil {
		return model.EmployeeStatus{}, eh.NewError(eh.ErrEmployeeStatusInsert, "DB Error: "+err.Error())
	}
	return es, err
}

// Get returns a single employeestatus by primary key
func (bl BLEmployeeStatus) Get(reqCtx common.RequestContext, employeeID int64) (es model.EmployeeStatus, err error) {
	es, err = dlGet(reqCtx, employeeID)

	if err != nil || es.ID == 0 {
		return model.EmployeeStatus{}, eh.NewErrorNotFound(eh.ErrEmployeeStatusDataNotFound, err, `EmployeeStatus not found: employeeID %d`, employeeID)
	}

	// Get the sources for the company
	srcs, err := companysource.BLCompanySource{}.GetForCompany(reqCtx, reqCtx.CompanyID())
	if err != nil {
		return model.EmployeeStatus{}, eh.WrapError(eh.ErrEmployeeStatusDataNotFound, err, `Sources not found for Company %d`, reqCtx.CompanyID)
	}

	//shows the sources
	es.Sources = es.ToSourceDetail(srcs)
	return
}

// GetAll will return all employeestatuss
func (bl BLEmployeeStatus) GetAll(reqCtx common.RequestContext, limit int64) (ess []model.EmployeeStatus, err error) {
	ess, err = dlGetAll(reqCtx, limit)
	if err != nil {
		return []model.EmployeeStatus{}, eh.NewError(eh.ErrEmployeeStatusDataNotFound, "DB Error: "+err.Error())
	}

	// Get the sources for the company
	srcs, err := companysource.BLCompanySource{}.GetForCompany(reqCtx, reqCtx.CompanyID())
	if err != nil {
		return []model.EmployeeStatus{}, eh.WrapError(eh.ErrEmployeeStatusDataNotFound, err, `Sources not found for Company %d`, reqCtx.CompanyID)
	}

	for idx, es := range ess {
		//shows the sources
		ess[idx].Sources = es.ToSourceDetail(srcs)
	}

	return
}

// Update updates a single employeestatus
func (bl BLEmployeeStatus) Update(reqCtx common.RequestContext, employeeID int64, source string, value bool, reference string) (model.EmployeeStatus, error) {

	s, err := bl.Get(reqCtx, employeeID)
	if err != nil {
		if eh.HasNoRowsError(err) {
			s = model.EmployeeStatus{EmployeeID: employeeID}
			s, err = bl.Create(reqCtx, s)
			if err != nil {
				return model.EmployeeStatus{}, eh.WrapError(eh.ErrEmployeeStatusUpdate, err, "")
			}
		}
	}

	result, err := dlUpdate(reqCtx, employeeID, source, value, reference)
	if err != nil {
		return model.EmployeeStatus{}, eh.NewError(eh.ErrEmployeeStatusUpdate, "DB Error: "+err.Error())
	}
	return result, err
}

// Delete marks a single employeestatus inactive
func (bl BLEmployeeStatus) Delete(reqCtx common.RequestContext, id int64) (string, error) {
	// todo: add validation here

	err := dlMarkInactive(reqCtx, id)
	if err != nil {
		return constant.ResultFail, eh.NewError(eh.ErrEmployeeStatusDelete, "DB Error: "+err.Error())
	}
	return constant.ResultOk, err
}

// HardDelete physically deletes a employeestatus, usually for testing
func (bl BLEmployeeStatus) HardDelete(reqCtx common.RequestContext, id int64) (string, error) {
	err := dlDelete(reqCtx, id)
	if err != nil {
		return constant.ResultFail, eh.NewError(eh.ErrEmployeeStatusDelete, "DB Error: "+err.Error())
	}
	return constant.ResultOk, err
}

// Search finds employeestatuss matching the term
// func (bl BLEmployeeStatus) Search(reqCtx common.RequestContext, term string, limit int64) (ess []model.EmployeeStatus, err error) {
// 	ess, err = dlSearch(reqCtx, term, limit)
// 	if err != nil {
// 		return []model.EmployeeStatus{}, eh.NewError(eh.ErrEmployeeStatusDataNotFound, "DB Error: "+err.Error())
// 	}
// 	return
// }

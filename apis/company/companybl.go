package company

import (
	"github.com/gkewl/pulsecheck/common"
	eh "github.com/gkewl/pulsecheck/errorhandler"

	"github.com/gkewl/pulsecheck/constant"
	"github.com/gkewl/pulsecheck/model"
)

// BizLogic is the interface for all company business logic methods
type BizLogic interface {
	Create(common.RequestContext, model.Company) (model.Company, error)
	Get(common.RequestContext, int) (model.Company, error)
	GetAll(common.RequestContext, int) ([]model.Company, error)
	Update(common.RequestContext, int, model.Company) (model.Company, error)
	Delete(common.RequestContext, int) (string, error)
	Search(common.RequestContext, string, int) ([]model.Company, error)
}

// BLCompany implements the company.BizLogic interface
type BLCompany struct {
}

// Create will insert a new company into the db
func (bl BLCompany) Create(reqCtx common.RequestContext, comp model.Company) (model.Company, error) {

	comp, err := dlCreate(reqCtx, comp)
	if err != nil {
		return model.Company{}, eh.NewError(eh.ErrCompanyInsert, "DB Error: "+err.Error())
	}
	return comp, err
}

// Get returns a single company by primary key
func (bl BLCompany) Get(reqCtx common.RequestContext, id int) (comp model.Company, err error) {
	comp, err = dlGet(reqCtx, id)

	if err != nil || comp.ID == 0 {
		return model.Company{}, eh.NewErrorNotFound(eh.ErrCompanyDataNotFound, err, `Company not found: id %d`, id)
	}
	return
}

// GetAll will return all companys
func (bl BLCompany) GetAll(reqCtx common.RequestContext, limit int) (comps []model.Company, err error) {
	comps, err = dlGetAll(reqCtx, limit)
	if err != nil {
		return []model.Company{}, eh.NewError(eh.ErrCompanyDataNotFound, "DB Error: "+err.Error())
	}
	return
}

// Update updates a single company
func (bl BLCompany) Update(reqCtx common.RequestContext, id int, comp model.Company) (model.Company, error) {
	// todo: add validation here

	result, err := dlUpdate(reqCtx, id, comp)
	if err != nil {
		return model.Company{}, eh.NewError(eh.ErrCompanyUpdate, "DB Error: "+err.Error())
	}
	return result, err
}

// Delete marks a single company inactive
func (bl BLCompany) Delete(reqCtx common.RequestContext, id int) (string, error) {
	// todo: add validation here

	err := dlMarkInactive(reqCtx, id)
	if err != nil {
		return constant.ResultFail, eh.NewError(eh.ErrCompanyDelete, "DB Error: "+err.Error())
	}
	return constant.ResultOk, err
}

// HardDelete physically deletes a company, usually for testing
func (bl BLCompany) HardDelete(reqCtx common.RequestContext, id int) (string, error) {
	err := dlDelete(reqCtx, id)
	if err != nil {
		return constant.ResultFail, eh.NewError(eh.ErrCompanyDelete, "DB Error: "+err.Error())
	}
	return constant.ResultOk, err
}

// Search finds companys matching the term
func (bl BLCompany) Search(reqCtx common.RequestContext, term string, limit int) (comps []model.Company, err error) {
	// comps, err = dlSearch(reqCtx, term, limit)
	// if err != nil {
	// 	return []model.Company{}, eh.NewError(eh.ErrCompanyDataNotFound, "DB Error: "+err.Error())
	// }
	return
}

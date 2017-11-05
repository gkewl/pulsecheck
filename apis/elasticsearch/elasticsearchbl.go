package elasticsearch

import (
	"github.com/gkewl/pulsecheck/common"
	eh "github.com/gkewl/pulsecheck/errorhandler"

	"github.com/gkewl/pulsecheck/constant"
	"github.com/gkewl/pulsecheck/model"
)

// BizLogic is the interface for all company business logic methods
type BizLogic interface {
	Search(common.RequestContext, string, int) ([]model.Company, error)
}

// BLCompany implements the company.BizLogic interface
type BLElasticSearch struct {
}

// Search finds companys matching the term
func (bl BLElasticSearch) Search(reqCtx common.RequestContext, emps []model.Employee ) (result []model.ElasticSearchResult, err error) {
	// comps, err = dlSearch(reqCtx, term, limit)
	// if err != nil {
	// 	return []model.Company{}, eh.NewError(eh.ErrCompanyDataNotFound, "DB Error: "+err.Error())
	// }
	return
}

// SearchOne finds one employee in eleastic search
func (bl BLElasticSearch) SearchOne(reqCtx common.RequestContext, emp model.Employee ) (result []model.ElasticSearchResult, err error) {
	
	// TBD get the sourec required to search
	var sources []string

	sources  = append(sources, "OIG")

	for _, src := range sources {

	}

	return
}


func (bl BLElasticSearch)  getSearchModel ( emp model.Employee, source string) 
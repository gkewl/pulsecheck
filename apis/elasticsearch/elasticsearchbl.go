package elasticsearch

import (
	"fmt"

	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/constant"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/logger"
	"github.com/gkewl/pulsecheck/model"
)

// BizLogic is the interface for all company business logic methods
type BizLogic interface {
	Search(common.RequestContext, []model.Employee) ([]model.ElasticSearchResult, error)
	SearchOIG(reqCtx common.RequestContext, oig model.OIGSearch)
}

// BLElasticSearch implements the company.BizLogic interface
type BLElasticSearch struct {
}

// Search finds companys matching the term
func (bl BLElasticSearch) Search(reqCtx common.RequestContext, emps []model.Employee) (result []model.ElasticSearchResult, err error) {
	// comps, err = dlSearch(reqCtx, term, limit)
	// if err != nil {
	// 	return []model.Company{}, eh.NewError(eh.ErrCompanyDataNotFound, "DB Error: "+err.Error())
	// }
	return
}

// SearchOne finds one employee in eleastic search
func (bl BLElasticSearch) SearchOne(reqCtx common.RequestContext, emp model.Employee, source string) (result []model.ElasticSearchResult, err error) {
	logger.LogInfo(fmt.Sprintf("Searching for employee %d  Source %s", emp.ID, source), reqCtx.Xid())
	input := bl.getSearchModel(reqCtx, emp, source)
	switch source {
	case constant.Source_OIG:
		result, err = bl.SearchOIG(reqCtx, input)
	default:
		logger.LogError(fmt.Sprintf("Source %s not implemented", source), reqCtx.Xid())
		err = eh.NewError(eh.ErrElasticSearchNotImplemented, "")
	}
	if len(result) > 0 {
		logger.LogInfo(fmt.Sprintf("employee %d found in Source %s  Result: %+v ", emp.ID, source, result), reqCtx.Xid())

	} else {
		logger.LogInfo(fmt.Sprintf("employee %d NOT found in Source %s ", emp.ID, source), reqCtx.Xid())
	}
	return
}

func (bl BLElasticSearch) getSearchModel(reqCtx common.RequestContext, emp model.Employee, source string) model.OIGSearch {
	switch source {
	case constant.Source_OIG:
		return emp.ToOIG()
	default:
		logger.LogError(fmt.Sprintf("Source %s not implemented", source), reqCtx.Xid())
	}
	return model.OIGSearch{}
}

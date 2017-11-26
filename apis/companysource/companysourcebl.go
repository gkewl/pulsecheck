package companysource

import (
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/constant"
)

// BizLogic is the interface for all company business logic methods
type BizLogic interface {
	GetForCompany(common.RequestContext, int) ([]string, error)
}

// BLCompanySource implements the companysource.BizLogic interface
type BLCompanySource struct {
}

// GetForCompany -
func (bl BLCompanySource) GetForCompany(reqCtx common.RequestContext, companyID int) ([]string, error) {
	return []string{constant.Source_OIG, constant.Source_OFAC, constant.Source_SAM}, nil
}

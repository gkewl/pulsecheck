package company

import (
	"net/http"

	"github.com/gkewl/pulsecheck/common"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
)

// TestingBizLogic will be set in testing to a MockBL
var TestingBizLogic BizLogic

type companyControllerFunc func(reqCtx common.RequestContext, companyBL BizLogic, companyInput model.Company) (interface{}, error)

// ControlWrapper extracts information from the request and calls the wrapped
// controller function
func ControlWrapper(handler companyControllerFunc) func(common.RequestContext) (interface{}, error) {
	return func(reqCtx common.RequestContext) (interface{}, error) {
		companyInterface := BizLogic(BLCompany{})
		if TestingBizLogic != nil {
			companyInterface = TestingBizLogic
		}
		comp := model.Company{}
		err := reqCtx.Scan("comp", &comp)
		if err != nil {
			return nil, eh.NewError(eh.ErrJsonDecodeFail, "JSON Error: "+err.Error())
		}
		return handler(reqCtx, companyInterface, comp)
	}
}

// GetRoutes returns all company-related routes
func GetRoutes() common.Routes {

	return common.Routes{

		common.Route{
			Name:           "SearchCompany",
			Method:         "GET",
			Pattern:        "/company/search",
			ControllerFunc: ControlWrapper(SearchCompany),
		},

		common.Route{
			Name:           "GetCompany",
			Method:         "GET",
			Pattern:        "/company/{id}",
			ControllerFunc: ControlWrapper(GetCompany),
		},
		common.Route{
			Name:           "GetCompanys",
			Method:         "GET",
			Pattern:        "/company",
			ControllerFunc: ControlWrapper(GetAllCompanys),
		},
		common.Route{
			Name:           "CreateCompany",
			Method:         "POST",
			Pattern:        "/company",
			ControllerFunc: ControlWrapper(CreateCompany),
			//	AuthRequired:   constant.Admin,
			NormalHttpCode: http.StatusCreated,
		},
		common.Route{
			Name:           "UpdateCompany",
			Method:         "PUT",
			Pattern:        "/company/{id}",
			ControllerFunc: ControlWrapper(UpdateCompany),
			//	AuthRequired:   constant.Admin,
		},
		common.Route{
			Name:           "DeleteCompany",
			Method:         "DELETE",
			Pattern:        "/company/{id}",
			ControllerFunc: ControlWrapper(DeleteCompany),
			//	AuthRequired:   constant.Admin,
		},
	}
}

// GetCompany gets a company by primary key
func GetCompany(reqCtx common.RequestContext, companyBL BizLogic, comp model.Company) (interface{}, error) {
	id := reqCtx.IntValue32("id", 0)
	return companyBL.Get(reqCtx, id)
}

// GetAllCompanys gets all companys
func GetAllCompanys(reqCtx common.RequestContext, companyBL BizLogic, comp model.Company) (interface{}, error) {
	limit := reqCtx.IntValue32("limit", 50)
	return companyBL.GetAll(reqCtx, limit)
}

// SearchCompany finds companys that match a term
func SearchCompany(reqCtx common.RequestContext, companyBL BizLogic, comp model.Company) (interface{}, error) {
	return companyBL.Search(reqCtx, reqCtx.Value("term", ""), reqCtx.IntValue32("limit", 50))
}

// CreateCompany creates a company
func CreateCompany(reqCtx common.RequestContext, companyBL BizLogic, comp model.Company) (interface{}, error) {
	return companyBL.Create(reqCtx, comp)
}

// UpdateCompany updates a single company
func UpdateCompany(reqCtx common.RequestContext, companyBL BizLogic, comp model.Company) (interface{}, error) {
	id := reqCtx.IntValue32("id", 0)
	return companyBL.Update(reqCtx, id, comp)
}

// DeleteCompany deletes a single company
func DeleteCompany(reqCtx common.RequestContext, companyBL BizLogic, comp model.Company) (interface{}, error) {
	id := reqCtx.IntValue32("id", 0)
	return companyBL.Delete(reqCtx, id)
}

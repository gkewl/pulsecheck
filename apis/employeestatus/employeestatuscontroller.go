package employeestatus

import (
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/constant"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
)

// TestingBizLogic will be set in testing to a MockBL
var TestingBizLogic BizLogic

type employeestatusControllerFunc func(reqCtx common.RequestContext, employeestatusBL BizLogic, employeestatusInput model.EmployeeStatus) (interface{}, error)

// ControlWrapper extracts information from the request and calls the wrapped
// controller function
func ControlWrapper(handler employeestatusControllerFunc) func(common.RequestContext) (interface{}, error) {
	return func(reqCtx common.RequestContext) (interface{}, error) {
		employeestatusInterface := BizLogic(BLEmployeeStatus{})
		if TestingBizLogic != nil {
			employeestatusInterface = TestingBizLogic
		}
		es := model.EmployeeStatus{}
		err := reqCtx.Scan("es", &es)
		if err != nil {
			return nil, eh.NewError(eh.ErrJsonDecodeFail, "JSON Error: "+err.Error())
		}
		return handler(reqCtx, employeestatusInterface, es)
	}
}

// GetRoutes returns all employeestatus-related routes
func GetRoutes() common.Routes {

	return common.Routes{

		common.Route{
			Name:           "GetEmployeeStatus",
			Method:         "GET",
			Pattern:        "/employeestatus/{id}",
			ControllerFunc: ControlWrapper(GetEmployeeStatus),
			AuthRequired:   constant.User,
		},
		common.Route{
			Name:           "GetEmployeeStatuss",
			Method:         "GET",
			Pattern:        "/employeestatus",
			ControllerFunc: ControlWrapper(GetAllEmployeeStatuss),
			AuthRequired:   constant.User,
		},
	}
}

// GetEmployeeStatus gets a employeestatus by primary key
func GetEmployeeStatus(reqCtx common.RequestContext, employeestatusBL BizLogic, es model.EmployeeStatus) (interface{}, error) {
	id := reqCtx.IntValue("id", 0)
	return employeestatusBL.Get(reqCtx, id)
}

// GetAllEmployeeStatuss gets all employeestatuss
func GetAllEmployeeStatuss(reqCtx common.RequestContext, employeestatusBL BizLogic, es model.EmployeeStatus) (interface{}, error) {
	limit := reqCtx.IntValue("limit", 50)
	return employeestatusBL.GetAll(reqCtx, limit)
}

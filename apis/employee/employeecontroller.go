package employee

import (
	"net/http"

	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/constant"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
)

// TestingBizLogic will be set in testing to a MockBL
var TestingBizLogic BizLogic

type employeeControllerFunc func(reqCtx common.RequestContext, employeeBL BizLogic, employeeInput model.Employee) (interface{}, error)

// ControlWrapper extracts information from the request and calls the wrapped
// controller function
func ControlWrapper(handler employeeControllerFunc) func(common.RequestContext) (interface{}, error) {
	return func(reqCtx common.RequestContext) (interface{}, error) {
		employeeInterface := getInterface()

		emp := model.Employee{}
		err := reqCtx.Scan("emp", &emp)
		if err != nil {
			return nil, eh.NewError(eh.ErrJsonDecodeFail, "JSON Error: "+err.Error())
		}
		return handler(reqCtx, employeeInterface, emp)
	}
}

func getInterface() BizLogic {
	employeeInterface := BizLogic(BLEmployee{})
	if TestingBizLogic != nil {
		employeeInterface = TestingBizLogic
	}
	return employeeInterface
}

// GetRoutes returns all employee-related routes
func GetRoutes() common.Routes {

	return common.Routes{

		// common.Route{
		// 	Name:           "SearchEmployee",
		// 	Method:         "GET",
		// 	Pattern:        "/employee/search",
		// 	ControllerFunc: ControlWrapper(SearchEmployee),
		// },

		common.Route{
			Name:           "GetEmployee",
			Method:         "GET",
			Pattern:        "/employee/{id}",
			ControllerFunc: ControlWrapper(GetEmployee),
			AuthRequired:   constant.User,
		},
		common.Route{
			Name:           "GetEmployees",
			Method:         "GET",
			Pattern:        "/employee",
			ControllerFunc: ControlWrapper(GetAllEmployees),
			AuthRequired:   constant.User,
		},
		common.Route{
			Name:           "CreateEmployee",
			Method:         "POST",
			Pattern:        "/employee",
			ControllerFunc: ControlWrapper(CreateEmployee),
			AuthRequired:   constant.User,
			NormalHttpCode: http.StatusCreated,
		},
		common.Route{
			Name:           "UploadEmployees",
			Method:         "POST",
			Pattern:        "/employee/upload",
			ControllerFunc: UploadEmployees,
			AuthRequired:   constant.User,
		},
		common.Route{
			Name:           "UpdateEmployee",
			Method:         "PUT",
			Pattern:        "/employee/{id}",
			ControllerFunc: ControlWrapper(UpdateEmployee),
			AuthRequired:   constant.User,
		},
		common.Route{
			Name:           "DeleteEmployee",
			Method:         "DELETE",
			Pattern:        "/employee/{id}",
			ControllerFunc: ControlWrapper(DeleteEmployee),
			AuthRequired:   constant.User,
		},
		common.Route{
			Name:           "SearchStatus",
			Method:         "POST",
			Pattern:        "/employee/searchstatus",
			ControllerFunc: ControlWrapper(SearchStatus),
			AuthRequired:   constant.User,
		},
	}
}

// GetEmployee gets a employee by primary key
func GetEmployee(reqCtx common.RequestContext, employeeBL BizLogic, emp model.Employee) (interface{}, error) {
	id := reqCtx.IntValue("id", 0)
	return employeeBL.Get(reqCtx, id)
}

// GetAllEmployees gets all employees
func GetAllEmployees(reqCtx common.RequestContext, employeeBL BizLogic, emp model.Employee) (interface{}, error) {
	limit := reqCtx.IntValue("limit", 50)
	return employeeBL.GetAll(reqCtx, limit)
}

// SearchEmployee finds employees that match a term
func SearchEmployee(reqCtx common.RequestContext, employeeBL BizLogic, emp model.Employee) (interface{}, error) {
	return employeeBL.Search(reqCtx, reqCtx.Value("term", ""), reqCtx.IntValue("limit", 1000))
}

// CreateEmployee creates a employee
func CreateEmployee(reqCtx common.RequestContext, employeeBL BizLogic, emp model.Employee) (interface{}, error) {
	return employeeBL.Create(reqCtx, emp)
}

// UpdateEmployee updates a single employee
func UpdateEmployee(reqCtx common.RequestContext, employeeBL BizLogic, emp model.Employee) (interface{}, error) {
	id := reqCtx.IntValue("id", 0)
	return employeeBL.Update(reqCtx, id, emp)
}

// DeleteEmployee deletes a single employee
func DeleteEmployee(reqCtx common.RequestContext, employeeBL BizLogic, emp model.Employee) (interface{}, error) {
	id := reqCtx.IntValue("id", 0)
	return employeeBL.Delete(reqCtx, id)
}

// SearchStatus finds employees that match a term
func SearchStatus(reqCtx common.RequestContext, employeeBL BizLogic, emp model.Employee) (interface{}, error) {
	return employeeBL.SearchStatus(reqCtx, emp)
}

// UploadEmployees uploads employees
func UploadEmployees(reqCtx common.RequestContext) (interface{}, error) {
	userInterface := getInterface()
	emps := []model.Employee{}
	err := reqCtx.Scan("emp", &emps)
	if err != nil {
		return nil, eh.NewError(eh.ErrJsonDecodeFail, "JSON Error: "+err.Error())
	}
	return userInterface.Upload(reqCtx, emps)
}

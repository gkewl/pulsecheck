package rroutes

import (
	"github.com/gkewl/pulsecheck/apis/company"
	"github.com/gkewl/pulsecheck/apis/employee"
	"github.com/gkewl/pulsecheck/apis/user"

	"github.com/gkewl/pulsecheck/authentication"

	"github.com/gkewl/pulsecheck/common"
)

var APIs = common.APIRoutes{

	authentication.GetRoutes(),
	company.GetRoutes(),
	user.GetRoutes(),
	employee.GetRoutes(),
}

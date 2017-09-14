package rroutes

import (
	"github.com/gkewl/pulsecheck/apis/actor"
	"github.com/gkewl/pulsecheck/authentication"
	"github.com/gkewl/pulsecheck/common"
)

var APIs = common.APIRoutes(

	authentication.GetRoutes(),
	actor.GetRoutes(),
)

package routehandler

import (
	"net/http"
	"pulsecheck/common"

	"github.com/gorilla/mux"
)

func NewRouter(ctx *common.AppContext, apis common.APIRoutes, subroute string) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	subrouter := router.PathPrefix(subroute).Subrouter()
	for _, api := range apis {
		for _, route := range api {
			var handler http.Handler

			handler = common.AppHandler{
				ctx,
				route.HandlerFunc,
			}

			subrouter.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(handler)

		}
	}

	return subrouter
}

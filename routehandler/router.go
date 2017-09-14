package routehandler

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"time"

	"github.com/gorilla/mux"

	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/config"
	"github.com/gkewl/pulsecheck/protocol"
)

func AttachProfiler(router *HTTPRouter) {
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/trace", pprof.Trace)

	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/debug/pprof/block", pprof.Handler("block"))

}

type HTTPRouter struct {
	*mux.Router

	actionRouteMap map[string]*common.Route
}

func (r *HTTPRouter) RouteForAction(action string) *common.Route {
	if val, present := r.actionRouteMap[action]; present {
		return val
	}
	return nil
}

func NewRouter(ctx *common.AppContext, apis common.APIRoutes, subroute string) *HTTPRouter {
	router := mux.NewRouter().StrictSlash(true)
	subrouter := &HTTPRouter{router.PathPrefix(subroute).Subrouter(), map[string]*common.Route{}}
	pathCheck := map[string]*common.Route{}
	timeout, err := time.ParseDuration(config.GetEnv(config.MOS_APP_TIMEOUT))
	if err != nil {
		fmt.Printf("FATAL ERROR: MOS_APP_TIMEOUT value '%s' could not be parsed\n", config.GetEnv(config.MOS_APP_TIMEOUT))
		os.Exit(1)
	}
	for _, api := range apis {
		for i := 0; i < len(api); i++ {
			route := &api[i]
			// Track all actions and prevent duplicates
			if subrouter.RouteForAction(route.Name) != nil {
				fmt.Printf("FATAL ERROR: duplicate route action %s %v %v\n", route.Name, *route, *subrouter.RouteForAction(route.Name))
				os.Exit(1)
			}
			subrouter.actionRouteMap[route.Name] = route

			// Check for dupe paths
			if prior, exists := pathCheck[route.Method+route.Pattern]; exists {
				fmt.Printf("FATAL ERROR: duplicate route pattern %s %v %v\n", route.Name, *route, *prior)
				os.Exit(1)
			}
			pathCheck[route.Method+route.Pattern] = route

			// Construct the actual HTTP handler for this action
			var handler http.Handler
			if route.ControllerFunc == nil {
				handler = common.AppHandler{
					ctx,
					route.HandlerFunc,
					route,
				}
			} else {
				handler = protocol.NewHTTPHandler{
					AppContext: ctx,
					Route:      route,
				}
			}
			handler = common.RecoverHandler(handler)
			if !route.Streaming {
				handler = http.TimeoutHandler(handler, timeout, "Request timed out")
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

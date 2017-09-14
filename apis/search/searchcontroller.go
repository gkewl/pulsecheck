package search

import (
	"github.com/gkewl/pulsecheck/common"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
	"github.com/gkewl/pulsecheck/utilities"
	"net/http"
)

func GetRoutes() common.Routes {

	return common.Routes{

		common.Route{
			Name:        "Search",
			Method:      "GET",
			Pattern:     "/search",
			HandlerFunc: SearchAny,
		},
	}
}

func Search(iSearch SearchInterface, entity string, term string, ctx *common.AppContext) ([]model.NameDescription, error) {

	result, err := iSearch.Search(entity, term, ctx)
	return result, err
}

func SearchAny(ctx *common.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {

	entity := r.FormValue("entity")
	input := r.FormValue("term")

	di := DBSearch{}

	output, err := Search(di, entity, input, ctx)

	if err != nil {
		return 0, eh.RespWithError(w, r, input, err)
	}
	utilities.WriteJSONStructuredResponse(r, w, http.StatusOK, output)
	return 0, nil
}

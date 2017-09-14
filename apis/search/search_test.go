package search_test

import (
	"fmt"
	"github.com/gkewl/pulsecheck/apis/search"
	"github.com/gkewl/pulsecheck/authentication"
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/config"
	"github.com/gkewl/pulsecheck/dbhandler"
	"github.com/gkewl/pulsecheck/model"
	"github.com/gkewl/pulsecheck/routehandler"
	"github.com/gkewl/pulsecheck/utilities"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	server *httptest.Server
	reader io.Reader //Ignore this for now
)

var apis = common.APIRoutes{
	search.GetRoutes(),
}
var ctx common.AppContext

func init() {

	config.LoadConfigurations()
	var err error
	ctx = common.AppContext{}
	ctx.Db, err = dbhandler.CreateConnection()
	if err != nil {

	}

	authBackend := authentication.InitJWTAuthenticationBackend()
	var role string = "Admin"
	token, _ = authBackend.GenerateToken("sspade", 1, "USER", role)
	router := routehandler.NewRouter(&ctx, apis, "/api/v1")
	server = httptest.NewServer(router) //Creating new server with the user handlers
}

func TestSearchModule_ActorAPI(t *testing.T) {

	term := "raj" //bmx%20metr"
	entity := "actor"
	apiURL := fmt.Sprintf("%s/api/v1/search?entity=%s&term=%s", server.URL, entity, term) //Grab the address for the API endpoint
	reader = strings.NewReader("")

	request, err := http.NewRequest("GET", apiURL, reader)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token.Token))
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 200 {
		t.Errorf("Expected: 200  Actual : %d", res.StatusCode)
	}

	var output []model.NameDescription

	utilities.ScanResponseObject(res.Body, &output)

	if len(output) == 0 {
		t.Error("Error in getting data for search term")
	}

}

func TestSearchModule_TaskAPI(t *testing.T) {

	term := "bmx" //bmx%20metr"
	entity := "task"
	apiURL := fmt.Sprintf("%s/api/v1/search?entity=%s&term=%s", server.URL, entity, term) //Grab the address for the API endpoint
	reader = strings.NewReader("")

	request, err := http.NewRequest("GET", apiURL, reader)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token.Token))
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 200 {
		t.Errorf("Expected: 200  Actual : %d", res.StatusCode)
	}

	var output []model.NameDescription

	utilities.ScanResponseObject(res.Body, &output)

	if len(output) == 0 {
		t.Error("Error in getting data for search term")
	}

}

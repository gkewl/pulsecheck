package actor_test

import (
	"net/http/httptest"

	"github.com/gkewl/pulsecheck/apis/actor"
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/config"
	"github.com/gkewl/pulsecheck/dbhandler"
	"github.com/gkewl/pulsecheck/protocol"
	"github.com/gkewl/pulsecheck/routehandler"
)

var (
	server     *httptest.Server
	router     *routehandler.HTTPRouter
	baseUrl    string
	mqttClient *protocol.MQTTClient
)

var apis = common.APIRoutes{
	actor.GetRoutes(),
}

func init() {
	config.LoadConfigurations()

	ctx := common.AppContext{}
	ctx.Db, _ = dbhandler.CreateConnection()

	router = routehandler.NewRouter(&ctx, apis, "/api/v1")
	server = httptest.NewServer(router) //Creating new server with the user handlers
	baseUrl = server.URL + "/api/v1"

	mqttClient = protocol.NewMQTTClient(&ctx, apis)
}

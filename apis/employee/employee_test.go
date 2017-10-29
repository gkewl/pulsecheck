package employee_test

import (
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gkewl/pulsecheck/apis/employee"

	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/config"
	"github.com/gkewl/pulsecheck/dbhandler"
	"github.com/gkewl/pulsecheck/protocol"
	"github.com/gkewl/pulsecheck/routehandler"
)

var (
	httpCaller protocol.HTTPTestCaller
	mqttCaller protocol.MQTTTestCaller
	server     *httptest.Server
	router     *routehandler.HTTPRouter
	baseUrl    string
	mqttClient *protocol.MQTTClient
)

var apis = common.APIRoutes{
	employee.GetRoutes(),
}

var _ = BeforeSuite(func() {
})

var _ = AfterSuite(func() {
})

func TestSuite(t *testing.T) {
	httpCaller = protocol.HTTPTestCaller{BaseURL: baseUrl, Router: router}
	mqttCaller = protocol.MQTTTestCaller{Client: mqttClient}

	RegisterFailHandler(Fail)
	RunSpecs(t, "SPARQ Employee Test Suite")
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

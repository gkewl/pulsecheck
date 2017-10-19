package authentication_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSuite(t *testing.T) {
	// httpCaller = protocol.HTTPTestCaller{BaseURL: baseUrl, Router: router}
	// mqttCaller = protocol.MQTTTestCaller{Client: mqttClient}

	RegisterFailHandler(Fail)
	RunSpecs(t, "SPARQ Authentication Test Suite")
}

var _ = BeforeSuite(func() {

})

var _ = AfterSuite(func() {

})

package protocol

import (
	"github.com/gkewl/pulsecheck/common"
)

// MQTTClient is a dummy client
type MQTTClient struct {
	Router *MQTTRouter
}

// NewMQTTClient returns a new dummy client
func NewMQTTClient(appCtx *common.AppContext, apis common.APIRoutes) *MQTTClient {
	mc := MQTTClient{Router: NewMQTTRouter(appCtx, apis)}
	return &mc
}

// Receive accepts an inbound request and returns the topic and payload for a
// response to publish
func (mc *MQTTClient) Receive(topic, payload string) (replyTopic, replyPayload string) {
	return mc.Router.Process(topic, payload)
}

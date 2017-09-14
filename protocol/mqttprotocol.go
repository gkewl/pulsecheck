package protocol

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gkewl/pulsecheck/constant"
	"strings"
	"time"

	auth "github.com/gkewl/pulsecheck/authentication"
	"github.com/gkewl/pulsecheck/common"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/logger"
)

// MQTTRouter is a singleton router for inbound MQTT requests
type MQTTRouter struct {
	appCtx         *common.AppContext
	actionRouteMap map[string]*common.Route
}

// Process handles processing of a single inbound MQTT request should be called
// in a go routine
func (r *MQTTRouter) Process(topic, payload string) (replyTopic, replyPayload string) {
	start := time.Now()
	reqCtx, err := NewMQTTRequestContext(r.appCtx, topic, payload)
	if err != nil {
		return "", ""
	}
	route := r.RouteForAction(reqCtx.Action())
	var result interface{}
	if route == nil {
		err = errors.New("unknown action") //todo: better error
	} else {
		if route.AuthRequired != constant.Guest {
			if auth.RoleIdForName(reqCtx.UserRole()) < route.AuthRequired {
				err = eh.NewError(eh.ErrUnAuthorizedUserForAPI, "")
			}
		}
		if err == nil {
			retry := true
			for retry {
				result, err = route.ControllerFunc(&reqCtx)
				retry, err = ProcessDeadlock(&reqCtx, err)
			}
		}
	}
	replyMap := map[string]interface{}{}
	replyMap["_meta_"] = reqCtx.Metadata()
	if err == nil {
		err = reqCtx.Tx().Commit()
	}
	if err != nil {
		_, content := StructuredError(err, reqCtx.Xid())
		replyMap["error"] = content.Error
		_ = reqCtx.Tx().Rollback()
	} else {
		replyMap["response"] = result
	}
	RunDeferredRequests(&reqCtx, (err == nil))
	fields := map[string]interface{}{
		"protocol": "mqtt",
		"action":   route.Name,
		"user":     reqCtx.UserName(),
	}
	for k, v := range *reqCtx.LogValues() {
		fields[k] = v
	}
	lm := logger.LogModel{
		Level:    logger.INFO,
		Msg:      "success",
		Xid:      reqCtx.Xid(),
		Fields:   fields,
		Duration: time.Since(start).Seconds(),
	}
	if err == nil {
		logger.Log(lm)
	} else {
		lm.Level = logger.ERROR
		lm.Err = err
		lm.Msg = "failure"
		if ee, ok := err.(eh.Error); ok {
			lm.Caller = ee.LocationStack("; ")
			lm.Msg = ee.DetailStack("; ")
		}
		logger.ErrorLog(lm)
	}
	replyTopic = reqCtx.ReplyTopic()
	replyBytes, err := json.Marshal(replyMap)
	if err != nil {
		logger.LogError("Failed to marshall json "+err.Error(), reqCtx.Xid())
	}
	replyPayload = string(replyBytes)
	return
}

// RouteForAction returns a Route object for the named action or nil
// if there is no Route for that name
func (r *MQTTRouter) RouteForAction(action string) *common.Route {
	if val, present := r.actionRouteMap[strings.ToLower(action)]; present {
		return val
	}
	return nil
}

// NewMQTTRouter returns a router that can dispatch inbound MQTT requests
func NewMQTTRouter(appCtx *common.AppContext, apis common.APIRoutes) (router *MQTTRouter) {
	router = &MQTTRouter{appCtx: appCtx, actionRouteMap: map[string]*common.Route{}}
	for _, api := range apis {
		for i := 0; i < len(api); i++ {
			route := &api[i]
			name := strings.ToLower(route.Name)
			// Track all actions and prevent duplicates
			if router.RouteForAction(name) != nil {
				logger.ErrorLog(logger.LogModel{
					Level: logger.ERROR,
					Msg:   fmt.Sprintf("Duplicate route action %s", route.Name),
				})
			}
			router.actionRouteMap[name] = route
		}
	}
	return
}

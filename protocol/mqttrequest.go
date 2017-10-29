package protocol

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	auth "github.com/gkewl/pulsecheck/authentication"
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/constant"
	"github.com/gkewl/pulsecheck/utilities"
	"github.com/gkewl/pulsecheck/xid"
)

const (
	metaKey = "_meta_"
)

// MQTTRequestContext implements the RequestContext interface using information from an
// inbound MQTT request
type MQTTRequestContext struct {
	common.RequestContextBase
	values   map[string]interface{}
	metaData map[string]interface{}
}

// NewMQTTRequestContext returns an initialized RequestContext ready for request processing
// or an error if the request contexts could not be parsed or there was a database error
func NewMQTTRequestContext(appCtx *common.AppContext, topic, payload string) (mrc MQTTRequestContext, err error) {
	mrc.AppCtx = appCtx
	if err = json.Unmarshal([]byte(payload), &mrc.values); err != nil {
		return
	}
	mrc.metaData = mrc.meta()
	mrc.Xnid = mrc.metaString("xid")
	if mrc.Xnid == "" {
		mrc.Xnid = xid.UniqueIdGenerator()
	}
	mrc.MaxDeadlockRetries = constant.MaxDeadlockRetries
	token := mrc.Token()
	if token != "" {
		mrc.Username, mrc.Userid, mrc.Userrole, mrc.Companyid = auth.GetUserInfoFromToken(token)
	} else {
		mrc.Userrole = auth.RoleName(constant.Guest)
	}
	mrc.Txn, err = appCtx.Db.Beginx()
	return
}

// meta returns the map of metadata from the payload
func (mrc MQTTRequestContext) meta() map[string]interface{} {
	if val, ok := mrc.values[metaKey]; ok {
		if valmap, ok := val.(map[string]interface{}); ok {
			return valmap
		}
	}
	return map[string]interface{}{}
}

// metastring returns a string value from the metadata
func (mrc MQTTRequestContext) metaString(name string) string {
	if val, ok := mrc.metaData[name]; ok {
		return fmt.Sprintf("%v", val)
	}
	return ""
}

// Metadata returns the metadata for this request
func (mrc MQTTRequestContext) Metadata() (md map[string]interface{}) {
	md = map[string]interface{}{}
	md["action"] = mrc.Action()
	md["token"] = mrc.Token()
	md["xid"] = mrc.Xid()
	return
}

// Action returns the action for this request
func (mrc MQTTRequestContext) Action() string {
	return mrc.metaString("action")
}

// Token returns the token presented with this request
func (mrc MQTTRequestContext) Token() string {
	return mrc.metaString("token")
}

// ReplyTopic returns the name of the topic to reply to
func (mrc MQTTRequestContext) ReplyTopic() string {
	return mrc.metaString("reply_topic")
}

// Value returns the named input variable as a string
func (mrc MQTTRequestContext) Value(name string, defValue string) string {
	var val interface{}
	var present bool
	if val, present = mrc.values[name]; present == false {
		return defValue
	}
	return fmt.Sprintf("%v", val)
}

// BoolValue returns the named input variable as an bool
func (mrc MQTTRequestContext) BoolValue(name string, defValue bool) bool {
	val := mrc.Value(name, "default")
	if val == "default" {
		return defValue
	}
	if iVal, err := strconv.ParseBool(val); err == nil {
		return iVal
	}
	return defValue
}

// IntValue returns the named input variable as an integer
func (mrc MQTTRequestContext) IntValue(name string, defValue int64) int64 {
	val := mrc.Value(name, "default")
	if val == "default" {
		return defValue
	}
	if iVal, err := strconv.ParseInt(val, 10, 64); err == nil {
		return iVal
	}
	return defValue
}

// IntValue32 returns the named input variable as an int64
func (mrc *MQTTRequestContext) IntValue32(name string, defValue int) int {
	val := mrc.Value(name, "default")
	if val == "default" {
		return defValue
	}
	if iVal, err := strconv.Atoi(val); err == nil {
		return iVal
	}
	return defValue
}

// FloatValue returns the named input variable as a float64
func (mrc MQTTRequestContext) FloatValue(name string, defValue float64) float64 {
	val := mrc.Value(name, "default")
	if val == "default" {
		return defValue
	}
	if fVal, err := strconv.ParseFloat(val, 64); err == nil {
		return fVal
	}
	return defValue
}

// Time Value returns the named input variable as a time.Time
func (mrc MQTTRequestContext) TimeValue(name string, defValue time.Time) time.Time {
	val := mrc.Value(name, "default")
	if val == "default" {
		return defValue
	}
	if parsedVal, err := utilities.ParseStringtoTime(val); err == nil {
		return parsedVal
	}
	return defValue
}

// Scan will unmarshal content from the payload into the destination
// Callers should pass the address of their destination struct
// If there is no data with key name in payload, no error is returned
func (mrc MQTTRequestContext) Scan(name string, dest interface{}) (err error) {
	// todo: see if there's a better way than re-marshalling and unmarshalling
	if val, present := mrc.values[name]; present {
		var valJSON []byte
		if valJSON, err = json.Marshal(val); err != nil {
			return err
		}
		return json.Unmarshal(valJSON, dest)
	}
	return nil
}

//RequestBody  todo: we have to implement for mttt
func (mrc MQTTRequestContext) RequestBody() []byte {
	return []byte{}
}

//RequestUploadFiles  todo: we have to implement for mttt
func (mrc MQTTRequestContext) RequestUploadFiles() []common.Upload {
	return []common.Upload{}
}

// ResetForRetry clears out deferred functions, log fields and rolls back
// the transaction
func (req *MQTTRequestContext) ResetForRetry() (err error) {
	req.Tx().Rollback()
	req.Txn, err = req.AppCtx.Db.Beginx()
	//	req.RequestContextBase.ClearDeferredRequests()
	//	req.LogValues = map[string]interface{}{}
	return
}

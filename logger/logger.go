package logger

import (
	"flag"
	"fmt"
	//	"github.com/juju/errors"
	"io/ioutil"
	"os"
	"runtime"
	//"time"

	log "github.com/Sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	if flag.Lookup("test.v") != nil { // go test sets this flag for us
		NoOutput()
	}
}

// NoOutput discards all logging
func NoOutput() {
	log.SetOutput(ioutil.Discard)
}

func OutputStdout() {
	log.SetOutput(os.Stdout)
}

func OutputStderr() {
	log.SetOutput(os.Stderr)
}

const FATAL = "FATAL"
const ERROR = "ERROR"
const WARNING = "WARNING"
const INFO = "INFO"
const DEBUG = "DEBUG"

type LogModel struct {
	Level    string
	Caller   string
	Msg      string
	Err      error
	Xid      string
	Duration float64
	Input    string
	Fields   map[string]interface{}
}

func addModelFields(fields *log.Fields, model *LogModel) {
	if model.Fields != nil {
		for k, v := range model.Fields {
			(*fields)[k] = v
		}
	}
}

func ErrorLog(lm LogModel) {
	errMsg := ""
	if lm.Err != nil {
		errMsg = lm.Err.Error()
	}
	fields := log.Fields{
		"caller":   lm.Caller,
		"message":  lm.Msg,
		"error":    errMsg,
		"input":    lm.Input,
		"xid":      lm.Xid,
		"duration": lm.Duration,
	}
	addModelFields(&fields, &lm)
	log.WithFields(fields).Error("")
}

func Log(lm LogModel) {
	errMsg := ""
	if lm.Err != nil {
		errMsg = lm.Err.Error()
	}
	fields := log.Fields{
		"caller":   lm.Caller,
		"message":  lm.Msg,
		"error":    errMsg,
		"xid":      lm.Xid,
		"duration": lm.Duration,
	}
	addModelFields(&fields, &lm)
	log.WithFields(fields).Info("")
}

func LogInfo(msg string, xid string) {

	//Get caller info
	_, fn, line, _ := runtime.Caller(1)

	caller := fmt.Sprintf("[%s:%d]", fn, line)
	log.WithFields(log.Fields{
		"Caller": caller,
		"Xid":    xid,
	}).Info(msg)
}

func LogDebug(msg string, xid string) {

	//Get caller info
	pc, fn, line, _ := runtime.Caller(1)
	caller := fmt.Sprintf("%s [%s:%d]", runtime.FuncForPC(pc).Name(), fn, line)

	log.WithFields(log.Fields{
		"Caller": caller,
		"Xid":    xid,
	}).Debug(msg)
}

func LogError(msg string, xid string) {

	//Get caller info
	pc, fn, line, _ := runtime.Caller(1)
	caller := fmt.Sprintf("%s [%s:%d]", runtime.FuncForPC(pc).Name(), fn, line)

	log.WithFields(log.Fields{
		"Caller": caller,
		"Xid":    xid,
	}).Error(msg)
}

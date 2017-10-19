package protocol

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gkewl/pulsecheck/logger"
)

func logRecover(why interface{}, msg, xid string, start time.Time, fields map[string]interface{}) {
	var caller string
	for levels := 2; caller == ""; levels++ {
		pc, fn, line, ok := runtime.Caller(levels)
		if !ok {
			break
		}
		funcName := runtime.FuncForPC(pc).Name()
		if !strings.HasPrefix(funcName, "runtime.") {
			caller = fmt.Sprintf("%s [%s:%d]", funcName, fn, line)
		}
	}
	var logErr error
	var ok bool
	if logErr, ok = why.(error); !ok {
		logErr = fmt.Errorf("panic: %v", why)
	}
	fields["stacktrace"] = strings.Split(string(debug.Stack()), "\n")
	logMsg := logger.LogModel{
		Caller:   caller,
		Level:    logger.ERROR,
		Msg:      msg,
		Err:      logErr,
		Xid:      xid,
		Duration: time.Now().Sub(start).Seconds(),
		Fields:   fields,
	}
	logger.ErrorLog(logMsg)
}

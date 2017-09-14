package protocol

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/gkewl/pulsecheck/common"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/logger"
)

// RunDeferredRequests executes deferred requests stored in the request
// context based on their conditions
func RunDeferredRequests(reqCtx common.RequestContext, succeeded bool) {
	requests := reqCtx.DeferredRequests()
	var wg sync.WaitGroup
	for i := 0; i < len(requests); i++ {
		r := requests[i]
		if r.Condition == common.DeferredAlways || (r.Condition == common.DeferredSuccess && succeeded) ||
			(r.Condition == common.DeferredFailure && !succeeded) {
			// Run synchronous requests in parallel and wait
			if r.Synchronous {
				wg.Add(1)
				go func() {
					defer wg.Done()
					runDeferredRequest(reqCtx.AppContext(), r)
				}()
			} else {
				// Fire and forget async requests
				go runDeferredRequest(reqCtx.AppContext(), r)
			}
		}
	}
	wg.Wait()
}

// ProcessDeadlock checks for and logs deadlocks. returns whether the request
// should be retried
func ProcessDeadlock(rc common.RequestContext, err error) (retry bool, newerr error) {
	if err != nil {
		newerr = err
		foundDeadlock := false
		if ehErr, ok := err.(eh.Error); ok {
			foundDeadlock = strings.Contains(ehErr.DetailStack(";"), eh.ErrDeadlockText)
			if foundDeadlock {
				newerr = eh.AddDetail(eh.NewErrorFromError(eh.ErrDBContention, ehErr), "Retry #%d retrying %v",
					rc.DeadlockRetryCount(), retry)
			}
		} else {
			foundDeadlock = strings.Contains(err.Error(), eh.ErrDeadlockText)
			if foundDeadlock {
				newerr = eh.NewError(eh.ErrDBContention, "Retry #%d retrying %v base error %s",
					rc.DeadlockRetryCount(), retry, err.Error())
			}
		}
		if foundDeadlock {
			retry = rc.BumpDeadlocks()
			lm := logger.LogModel{
				Msg:   eh.ErrDBContention.String(),
				Xid:   rc.Xid(),
				Input: string(rc.RequestBody()),
				Fields: map[string]interface{}{
					"deadlock_count": rc.DeadlockRetryCount(),
					"retrying":       retry,
				},
			}
			logger.ErrorLog(lm)
			if retry {
				time.Sleep(time.Duration(5*rand.Intn(5)) * time.Millisecond)
				rerr := rc.ResetForRetry()
				if rerr != nil {
					newerr = rerr
				}
			}
			rc.AddLogValue("deadlocks", rc.DeadlockRetryCount())
		}
	}
	return
}

// StructuredError returns a structured response for the given error
func StructuredError(err error, xid string) (statusCode int, content common.StructuredResponse) {
	statusCode = 500
	eherr, isnewerr := err.(eh.Error)
	if isnewerr {
		statusCode = eherr.NamedError.HTTPStatus()
		content = common.StructuredResponse{
			Xid:        xid,
			StatusCode: eherr.NamedError.HTTPStatus(),
			Error: &common.ErrorResponse{
				Code:      eherr.NamedError.Code(),
				Message:   eherr.NamedError.Description(),
				Details:   eherr.DetailStack("; "),
				Locations: eherr.LocationStack("; "),
			},
		}
	} else {
		content = common.StructuredResponse{
			Xid:        xid,
			StatusCode: statusCode,
			Error: &common.ErrorResponse{
				Message: err.Error(),
			},
		}
	}
	return
}

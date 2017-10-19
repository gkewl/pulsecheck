package utilities

import (
	"os"

	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/model"
)

// get the hostname from the os once when the binary starts
var hostname, err = os.Hostname()

// GetPingInfo takes in modelName,
// combines env info an create Ping info structure
func GetPingInfo(ctx *common.AppContext, moduleName string) model.Ping {

	var output = model.Ping{
		Status:    "OK",
		Version:   ctx.Version,
		BuildTime: ctx.BuildTime,
		Githash:   ctx.GitHash,
	}

	return output
}

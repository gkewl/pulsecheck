package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/handlers"

	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/config"
	//	"github.com/gkewl/pulsecheck/constant"
	"github.com/gkewl/pulsecheck/dbhandler"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/routehandler"
	"github.com/gkewl/pulsecheck/rroutes"
)

// go:generate swagger generate spec
// accept values from -ldflags
var (
	Version   = "undefined"
	BuildTime = "undefined"
	GitHash   = "undefined"
)
var log = logrus.New()

var ctx common.AppContext

func main() {
	fmt.Printf("Version    : %s\n", Version)
	fmt.Printf("Git Hash   : %s\n", GitHash)
	fmt.Printf("Build Time : %s\n", BuildTime)

	boolPtr := flag.Bool("version", false, "print version and exit")
	flag.Parse()
	//	fmt.Println("versionflag:", *boolPtr)

	if *boolPtr == true {
		return
	}
	config.LoadConfigurations()

	ctx := common.AppContext{}

	// pass to handlers in context
	ctx.Version = Version
	ctx.BuildTime = BuildTime
	ctx.GitHash = GitHash
	var err error
	ctx.Db, err = dbhandler.CreateConnection()
	if err != nil {
		// ------ JUST GO PANIC ------
		panic("Failed to connected to database: %s mysql ")
	} else {
		fmt.Println("Connected to database: mysql ")

	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)


}

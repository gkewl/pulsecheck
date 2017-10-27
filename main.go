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
		panic("Failed to connect to database: %s mysql ")
	} else {
		fmt.Println("Connected to database: mysql ")
	}

	log.Formatter = new(logrus.JSONFormatter)

	if ctx.Db != nil {
		defer ctx.Db.Close()
	}

	router := routehandler.NewRouter(&ctx, rroutes.APIs, "/api/v1")
	router.NotFoundHandler = http.HandlerFunc(eh.NotFound)

	routehandler.AttachProfiler(router)

	//Initialize concrete instances
	//Initialize()

	env := "DEV"
	timeout, _ := time.ParseDuration("60s")
	var httpServer *http.Server
	if env == "DEV" {
		httpServer = &http.Server{
			Addr: ":8080",
			Handler: handlers.CORS(
				handlers.AllowedOrigins([]string{"*"}),
				handlers.AllowedHeaders([]string{"Authorization", "Origin", "Content-Type", "X-Auth-Token"}),
				handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}))(router),
			ReadTimeout:    timeout,
			WriteTimeout:   timeout,
			MaxHeaderBytes: 1 << 20,
		}
	} else {
		httpServer = &http.Server{
			Addr: ":8080",
			Handler: handlers.CORS(
				handlers.AllowedOrigins([]string{"*"}),
				handlers.AllowedHeaders([]string{"Authorization", "Access-Control-Allow-Headers"}),
				handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}))(router),
			ReadTimeout:    timeout,
			WriteTimeout:   timeout,
			MaxHeaderBytes: 1 << 20,
		}
	}

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		sig := <-quitChannel
		log.Info("Received quit signal: ", sig)
		shutdownCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		httpServer.Shutdown(shutdownCtx)
		log.Info("Exiting application due to: ", sig)
		os.Exit(0)
	}()
	httpServer.ListenAndServe()
	<-make(chan bool) // wait for os.Exit() above
}

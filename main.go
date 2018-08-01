package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/vsukhin/booking/logging"
	"github.com/vsukhin/booking/persistence/sqldb"
	"github.com/vsukhin/booking/router"
)

const (
	// HostAddress is host address
	HostAddress = ""
	// ParameterNameHostAddress contains parameter host name
	ParameterNameHostAddress = "host"
	// PortHTTP is server port
	PortHTTP = 3000
	// ParameterNamePortHTTP contains parameter port name
	ParameterNamePortHTTP = "port"
	// Mode is service running mode
	Mode = logging.ModeDev
	// ParameterNameMode contains parameter mode name
	ParameterNameMode = "mode"
	// DBConnection is db connection string
	DBConnection = ""
	// ParameterNameDBConnection contains parameter db connection name
	ParameterNameDBConnection = "db"
)

var (
	host         = flag.String(ParameterNameHostAddress, HostAddress, "HTTP server host address")
	httpPort     = flag.Int(ParameterNamePortHTTP, PortHTTP, "HTTP server port")
	mode         = flag.String(ParameterNameMode, Mode, "Service running mode: dev, staging, prod")
	dbConnection = flag.String(ParameterNameDBConnection, DBConnection, "DB connection string")
)

func initParameters() []error {
	var err error
	var errs []error

	flag.Parse()

	envHost := os.Getenv("BOOKING_API_HTTP_HOST")
	if envHost != "" {
		*host = envHost
	}

	envPort := os.Getenv("BOOKING_API_HTTP_PORT")
	if envPort != "" {
		var value int

		value, err = strconv.Atoi(envPort)
		if err == nil {
			*httpPort = value
		} else {
			errs = append(errs, err)
		}
	}

	envMode := os.Getenv("BOOKING_API_MODE")
	if envMode != "" {
		*mode = envMode
	}

	envDBConnection := os.Getenv("BOOKING_API_DB")
	if envDBConnection != "" {
		*dbConnection = envDBConnection
	}

	return errs
}

func main() {
	errs := initParameters()

	logging.Log = logging.NewLogger()
	logging.Log.Info("Service is started at ", time.Now())
	if len(errs) != 0 {
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"errors": errs,
		}).Warn("Error reading values from env")
	}

	logging.Log.Init(*mode)

	isProduction := false
	if *mode == logging.ModeStaging || *mode == logging.ModeProd {
		isProduction = true
	}

	db, err := sqldb.NewDB(*dbConnection, nil, isProduction, sqldb.NewGorpLogger())
	if err != nil {
		os.Exit(1)
	}

	routerManager := router.NewManager(db)
	r := routerManager.CreateRouter(*mode)
	server := &http.Server{Addr: *host + ":" + strconv.Itoa(*httpPort), Handler: r}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logging.Log.WithFields(logging.DepthModerate, logging.Fields{
				"error": err,
				"host":  *host,
				"port":  *httpPort,
			}).Fatal("Error starting http service")
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	<-c
	close(c)

	err = server.Shutdown(context.Background())
	if err != nil {
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error": err,
			"host":  *host,
			"port":  *httpPort,
		}).Error("Error stopping http service")
	}

	logging.Log.Info("Service is stopped at ", time.Now())
}

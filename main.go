package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitoc "github.com/go-kit/kit/tracing/opencensus"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/7Maliko7/backend-trainee-assignment-2023/internal/config"
	"github.com/7Maliko7/backend-trainee-assignment-2023/internal/middleware"
	"github.com/7Maliko7/backend-trainee-assignment-2023/internal/service"
	segmentssvc "github.com/7Maliko7/backend-trainee-assignment-2023/internal/service/segments"
	"github.com/7Maliko7/backend-trainee-assignment-2023/internal/transport"
	httptransport "github.com/7Maliko7/backend-trainee-assignment-2023/internal/transport/http"
	"github.com/7Maliko7/backend-trainee-assignment-2023/pkg/db/driver/postgres"
	"github.com/7Maliko7/backend-trainee-assignment-2023/pkg/oc"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "c", "", "Custom config path")
	flag.Parse()
}

func main() {
	var logger log.Logger
	{
		logger = log.NewJSONLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"svc", "backend-trainee-assignment-2023",
			"ts", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	appConfig, err := config.New(configPath)
	if err != nil {
		level.Error(logger).Log(err.Error())
		os.Exit(1)
	}

	var db *sql.DB
	{
		db, err = sql.Open("postgres", appConfig.RWDB.ConnectionString)
		if err != nil {
			level.Error(logger).Log("msg", err.Error())
			os.Exit(1)
		}
	}
	defer db.Close()

	var svc service.SegmentsService
	{
		repository, err := postgres.New(db)
		if err != nil {
			level.Error(logger).Log("exit", err.Error())
			os.Exit(1)
		}

		svc = segmentssvc.NewService(repository, logger)
		svc = middleware.LoggingMiddleware(logger)(svc)
	}

	var endpoints transport.Endpoints
	{
		endpoints = transport.MakeEndpoints(svc)
		endpoints = transport.Endpoints{
			AddSegment:        oc.ServerEndpoint("Add")(endpoints.AddSegment),
			DeleteSegment:     oc.ServerEndpoint("Delete")(endpoints.DeleteSegment),
			UpdateUserSegment: oc.ServerEndpoint("Patch")(endpoints.UpdateUserSegment),
			GetSegments:       oc.ServerEndpoint("Get")(endpoints.GetSegments),
		}
	}

	var h http.Handler
	{
		ocTracing := kitoc.HTTPServerTrace()
		serverOptions := []kithttp.ServerOption{ocTracing}
		h = httptransport.NewService(endpoints, serverOptions, logger, appConfig)
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		level.Info(logger).Log("transport", "HTTP", "addr", appConfig.ListenAddress)
		server := &http.Server{
			Addr:    appConfig.ListenAddress,
			Handler: h,
		}
		errs <- server.ListenAndServe()
	}()

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	level.Error(logger).Log("exit", <-errs)
}

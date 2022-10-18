package main

import (
	"fmt"
	"loki-cf-logdrain/handlers"
	"os"
	"os/signal"

	zipkinReporter "github.com/openzipkin/zipkin-go/reporter"

	"github.com/spf13/viper"

	"github.com/labstack/echo/v4"

	"net/http"
	_ "net/http/pprof"

	"github.com/labstack/echo-contrib/zipkintracing"
	"github.com/openzipkin/zipkin-go"
	zipkinHttpReporter "github.com/openzipkin/zipkin-go/reporter/http"
)

var commit = "deadbeaf"
var release = "v1.2.2"
var buildVersion = release + "-" + commit

func main() {
	e := make(chan *echo.Echo, 1)
	os.Exit(realMain(e))
}

func realMain(echoChan chan<- *echo.Echo) int {

	viper.SetEnvPrefix("loki-cf-logdrain")
	viper.SetDefault("transport_url", "")
	viper.SetDefault("promtail_endpoint", "localhost:1514")
	viper.AutomaticEnv()

	transportURL := viper.GetString("transport_url")
	token := os.Getenv("TOKEN")

	// Echo framework
	e := echo.New()

	// Tracing
	endpoint, err := zipkin.NewEndpoint("loki-cf-logdrain", "")
	if err != nil {
		e.Logger.Fatalf("error creating zipkin endpoint: %s", err.Error())
	}
	reporter := zipkinReporter.NewNoopReporter()
	if transportURL != "" {
		reporter = zipkinHttpReporter.NewReporter(transportURL)
	}
	traceTags := make(map[string]string)
	traceTags["app"] = "loki-cf-logdrain"
	tracer, err := zipkin.NewTracer(reporter,
		zipkin.WithLocalEndpoint(endpoint),
		zipkin.WithTags(traceTags),
		zipkin.WithSampler(zipkin.AlwaysSample))

	// Middleware
	if err == nil {
		e.Use(zipkintracing.TraceServer(tracer))
	}
	healthHandler := handlers.HealthHandler{}
	e.GET("/health", healthHandler.Handler(tracer))
	e.GET("/api/version", handlers.VersionHandler(buildVersion))

	promtailEndpoint := viper.GetString("promtail_endpoint")
	syslogHandler, err := handlers.NewSyslogHandler(token, promtailEndpoint)
	if err != nil {
		fmt.Printf("syslogHandler: %v\n", err)
		return 8
	}
	e.POST("/syslog/drain/:token", syslogHandler.Handler(tracer))

	setupPprof()
	setupInterrupts()

	echoChan <- e
	exitCode := 0
	if err := e.Start(listenString()); err != nil {
		fmt.Printf("error: %v\n", err)
		exitCode = 6
	}
	return exitCode
}

func setupInterrupts() {
	// Setup a channel to receive a signal
	done := make(chan os.Signal, 1)

	// Notify this channel when a SIGINT is received
	signal.Notify(done, os.Interrupt)

	// Fire off a goroutine to loop until that channel receives a signal.
	// When a signal is received simply exit the program
	go func() {
		for range done {
			os.Exit(0)
		}
	}()
}

func setupPprof() {
	go func() {
		err := http.ListenAndServe("localhost:6060", nil)
		if err != nil {
		}
	}()
}

func listenString() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return (":" + port)
}

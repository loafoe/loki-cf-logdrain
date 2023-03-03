package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/loafoe/loki-cf-logdrain/handlers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/spf13/viper"

	"github.com/labstack/echo/v4"

	"net/http"
	_ "net/http/pprof"

	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var commit = "deadbeaf"
var release = "v1.2.2"
var buildVersion = release + "-" + commit

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func initProvider() (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceName("go-hello-world"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// If the OpenTelemetry Collector is running on a local cluster (minikube or
	// microk8s), it should be accessible through the NodePort service at the
	// `localhost:30080` endpoint. Otherwise, replace `localhost` with the
	// endpoint of your cluster. If you run the app inside k8s, then you can
	// probably connect directly to the service through dns.
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, os.Getenv("OTLP_ADDRESS"),
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider.Shutdown, nil
}

func main() {
	e := make(chan *echo.Echo, 1)
	os.Exit(realMain(e))
}

func realMain(echoChan chan<- *echo.Echo) int {
	ctx := context.Background()
	shutdown, err := initProvider()
	if err == nil {
		defer func() {
			if err := shutdown(ctx); err != nil {
				fmt.Printf("failed to shutdown TracerProvider: %v\n", err)
			}
		}()
	}

	viper.SetEnvPrefix("loki-cf-logdrain")
	viper.SetDefault("transport_url", "")
	viper.SetDefault("promtail_endpoint", "localhost:1514")
	viper.AutomaticEnv()

	token := os.Getenv("TOKEN")

	// Echo framework
	e := echo.New()

	// Tracing
	tracer := otel.Tracer("loki-cf-logdrain")

	e.Use(otelecho.Middleware("loki-cf-logdrain"))

	healthHandler := handlers.HealthHandler{}
	e.GET("/health", healthHandler.Handler(ctx, tracer))
	e.GET("/api/version", handlers.VersionHandler(buildVersion))

	promtailEndpoint := viper.GetString("promtail_endpoint")
	syslogHandler, err := handlers.NewSyslogHandler(token, promtailEndpoint)
	if err != nil {
		fmt.Printf("syslogHandler: %v\n", err)
		return 8
	}
	e.POST("/syslog/drain/:token", syslogHandler.Handler(ctx, tracer))

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
		_ = http.ListenAndServe("localhost:6060", nil)
	}()
}

func listenString() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return (":" + port)
}

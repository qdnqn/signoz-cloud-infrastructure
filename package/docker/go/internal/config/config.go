package config

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.9.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"time"
)

type Config struct {
	Environment  string `mapstructure:"ENVIRONMENT"`
	HelloMessage string `mapstructure:"HELLO_MESSAGE"`
	Collector    string `mapstructure:"COLLECTOR"`
	Logger       *zap.Logger
	Id           uuid.UUID
	Tracer       trace.Tracer
}

func LoadConfig() (config Config, err error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "default"
	}

	config.Environment = env
	config.Id = uuid.New()

	viper.SetConfigFile("configs/env-" + env)
	viper.SetConfigType("env")

	if err = viper.ReadInConfig(); err != nil {
		viper.SetDefault("ENVIRONMENT", "default")

		viper.SetConfigFile("configs/env-default")
		if err = viper.ReadInConfig(); err != nil {
			if err := viper.Unmarshal(&config); err != nil {
				panic("Failed to unmarshal default configuration: " + err.Error())
			}
		}
	} else {
		err = viper.Unmarshal(&config)

		if err != nil {
			panic("Failed to unmarshal default configuration: " + err.Error())
		}
	}

	setupTracing(config)
	config.Tracer = otel.Tracer("test-tracer")

	return
}

func setupTracing(config Config) (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String("BackendService"),
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
	conn, err := grpc.DialContext(ctx, config.Collector,
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

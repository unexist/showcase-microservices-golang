//
// @package Showcase-Microservices-Golang
//
// @file Todo main
// @copyright 2023-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc/encoding/gzip"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/unexist/showcase-microservices-golang/adapter"
	"github.com/unexist/showcase-microservices-golang/domain"
	"github.com/unexist/showcase-microservices-golang/infrastructure"

	"fmt"
	"log"
	"os"
)

func initTracer(ctx context.Context) *sdktrace.TracerProvider {
	var exporter sdktrace.SpanExporter
	var err error

	/* Create trace exporter */
	if "jaeger" == os.Getenv("TRACER") {
		exporter, err = otlptracegrpc.New(ctx,
			otlptracegrpc.WithEndpoint("localhost:14268"),
			otlptracegrpc.WithCompressor(gzip.Name),
		)
	} else {
		exporter, err = zipkin.New("http://localhost:9411/api/v2/spans",
			zipkin.WithLogger(log.New(os.Stderr, "todo-service",
				log.Ldate|log.Ltime|log.Llongfile)),
		)
	}

	if nil != err {
		log.Fatal(err)
	}

	/* Create processor */
	batcher := sdktrace.NewBatchSpanProcessor(exporter,
		sdktrace.WithMaxQueueSize(1000),
		sdktrace.WithMaxExportBatchSize(1000))

	/* Create resource */
	resource, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			attribute.String("service.name", "todo-service"),
			attribute.String("service.version", "1.0.0"),
		))
	if nil != err {
		log.Fatal(err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(resource),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(batcher),
	)
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{}))

	return provider
}

func main() {
	/* Init tracer */
	ctx := context.Background()

	provider := initTracer(ctx)
	defer func() {
		if err := provider.Shutdown(context.Background()); nil != err {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	/* Create business stuff */
	var todoRepository *infrastructure.TodoGormRepository

	todoRepository = infrastructure.NewTodoGormRepository()

	/* Create database connection */
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s host=localhost port=5432 sslmode=disable",
			os.Getenv("APP_DB_USERNAME"),
			os.Getenv("APP_DB_PASSWORD"),
			os.Getenv("APP_DB_NAME"))

	err := todoRepository.Open(connectionString)

	if nil != err {
		log.Fatal(err)
	}

	defer todoRepository.Close()

	todoService := domain.NewTodoService(todoRepository)
	todoResource := adapter.NewTodoResource(todoService)

	/* Finally start Gin */
	engine := gin.Default()

	engine.Use(otelgin.Middleware("todo-service"))

	todoResource.RegisterRoutes(engine)

	log.Fatal(http.ListenAndServe(":8080", engine))
}

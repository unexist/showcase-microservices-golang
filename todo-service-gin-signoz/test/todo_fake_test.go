//
// @package Showcase-Microservices-Golang
//
// @file Todo tests for fake repository
// @copyright 2023-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

////go:build fake

package test

import (
	"context"
	"log"

	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/unexist/showcase-microservices-golang/adapter"
	"github.com/unexist/showcase-microservices-golang/domain"
)

/* Test globals */
var engine *gin.Engine
var todoRepository *TodoFakeRepository

func initTracer(ctx context.Context) *sdktrace.TracerProvider {
	var exporter sdktrace.SpanExporter
	var err error

	/* Create in-memory trace exporter */
	exporter = tracetest.NewInMemoryExporter()

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

func TestMain(m *testing.M) {
	/* Init tracer */
	ctx := context.Background()

	provider := initTracer(ctx)
	defer func() {
		if err := provider.Shutdown(context.Background()); nil != err {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	/* Create business stuff */
	todoRepository = NewTodoFakeRepository()
	todoService := domain.NewTodoService(todoRepository)
	idService := domain.NewIdService()
	todoResource := adapter.NewTodoResource(todoService, idService)

	/* Finally start Gin */
	engine = gin.Default()

	engine.Use(otelgin.Middleware("todo-fake-service"))

	todoResource.RegisterRoutes(engine)

	retCode := m.Run()

	os.Exit(retCode)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()

	engine.ServeHTTP(recorder, req)

	return recorder
}

func checkResponseCode(t *testing.T, expected, actual int) {
	assert.Equal(t, expected, actual, "Expected different response code")
}

func TestEmptyTable(t *testing.T) {
	todoRepository.Clear()

	req, _ := http.NewRequest("GET", "/todo", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); "[]" != body {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentTodo(t *testing.T) {
	todoRepository.Clear()

	req, _ := http.NewRequest("GET", "/todo/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)

	assert.Equal(t, "Todo not found", m["error"],
		"Expected the 'error' key of the response to be set to 'Todo not found'")
}

func TestCreateTodo(t *testing.T) {
	todoRepository.Clear()

	var jsonStr = []byte(`{"title":"string", "description": "string"}`)

	req, _ := http.NewRequest("POST", "/todo", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	assert.Equal(t, 1.0, m["id"], "Expected todo ID to be '1'")
	assert.Equal(t, "string", m["title"], "Expected todo title to be 'string'")
	assert.Equal(t, "string", m["description"], "Expected todo description to be 'string'")
}

func TestGetTodo(t *testing.T) {
	todoRepository.Clear()
	addTodos(1)

	req, _ := http.NewRequest("GET", "/todo/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addTodos(count int) {
	if 1 > count {
		count = 1
	}

	todo := domain.Todo{}

	for i := 0; i < count; i++ {
		todo.ID = i
		todo.Title = "Todo " + strconv.Itoa(i)
		todo.Description = "string"

		todoRepository.CreateTodo(nil, &todo)
	}
}

func TestUpdateTodo(t *testing.T) {
	todoRepository.Clear()
	addTodos(1)

	req, _ := http.NewRequest("GET", "/todo/1", nil)
	response := executeRequest(req)

	var origTodo map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &origTodo)

	var jsonStr = []byte(`{"title":"new string", "description": "new string"}`)

	req, _ = http.NewRequest("PUT", "/todo/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	assert.Equal(t, origTodo["id"], m["id"], "Expected the id to remain the same")
	assert.NotEqual(t, origTodo["title"], m["title"], "Expected the title to change")
	assert.NotEqual(t, origTodo["description"], m["description"], "Expected the description to change")
}

func TestDeleteTodo(t *testing.T) {
	todoRepository.Clear()
	addTodos(1)

	req, _ := http.NewRequest("GET", "/todo/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/todo/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNoContent, response.Code)

	req, _ = http.NewRequest("GET", "/todo/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

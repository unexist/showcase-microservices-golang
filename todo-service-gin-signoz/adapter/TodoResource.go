//
// @package Showcase-Microservices-Golang
//
// @file Todo resource
// @copyright 2023-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package adapter

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/unexist/showcase-microservices-golang/docs"
	"github.com/unexist/showcase-microservices-golang/domain"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var (
	todoActionCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "todo_action_counter",
			Help: "Total number of times an action was called",
		},
		[]string{"action"},
	)
)

func init() {
	prometheus.MustRegister(todoActionCounter)
}

// @title OpenAPI for Todo showcase
// @version 1.0
// @description OpenAPI for Todo showcase

// @contact.name Christoph Kappel
// @contact.url https://unexist.dev
// @contact.email christoph@unexist.dev

// @license.name Apache 2.0
// @license.url https://www.apache.org/licenses/LICENSE-2.0

// @BasePath /todo

type TodoResource struct {
	service *domain.TodoService
}

func NewTodoResource(service *domain.TodoService) *TodoResource {
	return &TodoResource{
		service: service,
	}
}

// @Summary Get all todos
// @Description Get all todos
// @Accept json
// @Produce json
// @Tags Todo
// @Success 200 {array} string "List of todo"
// @Failure 500 {string} string "Server error"
// @Router /todo [get]
func (resource *TodoResource) getTodos(context *gin.Context) {
	tracer := otel.GetTracerProvider().Tracer("todo-resource")
	ctx, span := tracer.Start(context.Request.Context(), "get-todos",
		trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	todos, err := resource.service.GetTodos(ctx)

	todoActionCounter.WithLabelValues("getTodos").Inc()

	if nil != err {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		context.JSON(http.StatusOK, todos)
	}
}

// @Summary Create new todo
// @Description Create new todo
// @Accept json
// @Produce json
// @Tags Todo
// @Success 201 {string} string "New todo entry"
// @Failure 500 {string} string "Server error"
// @Router /todo [post]
func (resource *TodoResource) createTodo(context *gin.Context) {
	tracer := otel.GetTracerProvider().Tracer("todo-resource")
	ctx, span := tracer.Start(context.Request.Context(), "create-todo",
		trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	todoActionCounter.WithLabelValues("createTodo").Inc()

	var todo domain.Todo

	if nil == context.Bind(&todo) {
		if err := resource.service.CreateTodo(ctx, &todo); nil != err {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			return
		}
	} else {
		context.JSON(http.StatusBadRequest, "Invalid request payload")

		return
	}

	context.JSON(http.StatusCreated, todo)
}

// @Summary Get todo by id
// @Description Get todo by id
// @Produce json
// @Tags Todo
// @Param   id  path  int  true  "Todo ID"
// @Success 200 {string} string "Todo found"
// @Failure 404 {string} string "Todo not found"
// @Failure 500 {string} string "Server error"
// @Router /todo/{id} [get]
func (resource *TodoResource) getTodo(context *gin.Context) {
	tracer := otel.GetTracerProvider().Tracer("todo-resource")
	ctx, span := tracer.Start(context.Request.Context(), "get-todo",
		trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	todoActionCounter.WithLabelValues("getTodo").Inc()

	todoId, err := strconv.Atoi(context.Params.ByName("id"))

	if nil != err {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})

		return
	}

	todo, err := resource.service.GetTodo(ctx, todoId)

	if nil != err {
		if 0 == strings.Compare("Not found", err.Error()) {
			context.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	} else {
		context.JSON(http.StatusOK, todo)
	}
}

// @Summary Update todo by id
// @Description Update todo by id
// @Accept json
// @Produce json
// @Tags Todo
// @Param   id  path  int  true  "Todo ID"
// @Success 200 {string} string "List of todo"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Server error"
// @Router /todo/{id} [put]
func (resource *TodoResource) updateTodo(context *gin.Context) {
	tracer := otel.GetTracerProvider().Tracer("todo-resource")
	ctx, span := tracer.Start(context.Request.Context(), "update-todo",
		trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	todoActionCounter.WithLabelValues("updateTodo").Inc()

	todoId, err := strconv.Atoi(context.Params.ByName("id"))

	if nil != err {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})

		return
	}

	var todo domain.Todo

	if context.Bind(&todo) == nil {
		todo.ID = todoId

		if err := resource.service.UpdateTodo(ctx, &todo); nil != err {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			return
		}
	}

	context.JSON(http.StatusOK, todo)
}

// @Summary Delete todo by id
// @Description Delete todo by id
// @Produce json
// @Tags Todo
// @Param   id  path  int  true  "Todo ID"
// @Success 204 {string} string "Todo updated"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Server error"
// @Router /todo/{id} [delete]
func (resource *TodoResource) deleteTodo(context *gin.Context) {
	tracer := otel.GetTracerProvider().Tracer("todo-resource")
	ctx, span := tracer.Start(context.Request.Context(), "delete-todo",
		trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	todoActionCounter.WithLabelValues("deleteTodo").Inc()

	todoId, err := strconv.Atoi(context.Params.ByName("id"))

	if nil != err {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})

		return
	}

	if err := resource.service.DeleteTodo(ctx, todoId); nil != err {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	context.Status(http.StatusNoContent)
}

// Register REST routes on given engine
func (resource *TodoResource) RegisterRoutes(engine *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/"

	todo := engine.Group("/todo")
	{
		todo.GET("", resource.getTodos)
		todo.POST("", resource.createTodo)
		todo.GET("/:id", resource.getTodo)
		todo.PUT("/:id", resource.updateTodo)
		todo.DELETE("/:id", resource.deleteTodo)
	}

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

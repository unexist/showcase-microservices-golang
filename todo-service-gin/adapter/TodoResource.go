//
// @package Showcase-Microservices-Golang
//
// @file Todo main
// @copyright 2021-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package adapter

import (
	"net/http"
	"strconv"

	"github.com/unexist/showcase-microservices-golang/domain"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type TodoResource struct {
	service *domain.TodoService
}

func NewTodoResource(service *domain.TodoService) *TodoResource {
	return &TodoResource{
		service: service,
	}
}

func (resource *TodoResource) getTodos(context *gin.Context) {
	todos, err := resource.service.GetTodos()

	if nil != err {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		context.JSON(http.StatusOK, todos)
	}
}

func (resource *TodoResource) createTodo(context *gin.Context) {
	var todo domain.Todo

	if nil == context.Bind(&todo) {
		if err := resource.service.CreateTodo(&todo); nil != err {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			return
		}
	} else {
		context.JSON(http.StatusInternalServerError, "Invalid request payload")

		return
	}

	context.JSON(http.StatusOK, todo)
}

func (resource *TodoResource) getTodo(context *gin.Context) {
	todoId, err := strconv.Atoi(context.Params.ByName("id"))

	if nil != err {
		context.JSON(http.StatusBadRequest, "Invalid todo ID")

		return
	}

	todo, err := resource.service.GetTodo(todoId)

	if nil != err {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else if (domain.Todo{} == *todo) {
		context.JSON(http.StatusNotFound, "Todo not found")
	} else {
		context.JSON(http.StatusOK, todo)
	}
}

func (resource *TodoResource) updateTodo(context *gin.Context) {
	todoId, err := strconv.Atoi(context.Params.ByName("id"))

	if nil != err {
		context.JSON(http.StatusBadRequest, "Invalid todo ID")

		return
	}

	var todo domain.Todo

	if context.Bind(&todo) == nil {
		todo.ID = todoId

		if err := resource.service.UpdateTodo(&todo); nil != err {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			return
		}
	}

	context.JSON(http.StatusOK, todo)
}

func (resource *TodoResource) deleteTodo(context *gin.Context) {
	todoId, err := strconv.Atoi(context.Params.ByName("id"))

	if nil != err {
		context.JSON(http.StatusBadRequest, "Invalid todo ID")

		return
	}

	if _, err := resource.service.GetTodo(todoId); nil != err {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	context.Status(http.StatusNoContent)
}

func (resource *TodoResource) RegisterRoutes(engine *gin.Engine) {
	engine.GET("/todo", resource.getTodos)
	engine.POST("/todo", resource.createTodo)
	engine.GET("/todo/:id", resource.getTodo)
	engine.PUT("/todo/:id", resource.updateTodo)
	engine.DELETE("/todo/:id", resource.deleteTodo)
}

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
	"github.com/unexist/showcase-microservices-golang/domain"

	"database/sql"
	"fmt"
	"log"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type App struct {
	Engine *gin.Engine
	DB     *sql.DB
}

func (app *App) Initialize(user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	app.DB, err = sql.Open("postgres", connectionString)
	if nil != err {
		log.Fatal(err)
	}

	app.Engine = gin.Default()

	app.initializeRoutes()
}

func (app *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, app.Engine))
}

func (app *App) GetTodo(context *gin.Context) {
	id, err := strconv.Atoi(context.Params.ByName("id"))
	if nil != err {
		context.JSON(http.StatusBadRequest, "Invalid todo ID")

		return
	}

	todo := domain.Todo{ID: id}
	if err := todo.GetTodo(app.DB); nil != err {
		switch err {
		case sql.ErrNoRows:
			context.JSON(http.StatusNotFound, "Todo not found")
		default:
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	context.JSON(http.StatusOK, todo)
}

func (app *App) GetTodos(context *gin.Context) {
	count, _ := strconv.Atoi(context.PostForm("count"))
	start, _ := strconv.Atoi(context.PostForm("start"))

	if 10 < count || 1 > count {
		count = 10
	}
	if 0 > start {
		start = 0
	}

	todos, err := domain.GetTodos(app.DB, start, count)
	if nil != err {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	context.JSON(http.StatusOK, todos)
}

func (app *App) createTodo(context *gin.Context) {
	var todo domain.Todo

	if context.Bind(&todo) == nil {
		if err := todo.CreateTodo(app.DB); nil != err {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			return
		}
	} else {
		context.JSON(http.StatusInternalServerError, "Invalid request payload")

		return
	}

	context.JSON(http.StatusOK, todo)
}

func (app *App) UpdateTodo(context *gin.Context) {
	id, err := strconv.Atoi(context.Params.ByName("id"))
	if nil != err {
		context.JSON(http.StatusBadRequest, "Invalid todo ID")

		return
	}

	var todo domain.Todo

	if context.Bind(&todo) == nil {
		todo.ID = id

		if err := todo.UpdateTodo(app.DB); nil != err {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			return
		}
	}

	context.JSON(http.StatusOK, todo)
}

func (app *App) DeleteTodo(context *gin.Context) {
	id, err := strconv.Atoi(context.Params.ByName("id"))
	if nil != err {
		context.JSON(http.StatusBadRequest, "Invalid todo ID")

		return
	}

	todo := domain.Todo{ID: id}
	if err := todo.DeleteTodo(app.DB); nil != err {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	context.Status(http.StatusNoContent)
}

func (app *App) initializeRoutes() {
	app.Engine.GET("/todo", app.GetTodos)
	app.Engine.POST("/todo", app.createTodo)
	app.Engine.GET("/todo/:id", app.GetTodo)
	app.Engine.PUT("/todo/:id", app.UpdateTodo)
	app.Engine.DELETE("/todo/:id", app.DeleteTodo)
}

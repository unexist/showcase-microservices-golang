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

package main

import (
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

func (app *App) getTodo(context *gin.Context) {
	id, err := strconv.Atoi(context.Params.ByName("id"))
	if nil != err {
		context.JSON(http.StatusBadRequest, "Invalid todo ID")

		return
	}

	todo := Todo{ID: id}
	if err := todo.getTodo(app.DB); nil != err {
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

func (app *App) getTodos(context *gin.Context) {
	count, _ := strconv.Atoi(context.PostForm("count"))
	start, _ := strconv.Atoi(context.PostForm("start"))

	if 10 < count || 1 > count {
		count = 10
	}
	if 0 > start {
		start = 0
	}

	todos, err := getTodos(app.DB, start, count)
	if nil != err {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	context.JSON(http.StatusOK, todos)
}

func (app *App) createTodo(context *gin.Context) {
	var todo Todo

	if context.Bind(&todo) == nil {
		if err := todo.createTodo(app.DB); nil != err {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			return
		}
	} else {
		context.JSON(http.StatusInternalServerError, "Invalid request payload")

		return
	}

	context.JSON(http.StatusOK, todo)
}

func (app *App) updateTodo(context *gin.Context) {
	id, err := strconv.Atoi(context.Params.ByName("id"))
	if nil != err {
		context.JSON(http.StatusBadRequest, "Invalid todo ID")

		return
	}

	var todo Todo

	if context.Bind(&todo) == nil {
		todo.ID = id

		if err := todo.updateTodo(app.DB); nil != err {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

			return
		}
	}

	context.JSON(http.StatusOK, todo)
}

func (app *App) deleteTodo(context *gin.Context) {
	id, err := strconv.Atoi(context.Params.ByName("id"))
	if nil != err {
		context.JSON(http.StatusBadRequest, "Invalid todo ID")

		return
	}

	todo := Todo{ID: id}
	if err := todo.deleteTodo(app.DB); nil != err {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	context.Status(http.StatusNoContent)
}

func (app *App) initializeRoutes() {
	app.Engine.GET("/todo", app.getTodos)
	app.Engine.POST("/todo", app.createTodo)
	app.Engine.GET("/todo/:id", app.getTodo)
	app.Engine.PUT("/todo/:id", app.updateTodo)
	app.Engine.DELETE("/todo/:id", app.deleteTodo)
}

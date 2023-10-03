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
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/unexist/showcase-microservices-golang/adapter"
	"github.com/unexist/showcase-microservices-golang/domain"
	"github.com/unexist/showcase-microservices-golang/infrastructure"

	"database/sql"
	"fmt"
	"log"
	"os"
)

func main() {
	var database *sql.DB
	var engine *gin.Engine
	var err error

	/* Create database connection */
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("APP_DB_USERNAME"),
			os.Getenv("APP_DB_PASSWORD"),
			os.Getenv("APP_DB_NAME"))

	database, err = sql.Open("postgres", connectionString)
	if nil != err {
		log.Fatal(err)
	}

	/* Create business stuff */
	var todoRepository *infrastructure.TodoRepository
	var todoService *domain.TodoService
	var todoResource *adapter.TodoResource

	todoRepository = infrastructure.NewTodoRepository(database)
	todoService = domain.NewTodoService(todoRepository)
	todoResource = adapter.NewTodoResource(todoService)

	/* Finally start Gin */
	engine = gin.Default()

	todoResource.RegisterRoutes(engine)

	log.Fatal(http.ListenAndServe(":8080", engine))
}

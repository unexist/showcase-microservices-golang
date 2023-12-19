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
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/unexist/showcase-microservices-golang/adapter"
	"github.com/unexist/showcase-microservices-golang/domain"
	"github.com/unexist/showcase-microservices-golang/infrastructure"

	"fmt"
	"log"
	"os"
)

func main() {
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

	todoResource.RegisterRoutes(engine)

	log.Fatal(http.ListenAndServe(":8080", engine))
}

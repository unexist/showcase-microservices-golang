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
	"os"

	sqlxTransactor "github.com/Thiht/transactor/sqlx"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/unexist/showcase-microservices-golang/application"
	todoDomain "github.com/unexist/showcase-microservices-golang/domain/todo"
	userDomain "github.com/unexist/showcase-microservices-golang/domain/user"

	"github.com/gin-gonic/gin"

	"github.com/unexist/showcase-microservices-golang/adapter"
	"github.com/unexist/showcase-microservices-golang/infrastructure"

	"log"
)

func migrate(db *sqlx.DB, name string) {
	schema, err := os.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}

	db.MustExec(string(schema))
}

func main() {

	/* Create database connection */
	db, err := sqlx.Connect("sqlite3", "todo.db")

	if nil != err {
		log.Fatal(err)
	}

	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	/* Migrate */
	migrate(db, "../infrastructure/todos-with-user.sql")
	migrate(db, "../infrastructure/users.sql")

	/* Create transactor */
	transactor, dbGetter := sqlxTransactor.NewTransactor(
		db,
		sqlxTransactor.NestedTransactionsSavepoints,
	)

	/* Create business stuff */
	todoRepository := infrastructure.NewTodoSQLXRepository(dbGetter)
	userRepository := infrastructure.NewUserSQLXRepository(dbGetter)

	todoService := todoDomain.NewTodoService(todoRepository)
	userService := userDomain.NewUserService(userRepository)
	appService := application.NewTodoUserService(transactor, todoService, userService)

	userResource := adapter.NewUserResource(userService)
	todoResource := adapter.NewTodoResource(todoService, appService)

	/* Create middleware */
	authHandler := infrastructure.AuthUser(userService)

	/* Finally start Gin */
	engine := gin.Default()

	todoResource.RegisterRoutes(engine, authHandler)
	userResource.RegisterRoutes(engine, authHandler)

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	log.Fatal(http.ListenAndServe(":8080", engine))
}

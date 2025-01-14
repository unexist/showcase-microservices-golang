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

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/unexist/showcase-microservices-golang/application"
	todoDomain "github.com/unexist/showcase-microservices-golang/domain/todo"
	userDomain "github.com/unexist/showcase-microservices-golang/domain/user"

	"github.com/gin-gonic/gin"

	"github.com/unexist/showcase-microservices-golang/adapter"
	"github.com/unexist/showcase-microservices-golang/infrastructure"

	"log"
)

func main() {
	/* Create business stuff */
	var todoRepository *infrastructure.TodoGormRepository
	var userRepository *infrastructure.UserGormRepository

	/* Create database connection */
	db, err := gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})

	if nil != err {
		log.Fatal(err)
	}

	todoRepository = infrastructure.NewTodoGormRepository(db)
	userRepository = infrastructure.NewUserGormRepository(db)

	todoService := todoDomain.NewTodoService(todoRepository)
	userService := userDomain.NewUserService(userRepository)
	appService := application.NewTodoUserService(db, todoService, userService)

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

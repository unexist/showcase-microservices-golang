//
// @package Showcase-Microservices-Golang
//
// @file User resource
// @copyright 2023-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package adapter

import (
	"net/http"

	"github.com/unexist/showcase-microservices-golang/docs"
	"github.com/unexist/showcase-microservices-golang/domain/user"

	"github.com/gin-gonic/gin"
)

// @title OpenAPI for Todo showcase
// @version 1.0
// @description OpenAPI for Todo showcase

// @contact.name Christoph Kappel
// @contact.url https://unexist.dev
// @contact.email christoph@unexist.dev

// @license.name Apache 2.0
// @license.url https://www.apache.org/licenses/LICENSE-2.0

// @BasePath /user

type UserResource struct {
	service *domain.UserService
}

func NewUserResource(service *domain.UserService) *UserResource {
	return &UserResource{
		service: service,
	}
}

// @Summary Log user in
// @Description Log user in
// @Produce json
// @Tags User
// @Success 200 {string} string "User found"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Server error"
// @Router /user/login [post]
func (resource *UserResource) login(context *gin.Context) {
	defaultUser := domain.User{
		Name: "Default User",
	}

	if err := resource.service.CreateUser(context, &defaultUser); nil != err {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	context.JSON(http.StatusOK, gin.H{"token": defaultUser.Token})
}

// @Summary Get logged in user
// @Description Get logged in user
// @Produce json
// @Tags User
// @Success 200 {string} string "User found"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Server error"
// @Router /user/self [get]
func (resource *UserResource) getSelf(context *gin.Context) {
	user := context.MustGet("user").(*domain.User)

	context.JSON(http.StatusOK, user)
}

// Register REST routes on given engine
func (resource *UserResource) RegisterRoutes(engine *gin.Engine, authHandler gin.HandlerFunc) {
	docs.SwaggerInfo.BasePath = "/"

	user := engine.Group("/user")
	{
		user.GET("self", authHandler, resource.getSelf)
		user.POST("login", resource.login)
	}
}

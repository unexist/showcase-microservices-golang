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

//go:build fake

package test

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/unexist/showcase-microservices-golang/infrastructure"

	"os"
	"testing"

	"bytes"
	"github.com/unexist/showcase-microservices-golang/adapter"
	"github.com/unexist/showcase-microservices-golang/domain"
	"net/http"
	"net/http/httptest"
)

/* Test globals */
var (
	engine         *gin.Engine
	todoRepository *infrastructure.TodoFakeRepository
	userRepository *infrastructure.UserFakeRepository
)

func TestMain(m *testing.M) {
	/* Create business stuff */
	todoRepository = infrastructure.NewTodoFakeRepository()
	userRepository = infrastructure.NewUserFakeRepository()
	todoService := domain.NewTodoService(todoRepository)
	userService := domain.NewUserService(userRepository)
	todoResource := adapter.NewTodoResource(todoService)
	userResource := adapter.NewUserResource(userService)

	/* Finally start Gin */
	engine = gin.Default()

	todoResource.RegisterRoutes(engine)
	userResource.RegisterRoutes(engine)

	code := m.Run()

	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()

	engine.ServeHTTP(recorder, req)

	return recorder
}

func checkResponseCode(t *testing.T, expected, actual int) {
	assert.Equal(t, expected, actual, "Expected different response code")
}

func TestLogin(t *testing.T) {
	todoRepository.Clear()

	req, _ := http.NewRequest("POST", "/user/login", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

//
// @package Showcase-Microservices-Golang
//
// @file Todo tests for fake repository
// @copyright 2025-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package test

import (
	"encoding/json"

	sqlxTransactor "github.com/Thiht/transactor/sqlx"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/unexist/showcase-microservices-golang/application"
	todoDomain "github.com/unexist/showcase-microservices-golang/domain/todo"
	userDomain "github.com/unexist/showcase-microservices-golang/domain/user"
	"github.com/unexist/showcase-microservices-golang/infrastructure"

	"os"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/unexist/showcase-microservices-golang/adapter"
)

/* Test globals */
var (
	engine         *gin.Engine
	todoRepository *infrastructure.TodoSQLXRepository
	userRepository *infrastructure.UserSQLXRepository
)

func TestMain(m *testing.M) {
	transactor, dbGetter := sqlxTransactor.NewFakeTransactor(db)

	/* Create business stuff */
	todoRepository = infrastructure.NewTodoSQLXRepository()
	userRepository = infrastructure.NewUserSQLXRepository()

	todoService := todoDomain.NewTodoService(todoRepository)
	userService := userDomain.NewUserService(userRepository)
	appService := application.NewTodoUserService(todoService, userService)

	todoResource := adapter.NewTodoResource(todoService, appService)
	userResource := adapter.NewUserResource(userService)

	/* Create middleware */
	authHandler := infrastructure.AuthUser(userService)

	/* Finally start Gin */
	engine = gin.Default()

	todoResource.RegisterRoutes(engine, authHandler)
	userResource.RegisterRoutes(engine, authHandler)

	retCode := m.Run()

	os.Exit(retCode)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()

	engine.ServeHTTP(recorder, req)

	return recorder
}

func executeLogin(t *testing.T) string {
	req, _ := http.NewRequest("POST", "/user/login", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	return m["token"].(string)
}

func checkResponseCode(t *testing.T, expected, actual int) {
	assert.Equal(t, expected, actual, "Expected different response code")
}

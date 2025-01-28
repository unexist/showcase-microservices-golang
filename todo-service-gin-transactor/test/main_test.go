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
	"log"

	sqlxTransactor "github.com/Thiht/transactor/sqlx"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
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

func migrate(db *sqlx.DB, name string) {
	schema, err := os.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}

	db.MustExec(string(schema))
}

func TestMain(m *testing.M) {

	/* Create database connection */
	db, err := sqlx.Connect("sqlite3", ":memory:")

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
	migrate(db, "../../infrastructure/todos-with-user.sql")
	migrate(db, "../../infrastructure/users.sql")

	/* Create transactor */
	transactor, dbGetter := sqlxTransactor.NewFakeTransactor(db)

	/* Create business stuff */
	todoRepository = infrastructure.NewTodoSQLXRepository(dbGetter)
	userRepository = infrastructure.NewUserSQLXRepository(dbGetter)

	todoService := todoDomain.NewTodoService(todoRepository)
	userService := userDomain.NewUserService(userRepository)
	appService := application.NewTodoUserService(transactor, todoService, userService)

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

func executeLogin(t *testing.T) string {
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user/login", nil)
	engine.ServeHTTP(recorder, req)

	checkResponseCode(t, http.StatusOK, recorder.Code)

	var m map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &m)

	return m["token"].(string)
}

func checkResponseCode(t *testing.T, expected, actual int) {
	assert.Equal(t, expected, actual, "Expected different response code")
}

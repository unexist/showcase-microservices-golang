//
// @package Showcase-Microservices-Golang
//
// @file Todo godoc test main
// @copyright 2023-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

//go:build cucumber

package test

import (
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"testing"

	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"braces.dev/errtrace"

	"github.com/unexist/showcase-microservices-golang/adapter"
	"github.com/unexist/showcase-microservices-golang/domain"
)

/* Test globals */
var engine *gin.Engine
var todo *domain.Todo

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()

	engine.ServeHTTP(recorder, req)

	return recorder
}

type expectedAndActualAssertion func(t assert.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool

func assertExpectedAndActual(a expectedAndActualAssertion, expected, actual interface{}, msgAndArgs ...interface{}) error {
	var t asserter
	a(&t, expected, actual, msgAndArgs...)
	return errtrace.Wrap(t.err)
}

type asserter struct {
	err error
}

func (a *asserter) Errorf(format string, args ...interface{}) {
	a.err = fmt.Errorf(format, args...)
}

func givenCreateTodo() error {
	todo = &domain.Todo{}

	return nil
}

func whenSetTitle(title string) error {
	todo.Title = title

	return nil
}

func andSetDescription(description string) error {
	todo.Description = description

	return nil
}

func thenGetId(id float64) error {
	jsonStr, _ := json.Marshal(todo)

	req, _ := http.NewRequest("POST", "/todo", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	err := assertExpectedAndActual(
		assert.Equal, http.StatusCreated, response.Code, "Expected different response code")

	if nil != err {
		return errtrace.Wrap(err)
	}

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	err = assertExpectedAndActual(
		assert.Equal, id, m["id"], "Expected todo ID to be '%d'", id,
	)

	if nil != err {
		return errtrace.Wrap(err)
	}

	err = assertExpectedAndActual(
		assert.Equal, todo.Title, m["title"], "Expected todo title to be '%s'", todo.Title,
	)

	if nil != err {
		return errtrace.Wrap(err)
	}

	err = assertExpectedAndActual(
		assert.Equal, todo.Description, m["description"], "Expected todo description to be '%s'", todo.Description,
	)

	if nil != err {
		return errtrace.Wrap(err)
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		todo = nil
		return ctx, nil
	})

	ctx.Given(`^I create a todo$`, givenCreateTodo)
	ctx.When(`^its title is "([^"]*)"$`, whenSetTitle)
	ctx.When(`^its description is "([^"]*)"$`, andSetDescription)
	ctx.Then(`^its id should be (\d+)$`, thenGetId)
}

func TestMain(m *testing.M) {
	/* Create business stuff */
	var todoRepository *TodoFakeRepository

	todoRepository = NewTodoFakeRepository()
	todoService := domain.NewTodoService(todoRepository)
	todoResource := adapter.NewTodoResource(todoService)

	/* Finally start Gin */
	engine = gin.Default()

	todoResource.RegisterRoutes(engine)

	retCode := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format: "pretty",
			Paths:  []string{"features"},
			Output: colors.Colored(os.Stdout),
		},
	}.Run()

	os.Exit(retCode)
}

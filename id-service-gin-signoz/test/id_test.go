//
// @package Showcase-Microservices-Golang
//
// @file Id tests
// @copyright 2023-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package test

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"os"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/unexist/showcase-microservices-golang/adapter"
)

/* Test globals */
var engine *gin.Engine

func TestMain(m *testing.M) {
	idResource := adapter.NewIdResource()

	/* Finally start Gin */
	engine = gin.Default()

	idResource.RegisterRoutes(engine)

	retcode := m.Run()

	os.Exit(retcode)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()

	engine.ServeHTTP(recorder, req)

	return recorder
}

func checkResponseCode(t *testing.T, expected, actual int) {
	assert.Equal(t, expected, actual, "Expected different response code")
}

func TestGetTodo(t *testing.T) {
	req, _ := http.NewRequest("GET", "/id", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

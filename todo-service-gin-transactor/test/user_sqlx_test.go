//
// @package Showcase-Microservices-Golang
//
// @file User tests for sqlx repository
// @copyright 2023-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package test

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"net/http"

	"github.com/gin-gonic/gin"
)

func TestLogin(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	todoRepository.Clear(ctx)
	userRepository.Clear(ctx)

	req, _ := http.NewRequest("POST", "/user/login", nil)
	req.Header.Set("Content-Type", "application/json")

	engine.ServeHTTP(recorder, req)
	checkResponseCode(t, http.StatusOK, recorder.Code)
}

func TestSelf(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	todoRepository.Clear(ctx)
	userRepository.Clear(ctx)

	req, _ := http.NewRequest("GET", "/user/self", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", executeLogin(t)))

	engine.ServeHTTP(recorder, req)
	checkResponseCode(t, http.StatusOK, recorder.Code)
}

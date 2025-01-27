//
// @package Showcase-Microservices-Golang
//
// @file Todo tests for sqlx repository
// @copyright 2025-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	todoDomain "github.com/unexist/showcase-microservices-golang/domain/todo"

	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
)

func TestEmptyTable(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	todoRepository.Clear(ctx)
	userRepository.Clear(ctx)

	req, _ := http.NewRequest("GET", "/todo", nil)
	engine.ServeHTTP(recorder, req)

	checkResponseCode(t, http.StatusOK, recorder.Code)

	if body := recorder.Body.String(); "[]" != body {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentTodo(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	todoRepository.Clear(ctx)
	userRepository.Clear(ctx)

	req, _ := http.NewRequest("GET", "/todo/11", nil)
	engine.ServeHTTP(recorder, req)

	checkResponseCode(t, http.StatusNotFound, recorder.Code)

	var m map[string]string
	json.Unmarshal(recorder.Body.Bytes(), &m)

	assert.Equal(t, "Todo not found", m["error"],
		"Expected the 'error' key of the response to be set to 'Todo not found'")
}

func TestCreateTodo(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	todoRepository.Clear(ctx)
	userRepository.Clear(ctx)

	var jsonStr = []byte(`{"title":"string", "description": "string"}`)

	req, _ := http.NewRequest("POST", "/todo", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", executeLogin(t)))

	engine.ServeHTTP(recorder, req)
	checkResponseCode(t, http.StatusCreated, recorder.Code)

	var m map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &m)

	assert.Equal(t, 1.0, m["id"], "Expected todo ID to be '1'")
	assert.Equal(t, "string", m["title"], "Expected todo title to be 'string'")
	assert.Equal(t, "string", m["description"], "Expected todo description to be 'string'")
}

func TestCreateTodoAnon(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	todoRepository.Clear(ctx)
	userRepository.Clear(ctx)

	var jsonStr = []byte(`{"title":"string", "description": "string"}`)

	req, _ := http.NewRequest("POST", "/todo/anon", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	engine.ServeHTTP(recorder, req)
	checkResponseCode(t, http.StatusCreated, recorder.Code)

	var m map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &m)

	assert.Equal(t, 1.0, m["id"], "Expected todo ID to be '1'")
	assert.Equal(t, "string", m["title"], "Expected todo title to be 'string'")
	assert.Equal(t, "string", m["description"], "Expected todo description to be 'string'")
}

func TestGetTodo(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	todoRepository.Clear(ctx)
	userRepository.Clear(ctx)
	addTodos(ctx, 1)

	req, _ := http.NewRequest("GET", "/todo/1", nil)
	engine.ServeHTTP(recorder, req)

	checkResponseCode(t, http.StatusOK, recorder.Code)
}

func addTodos(ctx context.Context, count int) {
	if 1 > count {
		count = 1
	}

	todo := todoDomain.Todo{}

	for i := 0; i < count; i++ {
		todo.ID = i
		todo.Title = "Todo " + strconv.Itoa(i)
		todo.Description = "string"

		todoRepository.CreateTodo(ctx, &todo)
	}
}

func TestUpdateTodo(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	todoRepository.Clear(ctx)
	userRepository.Clear(ctx)
	addTodos(ctx, 1)

	req, _ := http.NewRequest("GET", "/todo/1", nil)
	engine.ServeHTTP(recorder, req)

	var origTodo map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &origTodo)

	var jsonStr = []byte(`{"title":"new string", "description": "new string"}`)

	req, _ = http.NewRequest("PUT", "/todo/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", executeLogin(t)))

	engine.ServeHTTP(recorder, req)
	checkResponseCode(t, http.StatusOK, recorder.Code)

	var m map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &m)

	assert.Equal(t, origTodo["id"], m["id"], "Expected the id to remain the same")
	assert.NotEqual(t, origTodo["title"], m["title"], "Expected the title to change")
	assert.NotEqual(t, origTodo["description"], m["description"], "Expected the description to change")
}

func TestDeleteTodo(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	todoRepository.Clear(ctx)
	userRepository.Clear(ctx)
	addTodos(ctx, 1)

	req, _ := http.NewRequest("GET", "/todo/1", nil)
	engine.ServeHTTP(recorder, req)
	checkResponseCode(t, http.StatusOK, recorder.Code)

	req, _ = http.NewRequest("DELETE", "/todo/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", executeLogin(t)))
	engine.ServeHTTP(recorder, req)
	checkResponseCode(t, http.StatusNoContent, recorder.Code)

	req, _ = http.NewRequest("GET", "/todo/1", nil)
	engine.ServeHTTP(recorder, req)
	checkResponseCode(t, http.StatusNotFound, recorder.Code)
}

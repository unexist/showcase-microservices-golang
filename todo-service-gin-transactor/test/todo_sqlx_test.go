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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	todoDomain "github.com/unexist/showcase-microservices-golang/domain/todo"

	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
)

func TestEmptyTable(t *testing.T) {
	todoRepository.Clear()
	userRepository.Clear()

	req, _ := http.NewRequest("GET", "/todo", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); "[]" != body {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentTodo(t *testing.T) {
	todoRepository.Clear()
	userRepository.Clear()

	req, _ := http.NewRequest("GET", "/todo/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)

	assert.Equal(t, "Todo not found", m["error"],
		"Expected the 'error' key of the response to be set to 'Todo not found'")
}

func TestCreateTodo(t *testing.T) {
	todoRepository.Clear()
	userRepository.Clear()

	var jsonStr = []byte(`{"title":"string", "description": "string"}`)

	req, _ := http.NewRequest("POST", "/todo", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", executeLogin(t)))

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	assert.Equal(t, 1.0, m["id"], "Expected todo ID to be '1'")
	assert.Equal(t, "string", m["title"], "Expected todo title to be 'string'")
	assert.Equal(t, "string", m["description"], "Expected todo description to be 'string'")
}

func TestCreateTodoAnon(t *testing.T) {
	todoRepository.Clear()
	userRepository.Clear()

	var jsonStr = []byte(`{"title":"string", "description": "string"}`)

	req, _ := http.NewRequest("POST", "/todo/anon", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	assert.Equal(t, 1.0, m["id"], "Expected todo ID to be '1'")
	assert.Equal(t, "string", m["title"], "Expected todo title to be 'string'")
	assert.Equal(t, "string", m["description"], "Expected todo description to be 'string'")
}

func TestGetTodo(t *testing.T) {
	todoRepository.Clear()
	userRepository.Clear()
	addTodos(1)

	req, _ := http.NewRequest("GET", "/todo/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addTodos(count int) {
	if 1 > count {
		count = 1
	}

	todo := todoDomain.Todo{}

	for i := 0; i < count; i++ {
		todo.ID = i
		todo.Title = "Todo " + strconv.Itoa(i)
		todo.Description = "string"

		todoRepository.CreateTodo(&todo)
	}
}

func TestUpdateTodo(t *testing.T) {
	todoRepository.Clear()
	userRepository.Clear()
	addTodos(1)

	req, _ := http.NewRequest("GET", "/todo/1", nil)
	response := executeRequest(req)

	var origTodo map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &origTodo)

	var jsonStr = []byte(`{"title":"new string", "description": "new string"}`)

	req, _ = http.NewRequest("PUT", "/todo/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", executeLogin(t)))

	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	assert.Equal(t, origTodo["id"], m["id"], "Expected the id to remain the same")
	assert.NotEqual(t, origTodo["title"], m["title"], "Expected the title to change")
	assert.NotEqual(t, origTodo["description"], m["description"], "Expected the description to change")
}

func TestDeleteTodo(t *testing.T) {
	todoRepository.Clear()
	userRepository.Clear()
	addTodos(1)

	req, _ := http.NewRequest("GET", "/todo/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/todo/1", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", executeLogin(t)))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNoContent, response.Code)

	req, _ = http.NewRequest("GET", "/todo/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

//
// @package Showcase-Microservices-Golang
//
// @file Todo test main
// @copyright 2021-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package main_test

import (
	"log"
	"os"
	"testing"

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
)

var app main.App

func TestMain(m *testing.M) {
	app = main.App{}

	app.Initialize(
		os.Getenv("TEST_DB_USERNAME"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"))

	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func ensureTableExists() {
	if _, err := app.DB.Exec(tableCreationQuery); nil != err {
		log.Fatal(err)
	}
}

func clearTable() {
	app.DB.Exec("DELETE FROM todos")
	app.DB.Exec("ALTER SEQUENCE todos_id_seq RESTART WITH 1")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS todos
(
    id SERIAL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    CONSTRAINT todos_pkey PRIMARY KEY (id)
)`

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/todo", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); "[]" != body {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()

	app.Router.ServeHTTP(recorder, req)

	return recorder
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentTodo(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/todo/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if "Todo not found" != m["error"] {
		t.Errorf("Expected the 'error' key of the response to be set to 'Todo not found'. Got '%s'", m["error"])
	}
}

func TestCreateTodo(t *testing.T) {
	clearTable()

	var jsonStr = []byte(`{"title":"string", "description": "string"}`)

	req, _ := http.NewRequest("POST", "/todo", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if "string" != m["name"] {
		t.Errorf("Expected todo title to be 'string'. Got '%v'", m["title"])
	}

	if "string" != m["description"] {
		t.Errorf("Expected todo description to be 'string'. Got '%v'",
			m["description"])
	}

	if 1.0 != m["id"] {
		t.Errorf("Expected todo ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetTodo(t *testing.T) {
	clearTable()
	addTodos(1)

	req, _ := http.NewRequest("GET", "/todo/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addTodos(count int) {
	if 1 > count {
		count = 1
	}

	for i := 0; i < count; i++ {
		app.DB.Exec("INSERT INTO todos(title, description) VALUES($1, $2)",
			"Todo "+strconv.Itoa(i), "string")
	}
}

func TestUpdateTodo(t *testing.T) {
	clearTable()
	addTodos(1)

	req, _ := http.NewRequest("GET", "/todo/1", nil)
	response := executeRequest(req)

	var origTodo map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &origTodo)

	var jsonStr = []byte(`{"title":"new string", "description": "string"}`)

	req, _ = http.NewRequest("PUT", "/todo/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != origTodo["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v",
			origTodo["id"], m["id"])
	}

	if m["title"] == origTodo["title"] {
		t.Errorf("Expected the title to change from '%v' to '%v'. Got '%v'",
			origTodo["title"], m["title"], m["title"])
	}

	if m["description"] == origTodo["description"] {
		t.Errorf("Expected the description to change from '%v' to '%v'. Got '%v'",
			origTodo["description"], m["description"], m["description"])
	}
}

func TestDeleteTodo(t *testing.T) {
	clearTable()
	addTodos(1)

	req, _ := http.NewRequest("GET", "/todo/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/todo/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/todo/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

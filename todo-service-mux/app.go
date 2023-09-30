//
// @package Showcase-Microservices-Golang
//
// @file Todo main
// @copyright 2021-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package main

import (
	"database/sql"

	"fmt"
	"log"

	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (app *App) Initialize(user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	app.DB, err = sql.Open("postgres", connectionString)
	if nil != err {
		log.Fatal(err)
	}

	app.Router = mux.NewRouter()

	app.initializeRoutes()
}

func (app *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, app.Router))
}

// tom: these are added later
func (app *App) getTodo(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])
	if nil != err {
		respondWithError(writer, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	todo := Todo{ID: id}
	if err := todo.getTodo(app.DB); nil != err {
		switch err {
		case sql.ErrNoRows:
			respondWithError(writer, http.StatusNotFound, "Todo not found")
		default:
			respondWithError(writer, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(writer, http.StatusOK, todo)
}

func respondWithError(writer http.ResponseWriter, code int, message string) {
	respondWithJSON(writer, code, map[string]string{"error": message})
}

func respondWithJSON(writer http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	writer.Write(response)
}

func (app *App) getTodos(writer http.ResponseWriter, request *http.Request) {
	count, _ := strconv.Atoi(request.FormValue("count"))
	start, _ := strconv.Atoi(request.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	todos, err := getTodos(app.DB, start, count)
	if nil != err {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(writer, http.StatusOK, todos)
}

func (app *App) createTodo(writer http.ResponseWriter, request *http.Request) {
	var todo Todo
	decoder := json.NewDecoder(request.Body)

	if err := decoder.Decode(&todo); nil != err {
		respondWithError(writer, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer request.Body.Close()

	if err := todo.createTodo(app.DB); nil != err {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(writer, http.StatusCreated, todo)
}

func (app *App) updateTodo(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])
	if nil != err {
		respondWithError(writer, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	var todo Todo
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&todo); nil != err {
		respondWithError(writer, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer request.Body.Close()
	todo.ID = id

	if err := todo.updateTodo(app.DB); nil != err {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(writer, http.StatusOK, todo)
}

func (app *App) deleteTodo(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])
	if nil != err {
		respondWithError(writer, http.StatusBadRequest, "Invalid Todo ID")
		return
	}

	todo := Todo{ID: id}
	if err := todo.deleteTodo(app.DB); nil != err {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(writer, http.StatusOK, map[string]string{"result": "success"})
}

func (app *App) initializeRoutes() {
	app.Router.HandleFunc("/todo", app.getTodos).Methods("GET")
	app.Router.HandleFunc("/todo", app.createTodo).Methods("POST")
	app.Router.HandleFunc("/todo/{id:[0-9]+}", app.getTodo).Methods("GET")
	app.Router.HandleFunc("/todo/{id:[0-9]+}", app.updateTodo).Methods("PUT")
	app.Router.HandleFunc("/todo/{id:[0-9]+}", app.deleteTodo).Methods("DELETE")
}

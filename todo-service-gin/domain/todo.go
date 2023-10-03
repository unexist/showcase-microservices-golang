//
// @package Showcase-Microservices-Golang
//
// @file Todo model
// @copyright 2021-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package domain

import (
	"database/sql"
)

type Todo struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (todo *Todo) GetTodo(db *sql.DB) error {
	return db.QueryRow("SELECT title, description FROM todos WHERE id=$1",
		todo.ID).Scan(&todo.Title, &todo.Description)
}

func (todo *Todo) UpdateTodo(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE todos SET title=$1, description=$2 WHERE id=$3",
			todo.Title, todo.Description, todo.ID)

	return err
}

func (todo *Todo) DeleteTodo(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM todos WHERE id=$1", todo.ID)

	return err
}

func (todo *Todo) CreateTodo(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO todos(title, description) VALUES($1, $2) RETURNING id",
		todo.Title, todo.Description).Scan(&todo.ID)

	if nil != err {
		return err
	}

	return nil
}

func GetTodos(db *sql.DB, start, count int) ([]Todo, error) {
	rows, err := db.Query(
		"SELECT id, title, description FROM todos LIMIT $1 OFFSET $2",
		count, start)

	if nil != err {
		return nil, err
	}

	defer rows.Close()

	todos := []Todo{}

	for rows.Next() {
		var todo Todo

		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description); nil != err {
			return nil, err
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

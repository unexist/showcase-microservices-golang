//
// @package Showcase-Microservices-Golang
//
// @file Todo SQL repository
// @copyright 2023-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package infrastructure

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"

	"braces.dev/errtrace"
	"github.com/unexist/showcase-microservices-golang/domain"
)

type TodoSQLRepository struct {
	database *sql.DB
}

func NewTodoSQLRepository() *TodoSQLRepository {
	return &TodoSQLRepository{}
}

func (repository *TodoSQLRepository) Open(connectionString string) error {
	var err error

	repository.database, err = sql.Open("postgres", connectionString)

	if nil != err {
		return errtrace.Wrap(err)
	}

	return nil
}

func (repository *TodoSQLRepository) GetTodos() ([]domain.Todo, error) {
	rows, err := repository.database.Query(
		"SELECT id, title, description FROM todos")

	if nil != err {
		return nil, errtrace.Wrap(err)
	}

	defer rows.Close()

	todos := []domain.Todo{}

	for rows.Next() {
		var todo domain.Todo

		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description); nil != err {
			return nil, errtrace.Wrap(err)
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

func (repository *TodoSQLRepository) CreateTodo(todo *domain.Todo) error {
	return errtrace.Wrap(repository.database.QueryRow(
		"INSERT INTO todos(title, description) VALUES($1, $2) RETURNING id",
		todo.Title, todo.Description).Scan(&todo.ID))
}

func (repository *TodoSQLRepository) GetTodo(todoId int) (*domain.Todo, error) {
	todo := domain.Todo{}

	err := repository.database.QueryRow("SELECT id, title, description FROM todos WHERE id=$1",
		todoId).Scan(&todo.ID, &todo.Title, &todo.Description)

	if errors.Is(err, sql.ErrNoRows) {
		err = errors.New("Not found")
	}

	return &todo, errtrace.Wrap(err)
}

func (repository *TodoSQLRepository) UpdateTodo(todo *domain.Todo) error {
	_, err :=
		repository.database.Exec("UPDATE todos SET title=$1, description=$2 WHERE id=$3",
			todo.Title, todo.Description, todo.ID)

	return errtrace.Wrap(err)
}

func (repository *TodoSQLRepository) DeleteTodo(todoId int) error {
	_, err := repository.database.Exec("DELETE FROM todos WHERE id=$1", todoId)

	return errtrace.Wrap(err)
}

func (repository *TodoSQLRepository) Close() error {
	return errtrace.Wrap(repository.database.Close())
}

func (repository *TodoSQLRepository) Clear() error {
	_, err := repository.database.Exec("DELETE FROM todos; ALTER SEQUENCE todos_id_seq RESTART WITH 1")

	return errtrace.Wrap(err)
}

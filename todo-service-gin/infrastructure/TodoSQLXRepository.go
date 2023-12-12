//
// @package Showcase-Microservices-Golang
//
// @file Todo SQLX repository
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

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"braces.dev/errtrace"
	"github.com/unexist/showcase-microservices-golang/domain"
)

type TodoSQLXRepository struct {
	database *sqlx.DB
}

func NewTodoSQLXRepository() *TodoSQLXRepository {
	return &TodoSQLXRepository{}
}

func (repository *TodoSQLXRepository) Open(connectionString string) error {
	var err error

	repository.database, err = sqlx.Connect("postgres", connectionString)

	if nil != err {
		return errtrace.Wrap(err)
	}

	return nil
}

func (repository *TodoSQLXRepository) GetTodos() ([]domain.Todo, error) {
	rows, err := repository.database.Queryx(
		"SELECT id, title, description FROM todos")

	if nil != err {
		return nil, errtrace.Wrap(err)
	}

	defer rows.Close()

	todos := []domain.Todo{}

	for rows.Next() {
		var todo domain.Todo

		if err := rows.StructScan(&todo); nil != err {
			return nil, errtrace.Wrap(err)
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

func (repository *TodoSQLXRepository) CreateTodo(todo *domain.Todo) error {
	return errtrace.Wrap(repository.database.QueryRow(
		"INSERT INTO todos(title, description) VALUES($1, $2) RETURNING id",
		todo.Title, todo.Description).Scan(&todo.ID))
}

func (repository *TodoSQLXRepository) GetTodo(todoId int) (*domain.Todo, error) {
	todo := domain.Todo{}

	err := repository.database.QueryRowx("SELECT id, title, description FROM todos WHERE id=$1",
		todoId).StructScan(&todo)

	if errors.Is(err, sql.ErrNoRows) {
		err = errors.New("Not found")
	}

	return &todo, errtrace.Wrap(err)
}

func (repository *TodoSQLXRepository) UpdateTodo(todo *domain.Todo) error {
	_, err :=
		repository.database.Exec("UPDATE todos SET title=$1, description=$2 WHERE id=$3",
			todo.Title, todo.Description, todo.ID)

	return errtrace.Wrap(err)
}

func (repository *TodoSQLXRepository) DeleteTodo(todoId int) error {
	_, err := repository.database.Exec("DELETE FROM todos WHERE id=$1", todoId)

	return errtrace.Wrap(err)
}

func (repository *TodoSQLXRepository) Close() error {
	return errtrace.Wrap(repository.database.Close())
}

func (repository *TodoSQLXRepository) Clear() error {
	_, err := repository.database.Exec("DELETE FROM todos; ALTER SEQUENCE todos_id_seq RESTART WITH 1")

	return errtrace.Wrap(err)
}

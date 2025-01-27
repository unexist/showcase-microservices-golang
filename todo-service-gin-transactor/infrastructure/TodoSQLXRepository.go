//
// @package Showcase-Microservices-Golang
//
// @file Todo SQLX repository
// @copyright 2025-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package infrastructure

import (
	"context"
	"database/sql"
	"errors"

	"braces.dev/errtrace"

	"github.com/unexist/showcase-microservices-golang/domain/todo"

	sqlxTransactor "github.com/Thiht/transactor/sqlx"
)

type TodoSQLXRepository struct {
	dbGetter sqlxTransactor.DBGetter
}

func NewTodoSQLXRepository(dbGetter sqlxTransactor.DBGetter) *TodoSQLXRepository {
	return &TodoSQLXRepository{
		dbGetter: dbGetter,
	}
}

func (repository *TodoSQLXRepository) GetTodos(ctx context.Context) ([]domain.Todo, error) {
	rows, err := repository.dbGetter(ctx).Queryx(
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

func (repository *TodoSQLXRepository) CreateTodo(ctx context.Context, todo *domain.Todo) error {
	return errtrace.Wrap(repository.dbGetter(ctx).QueryRow(
		"INSERT INTO todos(title, description) VALUES($1, $2) RETURNING id",
		todo.Title, todo.Description).Scan(&todo.ID))
}

func (repository *TodoSQLXRepository) GetTodo(ctx context.Context, todoId int) (*domain.Todo, error) {
	todo := domain.Todo{}

	err := repository.dbGetter(ctx).QueryRowx("SELECT id, title, description FROM todos WHERE id=$1",
		todoId).StructScan(&todo)

	if errors.Is(err, sql.ErrNoRows) {
		err = errors.New("Not found")
	}

	return &todo, errtrace.Wrap(err)
}

func (repository *TodoSQLXRepository) UpdateTodo(ctx context.Context, todo *domain.Todo) error {
	_, err :=
		repository.dbGetter(ctx).Exec("UPDATE todos SET title=$1, description=$2 WHERE id=$3",
			todo.Title, todo.Description, todo.ID)

	return errtrace.Wrap(err)
}

func (repository *TodoSQLXRepository) DeleteTodo(ctx context.Context, todoId int) error {
	_, err := repository.dbGetter(ctx).Exec("DELETE FROM todos WHERE id=$1", todoId)

	return errtrace.Wrap(err)
}

func (repository *TodoSQLXRepository) Clear(ctx context.Context) error {
	_, err := repository.dbGetter(ctx).Exec("DELETE FROM todos")
	_, err = repository.dbGetter(ctx).Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name = 'todos'")

	return errtrace.Wrap(err)
}

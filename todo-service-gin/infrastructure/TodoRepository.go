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

package infrastructure

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/unexist/showcase-microservices-golang/domain"
)

type TodoRepository struct {
	database *sql.DB
}

func NewTodoRepository(database *sql.DB) *TodoRepository {
	return &TodoRepository{
		database: database,
	}
}

func (repository *TodoRepository) GetTodos() ([]domain.Todo, error) {
	rows, err := repository.database.Query(
		"SELECT id, title, description FROM todos")

	if nil != err {
		return nil, err
	}

	defer rows.Close()

	todos := []domain.Todo{}

	for rows.Next() {
		var todo domain.Todo

		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description); nil != err {
			return nil, err
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

func (repository *TodoRepository) CreateTodo(todo *domain.Todo) error {
	return repository.database.QueryRow(
		"INSERT INTO todos(title, description) VALUES($1, $2) RETURNING id",
		todo.Title, todo.Description).Scan(&todo.ID)
}

func (repository *TodoRepository) GetTodo(todoId int) (*domain.Todo, error) {
	todo := domain.Todo{}

	err := repository.database.QueryRow("SELECT title, description FROM todos WHERE id=$1",
		todoId).Scan(&todo.Title, &todo.Description)

	return &todo, err
}

func (repository *TodoRepository) UpdateTodo(todo *domain.Todo) error {
	_, err :=
		repository.database.Exec("UPDATE todos SET title=$1, description=$2 WHERE id=$3",
			todo.Title, todo.Description, todo.ID)

	return err
}

func (repository *TodoRepository) DeleteTodo(todoId int) error {
	_, err := repository.database.Exec("DELETE FROM todos WHERE id=$1", todoId)

	return err
}

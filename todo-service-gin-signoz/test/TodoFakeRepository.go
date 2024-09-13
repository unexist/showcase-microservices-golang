//
// @package Showcase-Microservices-Golang
//
// @file Todo fake repository
// @copyright 2023-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package test

import (
	"context"
	"errors"

	"braces.dev/errtrace"

	"github.com/unexist/showcase-microservices-golang/domain"
)

type TodoFakeRepository struct {
	todos []domain.Todo
}

func NewTodoFakeRepository() *TodoFakeRepository {
	return &TodoFakeRepository{
		todos: make([]domain.Todo, 0),
	}
}

func (repository *TodoFakeRepository) Open(_ string) error {
	return nil
}

func (repository *TodoFakeRepository) GetTodos(_ context.Context) ([]domain.Todo, error) {
	return repository.todos, nil
}

func (repository *TodoFakeRepository) CreateTodo(_ context.Context, todo *domain.Todo) error {
	newTodo := domain.Todo{
		ID:          len(repository.todos) + 1,
		Title:       todo.Title,
		Description: todo.Description,
	}

	todo.ID = newTodo.ID

	repository.todos = append(repository.todos, newTodo)

	return nil
}

func (repository *TodoFakeRepository) GetTodo(_ context.Context, todoId int) (*domain.Todo, error) {
	for i := 0; i < len(repository.todos); i++ {
		if repository.todos[i].ID == todoId {
			return &repository.todos[i], nil
		}
	}

	return nil, errtrace.Wrap(errors.New("Not found"))
}

func (repository *TodoFakeRepository) UpdateTodo(_ context.Context, todo *domain.Todo) error {
	for i := 0; i < len(repository.todos); i++ {
		if repository.todos[i].ID == todo.ID {
			repository.todos[i].Title = todo.Title
			repository.todos[i].Description = todo.Description

			return nil
		}
	}

	return errtrace.Wrap(errors.New("Not found"))
}

func (repository *TodoFakeRepository) DeleteTodo(_ context.Context, todoId int) error {
	for i := 0; i < len(repository.todos); i++ {
		if repository.todos[i].ID == todoId {
			repository.todos = removeIndex(repository.todos, i)

			return nil
		}
	}

	return errtrace.Wrap(errors.New("Not found"))
}

func (repository *TodoFakeRepository) Close() error {
	return nil
}

func (repository *TodoFakeRepository) Clear() error {
	repository.todos = make([]domain.Todo, 0)

	return nil
}

func removeIndex(s []domain.Todo, idx int) []domain.Todo {
	return append(s[:idx], s[idx+1:]...)
}

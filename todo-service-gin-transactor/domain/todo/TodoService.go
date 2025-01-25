//
// @package Showcase-Microservices-Golang
//
// @file Todo service
// @copyright 2025-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package domain

import (
	"context"

	"braces.dev/errtrace"
)

type TodoService struct {
	repository TodoRepository
}

func NewTodoService(repository TodoRepository) *TodoService {
	return &TodoService{
		repository: repository,
	}
}

func (service *TodoService) GetTodos(ctx context.Context) ([]Todo, error) {
	return errtrace.Wrap2(service.repository.GetTodos(ctx))
}

func (service *TodoService) CreateTodo(ctx context.Context, todo *Todo) error {
	return errtrace.Wrap(service.repository.CreateTodo(ctx, todo))
}

func (service *TodoService) GetTodo(ctx context.Context, todoId int) (*Todo, error) {
	return errtrace.Wrap2(service.repository.GetTodo(ctx, todoId))
}

func (service *TodoService) UpdateTodo(ctx context.Context, todo *Todo) error {
	return errtrace.Wrap(service.repository.UpdateTodo(ctx, todo))
}

func (service *TodoService) DeleteTodo(ctx context.Context, todoId int) error {
	return errtrace.Wrap(service.repository.DeleteTodo(ctx, todoId))
}

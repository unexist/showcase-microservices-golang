//
// @package Showcase-Microservices-Golang
//
// @file Todo service
// @copyright 2023-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package domain

import (
	"context"

	"braces.dev/errtrace"
	"go.opentelemetry.io/otel"
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
	tracer := otel.GetTracerProvider().Tracer("todo-service")
	_, span := tracer.Start(ctx, "get-todos")
	defer span.End()

	return errtrace.Wrap2(service.repository.GetTodos())
}

func (service *TodoService) CreateTodo(todo *Todo, ctx context.Context) error {
	tracer := otel.GetTracerProvider().Tracer("todo-service")
	_, span := tracer.Start(ctx, "create-todo")
	defer span.End()

	return errtrace.Wrap(service.repository.CreateTodo(todo))
}

func (service *TodoService) GetTodo(todoId int, ctx context.Context) (*Todo, error) {
	tracer := otel.GetTracerProvider().Tracer("todo-service")
	_, span := tracer.Start(ctx, "get-todo")
	defer span.End()

	return errtrace.Wrap2(service.repository.GetTodo(todoId))
}

func (service *TodoService) UpdateTodo(todo *Todo, ctx context.Context) error {
	tracer := otel.GetTracerProvider().Tracer("todo-service")
	_, span := tracer.Start(ctx, "update-todo")
	defer span.End()

	return errtrace.Wrap(service.repository.UpdateTodo(todo))
}

func (service *TodoService) DeleteTodo(todoId int, ctx context.Context) error {
	tracer := otel.GetTracerProvider().Tracer("todo-service")
	_, span := tracer.Start(ctx, "update-todo")
	defer span.End()

	return errtrace.Wrap(service.repository.DeleteTodo(todoId))
}

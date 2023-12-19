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

	todo, err := service.repository.GetTodos()

	span.End()

	return errtrace.Wrap2(todo, err)
}

func (service *TodoService) CreateTodo(todo *Todo, ctx context.Context) error {
	tracer := otel.GetTracerProvider().Tracer("todo-service")
	_, span := tracer.Start(ctx, "create-todo")

	err := service.repository.CreateTodo(todo)

	span.End()

	return errtrace.Wrap(err)
}

func (service *TodoService) GetTodo(todoId int, ctx context.Context) (*Todo, error) {
	tracer := otel.GetTracerProvider().Tracer("todo-service")
	_, span := tracer.Start(ctx, "get-todo")

	todo, err := service.repository.GetTodo(todoId)

	span.End()

	return errtrace.Wrap2(todo, err)
}

func (service *TodoService) UpdateTodo(todo *Todo, ctx context.Context) error {
	tracer := otel.GetTracerProvider().Tracer("todo-service")
	_, span := tracer.Start(ctx, "update-todo")

	err := service.repository.UpdateTodo(todo)

	span.End()

	return errtrace.Wrap(err)
}

func (service *TodoService) DeleteTodo(todoId int, ctx context.Context) error {
	tracer := otel.GetTracerProvider().Tracer("todo-service")
	_, span := tracer.Start(ctx, "update-todo")

	err := service.repository.DeleteTodo(todoId)

	span.End()

	return errtrace.Wrap(err)
}

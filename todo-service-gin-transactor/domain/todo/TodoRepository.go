//
// @package Showcase-Microservices-Golang
//
// @file Todo repository interface
// @copyright 2025-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package domain

import (
	"context"
)

type TodoRepository interface {

	// Get all todos stored by this repository
	GetTodos(ctx context.Context) ([]Todo, error)

	// Create new todo based on given values
	CreateTodo(ctx context.Context, todo *Todo) error

	// Get todo entry with given id
	GetTodo(ctx context.Context, todoId int) (*Todo, error)

	// Update todo entry with given id
	UpdateTodo(ctx context.Context, todo *Todo) error

	// Delete todo entry with given id
	DeleteTodo(ctx context.Context, todoId int) error

	// Clear table
	Clear(ctx context.Context) error
}

//
// @package Showcase-Microservices-Golang
//
// @file Todo repository interface
// @copyright 2023-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package domain

type TodoRepository interface {
	// Open connection to database
	Open(connectionString string) error

	// Get all todos stored by this repository
	GetTodos() ([]Todo, error)

	// Create new todo based on given values
	CreateTodo(todo *Todo) error

	// Get todo entry with given id
	GetTodo(todoId int) (*Todo, error)

	// Update todo entry with given id
	UpdateTodo(todo *Todo) error

	// Delete todo entry with given id
	DeleteTodo(todoId int) error

	// Clear table
	Clear() error

	// Close database connection
	Close() error
}

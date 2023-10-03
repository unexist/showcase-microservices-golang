//
// @package Showcase-Microservices-Golang
//
// @file Todo repository interface
// @copyright 2021-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package domain

type TodoRepository interface {
	GetTodos() ([]Todo, error)
	CreateTodo(todo *Todo) error
	GetTodo(todoId int) (*Todo, error)
	UpdateTodo(todo *Todo) error
	DeleteTodo(todoId int) error
}

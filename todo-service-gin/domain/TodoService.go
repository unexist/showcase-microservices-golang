//
// @package Showcase-Microservices-Golang
//
// @file Todo service
// @copyright 2021-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package domain

type TodoService struct {
	repository TodoRepository
}

func NewTodoService(repository TodoRepository) *TodoService {
	return &TodoService{
		repository: repository,
	}
}

func (service *TodoService) GetTodos() ([]Todo, error) {
	return service.repository.GetTodos()
}

func (service *TodoService) CreateTodo(todo *Todo) error {
	return service.repository.CreateTodo(todo)
}

func (service *TodoService) GetTodo(todoId int) (*Todo, error) {
	return service.repository.GetTodo(todoId)
}

func (service *TodoService) UpdateTodo(todo *Todo) error {
	return service.repository.UpdateTodo(todo)
}

func (service *TodoService) DeleteTodo(todoId int) error {
	return service.repository.DeleteTodo(todoId)
}

//
// @package Showcase-Microservices-Golang
//
// @file User resource
// @copyright 2023-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package application

import (
	todoDomain "github.com/unexist/showcase-microservices-golang/domain/todo"
	userDomain "github.com/unexist/showcase-microservices-golang/domain/user"
)

type TodoUserService struct {
	todoService *todoDomain.TodoService
	userService *userDomain.UserService
}

func NewTodoUserService(todoService *todoDomain.TodoService, userService *userDomain.UserService) *TodoUserService {
	return &TodoUserService{
		todoService: todoService,
		userService: userService,
	}
}

func (service *TodoUserService) CreateAnonTodo(todo *todoDomain.Todo) error {
	anonUser := userDomain.User{
		Name: "Anon User",
	}

	/* Start unit-of-work */
	if err := service.userService.CreateUser(&anonUser); nil != err {
		/* Roll back */
		return err
	}

	todo.UserID = anonUser.ID

	if err := service.todoService.CreateTodo(todo); nil != err {
		/* Roll back */
		return err
	}
	/* End unit-of-work */

	return nil
}

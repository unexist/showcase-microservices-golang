//
// @package Showcase-Microservices-Golang
//
// @file User resource
// @copyright 2025-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package application

import (
	"context"

	"github.com/Thiht/transactor"

	todoDomain "github.com/unexist/showcase-microservices-golang/domain/todo"
	userDomain "github.com/unexist/showcase-microservices-golang/domain/user"
)

type TodoUserService struct {
	transactor  transactor.Transactor
	todoService *todoDomain.TodoService
	userService *userDomain.UserService
}

func NewTodoUserService(transactor transactor.Transactor, todoService *todoDomain.TodoService,
	userService *userDomain.UserService) *TodoUserService {
	return &TodoUserService{
		transactor:  transactor,
		todoService: todoService,
		userService: userService,
	}
}

func (service *TodoUserService) CreateAnonTodo(todo *todoDomain.Todo) error {
	anonUser := userDomain.User{
		Name: "Anon User",
	}

	/* Start unit-of-work */
	return service.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := service.userService.CreateUser(&anonUser); nil != err {
			/* Roll back */
			return err
		}

		todo.UserID = anonUser.ID

		if err := service.todoService.CreateTodo(todo); nil != err {
			/* Roll back */
			return err
		}

		return nil
	})
	/* End unit-of-work */
}

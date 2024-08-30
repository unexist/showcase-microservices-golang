//
// @package Showcase-Microservices-Golang
//
// @file Todo SQL repository
// @copyright 2023-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package infrastructure

import (
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"braces.dev/errtrace"
	"github.com/unexist/showcase-microservices-golang/domain"
)

type TodoGormRepository struct {
	database *gorm.DB
}

func NewTodoGormRepository() *TodoGormRepository {
	return &TodoGormRepository{}
}

func (repository *TodoGormRepository) Open(connectionString string) (err error) {
	repository.database, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if nil != err {
		err = errtrace.Wrap(err)
		return
	}

	repository.database.AutoMigrate(&domain.Todo{})

	err = errtrace.Wrap(err)
	return
}

func (repository *TodoGormRepository) GetTodos() ([]domain.Todo, error) {
	todos := []domain.Todo{}

	result := repository.database.Find(&todos)

	if nil != result.Error {
		return nil, errtrace.Wrap(result.Error)
	}

	return todos, nil
}

func (repository *TodoGormRepository) CreateTodo(todo *domain.Todo) error {
	result := repository.database.Create(todo)

	if nil != result.Error {
		return errtrace.Wrap(result.Error)
	}

	return nil
}

func (repository *TodoGormRepository) GetTodo(todoId int) (*domain.Todo, error) {
	var err error

	todo := domain.Todo{ID: todoId}

	result := repository.database.First(&todo)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err = errors.New("Not found")
	} else {
		err = nil
	}

	return &todo, errtrace.Wrap(err)
}

func (repository *TodoGormRepository) UpdateTodo(todo *domain.Todo) (err error) {
	result := repository.database.Save(todo)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err = errors.New("Not found")
	} else {
		err = nil
	}

	err = errtrace.Wrap(err)
	return
}

func (repository *TodoGormRepository) DeleteTodo(todoId int) (err error) {
	result := repository.database.Delete(&domain.Todo{}, todoId)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err = errors.New("Not found")
	} else {
		err = nil
	}

	err = errtrace.Wrap(err)
	return
}

func (repository *TodoGormRepository) Clear() error {
	result := repository.database.Exec(
		"DELETE FROM todos; ALTER SEQUENCE todos_id_seq RESTART WITH 1")

	if nil != result.Error {
		return errtrace.Wrap(result.Error)
	}

	return nil
}

func (repository *TodoGormRepository) Close() error {
	return nil
}

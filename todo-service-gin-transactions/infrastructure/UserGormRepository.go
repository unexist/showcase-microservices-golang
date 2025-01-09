//
// @package Showcase-Microservices-Golang
//
// @file User fake repository
// @copyright 2023-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package infrastructure

import (
	"errors"

	uuid2 "github.com/hashicorp/go-uuid"
	"gorm.io/gorm"

	"braces.dev/errtrace"

	"github.com/unexist/showcase-microservices-golang/domain/user"
)

type UserGormRepository struct {
	database *gorm.DB
}

func NewUserGormRepository(db *gorm.DB) *UserGormRepository {
	return &UserGormRepository{
		database: db,
	}
}

func (repository *UserGormRepository) GetUsers() ([]domain.User, error) {
	todos := []domain.User{}

	result := repository.database.Find(&todos)

	if nil != result.Error {
		return nil, errtrace.Wrap(result.Error)
	}

	return todos, nil
}

func (repository *UserGormRepository) CreateUser(user *domain.User) error {
	uuid, _ := uuid2.GenerateUUID()

	user.Token = uuid

	result := repository.database.Create(user)

	if nil != result.Error {
		return errtrace.Wrap(result.Error)
	}

	return nil
}

func (repository *UserGormRepository) GetUser(userId int) (*domain.User, error) {
	var err error

	user := domain.User{ID: userId}

	result := repository.database.First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err = errors.New("Not found")
	} else {
		err = nil
	}

	return &user, errtrace.Wrap(err)
}

func (repository *UserGormRepository) GetUserByToken(token string) (*domain.User, error) {
	var err error

	user := domain.User{Token: token}

	result := repository.database.First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err = errors.New("Not found")
	} else {
		err = nil
	}

	return &user, errtrace.Wrap(err)
}

func (repository *UserGormRepository) Clear() error {
	result := repository.database.Exec(
		"DELETE FROM users; ALTER SEQUENCE todos_id_seq RESTART WITH 1")

	if nil != result.Error {
		return errtrace.Wrap(result.Error)
	}

	return nil
}

func (repository *UserGormRepository) Close() error {
	return nil
}

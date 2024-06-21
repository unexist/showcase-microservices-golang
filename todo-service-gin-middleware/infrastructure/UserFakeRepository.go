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

	"braces.dev/errtrace"
	"github.com/unexist/showcase-microservices-golang/domain/user"
)

type UserFakeRepository struct {
	users []domain.User
}

func NewUserFakeRepository() *UserFakeRepository {
	return &UserFakeRepository{
		users: make([]domain.User, 0),
	}
}

func (repository *UserFakeRepository) Open(_ string) error {
	return nil
}

func (repository *UserFakeRepository) CreateUser(user *domain.User) error {
	uuid, _ := uuid2.GenerateUUID()

	newUser := domain.User{
		ID:    len(repository.users) + 1,
		Name:  user.Name,
		Token: uuid,
	}

	user.ID = newUser.ID
	user.Token = newUser.Token

	repository.users = append(repository.users, newUser)

	return nil
}

func (repository *UserFakeRepository) GetUser(userId int) (*domain.User, error) {
	for i := 0; i < len(repository.users); i++ {
		if repository.users[i].ID == userId {
			return &repository.users[i], nil
		}
	}

	return nil, errtrace.Wrap(errors.New("not found"))
}

func (repository *UserFakeRepository) GetUserByToken(token string) (*domain.User, error) {
	for i := 0; i < len(repository.users); i++ {
		if repository.users[i].Token == token {
			return &repository.users[i], nil
		}
	}

	return nil, errtrace.Wrap(errors.New("not found"))
}

func (repository *UserFakeRepository) Close() error {
	return nil
}

func (repository *UserFakeRepository) Clear() error {
	repository.users = make([]domain.User, 0)

	return nil
}

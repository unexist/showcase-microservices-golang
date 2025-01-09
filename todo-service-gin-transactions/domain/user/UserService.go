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

import "braces.dev/errtrace"

type UserService struct {
	repository UserRepository
}

func NewUserService(repository UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (service *UserService) CreateUser(user *User) error {
	return errtrace.Wrap(service.repository.CreateUser(user))
}

func (service *UserService) GetUser(userId int) (*User, error) {
	return errtrace.Wrap2(service.repository.GetUser(userId))
}

func (service *UserService) ValidateToken(token string) (*User, error) {
	/* <Insert fancy token validate mechanic here> */
	return errtrace.Wrap2(service.repository.GetUserByToken(token))
}

//
// @package Showcase-Microservices-Golang
//
// @file User repository interface
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

type UserRepository interface {

	// Create new user based on given values
	CreateUser(ctx context.Context, user *User) error

	// Get user with given id
	GetUser(ctx context.Context, userId int) (*User, error)

	// Get user with given token
	GetUserByToken(ctx context.Context, token string) (*User, error)

	// Clear table
	Clear(ctx context.Context) error
}

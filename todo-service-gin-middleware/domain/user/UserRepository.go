//
// @package Showcase-Microservices-Golang
//
// @file User repository interface
// @copyright 2023-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package domain

type UserRepository interface {
	// Open connection to database
	Open(connectionString string) error

	// Create new user based on given values
	CreateUser(user *User) error

	// Get user with given id
	GetUser(userId int) (*User, error)

	// Get user with given token
	GetUserByToken(token string) (*User, error)

	// Clear table
	Clear() error

	// Close database connection
	Close() error
}

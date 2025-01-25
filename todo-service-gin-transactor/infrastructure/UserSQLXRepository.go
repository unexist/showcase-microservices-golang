//
// @package Showcase-Microservices-Golang
//
// @file User SQLX repository
// @copyright 2025-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package infrastructure

import (
	"context"
	"database/sql"
	"errors"

	"braces.dev/errtrace"
	uuid2 "github.com/hashicorp/go-uuid"

	"github.com/unexist/showcase-microservices-golang/domain/user"

	sqlxTransactor "github.com/Thiht/transactor/sqlx"
)

type UserSQLXRepository struct {
	dbGetter sqlxTransactor.DBGetter
}

func NewUserSQLXRepository(dbGetter sqlxTransactor.DBGetter) *UserSQLXRepository {
	return &UserSQLXRepository{
		dbGetter: dbGetter,
	}
}

func (repository *UserSQLXRepository) GetUsers(ctx context.Context) ([]domain.User, error) {
	rows, err := repository.dbGetter(ctx).Queryx(
		"SELECT id, name, token FROM users")

	if nil != err {
		return nil, errtrace.Wrap(err)
	}

	defer rows.Close()

	users := []domain.User{}

	for rows.Next() {
		var user domain.User

		if err := rows.StructScan(&user); nil != err {
			return nil, errtrace.Wrap(err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (repository *UserSQLXRepository) CreateUser(ctx context.Context, user *domain.User) error {
	uuid, _ := uuid2.GenerateUUID()

	return errtrace.Wrap(repository.dbGetter(ctx).QueryRow(
		"INSERT INTO users(name, token) VALUES($1, $2) RETURNING id",
		user.Name, uuid).Scan(&user.ID))
}

func (repository *UserSQLXRepository) GetUser(ctx context.Context, userId int) (*domain.User, error) {
	user := domain.User{}

	err := repository.dbGetter(ctx).QueryRowx("SELECT id, name, token FROM users WHERE id=$1",
		userId).StructScan(&user)

	if errors.Is(err, sql.ErrNoRows) {
		err = errors.New("Not found")
	}

	return &user, errtrace.Wrap(err)
}

func (repository *UserSQLXRepository) GetUserByToken(ctx context.Context, token string) (*domain.User, error) {
	user := domain.User{}

	err := repository.dbGetter(ctx).QueryRowx("SELECT id, name, token FROM users WHERE token$1",
		token).StructScan(&user)

	if errors.Is(err, sql.ErrNoRows) {
		err = errors.New("Not found")
	}

	return &user, errtrace.Wrap(err)
}

func (repository *UserSQLXRepository) Clear(ctx context.Context) error {
	_, err := repository.dbGetter(ctx).Exec("DELETE FROM users; ALTER SEQUENCE users_id_seq RESTART WITH 1")

	return errtrace.Wrap(err)
}

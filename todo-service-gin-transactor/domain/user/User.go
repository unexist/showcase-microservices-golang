//
// @package Showcase-Microservices-Golang
//
// @file User model
// @copyright 2025-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package domain

import (
	"fmt"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

func (user User) String() string {
	return fmt.Sprintf("ID: %s\nName: %s\nToken: %s",
		user.ID, user.Name, user.Token)
}

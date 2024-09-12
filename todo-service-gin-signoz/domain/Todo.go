//
// @package Showcase-Microservices-Golang
//
// @file Todo model
// @copyright 2023-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package domain

import (
	"fmt"
)

type Todo struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	UUID        string `json:"uuid"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (todo Todo) String() string {
	return fmt.Sprintf("ID: %s\nUUID: %s\nTitle: %s\nDescription: %s",
		todo.ID, todo.UUID, todo.Title, todo.Description)
}

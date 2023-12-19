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
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (todo Todo) String() string {
	return fmt.Sprintf("ID: %s\nTitle: %s\nDescription: %s",
		todo.ID, todo.Title, todo.Description)
}

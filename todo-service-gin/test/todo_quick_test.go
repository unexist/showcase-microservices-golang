//
// @package Showcase-Microservices-Golang
//
// @file Todo tests with property-based testing
// @copyright 2023-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

//go:build quick

package test

import (
	"testing"
	"testing/quick"

	"github.com/unexist/showcase-microservices-golang/domain"
)

func Test_Property_Todo(t *testing.T) {
	property := func(id int, title, description string) bool {
		todo := domain.Todo{
			ID:          id,
			Title:       title,
			Description: description,
		}

		return 0 != todo.ID
	}

	if err := quick.Check(property, &quick.Config{
		MaxCount: 10,
	}); nil != err {
		t.Error(err)
	}
}

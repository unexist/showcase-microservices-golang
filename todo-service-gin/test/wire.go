//
// @package Showcase-Microservices-Golang
//
// @file Todo wire
// @copyright 2023-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

//go:build wireinject

package test

import (
	"github.com/google/wire"

	"github.com/unexist/showcase-microservices-golang/adapter"
	"github.com/unexist/showcase-microservices-golang/domain"
)

func InitializeResource() *adapter.TodoResource {
	panic(wire.Build(
		// Set binding for concrete interface
		wire.Bind(new(domain.TodoRepository), new(*TodoFakeRepository)),

		// Provider
		NewTodoFakeRepository,
		domain.NewTodoService,
		adapter.NewTodoResource))
}

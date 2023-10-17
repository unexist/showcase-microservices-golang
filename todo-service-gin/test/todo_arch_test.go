//
// @package Showcase-Microservices-Golang
//
// @file Todo tests with wire
// @copyright 2023-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

//go:build arch

package test

import (
	"testing"

	"github.com/datosh/gau"
)

func Test_Layer_Adapter(t *testing.T) {
	gau.Packages(t, "github.com/unexist/showcase-microservices-golang/...").That().
		ResideIn("github.com/unexist/showcase-microservices-golang/adapter").
		Should().DirectlyDependOn("github.com/unexist/showcase-microservices-golang/domain")
}

func Test_Layer_Domain(t *testing.T) {
	gau.Packages(t, "github.com/unexist/showcase-microservices-golang/...").That().
		ResideIn("github.com/unexist/showcase-microservices-golang/domain").
		ShouldNot().DirectlyDependOn("github.com/unexist/showcase-microservices-golang/infrastructure")
}

func Test_Layer_Infrastructure(t *testing.T) {
	gau.Packages(t, "github.com/unexist/showcase-microservices-golang/...").That().
		ResideIn("github.com/unexist/showcase-microservices-golang/infrastructure").
		ShouldNot().DirectlyDependOn("github.com/unexist/showcase-microservices-golang/adapter")
}

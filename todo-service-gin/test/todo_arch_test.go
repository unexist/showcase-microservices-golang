//
// @package Showcase-Microservices-Golang
//
// @file Todo tests with gau
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
	pkgs := gau.Packages(t, "github.com/unexist/showcase-microservices-golang/...").That().
		ResideIn("github.com/unexist/showcase-microservices-golang/adapter")

	pkgs.Should().DirectlyDependOn("github.com/unexist/showcase-microservices-golang/domain")
	//pkgs.ShouldNot().DirectlyDependOn("github.com/unexist/showcase-microservices-golang/infrastructure")
}

func Test_Layer_Domain(t *testing.T) {
	pkgs := gau.Packages(t, "github.com/unexist/showcase-microservices-golang/...").That().
		ResideIn("github.com/unexist/showcase-microservices-golang/domain")

	pkgs.ShouldNot().DirectlyDependOn("github.com/unexist/showcase-microservices-golang/infrastructure")
	pkgs.ShouldNot().DirectlyDependOn("github.com/unexist/showcase-microservices-golang/adapter")
}

func Test_Layer_Infrastructure(t *testing.T) {
	pkgs := gau.Packages(t, "github.com/unexist/showcase-microservices-golang/...").That().
		ResideIn("github.com/unexist/showcase-microservices-golang/infrastructure")

	pkgs.ShouldNot().DirectlyDependOn("github.com/unexist/showcase-microservices-golang/adapter")
	pkgs.Should().DirectlyDependOn("github.com/unexist/showcase-microservices-golang/domain")
}

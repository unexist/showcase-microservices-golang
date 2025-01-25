//
// @package Showcase-Microservices-Golang
//
// @file User tests for Gorm repository
// @copyright 2023-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package test

import (
	"fmt"
	"testing"

	"net/http"
)

func TestLogin(t *testing.T) {
	todoRepository.Clear()
	userRepository.Clear()

	req, _ := http.NewRequest("POST", "/user/login", nil)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestSelf(t *testing.T) {
	todoRepository.Clear()
	userRepository.Clear()

	req, _ := http.NewRequest("GET", "/user/self", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", executeLogin(t)))

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

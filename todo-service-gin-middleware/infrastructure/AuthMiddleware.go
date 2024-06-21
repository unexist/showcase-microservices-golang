//
// @package Showcase-Microservices-Golang
//
// @file Authentication middleware
// @copyright 2023-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package infrastructure

import (
	"github.com/gin-gonic/gin"
	"github.com/unexist/showcase-microservices-golang/domain/user"
	"net/http"
	"strings"
)

type authHeader struct {
	value string `header:"Authorization"`
}

func AuthUser(userService *domain.UserService) gin.HandlerFunc {
	return func(context *gin.Context) {
		header := authHeader{}

		if err := context.ShouldBindHeader(&header); nil != err {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			context.Abort()

			return
		}

		print("foobar= " + header.value)

		bearerHeader := strings.Split(header.value, "Bearer ")

		if 2 > len(bearerHeader) {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Must provide Authorization header with format `Bearer {token}`"})
			context.Abort()

			return
		}

		user, err := userService.ValidateToken(bearerHeader[1])

		if nil != err {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid"})
			context.Abort()

			return
		}

		context.Set("user", user)

		context.Next()
	}
}

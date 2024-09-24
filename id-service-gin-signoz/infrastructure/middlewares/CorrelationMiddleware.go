//
// @package Showcase-Microservices-Golang
//
// @file Correlation middleware
// @copyright 2024-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-uuid"
)

func CorrelationMiddleware() gin.HandlerFunc {
	return addCorrelationId
}

func addCorrelationId(c *gin.Context) {
	correlationId := c.Request.Header.Get("CorrelationId")

	if "" == strings.TrimSpace(correlationId) {
		correlationId, _ = uuid.GenerateUUID()
	}

	c.Set("correlation_id", correlationId)
	c.Request.Header.Add("CorrelationId", correlationId)

	c.Next()
}

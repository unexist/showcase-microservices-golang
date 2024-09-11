//
// @package Showcase-Microservices-Golang
//
// @file Todo resource
// @copyright 2023-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package infrastructure

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	todoHttpStatusCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "todo_http_status_counter",
			Help: "Total number of requests with each status code",
		},
		[]string{"item_type"},
	)
)

func init() {
	prometheus.MustRegister(todoHttpStatusCounter)
}

func HttpStatusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		statusCode := c.Writer.Status()

		if 200 <= statusCode && statusCode <= 299 {
			todoHttpStatusCounter.WithLabelValues(fmt.Sprintf("HTTP%d", statusCode)).Inc()
		} else if 400 <= statusCode && statusCode <= 499 {
			todoHttpStatusCounter.WithLabelValues(fmt.Sprintf("HTTP%d", statusCode)).Inc()
		} else if 500 <= statusCode && statusCode <= 599 {
			todoHttpStatusCounter.WithLabelValues(fmt.Sprintf("HTTP%d", statusCode)).Inc()
		}

		c.Next()
	}
}

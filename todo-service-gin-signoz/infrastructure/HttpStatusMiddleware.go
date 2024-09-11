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
		[]string{"code"},
	)
)

func init() {
	prometheus.MustRegister(todoHttpStatusCounter)
}

func HttpStatusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		statusCode := c.Writer.Status()

		if 200 <= statusCode && 299 >= statusCode {
			todoHttpStatusCounter.WithLabelValues(fmt.Sprintf("%d", statusCode)).Inc()
		} else if 400 <= statusCode && 499 >= statusCode {
			todoHttpStatusCounter.WithLabelValues(fmt.Sprintf("%d", statusCode)).Inc()
		} else if 500 <= statusCode && 599 >= statusCode {
			todoHttpStatusCounter.WithLabelValues(fmt.Sprintf("%d", statusCode)).Inc()
		}

		c.Next()
	}
}

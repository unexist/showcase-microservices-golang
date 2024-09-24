//
// @package Showcase-Microservices-Golang
//
// @file HTTP status middleware
// @copyright 2024-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package middlewares

import (
	"fmt"
	"time"

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
	todoHttpLatency = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "todo_http_latency",
			Help: "Total time taken for HTTP requests",
		},
	)
)

func init() {
	prometheus.MustRegister(todoHttpStatusCounter, todoHttpLatency)
}

func HttpStatusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		statusCode := c.Writer.Status()

		if 200 <= statusCode && 299 >= statusCode {
			todoHttpStatusCounter.WithLabelValues(fmt.Sprintf("%d", statusCode)).Inc()
		} else if 400 <= statusCode && 499 >= statusCode {
			todoHttpStatusCounter.WithLabelValues(fmt.Sprintf("%d", statusCode)).Inc()
		} else if 500 <= statusCode && 599 >= statusCode {
			todoHttpStatusCounter.WithLabelValues(fmt.Sprintf("%d", statusCode)).Inc()
		}

		latency := time.Now().Sub(startTime)

		if latency > time.Minute {
			latency = latency.Truncate(time.Second)
		}

		todoHttpLatency.Set(float64(latency.Milliseconds()))
	}
}

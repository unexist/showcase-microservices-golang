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

package infrastructure

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	idHttpStatusCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "id_http_status_counter",
			Help: "Total number of requests with each status code",
		},
		[]string{"code"},
	)
	idHttpLatency = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "id_http_latency",
			Help: "Total time taken for HTTP requests",
		},
	)
)

func init() {
	prometheus.MustRegister(idHttpStatusCounter, idHttpLatency)
}

func HttpStatusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		statusCode := c.Writer.Status()

		if 200 <= statusCode && 299 >= statusCode {
			idHttpStatusCounter.WithLabelValues(fmt.Sprintf("%d", statusCode)).Inc()
		} else if 400 <= statusCode && 499 >= statusCode {
			idHttpStatusCounter.WithLabelValues(fmt.Sprintf("%d", statusCode)).Inc()
		} else if 500 <= statusCode && 599 >= statusCode {
			idHttpStatusCounter.WithLabelValues(fmt.Sprintf("%d", statusCode)).Inc()
		}

		latency := time.Now().Sub(startTime)

		if latency > time.Minute {
			latency = latency.Truncate(time.Second)
		}

		idHttpLatency.Set(float64(latency.Milliseconds()))
	}
}

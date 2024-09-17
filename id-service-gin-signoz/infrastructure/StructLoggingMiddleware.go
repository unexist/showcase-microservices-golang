//
// @package Showcase-Microservices-Golang
//
// @file Structured Logging middleware
// @copyright 2024-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package infrastructure

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func DefaultStructuredLogger() gin.HandlerFunc {
	return StructLoggingMiddleware(&log.Logger)
}

func StructLoggingMiddleware(logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now() // Start timer
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		param := gin.LogFormatterParams{}

		param.TimeStamp = time.Now() // Stop timer
		param.Latency = param.TimeStamp.Sub(startTime)

		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		param.BodySize = c.Writer.Size()

		if "" != raw {
			path = path + "?" + raw
		}
		param.Path = path

		// Log using the params
		var logEvent *zerolog.Event

		if 500 <= c.Writer.Status() {
			logEvent = logger.Error()
		} else {
			logEvent = logger.Info()
		}

		logEvent.Str("client_id", param.ClientIP).
			Str("method", param.Method).
			Int("status_code", param.StatusCode).
			Int("body_size", param.BodySize).
			Str("path", param.Path).
			Str("latency", param.Latency.String()).
			Msg(param.ErrorMessage)
	}
}

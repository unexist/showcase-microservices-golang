//
// @package Showcase-Microservices-Golang
//
// @file Todo resource
// @copyright 2024-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package adapter

import (
	"net/http"

	"github.com/hashicorp/go-uuid"
	"github.com/prometheus/client_golang/prometheus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/gin-gonic/gin"

	"github.com/unexist/showcase-microservices-golang/docs"
)

var (
	idActionCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "id_action_counter",
			Help: "Total number of times an action was called",
		},
		[]string{"action"},
	)
)

func init() {
	prometheus.MustRegister(idActionCounter)
}

// @title OpenAPI for Todo showcase
// @version 1.0
// @description OpenAPI for Todo showcase

// @contact.name Christoph Kappel
// @contact.url https://unexist.dev
// @contact.email christoph@unexist.dev

// @license.name Apache 2.0
// @license.url https://www.apache.org/licenses/LICENSE-2.0

// @BasePath /id

type IdResource struct {
}

func NewIdResource() *IdResource {
	return &IdResource{}
}

// @Summary Get new id
// @Description Get new id
// @Produce json
// @Tags Id
// @Success 200 {string} string "New ID"
// @Failure 500 {string} string "Server error"
// @Router /id [get]
func (resource *IdResource) getId(context *gin.Context) {
	tracer := otel.GetTracerProvider().Tracer("id-resource")
	_, span := tracer.Start(context.Request.Context(), "get-id",
		trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	idActionCounter.WithLabelValues("getId").Inc()

	newUUID, err := uuid.GenerateUUID()

	if nil != err {
		context.JSON(http.StatusTeapot, gin.H{"error": "Cannot generate ID"})

		span.SetStatus(http.StatusTeapot, "UUID failed")
		span.RecordError(err)

		return
	}

	span.SetStatus(http.StatusCreated, "UUID created")
	span.SetAttributes(attribute.String("uuid", newUUID))

	context.JSON(http.StatusCreated, gin.H{"id": newUUID})
}

// Register REST routes on given engine
func (resource *IdResource) RegisterRoutes(engine *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/"

	todo := engine.Group("/id")
	{
		todo.GET("", resource.getId)
	}

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

//
// @package Showcase-Microservices-Golang
//
// @file Id service
// @copyright 2024-present Christoph Kappel <christoph@unexist.dev>
// @version $Id$
//
// This program can be distributed under the terms of the Apache License v2.0.
// See the file LICENSE for details.
//

package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"

	"github.com/unexist/showcase-microservices-golang/infrastructure/utils"
)

type IdService struct{}

type IdServiceReply struct {
	UUID string `json:"uuid"`
}

func NewIdService() *IdService {
	return &IdService{}
}

func (service *IdService) GetId(ctx context.Context) (string, error) {
	tracer := otel.GetTracerProvider().Tracer("todo-service")
	_, span := tracer.Start(ctx, "get-id")
	defer span.End()

	response, err := otelhttp.Get(ctx, fmt.Sprintf("%s/id",
		utils.GetEnvOrDefault("APP_ID_LISTEN_HOST_PORT", ":8081")))

	if err != nil {
		return "", err
	}

	jsonStr, _ := io.ReadAll(response.Body)

	var reply IdServiceReply

	err = json.Unmarshal([]byte(jsonStr), &reply)

	if err != nil {
		return "", err
	}

	return reply.UUID, nil
}

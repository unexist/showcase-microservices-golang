.DEFAULT_GOAL := build
.ONESHELL:
.PHONY: test

PG_USER := postgres
PG_PASS := postgres

define JSON_TODO
curl -X 'POST' \
  'http://localhost:8080/todo' \
  -H 'accept: */*' \
  -H 'Content-Type: application/json' \
  -d '{
  "description": "string",
  "done": true,
  "title": "string"
}'
endef
export JSON_TODO

# Tools
todo:
	@echo $$JSON_TODO | bash

list:
	@curl -X 'GET' 'http://localhost:8080/todo' -H 'accept: */*' | jq .

open-swagger:
	open http://localhost:8080/swagger/index.html

open-prometheus:
	open http://localhost:9090

open-signoz:
	open http://localhost:3301

# Build
build-gin-signoz:
	@$(SHELL) -c "cd todo-service-gin-signoz; GO111MODULE=on; go mod download; go build -o $(BINARY)"

# Run
run-gin-signoz:
	@$(SHELL) -c "cd todo-service-gin-signoz; APP_DB_USERNAME=$(PG_USER) APP_DB_PASSWORD=$(PG_PASS) APP_DB_NAME=postgres ./$(BINARY)"

# Test
hurl:
	@hurl --color --test todo.hurl

slumber:
	@slumber ./slumber.yml

# Modules
ifneq (,$(findstring gin-tracing,$(MAKECMDGOALS)))
-include todo-service-gin-tracing/Makefile
endif

ifneq (,$(findstring mux,$(MAKECMDGOALS)))
-include todo-service-mux/Makefile
endif

install:
	go install braces.dev/errtrace/cmd/errtrace@latest
	go install golang.org/x/tools/cmd/deadcode@latest
	go install github.com/swaggo/swag/cmd/swag@latest

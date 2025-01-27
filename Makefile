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

# Helper
todo:
	@echo $$JSON_TODO | bash

list:
	@curl -X 'GET' 'http://localhost:8080/todo' -H 'accept: */*' | jq .

open-swagger:
	open http://localhost:8080/swagger/index.html

open-prometheus:
	open http://localhost:9090

# Test
hurl-todo:
	@hurl --color --test hurl/todo.hurl

hurl-user:
	@hurl --color --test hurl/user.hurl

slumber:
	@slumber ./slumber.yml

# Modules
ifneq (,$(findstring gin,$(MAKECMDGOALS)))
-include todo-service-gin/Makefile
endif

ifneq (,$(findstring metrics,$(MAKECMDGOALS)))
-include todo-service-gin-metrics/Makefile
endif

ifneq (,$(findstring middleware,$(MAKECMDGOALS)))
-include todo-service-gin-middleware/Makefile
endif

ifneq (,$(findstring tracing,$(MAKECMDGOALS)))
-include todo-service-gin-tracing/Makefile
endif

ifneq (,$(findstring transactions,$(MAKECMDGOALS)))
-include todo-service-gin-transactions/Makefile
endif

ifneq (,$(findstring transactor,$(MAKECMDGOALS)))
-include todo-service-gin-transactor/Makefile
endif

ifneq (,$(findstring mux,$(MAKECMDGOALS)))
-include todo-service-mux/Makefile
endif

ifneq (,$(findstring infra,$(MAKECMDGOALS)))
-include infrastructure/Makefile
endif

install:
	go install braces.dev/errtrace/cmd/errtrace@latest
	go install golang.org/x/tools/cmd/deadcode@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/kisielk/godepgraph@latest


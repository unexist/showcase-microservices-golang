.DEFAULT_GOAL := build
.ONESHELL:

BINARY := todo-service.bin

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

swagger:
	@$(SHELL) -c "cd todo-service-gin; swag init"

open-swagger:
	open http://localhost:8080/swagger/index.html

kat-test:
	@kcat -t todo_created -b localhost:9092 -P

kat-listen:
	@kcat -t todo_created -b localhost:9092 -C

# Build
build-mux:
	@$(SHELL) -c  "cd todo-service-mux; GO111MODULE=on GOFLAGS=-mod=vendor; go mod download; go build -o $(BINARY)"

build-gin:
	@$(SHELL) -c "cd todo-service-gin; GO111MODULE=on; go mod download; go build -o $(BINARY)"

build-gin-tracing:
	@$(SHELL) -c "cd todo-service-gin-tracing; GO111MODULE=on; go mod download; go build -o $(BINARY)"

# Analysis
vet-mux:
	@$(SHELL) -c "cd todo-service-mux; go vet"

vet-gin:
	@$(SHELL) -c "cd todo-service-gin; go vet"

wire-gin:
	@$(SHELL) -c  "cd todo-service-gin/test; wire"

# Run
run-mux:
	@$(SHELL) -c  "cd todo-service-mux; APP_DB_USERNAME=$(PG_USER) APP_DB_PASSWORD=$(PG_PASS) APP_DB_NAME=postgres ./$(BINARY)"

run-gin:
	@$(SHELL) -c "cd todo-service-gin; APP_DB_USERNAME=$(PG_USER) APP_DB_PASSWORD=$(PG_PASS) APP_DB_NAME=postgres ./$(BINARY)"

run-gin-jaeger:
	@$(SHELL) -c "cd todo-service-gin-tracing; APP_DB_USERNAME=$(PG_USER) APP_DB_PASSWORD=$(PG_PASS) APP_DB_NAME=postgres TRACER=jaeger ./$(BINARY)"

run-gin-zipkin:
	@$(SHELL) -c "cd todo-service-gin-tracing; APP_DB_USERNAME=$(PG_USER) APP_DB_PASSWORD=$(PG_PASS) APP_DB_NAME=postgres TRACER=zipkin ./$(BINARY)"

# Tests
test-mux:
	@$(SHELL) -c "cd todo-service-mux; go test -v"

test-fake-gin:
	@$(SHELL) -c "cd todo-service-gin; go test -v -tags=fake ./test"

test-cucumber-gin:
	@$(SHELL) -c "cd todo-service-gin; go test -v -tags=cucumber ./test"

test-gorm-gin:
	@$(SHELL) -c "cd todo-service-gin; TEST_DB_USERNAME=$(PG_USER) TEST_DB_PASSWORD=$(PG_PASS) TEST_DB_NAME=postgres go test -v -tags=gorm ./test"

test-sqlx-gin:
	@$(SHELL) -c "cd todo-service-gin; TEST_DB_USERNAME=$(PG_USER) TEST_DB_PASSWORD=$(PG_PASS) TEST_DB_NAME=postgres go test -v -tags=sqlx ./test"

test-arch-gin:
	@$(SHELL) -c "cd todo-service-gin; go test -v -tags=arch ./test"

# Helper
clean:
	rm -rf todo-service-mux/$(BINARY)
	rm -rf todo-service-gin/$(BINARY)
	rm -rf todo-service-gin-tracing/$(BINARY)

install:
	go install braces.dev/errtrace/cmd/errtrace@latest
	go install golang.org/x/tools/cmd/deadcode@latest
	go install github.com/swaggo/swag/cmd/swag@latest

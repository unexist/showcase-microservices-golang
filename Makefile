.DEFAULT_GOAL := build
.ONESHELL:

BINARY := todo-service.bin
PODNAME := showcase
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

psql:
	@PGPASSWORD=$(PG_PASS) psql -h 127.0.0.1 -U $(PG_USER)

schema:
	@PGPASSWORD=$(PG_PASS) psql -h 127.0.0.1 -U $(PG_USER) -f ./schema.sql

swagger:
	@$(SHELL) -c "cd todo-service-gin; swag init"

open-swagger:
	open http://localhost:8080/swagger/index.html

kat-test:
	@kcat -t todo_created -b localhost:9092 -P

kat-listen:
	@kcat -t todo_created -b localhost:9092 -C

# Podman
pd-machine-init:
	podman machine init --memory=8192 --cpus=2 --disk-size=20

pd-machine-start:
	podman machine start

pd-machine-stop:
	podman machine stop

pd-machine-rm:
	@podman machine rm

pd-machine-recreate: pd-machine-rm pd-machine-init pd-machine-start

pd-pod-create:
	@podman pod create -n $(PODNAME) --network bridge \
		-p 5432:5432 \
		-p 9092:9092 \
		-p 9411:9411 \
		-p 4318:4318

pd-pod-rm:
	podman pod rm -f $(PODNAME)

pd-pod-recreate: pd-pod-rm pd-pod-create

pd-postgres:
	@podman run -dit --name postgres --pod=$(PODNAME) \
		-e POSTGRES_USER=$(PG_USER) \
		-e POSTGRES_PASSWORD=$(PG_PASS) \
		postgres:latest

pd-zipkin:
	@podman run -dit --name zipkin --pod=$(PODNAME) openzipkin/zipkin

pd-jaeger:
	@podman run -dit --name jaeger --pod=$(PODNAME) jaegertracing/all-in-one

pd-redpanda:
	@podman run -dit --name redpanda --pod=$(PODNAME) vectorized/redpanda

# Build
build-mux:
	@$(SHELL) -c  "cd todo-service-mux; GO111MODULE=on GOFLAGS=-mod=vendor; go mod download; go build -o $(BINARY)"

vet-mux:
	@$(SHELL) -c "cd todo-service-mux; go vet"

run-mux:
	#source env-sample
	@$(SHELL) -c  "cd todo-service-mux; APP_DB_USERNAME=$(PG_USER) APP_DB_PASSWORD=$(PG_PASS) APP_DB_NAME=postgres ./$(BINARY)"

test-mux:
	#source env-test
	@$(SHELL) -c "cd todo-service-mux; go test -v"

build-gin:
	@$(SHELL) -c "cd todo-service-gin; GO111MODULE=on; go mod download; go build -o $(BINARY)"

vet-gin:
	@$(SHELL) -c "cd todo-service-gin; go vet"

wire-gin:
	@$(SHELL) -c  "cd todo-service-gin/test; wire"

run-gin:
	#source env-sample
	@$(SHELL) -c "cd todo-service-gin; APP_DB_USERNAME=$(PG_USER) APP_DB_PASSWORD=$(PG_PASS) APP_DB_NAME=postgres ./$(BINARY)"

run-trace-gin:
	#source env-sample
	@$(SHELL) -c "cd todo-service-gin; APP_DB_USERNAME=$(PG_USER) APP_DB_PASSWORD=$(PG_PASS) APP_DB_NAME=postgres ./$(BINARY)"

test-fake-gin:
	#source env-test
	@$(SHELL) -c "cd todo-service-gin; go test -v -tags=fake ./test"

test-cucumber-gin:
	#source env-test
	@$(SHELL) -c "cd todo-service-gin; go test -v -tags=cucumber ./test"

test-gorm-gin:
	#source env-test
	@$(SHELL) -c "cd todo-service-gin; TEST_DB_USERNAME=$(PG_USER) TEST_DB_PASSWORD=$(PG_PASS) TEST_DB_NAME=postgres go test -v -tags=gorm ./test"

test-sqlx-gin:
	#source env-test
	@$(SHELL) -c "cd todo-service-gin; TEST_DB_USERNAME=$(PG_USER) TEST_DB_PASSWORD=$(PG_PASS) TEST_DB_NAME=postgres go test -v -tags=sqlx ./test"

test-arch-gin:
	@$(SHELL) -c "cd todo-service-gin; go test -v -tags=arch ./test"

clear:
	rm -rf todo-service-mux/$(BINARY)
	rm -rf todo-service-gin/$(BINARY)

install:
	go install braces.dev/errtrace/cmd/errtrace@latest
	go install golang.org/x/tools/cmd/deadcode@latest

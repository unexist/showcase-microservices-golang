.DEFAULT_GOAL := build

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
	@PGPASSWORD=$(PG_PASS) psql -h localhost -U $(PG_USER)

# Podman
pd-machine-init:
	podman machine init --memory=8192 --cpus=2 --disk-size=20

pd-machine-start:
	podman machine start

pd-machine-rm:
	@podman machine rm

pd-machine-recreate: pd-machine-rm pd-machine-init pd-machine-start

pd-pod-create:
	@podman pod create -n $(PODNAME) --network bridge \
		-p 5432:5432

pd-pod-rm:
	podman pod rm -f $(PODNAME)

pd-pod-recreate: pd-pod-rm pd-pod-create

pd-postgres:
	@podman run -dit --name postgres --pod=$(PODNAME) \
		-e POSTGRES_USER=$(PG_USER) \
		-e POSTGRES_PASSWORD=$(PG_PASS) \
		postgres:latest

# Build
build-mux:
	cd todo-service-mux
	export GO111MODULE=on
	export GOFLAGS=-mod=vendor
	go mod download
	go build -o $(BINARY)

run-mux:
	cd todo-service-mux
	./$(BINARY)

test-mux:
	cd todo-service-mux
	go test -v

build-gin:
	cd todo-service-gin
	export GO111MODULE=on
	export GOFLAGS=-mod=vendor
	go mod download
	go build -o $(BINARY)

run-gin:
	cd todo-service-gin
	./$(BINARY)

test-gin:
	cd todo-service-gin
	go test -v

clear:
	rm -rf todo-service-mux/$(BINARY)
	rm -rf todo-service-gin/$(BINARY)

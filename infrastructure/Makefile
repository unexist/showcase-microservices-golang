.DEFAULT_GOAL := build
.ONESHELL:

PODNAME := showcase
PG_USER := postgres
PG_PASS := postgres

psql:
	@PGPASSWORD=$(PG_PASS) psql -h 127.0.0.1 -U $(PG_USER)

psql-schema:
	@PGPASSWORD=$(PG_PASS) psql -h 127.0.0.1 -U $(PG_USER) -f ./../schema.sql

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
		-p 14268:14268

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
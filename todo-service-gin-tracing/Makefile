BINARY := todo-service.bin

# Dummy for parent Makefile
tracing:
	$(NOECHO) $(NOOP)

# Swagger
swagger:
	@$(SHELL) -c "cd todo-service-gin-tracing; swag init"

# Build
build:
	@$(SHELL) -c "cd todo-service-gin-tracing; GO111MODULE=on; go mod download; go build -o $(BINARY)"

# Run
run-jaeger:
	@$(SHELL) -c "cd todo-service-gin-tracing; APP_DB_USERNAME=$(PG_USER) APP_DB_PASSWORD=$(PG_PASS) APP_DB_NAME=postgres TRACER=jaeger ./$(BINARY)"

run-zipkin:
	@$(SHELL) -c "cd todo-service-gin-tracing; APP_DB_USERNAME=$(PG_USER) APP_DB_PASSWORD=$(PG_PASS) APP_DB_NAME=postgres TRACER=zipkin ./$(BINARY)"

clean:
	rm -rf todo-service-gin-tracing/$(BINARY)

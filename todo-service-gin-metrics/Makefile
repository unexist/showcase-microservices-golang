BINARY := todo-service.bin

# Dummy for parent Makefile
metrics:
	$(NOECHO) $(NOOP)

# Swagger
swagger:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin-metrics; swag init"

# Build
build:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin-metrics; GO111MODULE=on; go mod download; go build -o $(BINARY)"

# Run
run:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin-metrics; APP_DB_USERNAME=$(PG_USER) APP_DB_PASSWORD=$(PG_PASS) APP_DB_NAME=postgres ./$(BINARY)"

# Helper
clean:
	rm -rf $(CURDIR)/todo-service-gin-metrics/$(BINARY)

BINARY := todo-service.bin

# Dummy for parent Makefile
todosignoz:
	$(NOECHO) $(NOOP)

# Build
build:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin-signoz; GO111MODULE=on; go mod download; go build -o $(BINARY)"

# Run
run:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin-signoz; APP_DB_USERNAME=$(PG_USER) APP_DB_PASSWORD=$(PG_PASS) APP_DB_NAME=postgres ./$(BINARY)"

# Helper
clean:
	rm -rf $(CURDIR)/todo-service-gin-signoz/$(BINARY)

# Container
container:
	@docker build -t todo-todo-service -f $(CURDIR)/todo-service-gin-signoz/Containerfile $(CURDIR)/todo-service-gin-signoz
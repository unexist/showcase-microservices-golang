BINARY := todo-service.bin

# Dummy for parent Makefile
todosignoz:
	$(NOECHO) $(NOOP)

# Swagger
swagger:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin-signoz; swag init"

# Build
build:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin-signoz; GO111MODULE=on; go mod download; go build -o $(BINARY)"

# Run
run: build
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin-signoz; APP_DB_USERNAME=$(PG_USER) APP_DB_PASSWORD=$(PG_PASS) APP_DB_NAME=postgres ./$(BINARY)"

# Analysis
vet:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin-signoz; go vet"

# Tests
test-fake:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin-signoz; go test -v -tags=fake ./test"

test-gorm:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin-signoz; TEST_DB_USERNAME=$(PG_USER) TEST_DB_PASSWORD=$(PG_PASS) TEST_DB_NAME=postgres go test -v -tags=gorm ./test"

# Helper
clean:
	rm -rf $(CURDIR)/todo-service-gin-signoz/$(BINARY)

# Container
container:
	@docker build -t todo-todo-service -f $(CURDIR)/todo-service-gin-signoz/Containerfile $(CURDIR)/todo-service-gin-signoz
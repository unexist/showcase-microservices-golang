BINARY := todo-service.bin

# Dummy for parent Makefile
gin:
	$(NOECHO) $(NOOP)

# Swagger
swagger:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin; swag init"

# Build
build:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin; GO111MODULE=on; go mod download; go build -o $(BINARY)"

# Analysis
vet:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin; go vet"

wire:
	@$(SHELL) -c  "cd $(CURDIR)/todo-service-gin/test; wire"

# Run
run:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin; APP_DB_USERNAME=$(PG_USER) APP_DB_PASSWORD=$(PG_PASS) APP_DB_NAME=postgres ./$(BINARY)"

# Tests
test-fake:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin; go test -v -tags=fake ./test"

test-cucumber:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin; go test -v -tags=cucumber ./test"

test-gorm:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin; TEST_DB_USERNAME=$(PG_USER) TEST_DB_PASSWORD=$(PG_PASS) TEST_DB_NAME=postgres go test -v -tags=gorm ./test"

test-sqlx:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin; TEST_DB_USERNAME=$(PG_USER) TEST_DB_PASSWORD=$(PG_PASS) TEST_DB_NAME=postgres go test -v -tags=sqlx ./test"

test-arch:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin; go test -v -tags=arch ./test"

# Helper
clean:
	rm -rf $(CURDIR)/todo-service-gin/$(BINARY)

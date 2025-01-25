BINARY := todo-service.bin

# Dummy for parent Makefile
transactor:
	$(NOECHO) $(NOOP)

# Swagger
swagger:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin-transactor; swag init"

# Build
build:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin-transactor; GO111MODULE=on; go mod download; go build -o $(BINARY)"

# Analysis
vet:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin-transactor; go vet"

# Run
run: build
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin-transactor; ./$(BINARY)"

# Tests
test:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin-transactor; go test -count=1 -v ./test"

# Helper
clear:
	rm -rf $(CURDIR)/todo-service-gin-transactor/$(BINARY)

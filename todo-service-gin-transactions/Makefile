BINARY := todo-service.bin

# Dummy for parent Makefile
transactions:
	$(NOECHO) $(NOOP)

# Swagger
swagger:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin-transactions; swag init"

# Build
build:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin-transactions; GO111MODULE=on; go mod download; go build -o $(BINARY)"

# Analysis
vet:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin-transactions; go vet"

# Run
run: build
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin-transactions; ./$(BINARY)"

# Tests
test:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin-transactions; go test -count=1 -v ./test"

# Helper
clear:
	rm -rf $(CURDIR)/todo-service-gin-transactions/$(BINARY)

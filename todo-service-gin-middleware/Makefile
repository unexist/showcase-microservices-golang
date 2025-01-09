BINARY := todo-service.bin

# Dummy for parent Makefile
middleware:
	$(NOECHO) $(NOOP)

# Swagger
swagger:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin-middleware; swag init"

# Build
build:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin-middleware; GO111MODULE=on; go mod download; go build -o $(BINARY)"

# Analysis
vet:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin-middleware; go vet"

# Run
run: build
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin-middleware; ./$(BINARY)"

# Tests
test:
	@$(SHELL) -c "cd $(CURDIR)/todo-service-gin-middleware; go test -count=1 -v ./test"

# Helper
clear:
	rm -rf $(CURDIR)/todo-service-gin-middleware/$(BINARY)

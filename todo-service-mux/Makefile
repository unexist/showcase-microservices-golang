BINARY := todo-service.bin

# Dummy for parent Makefile
mux:
	$(NOECHO) $(NOOP)

# Build
build:
	@$(SHELL) -c  "cd todo-service-mux; GO111MODULE=on GOFLAGS=-mod=vendor; go mod download; go build -o $(BINARY)"

# Analysis
vet:
	@$(SHELL) -c "cd todo-service-mux; go vet"

# Run
run:
	@$(SHELL) -c  "cd todo-service-mux; APP_DB_USERNAME=$(PG_USER) APP_DB_PASSWORD=$(PG_PASS) APP_DB_NAME=postgres ./$(BINARY)"

# Tests
test:
	@$(SHELL) -c "cd todo-service-mux; go test -v"

# Helper
clean:
	rm -rf todo-service-mux/$(BINARY)

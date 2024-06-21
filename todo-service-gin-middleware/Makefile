.DEFAULT_GOAL := build-gin
.ONESHELL:

EXE_GIN := todo-service.bin

swagger:
	@$(SHELL) -c "swag init"

# Build
build:
	@$(SHELL) -c "GO111MODULE=on; go mod download; go build -o $(EXE_GIN)"

# Analysis
vet:
	@$(SHELL) -c "go vet"

# Run
run: build
	@$(SHELL) -c "./$(EXE_GIN)"

# Tests
test:
	@$(SHELL) -c "go test -count=1 -v ./test"

# Helper
clear:
	rm -rf ./$(EXE_GIN)

install:
	go install github.com/swaggo/swag/cmd/swag@latest
BINARY := id-service.bin

# Dummy for parent Makefile
idsignoz:
	$(NOECHO) $(NOOP)

# Swagger
swagger:
	@$(SHELL) -c "cd $(CURDIR)/id-service-gin-signoz; swag init"

# Build
build:
	@$(SHELL) -c "cd $(CURDIR)/id-service-gin-signoz; GO111MODULE=on; go mod download; go build -o $(BINARY)"

# Run
run: build
	@$(SHELL) -c "cd $(CURDIR)/id-service-gin-signoz; ./$(BINARY)"

# Analysis
vet:
	@$(SHELL) -c "cd $(CURDIR)/id-service-gin-signoz; go vet"

# Helper
clean:
	rm -rf $(CURDIR)/todo-service-gin-signoz/$(BINARY)

# Container
container:
	@docker build -t todo-id-service -f $(CURDIR)/id-service-gin-signoz/Containerfile $(CURDIR)/id-service-gin-signoz

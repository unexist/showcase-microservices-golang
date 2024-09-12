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
run:
	@$(SHELL) -c "cd $(CURDIR)/id-service-gin-signoz; APP_DB_USERNAME=$(PG_USER) APP_DB_PASSWORD=$(PG_PASS) APP_DB_NAME=postgres ./$(BINARY)"

# Container
container:
	@docker build -t todo-id-service -f $(CURDIR)/id-service-gin-signoz/Containerfile $(CURDIR)/id-service-gin-signoz

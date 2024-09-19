BIN_DIR = bin
BIN_FILE = api-server
CMD_DIR = cmd
CMD_FILE = main.go
GOOSE_SETTINGS = GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=postgres dbname=todo_api sslmode=disable" GOOSE_MIGRATION_DIR="repo/migrations"

all: build

run: build
	@./$(BIN_DIR)/$(BIN_FILE)

build:
	@echo "> start building the server..."
	@go build -o $(BIN_DIR)/$(BIN_FILE) $(CMD_DIR)/$(CMD_FILE)
	@echo "> finished building"

clean:
	@echo "> cleaning bin dir..."
	@rm -rf $(BIN_DIR)
	@echo "> finished cleaning"

up:
	$(GOOSE_SETTINGS) goose up

down:
	$(GOOSE_SETTINGS) goose down

reset:
	$(GOOSE_SETTINGS) goose reset

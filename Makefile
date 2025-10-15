APP_NAME := go-gin-todo-api
APP_VERSION := 0.1.0

BIN_DIR := bin

.PHONY: run build clean

# default: run in watch mode
run:
	@echo "Running $(APP_NAME) in watch mode..."
	air

build:
	@echo "Building $(APP_NAME) version $(APP_VERSION)..."
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) .

clean:
	@rm -rf $(BIN_DIR)
	@echo "Cleaned up..."

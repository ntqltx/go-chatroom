BINARY_DIR = bin
SERVER_BIN = $(BINARY_DIR)/server
CLIENT_BIN = $(BINARY_DIR)/client

.PHONY: all build server client clean

all: build

build: $(SERVER_BIN) $(CLIENT_BIN)
	@echo "Binaries built successfully"

$(SERVER_BIN):
	@mkdir -p $(BINARY_DIR)
	@go build -o $@ ./server

$(CLIENT_BIN):
	@mkdir -p $(BINARY_DIR)
	@go build -o $@ ./client

server:
	@-go run ./server; exit 0

client:
	@-go run ./client; exit 0

clean:
	@rm -rf $(BINARY_DIR)
	@echo "Cleaned binaries"

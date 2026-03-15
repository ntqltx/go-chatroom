BINARY_DIR = bin
SERVER_BIN = $(BINARY_DIR)/server
CLIENT_BIN = $(BINARY_DIR)/client

.PHONY: all build server client clean

build: $(SERVER_BIN) $(CLIENT_BIN)
	@echo "Binaries built successfully"

$(SERVER_BIN):
	@mkdir -p $(BINARY_DIR)
	@go build -o $@ ./server

$(CLIENT_BIN):
	@mkdir -p $(BINARY_DIR)
	@go build -o $@ ./client

server:
	@if echo "$(ARGS)" | grep -qE "\-\-verbose|\-\-log"; then \
		stty -echoctl; \
		go run ./server $(ARGS); \
		stty echoctl; \
	else \
		go run ./server $(ARGS) & \
		sleep 0.3; \
	fi; exit 0

client:
	@-go run ./client; exit 0

clean:
	@rm -rf $(BINARY_DIR)
	@echo "Cleaned binaries"

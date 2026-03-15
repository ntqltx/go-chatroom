BINARY_DIR = bin
SERVER_BIN = $(BINARY_DIR)/server
CLIENT_BIN = $(BINARY_DIR)/client

.PHONY: all build kill clean

build: clean $(SERVER_BIN) $(CLIENT_BIN)

$(SERVER_BIN):
	@mkdir -p $(BINARY_DIR)
	@go build -o $@ ./server
	@chmod +x $(SERVER_BIN)

$(CLIENT_BIN):
	@mkdir -p $(BINARY_DIR)
	@go build -o $@ ./client
	@chmod +x $(CLIENT_BIN)

kill:
	@echo ""; \
	read -p "Port (default 8080): " port; \
	port=$${port:-8080}; \
	pid=$$(lsof -t -i :$$port); \
	if [ -z "$$pid" ]; then \
		echo "Nothing running on port $$port"; \
	else \
		kill -15 $$pid && echo "Server on port $$port stopped" \
		|| echo "Failed to stop server"; \
	fi

clean:
	@rm -rf $(BINARY_DIR)

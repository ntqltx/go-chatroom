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
		cmd=$$(ps -o command= -p $$pid | awk '{print $$1}' | xargs basename 2>/dev/null); \
		if echo "$$cmd" | grep -qE "^server$$"; then \
			kill -15 -- $$pid && echo "Server successfully stopped" || echo "Failed to stop server"; \
		else \
			echo "Process on port $$port is '$$cmd', not a server"; \
		fi \
	fi

clean:
	@rm -rf $(BINARY_DIR)

.PHONY: build proto run clean test clean-db run-db-mode

APP_NAME = account-cred-manager-go-grpc
APP_DB_PATH = ./badger.db
SERVER_PATH = ./cmd/server
PROTO_DIR = ./api/proto
PROTO_OUTPUT_DIR = ./api/proto/v1
BUILD_DIR = ./build

build: proto
	@echo "Building $(APP_NAME)"
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(SERVER_PATH)

proto:
	@echo "Generating proto code"
	@mkdir -p ${PROTO_OUTPUT_DIR}
	protoc --proto_path=$(PROTO_DIR)/v1 \
		   --go_out=$(PROTO_OUTPUT_DIR) --go_opt=paths=source_relative \
		   --go-grpc_out=$(PROTO_OUTPUT_DIR) --go-grpc_opt=paths=source_relative \
		   $(PROTO_DIR)/v1/*.proto

run: build
	@echo "Running $(APP_NAME) in Memory mode"
	@GRPC_PORT="50051" $(BUILD_DIR)/$(APP_NAME)

run-db-mode: build
	@echo "Running $(APP_NAME) in DB mode"
	@GRPC_PORT="50051" STORAGE_MODE=DB $(BUILD_DIR)/$(APP_NAME)

test:
	@echo "Running test"
	go clean --testcache
	go test ./test/ -v

clean-db:
	@echo "Cleaning DB"
	rm -rf $(APP_DB_PATH)

clean:
	@echo "Cleaning up"
	rm -rf $(BUILD_DIR)
	rm $(PROTO_OUTPUT_DIR)/*.pb.go
	go clean ./...
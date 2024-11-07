.PHONY: all build clean test lint run help

# 项目信息
PROJECT_NAME := lingjian
MAIN_SERVER := cmd/server/main.go
MAIN_WORKER := cmd/worker/main.go
BUILD_DIR := build
BINARY_SERVER := $(BUILD_DIR)/server
BINARY_WORKER := $(BUILD_DIR)/worker

# Go相关配置
GO := go
GOFLAGS := -v
LDFLAGS := -s -w

# 数据库配置
DB_USER := root
DB_PASS := MtNNJasQv5GptQz
DB_NAME := lingjian
DB_TEST_NAME := lingjian_test
DB_HOST := localhost
DB_PORT := 3306

# 默认目标
all: lint test build

# 构建
build: $(BINARY_SERVER) $(BINARY_WORKER)

$(BINARY_SERVER):
	@echo "Building server..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_SERVER) $(MAIN_SERVER)

$(BINARY_WORKER):
	@echo "Building worker..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_WORKER) $(MAIN_WORKER)

# 清理
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean -i ./...

# 测试相关命令
test: init-test-db test-unit test-integration

# 初始化测试数据库
init-test-db:
	@echo "Initializing test database..."
	@mysql -h $(DB_HOST) -P $(DB_PORT) -u $(DB_USER) -p$(DB_PASS) < internal/test/init_test_db.sql

# 单元测试
test-unit:
	@echo "Running unit tests..."
	@go test -v -short ./...

# 集成测试
test-integration:
	@echo "Running integration tests..."
	@go test -v -run Integration ./...

# 测试覆盖率
test-coverage:
	@echo "Running tests with coverage..."
	@mkdir -p $(BUILD_DIR)
	@go test -v -coverprofile=$(BUILD_DIR)/coverage.out ./...
	@go tool cover -html=$(BUILD_DIR)/coverage.out -o $(BUILD_DIR)/coverage.html
	@echo "Coverage report generated at $(BUILD_DIR)/coverage.html"

# 性能测试
test-bench:
	@echo "Running benchmark tests..."
	@go test -v -bench=. -benchmem ./...

# 代码检查
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint is not installed. Installing..."; \
		go install github.com/golangci/golangci/cmd/golangci-lint@latest; \
		golangci-lint run ./...; \
	fi

# 运行服务器
run-server: build
	@echo "Running server..."
	@$(BINARY_SERVER)

# 运行worker
run-worker: build
	@echo "Running worker..."
	@$(BINARY_WORKER)

# 初始化数据库
init-db:
	@echo "Initializing database..."
	@mysql -h $(DB_HOST) -P $(DB_PORT) -u $(DB_USER) -p$(DB_PASS) -e "CREATE DATABASE IF NOT EXISTS $(DB_NAME) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;"
	@mysql -h $(DB_HOST) -P $(DB_PORT) -u $(DB_USER) -p$(DB_PASS) $(DB_NAME) < internal/model/schema.sql
	@mysql -h $(DB_HOST) -P $(DB_PORT) -u $(DB_USER) -p$(DB_PASS) $(DB_NAME) < internal/model/schema_application.sql
	@mysql -h $(DB_HOST) -P $(DB_PORT) -u $(DB_USER) -p$(DB_PASS) $(DB_NAME) < internal/model/schema_config.sql
	@mysql -h $(DB_HOST) -P $(DB_PORT) -u $(DB_USER) -p$(DB_PASS) $(DB_NAME) < internal/model/schema_task.sql

# 生成API文档
docs:
	@echo "Generating API documentation..."
	@if command -v swag >/dev/null; then \
		swag init -g cmd/server/main.go; \
	else \
		echo "swag is not installed. Installing..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
		$(shell go env GOPATH)/bin/swag init -g cmd/server/main.go; \
	fi
	@echo "API documentation generated. Visit http://localhost:8081/docs when server is running."

# 依赖管理
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# 格式化代码
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# 检查代码
vet:
	@echo "Vetting code..."
	@go vet ./...

# 开发模式运行（支持热重载）
dev-server: docs
	@if command -v air >/dev/null; then \
		air -c .air.toml; \
	else \
		echo "air is not installed. Installing..."; \
		GOBIN=$(shell go env GOPATH)/bin go install github.com/air-verse/air@latest; \
		export PATH=$(shell go env GOPATH)/bin:$$PATH; \
		air -c .air.toml; \
	fi

# 帮助信息
help:
	@echo "Make commands:"
	@echo "  all              - Run lint, test, and build"
	@echo "  build            - Build server and worker binaries"
	@echo "  clean            - Remove build artifacts"
	@echo "  test             - Run all tests"
	@echo "  test-unit        - Run unit tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  test-coverage    - Run tests with coverage report"
	@echo "  test-bench       - Run benchmark tests"
	@echo "  lint             - Run linter"
	@echo "  run-server       - Build and run server"
	@echo "  run-worker       - Build and run worker"
	@echo "  init-db          - Initialize database"
	@echo "  init-test-db     - Initialize test database"
	@echo "  docs             - Generate API documentation"
	@echo "  deps             - Download and tidy dependencies"
	@echo "  fmt              - Format code"
	@echo "  vet              - Check code for common errors"
	@echo "  dev-server       - Run server in development mode with hot reload"

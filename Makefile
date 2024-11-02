# Go参数
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# 项目参数
BINARY_NAME=lingjian
MAIN_FILE=cmd/server/main.go
CONFIG_FILE=config/config.yaml
CONFIG_EXAMPLE=config/config.yaml.example

# 数据库脚本
SCHEMA_FILE=scripts/schema.sql
INIT_DATA_FILE=scripts/init_data.sql

.PHONY: all build run clean test deps setup help

all: clean build

# 构建项目
build:
	@echo "构建项目..."
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_FILE)

# 运行项目
run:
	@echo "运行项目..."
	$(GORUN) $(MAIN_FILE)

# 清理构建产物
clean:
	@echo "清理构建产物..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe

# 运行测试
test:
	@echo "运行测试..."
	$(GOTEST) -v ./...

# 更新依赖
deps:
	@echo "更新依赖..."
	$(GOGET) -u ./...
	$(GOMOD) tidy

# 初始化开发环境
setup:
	@echo "初始化开发环境..."
	@if [ ! -f $(CONFIG_FILE) ]; then \
		cp $(CONFIG_EXAMPLE) $(CONFIG_FILE); \
		echo "创建配置文件 $(CONFIG_FILE)"; \
	fi
	@mkdir -p logs uploads
	@echo "创建必要的目录"
	@chmod +x scripts/setup_dev.sh
	@./scripts/setup_dev.sh

# 生成API文档
swagger:
	@echo "生成Swagger文档..."
	swag init -g $(MAIN_FILE) -o api/docs

# 帮助信息
help:
	@echo "可用的命令："
	@echo "  make build    - 构建项目"
	@echo "  make run      - 运行项目"
	@echo "  make clean    - 清理构建产物"
	@echo "  make test     - 运行测试"
	@echo "  make deps     - 更新依赖"
	@echo "  make setup    - 初始化开发环境"
	@echo "  make swagger  - 生成API文档"
	@echo "  make help     - 显示帮助信息"

# 默认目标
.DEFAULT_GOAL := help

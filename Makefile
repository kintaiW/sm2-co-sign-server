# Makefile for sm2-co-sign-server

# 项目名称
PROJECT_NAME := sm2-co-sign-server

# 构建输出目录
BUILD_DIR := ./bin

# 主入口文件
MAIN_FILE := ./cmd/server/main.go

# Go 命令
GO := go

# 构建标志
BUILD_FLAGS := -ldflags="-w -s"

# 目标
all: build

# 构建
build: 
	@echo "Building $(PROJECT_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@$(GO) build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME) $(MAIN_FILE)
	@echo "Build completed: $(BUILD_DIR)/$(PROJECT_NAME)"

# 运行
run: build
	@echo "Running $(PROJECT_NAME)..."
	@$(BUILD_DIR)/$(PROJECT_NAME)

# 测试
test:
	@echo "Running tests..."
	@$(GO) test ./...

# 清理
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@$(GO) clean
	@echo "Clean completed"

# 依赖
deps:
	@echo "Downloading dependencies..."
	@$(GO) mod tidy
	@echo "Dependencies downloaded"

# 生成 Docker 镜像
docker:
	@echo "Building Docker image..."
	@docker build -t $(PROJECT_NAME) .
	@echo "Docker image built"

# 帮助
help:
	@echo "Makefile for $(PROJECT_NAME)"
	@echo ""
	@echo "Targets:"
	@echo "  all       - Build the project"
	@echo "  build     - Build the binary"
	@echo "  run       - Build and run the binary"
	@echo "  test      - Run tests"
	@echo "  clean     - Clean build artifacts"
	@echo "  deps      - Download dependencies"
	@echo "  docker    - Build Docker image"
	@echo "  help      - Show this help"

.PHONY: all build run test clean deps docker help

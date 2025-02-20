# 项目配置
BIN_NAME := gamelife
BIN_DIR := bin
SRC_DIR := ./...

# 构建参数
BUILD_FLAGS := -v
LD_FLAGS := -s -w

.PHONY: all build run clean test deps install

all: build

build:
	@echo "Building $(BIN_NAME)..."
	@mkdir -p $(BIN_DIR)
	go build $(BUILD_FLAGS) -ldflags="$(LD_FLAGS)" -o $(BIN_DIR)/$(BIN_NAME) ./cmd

run: build
	@echo "Running $(BIN_NAME)..."
	@$(BIN_DIR)/$(BIN_NAME) $(ARGS)

clean:
	@echo "Cleaning..."
	@rm -rf $(BIN_DIR)
	@rm -f *.png *.jpg *.mp4

test:
	@echo "Running tests..."
	@go test $(SRC_DIR) -v

deps:
	@echo "Downloading dependencies..."
	@go mod download

install: build
	@echo "Installing to GOPATH/bin..."
	@cp $(BIN_DIR)/$(BIN_NAME) $(GOPATH)/bin/$(BIN_NAME)

# 交叉编译目标
build-linux:
	@GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(BIN_NAME)-linux ./cmd

build-windows:
	@GOOS=windows GOARCH=amd64 go build -o $(BIN_DIR)/$(BIN_NAME).exe ./cmd

# 开发工具检查
check-ffmpeg:
	@which ffmpeg || (echo "Error: ffmpeg is required for video generation" && exit 1)
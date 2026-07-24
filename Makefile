.PHONY: build clean test frontend install uninstall

VERSION?=1.0.0
BUILD_TIME=$(shell date +%Y-%m-%d_%H:%M:%S)
GOOS?=linux
GOARCH?=amd64
LDFLAGS=-ldflags "-s -w -X github.com/openforge-maintain/openforge-maintain/pkg/utils.Version=$(VERSION) -X github.com/openforge-maintain/openforge-maintain/pkg/utils.BuildTime=$(BUILD_TIME)"

# ==================== 构建 ====================

# 构建所有二进制
build: build-core build-agent

# 构建 Core 服务
build-core:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o bin/core ./cmd/core

# 构建 Agent 服务
build-agent:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o bin/agent ./cmd/agent

# 交叉编译 (当前平台 + arm64)
build-all:
	@echo "Building for linux/amd64..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/core-linux-amd64 ./cmd/core
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/agent-linux-amd64 ./cmd/agent
	@echo "Building for linux/arm64..."
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/core-linux-arm64 ./cmd/core
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/agent-linux-arm64 ./cmd/agent
	@echo "Done."

# 前端构建
frontend:
	cd frontend && npm install && npm run build

# ==================== 开发 ====================

# 开发模式运行 Core
run-core:
	go run ./cmd/core

# 开发模式运行 Agent
run-agent:
	go run ./cmd/agent

# 前端开发模式
dev-frontend:
	cd frontend && npm install && npm run dev

# 同时启动后端和前端开发模式
dev:
	@echo "Starting Core service on :9999..."
	@go run ./cmd/core &
	@echo "Starting Agent service on :10000..."
	@go run ./cmd/agent &
	@echo "Starting frontend dev server..."
	@cd frontend && npm run dev

# 数据库迁移
migrate-core:
	go run ./cmd/core --migrate

migrate-agent:
	go run ./cmd/agent --migrate

# ==================== 测试 ====================

# 运行所有测试
test:
	go test ./... -v -race -coverprofile=coverage.out

# 运行后端测试
test-core:
	go test ./core/... -v

test-agent:
	go test ./agent/... -v

test-pkg:
	go test ./pkg/... -v

# 前端测试
test-frontend:
	cd frontend && npm run test

# ==================== 代码质量 ====================

# 代码格式化
fmt:
	go fmt ./...

# 代码检查
lint:
	golangci-lint run ./...

# 依赖检查
tidy:
	go mod tidy

# ==================== 部署 ====================

# 打包发布
release: build-all
	@mkdir -p releases
	tar czf releases/openforge-maintain-linux-amd64.tar.gz -C bin core-linux-amd64 agent-linux-amd64
	tar czf releases/openforge-maintain-linux-arm64.tar.gz -C bin core-linux-arm64 agent-linux-arm64
	@echo "Release packages created in releases/"

# 安装
install: build
	sudo bash deploy/install.sh

# 卸载
uninstall:
	sudo bash deploy/install.sh uninstall

# ==================== 清理 ====================

clean:
	rm -rf bin/ releases/
	cd frontend && rm -rf node_modules dist

# ==================== Docker ====================

docker-build:
	docker build -t openforge-maintain:$(VERSION) .

docker-run:
	docker run -d \
		--name openforge-maintain \
		-p 9999:9999 \
		-p 10000:10000 \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v /opt/openforge-maintain/data:/opt/openforge-maintain/data \
		openforge-maintain:$(VERSION)

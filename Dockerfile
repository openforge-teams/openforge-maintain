# ==================== Build Stage ====================
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm install --registry https://registry.npmmirror.com
COPY frontend/ ./
RUN npm run build

# ==================== Final Stage ====================
FROM golang:1.22-alpine AS builder

WORKDIR /app

# 安装构建依赖
RUN apk add --no-cache gcc musl-dev

# 复制 Go 模块文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建二进制
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o /app/bin/core ./cmd/core

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o /app/bin/agent ./cmd/agent

# ==================== Runtime Stage ====================
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata sqlite-libs docker-cli

WORKDIR /opt/openforge-maintain

# 复制二进制
COPY --from=builder /app/bin/core /opt/openforge-maintain/bin/core
COPY --from=builder /app/bin/agent /opt/openforge-maintain/bin/agent

# 复制前端资源
COPY --from=frontend-builder /app/frontend/dist /opt/openforge-maintain/frontend

# 创建数据目录
RUN mkdir -p /opt/openforge-maintain/data /opt/openforge-maintain/frontend /var/log/openforge-maintain

# 复制部署文件
COPY deploy/systemd/*.service /etc/systemd/system/

# 环境变量
ENV MAINTAIN_PORT=9999
ENV AGENT_PORT=10000
ENV MAINTAIN_DB_PATH=/opt/openforge-maintain/data/core.db
ENV AGENT_DB_PATH=/opt/openforge-maintain/data/agent.db
ENV DOCKER_HOST=unix:///var/run/docker.sock
ENV GIN_MODE=release

# 暴露端口
EXPOSE 9999 10000

# 入口
CMD ["/bin/sh", "-c", "/opt/openforge-maintain/bin/core & /opt/openforge-maintain/bin/agent & wait"]

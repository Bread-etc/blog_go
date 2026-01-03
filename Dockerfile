# Build Stage
FROM golang:1.25 AS builder

# 设置环境变量
ENV GO111MODULE=on \
  GOPROXY=https://goproxy.io,direct \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

WORKDIR /app

# 缓存依赖
COPY go.mod go.sum ./
RUN go mod download

# 编译
COPY . .
RUN go build -ldflags="-s -w" -o blog-server .

# Run Stage
FROM debian:bookworm-slim

WORKDIR /app

# 安装必要的系统证书 & 清理缓存 / 设置时区
RUN apt-get update && apt-get install -y ca-certificates tzdata && \
  rm -rf /var/lib/apt/lists/*

ENV TZ=Asia/Shanghai

# 复制二进制和配置文件
COPY --from=builder /app/blog-server .
COPY config/config.yaml config/config.yaml

RUN mkdir logs

EXPOSE 8080

CMD ["./blog-server"]
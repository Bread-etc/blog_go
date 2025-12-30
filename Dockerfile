# 1. 编译 (Builder)
FROM golang:1.25 AS builder

ENV GO111MODULE=on \
  GOPROXY=https://goproxy.cn,direct

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

# 拷贝所有源码
COPY . .

# CGO_ENABLED=0: 禁用 CGO，生成纯静态二进制文件，兼容性最好
# GOOS=linux: 目标是 Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o blog-server main.go

# 2. 运行阶段 (Runner)
FROM debian:bookworm-slim

# 安装证书
RUN apt-get update && apt-get install -y ca-certificates tzdata && rm -rf /var/lib/apt/lists/*
ENV TZ=Asia/Shanghai

WORKDIR /app

# 从 builder 拷贝编译好的程序
COPY --from=builder /build/blog-server .
# 拷贝配置文件
RUN mkdir config
COPY config/config.yaml config/

# 创建日志目录
RUN mkdir logs

EXPOSE 8080

ENTRYPOINT [ "./blog-server" ]
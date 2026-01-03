# blog_go

基于 `Go` + `Gin` + `GORM` + `MySQL` 构建的博客后端，本项目集成了`CI/CD`流程、`Docker`容器化部署方案。

## 🛠 技术栈

- **语言**: Golang 1.25
- **Web 框架**: Gin
- **ORM**: GORM (MySQL 8.4)
- **日志**: Zap + Lumberjack
- **配置管理**: Viper
- **测试**: Testify + SQLite (In-memory)
- **部署**: Docker, Caddy, GitHub Actions

## 📂 目录结构

```text
├── config/             # 配置文件模板 (config.yaml)
├── controller/         # 控制器层
├── docs/               # 项目文档
├── middleware/         # 中间件
├── model/              # 数据库模型
├── pkg/                # 公共工具包
├── router/             # 路由定义
├── services/           # 业务逻辑层 & 单元测试
├── .github/workflows/  # CI/CD 配置
├── Caddyfile           # Caddy 反向代理配置
├── Dockerfile          # docker镜像配置
├── docker-compose.yml  # 容器编排配置
└── main.go             # 入口
```

## 🚀 快速开始

### 前置要求

- Go 1.25+
- MySQL 8.0+

### 1. 克隆项目

```bash
git clone https://github.com/your/repo.git
cd blog_go
```

### 2. 配置环境

在项目根目录创建 `.env` 文件 (可参考 config.yaml 结构，或者直接修改 config/config.yaml 用于本地调试)
**注意**: 不要将包含真实密码的 `.env` 提交到版本控制

### 3. 运行程序

```bash
# 下载依赖
go mod tidy

# 启动服务
go run main.go
```

服务默认运行在 `http://localhost:8080`

### 4. 运行单元测试

```bash
# 运行 Service 层的单元测试
go test -v ./services/...
```

## 🐳 Docker 部署 (生产环境)

本项目使用 Docker Compose 进行一键部署，自动包含 **MySQL** 和 **Caddy** (反向代理)

```bash
# 构建并启动所有服务
docker compose up -d --build
```

启动后，Caddy 会自动为你的域名申请 HTTPS 证书

- API 接口: `https://your-domain.com/api`
- 健康检查: `https://your-domain.com/api/health`

## 🔄 CI/CD 工作流

本项目配置了 **GitHub Actions** 实现自动化部署

1.  **自动测试 (Test)**: 修改版本变更文档`CHANGELOG.md`会触发单元测试，确保代码质量
2.  **自动发布 (Deploy)**: 只有 **Master** 分支且提交信息时，才会执行构建与部署

示例：

```bash
git commit -m "build(version): release v1.0.0 - initial launch"
git push origin master
```

触发后，Actions 流水线会自动：

1.  运行单元测试 (PASS 后继续)
2.  构建 Docker 镜像并推送至 Docker Hub
3.  SSH 连接服务器拉取新镜像并重启服务

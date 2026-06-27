# blog_go

基于 `Go` + `Gin` + `GORM` + `MySQL` 构建的个人博客后端 API。项目已完成核心后端功能重构，当前重点是为前端开发提供稳定、可测试、文档清晰的接口契约。

## 技术栈

- Go 1.25
- Gin
- GORM + MySQL 8.4
- JWT 鉴权
- RSA 非对称加密登录
- Zap 日志
- Viper 配置管理
- Testify + SQLite 内存数据库单元测试
- Bruno CLI 接口自动化测试
- Docker Compose + Caddy 部署
- GitHub Actions CI/CD

## 功能模块

- Auth：公钥获取、登录、当前用户信息、修改密码
- Post：文章 CRUD、Slug 详情、浏览量、分页、关键词、分类和标签过滤
- Category：分类 CRUD、唯一性校验、删除引用限制
- Tag：标签 CRUD、唯一性校验、删除引用限制
- Link：友情链接 CRUD 和排序
- Config：站点配置读取和更新
- Dashboard：统计卡片、热门文章排行
- Health：数据库健康检查

## 项目结构

```text
├── bruno/              # Bruno 接口集合和自动化测试
├── config/             # 配置文件
├── controller/         # HTTP 控制层
├── docs/               # API 文档与版本记录
├── dto/                # 请求/响应 DTO
├── middleware/         # Gin 中间件
├── model/              # GORM 数据模型
├── pkg/                # 公共工具包
├── router/             # 路由定义
├── services/           # 业务逻辑层和单元测试
├── .github/workflows/  # GitHub Actions
├── Caddyfile           # Caddy 反向代理配置
├── Dockerfile          # 后端镜像构建
├── docker-compose.yml  # 生产编排
└── main.go             # 应用入口
```

## 本地启动

前置要求：

- Go 1.25+
- MySQL 8.0+

配置 `config/config.yaml`，然后启动服务：

```powershell
go mod download
go run main.go
```

默认地址：

```text
http://localhost:8080
```

健康检查：

```text
GET http://localhost:8080/api/health
```

## 测试

运行 Go 单元测试：

```powershell
go test ./...
```

运行 Bruno 接口自动化测试：

```powershell
cd bruno/Blog-gin
bru run --env Local --sandbox developer --exclude-tags manual
```

说明：

- Bruno 默认覆盖主流程、负向用例和边界用例。
- `manual` 标签用于改密等高风险接口，默认排除。
- 当前接口调试与契约验证以 Bruno 和 `docs/API_REFERENCE.md` 为准，暂不引入 OpenAPI/Swagger。

## 接口文档

- API 参考：[docs/API_REFERENCE.md](docs/API_REFERENCE.md)
- Bruno 使用说明：[bruno/Blog-gin/README.md](bruno/Blog-gin/README.md)
- 版本记录：[docs/CHANGELOG.md](docs/CHANGELOG.md)

## Docker 部署

生产环境使用 Docker Compose 编排后端、MySQL 和 Caddy：

```powershell
docker compose up -d
```

服务器部署目录需要提前准备 `.env`，该文件不提交到仓库：

```env
DB_ROOT_PASSWORD=your_root_password
DB_NAME=blog
DB_USER=blog_user
DB_PASSWORD=your_db_password
JWT_SECRET=your_jwt_secret
TZ=Asia/Shanghai
```

后端容器会通过 Compose 注入的环境变量覆盖 `config/config.yaml` 中的数据库配置。`keys/`、`logs/`、`frontend/dist/` 等目录会作为持久化或静态资源目录使用。

## CI/CD

GitHub Actions 工作流位于 `.github/workflows/deploy.yml`。

触发方式：

- 推送到 `master` 且修改了 `docs/CHANGELOG.md`
- 手动执行 `workflow_dispatch`

流水线步骤：

1. 运行 `go test -v ./...`
2. 构建并推送 Docker 镜像到 Docker Hub
3. 通过 SSH 登录服务器，拉取新镜像并执行 `docker compose up -d`

生产服务器必须提前在 `/opt/blog_go/.env` 准备真实环境变量，否则部署会主动失败并提示缺少 `.env`。

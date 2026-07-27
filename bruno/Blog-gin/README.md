# Blog-gin Bruno 接口测试集合

这是当前 Go 博客后端的 Bruno collection，用于配合 Bruno Desktop 或 Bruno CLI 做接口调试与自动化测试。

## 环境要求

- 后端服务需要运行在 `http://localhost:8080`。
- MySQL 需要指向可写的本地或测试数据库。
- 默认管理员账号需要可用，当前默认值为 `admin / admin`。
- 当前已验证的 Bruno CLI 版本：`3.5.0`。

## 本地环境变量

环境文件位于 `environments/Local.bru`。

常用配置如下：

```bru
vars {
  baseUrl: http://localhost:8080
  adminUsername: admin
  adminPassword: admin
  newAdminPassword:
}
```

说明：

- `baseUrl` 是本地后端地址，默认使用 `http://localhost:8080`。
- `adminUsername` 和 `adminPassword` 用于登录并获取 JWT token。
- `newAdminPassword` 只用于手动改密请求，默认不要填写。
- `token`、`runId`、`categoryId`、`postId` 等变量会在运行过程中由脚本自动写入，不需要手动维护。

## 运行方式

进入 collection 目录后执行：

```powershell
cd bruno/Blog-gin
bru run --env Local --sandbox developer --exclude-tags manual
```

必须使用 `--sandbox developer`，因为登录接口会在 pre-request script 中调用 Node.js 的 `crypto` 模块，对密码做 RSA 公钥加密。

必须保留 `--exclude-tags manual`，因为 `01_Auth/04_change_password_manual.bru` 会修改真实管理员密码，默认不应该进入自动化测试。

## 测试目录说明

- `00_System`：健康检查，并初始化本轮测试的 `runId`。
- `01_Auth`：获取公钥、登录、获取当前用户信息，以及手动改密请求。
- `02_Category`：分类 CRUD 主流程。
- `03_Tag`：标签 CRUD 主流程。
- `04_Post`：文章依赖数据创建、文章 CRUD、详情与浏览量拆分、`tagIds` 重复参数和 CSV 参数查询、依赖数据清理。
- `05_Link`：友情链接 CRUD 主流程。
- `06_Config`：站点配置获取、更新、再次获取。
- `07_Dashboard`：后台统计卡片和热门文章排行。
- `90_Negative`：基础负向测试，例如未登录、缺少必填字段、非法 URL、非法邮箱等。
- `91_Boundary`：边界与健壮性测试，例如重复名称/slug、长度超限、关联不存在、删除被引用资源、PUT 缺字段、分页上限、重复 `tagIds`、鉴权异常等。

## 当前自动化覆盖范围

默认命令会执行主流程、基础负向测试和边界测试。

覆盖类型包括：

- 正常 CRUD：Category、Tag、Post、Link、Config。
- 鉴权流程：登录、token 复用、未携带 token、非法 token。
- 参数校验：缺少必填字段、字段类型错误、长度超限、非法 URL、非法 email、非法分页参数。
- 业务冲突：重复分类名称、重复分类 slug、重复标签名称、重复标签 slug、重复文章 slug。
- 关联约束：文章分类不存在、文章标签不存在、分类被文章引用时禁止删除、标签被文章引用时禁止删除。
- 列表查询：分页参数、关键词长度、`tagIds` 重复参数、`tagIds` CSV 参数、重复/空 tag id。
- 浏览量：详情查询不自增、显式记录浏览后只增加一次、不存在的文章返回 404。
- Dashboard：热门文章 limit 最小值和最大值校验。

## 测试数据策略

测试数据尽量通过接口创建，不依赖数据库已有业务数据。

主流程和边界流程都会创建自己的分类、标签、文章和友链，并在流程结束时尽量通过接口清理。

需要注意：

- 如果测试中途被强制中断，可能会留下以 `Auto` 或 `Boundary` 开头的测试数据。
- `06_Config/02_update_config.bru` 会修改单例站点配置，目前不会自动恢复执行前的原始配置。后续如果需要完全无副作用，可以补充“读取原配置 -> 更新测试配置 -> 恢复原配置”的闭环。

## 手动改密请求

`01_Auth/04_change_password_manual.bru` 被标记为 `manual`，不会被默认命令执行。

如果确实需要测试改密：

1. 在 `environments/Local.bru` 中设置 `newAdminPassword`。
2. 先执行获取公钥和登录请求，确保 `publicKey` 与 `token` 已写入运行时变量。
3. 单独执行 `01_Auth/04_change_password_manual.bru`。
4. 改密成功后，记得同步更新 `adminPassword`，否则后续自动化登录会失败。

## 生成报告

可以手动指定 JSON 报告路径：

```powershell
bru run --env Local --sandbox developer --exclude-tags manual --output reports/bruno-run-local.json
```

如果需要 HTML 或 JUnit 报告，可以使用 Bruno CLI 的 reporter 参数，例如：

```powershell
bru run --env Local --sandbox developer --exclude-tags manual --reporter-html reports/bruno-run-local.html --reporter-junit reports/bruno-run-local.xml
```

## 推荐使用方式

开发阶段建议按目录小步执行：

```powershell
bru run 01_Auth --env Local --sandbox developer
bru run 04_Post --env Local --sandbox developer
bru run 91_Boundary --env Local --sandbox developer
```

发版前建议执行完整集合：

```powershell
bru run --env Local --sandbox developer --exclude-tags manual
```

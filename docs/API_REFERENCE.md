# API 接口参考

> 本文档对应当前后端接口契约。字段命名统一使用 `camelCase`，数据库字段仍使用 `snake_case`。

## 基础约定

- 本地地址：`http://localhost:8080`
- API 前缀：`/api`
- 鉴权方式：`Authorization: Bearer <token>`
- 成功响应：`errorCode` 为空字符串
- 失败响应：`data` 固定为 `null`

```json
{
  "code": 200,
  "errorCode": "",
  "message": "success",
  "data": {}
}
```

```json
{
  "code": 400,
  "errorCode": "INVALID_PARAMS",
  "message": "invalid request parameters",
  "data": null
}
```

通用错误码：

| errorCode              | HTTP | 说明                            |
| ---------------------- | ---: | ------------------------------- |
| `INVALID_PARAMS`       |  400 | 参数错误或 JSON 类型错误        |
| `UNAUTHORIZED`         |  401 | 未登录、Token 缺失或 Token 无效 |
| `FORBIDDEN`            |  403 | 无权限                          |
| `NOT_FOUND`            |  404 | 资源不存在                      |
| `CONFLICT`             |  409 | 资源冲突                        |
| `INTERNAL_ERROR`       |  500 | 内部错误                        |
| `DATABASE_UNAVAILABLE` |  500 | 数据库不可用或数据库操作失败    |

## Auth

| 方法   | 路径                        | 鉴权 | 说明          |
| ------ | --------------------------- | ---- | ------------- |
| `GET`  | `/api/auth/public-key`      | 否   | 获取 RSA 公钥 |
| `POST` | `/api/auth/login`           | 否   | 登录          |
| `GET`  | `/api/auth/profile`         | 是   | 当前用户信息  |
| `POST` | `/api/auth/change-password` | 是   | 修改密码      |

登录请求：

```json
{
  "username": "admin",
  "password": "<RSA 公钥加密后的 Base64 密文>"
}
```

登录响应 `data`：

```json
{
  "token": "<jwt>",
  "user": {
    "id": "uuid",
    "username": "admin",
    "role": "admin"
  }
}
```

改密请求：

```json
{
  "oldPassword": "<RSA 公钥加密后的旧密码>",
  "newPassword": "<RSA 公钥加密后的新密码>"
}
```

主要错误码：`USERNAME_REQUIRED`、`USERNAME_TOO_LONG`、`PASSWORD_REQUIRED`、`OLD_PASSWORD_REQUIRED`、`NEW_PASSWORD_REQUIRED`、`INVALID_PASSWORD_ENCRYPTION`、`INVALID_CREDENTIALS`、`INCORRECT_OLD_PASSWORD`、`UNAUTHORIZED`。

## Post

| 方法     | 路径                   | 鉴权 | 说明         |
| -------- | ---------------------- | ---- | ------------ |
| `GET`    | `/api/posts`           | 否   | 文章列表     |
| `GET`    | `/api/posts/:slug`     | 否   | 文章详情     |
| `POST`   | `/api/posts/:id/views` | 否   | 浏览量加一   |
| `POST`   | `/api/posts`           | 是   | 创建文章     |
| `PUT`    | `/api/posts/:id`       | 是   | 全量更新文章 |
| `DELETE` | `/api/posts/:id`       | 是   | 删除文章     |

列表 Query：

| 参数          | 类型     | 默认值 | 约束                       |
| ------------- | -------- | ------ | -------------------------- |
| `page`        | number   | `1`    | 最小 `1`                   |
| `pageSize`    | number   | `10`   | `1-50`                     |
| `keyword`     | string   | -      | 最长 `100`，搜索标题和摘要 |
| `categoryId`  | string   | -      | 分类 ID                    |
| `tagIds`      | string[] | -      | 支持重复参数和 CSV         |
| `isPublished` | boolean  | -      | 发布状态；不传时查询全部   |

前台文章列表应显式传递 `isPublished=true`；后台管理列表不传时查询全部文章，传递 `false` 时仅查询未发布文章。

`tagIds` 示例：

```http
GET /api/posts?tagIds=tag-a&tagIds=tag-b
GET /api/posts?tagIds=tag-a,tag-b
GET /api/posts?tagIds=tag-a&tagIds=tag-b,tag-c
```

获取文章详情不会自动增加浏览量。前台成功展示文章后，调用以下无请求体接口记录一次浏览；后台编辑文章时无需调用：

```http
POST /api/posts/:id/views
```

成功响应中的 `data` 为 `null`。

创建/更新请求：

```json
{
  "title": "文章标题",
  "content": "正文",
  "summary": "摘要",
  "slug": "post-slug",
  "cover": "https://example.com/cover.png",
  "categoryId": "category-uuid",
  "tagIds": ["tag-uuid"],
  "isPublished": true
}
```

字段约束：

| 字段          | 约束                                                      |
| ------------- | --------------------------------------------------------- |
| `title`       | 必填，最长 `255`                                          |
| `content`     | 必填                                                      |
| `summary`     | 可为空，最长 `500`                                        |
| `slug`        | 必填，最长 `255`，唯一                                    |
| `cover`       | 可为空，最长 `255`                                        |
| `categoryId`  | 必填，必须存在                                            |
| `tagIds`      | 必填，至少 1 个，不能重复，不能包含空字符串，标签必须存在 |
| `isPublished` | 必填，允许 `true/false`                                   |

列表响应 `data`：

```json
{
  "list": [
    {
      "id": "uuid",
      "title": "文章标题",
      "summary": "摘要",
      "slug": "post-slug",
      "cover": "",
      "views": 0,
      "isPublished": true,
      "createdAt": "2026-06-27T12:00:00+08:00",
      "updatedAt": "2026-06-27T12:00:00+08:00",
      "category": {
        "id": "uuid",
        "name": "Go",
        "slug": "go"
      },
      "tags": [
        {
          "id": "uuid",
          "name": "Gin",
          "slug": "gin"
        }
      ]
    }
  ],
  "total": 1,
  "page": 1,
  "pageSize": 10
}
```

主要错误码：`POST_NOT_FOUND`、`POST_SLUG_EXISTS`、`POST_CATEGORY_NOT_FOUND`、`POST_TAG_NOT_FOUND`、`POST_TITLE_REQUIRED`、`POST_CONTENT_REQUIRED`、`POST_SLUG_REQUIRED`、`POST_CATEGORY_ID_REQUIRED`、`POST_TAG_IDS_REQUIRED`、`POST_IS_PUBLISHED_REQUIRED`、`POST_TITLE_TOO_LONG`、`POST_SUMMARY_TOO_LONG`、`POST_SLUG_TOO_LONG`、`POST_COVER_TOO_LONG`、`POST_KEYWORD_TOO_LONG`、`POST_TAG_IDS_INVALID`、`POST_PAGE_INVALID`、`POST_PAGE_SIZE_INVALID`。

## Category

| 方法     | 路径                  | 鉴权 | 说明         |
| -------- | --------------------- | ---- | ------------ |
| `GET`    | `/api/categories`     | 否   | 分类列表     |
| `POST`   | `/api/categories`     | 是   | 创建分类     |
| `PUT`    | `/api/categories/:id` | 是   | 全量更新分类 |
| `DELETE` | `/api/categories/:id` | 是   | 删除分类     |

创建/更新请求：

```json
{
  "name": "Go",
  "slug": "go"
}
```

约束：`name` 必填、最长 `50`、唯一；`slug` 必填、最长 `100`、唯一。

列表响应 `data`：

```json
[
  {
    "id": "uuid",
    "name": "Go",
    "slug": "go"
  }
]
```

主要错误码：`CATEGORY_NOT_FOUND`、`CATEGORY_NAME_EXISTS`、`CATEGORY_SLUG_EXISTS`、`CATEGORY_IN_USE`、`CATEGORY_NAME_REQUIRED`、`CATEGORY_SLUG_REQUIRED`、`CATEGORY_NAME_TOO_LONG`、`CATEGORY_SLUG_TOO_LONG`。

## Tag

| 方法     | 路径            | 鉴权 | 说明         |
| -------- | --------------- | ---- | ------------ |
| `GET`    | `/api/tags`     | 否   | 标签列表     |
| `POST`   | `/api/tags`     | 是   | 创建标签     |
| `PUT`    | `/api/tags/:id` | 是   | 全量更新标签 |
| `DELETE` | `/api/tags/:id` | 是   | 删除标签     |

创建/更新请求：

```json
{
  "name": "Gin",
  "slug": "gin"
}
```

约束：`name` 必填、最长 `50`、唯一；`slug` 必填、最长 `100`、唯一。

主要错误码：`TAG_NOT_FOUND`、`TAG_NAME_EXISTS`、`TAG_SLUG_EXISTS`、`TAG_IN_USE`、`TAG_NAME_REQUIRED`、`TAG_SLUG_REQUIRED`、`TAG_NAME_TOO_LONG`、`TAG_SLUG_TOO_LONG`。

## Link

| 方法     | 路径             | 鉴权 | 说明                                           |
| -------- | ---------------- | ---- | ---------------------------------------------- |
| `GET`    | `/api/links`     | 否   | 友链列表，按 `sort DESC, created_at DESC` 排序 |
| `POST`   | `/api/links`     | 是   | 创建友链                                       |
| `PUT`    | `/api/links/:id` | 是   | 全量更新友链                                   |
| `DELETE` | `/api/links/:id` | 是   | 删除友链                                       |

创建/更新请求：

```json
{
  "name": "Example",
  "url": "https://example.com",
  "description": "描述",
  "sort": 10
}
```

约束：

| 字段          | 约束                       |
| ------------- | -------------------------- |
| `name`        | 必填，最长 `50`            |
| `url`         | 必填，合法 URL，最长 `255` |
| `description` | 可为空，最长 `255`         |
| `sort`        | 最小 `0`，数值越大越靠前   |

主要错误码：`LINK_NOT_FOUND`、`LINK_NAME_REQUIRED`、`LINK_URL_REQUIRED`、`LINK_URL_INVALID`、`LINK_SORT_INVALID`、`LINK_NAME_TOO_LONG`、`LINK_URL_TOO_LONG`、`LINK_DESCRIPTION_TOO_LONG`。

## Site Config

| 方法  | 路径          | 鉴权 | 说明               |
| ----- | ------------- | ---- | ------------------ |
| `GET` | `/api/config` | 否   | 获取站点配置       |
| `PUT` | `/api/config` | 是   | 创建或更新站点配置 |

请求/响应字段：

```json
{
  "title": "站点标题",
  "subtitle": "副标题",
  "description": "站点描述",
  "keywords": "go,blog",
  "author": "作者",
  "email": "name@example.com",
  "githubUrl": "https://github.com/example"
}
```

约束：

| 字段          | 约束                                       |
| ------------- | ------------------------------------------ |
| `title`       | 必填，最长 `100`                           |
| `subtitle`    | 最长 `255`                                 |
| `description` | 最长 `1000`                                |
| `keywords`    | 最长 `255`                                 |
| `author`      | 最长 `50`                                  |
| `email`       | 可为空；非空时必须是合法 email，最长 `100` |
| `githubUrl`   | 可为空；非空时必须是合法 URL，最长 `255`   |

主要错误码：`CONFIG_TITLE_REQUIRED`、`CONFIG_TITLE_TOO_LONG`、`CONFIG_SUBTITLE_TOO_LONG`、`CONFIG_DESCRIPTION_TOO_LONG`、`CONFIG_KEYWORDS_TOO_LONG`、`CONFIG_AUTHOR_TOO_LONG`、`CONFIG_EMAIL_INVALID`、`CONFIG_EMAIL_TOO_LONG`、`CONFIG_GITHUB_URL_INVALID`、`CONFIG_GITHUB_URL_TOO_LONG`。

## Dashboard

| 方法  | 路径                       | 鉴权 | 说明         |
| ----- | -------------------------- | ---- | ------------ |
| `GET` | `/api/dashboard/stats`     | 是   | 后台统计卡片 |
| `GET` | `/api/dashboard/top-posts` | 是   | 热门文章排行 |

`GET /api/dashboard/top-posts` Query：

| 参数    | 类型   | 默认值 | 约束   |
| ------- | ------ | ------ | ------ |
| `limit` | number | `5`    | `1-10` |

统计响应 `data`：

```json
{
  "posts": {
    "total": 10,
    "moMGrowth": 20
  },
  "categories": {
    "total": 3,
    "moMGrowth": 0
  },
  "tags": {
    "total": 8,
    "moMGrowth": 100
  },
  "links": {
    "total": 5,
    "moMGrowth": -10
  },
  "totalViews": 123
}
```

主要错误码：`AGGREGATE_TOP_POSTS_LIMIT_INVALID`、`DATABASE_UNAVAILABLE`。

## Health

```http
GET /api/health
```

成功响应 `data`：

```json
{
  "status": "UP",
  "database": "connected"
}
```

数据库不可用时返回 `DATABASE_UNAVAILABLE`。

## Bruno 自动化测试

当前 Bruno collection 位于 `bruno/Blog-gin`，默认覆盖 `84` 个请求：

- 主流程：Health、Auth、Category、Tag、Post、Link、Config、Dashboard。
- 基础负向测试：未登录、缺少必填字段、非法 URL、非法 email、非法分页等。
- 边界测试：重复名称/slug、长度超限、关联不存在、删除被引用资源、PUT 缺字段、重复或空 `tagIds`、非法 token。

运行命令：

```powershell
cd bruno/Blog-gin
bru run --env Local --sandbox developer --exclude-tags manual
```

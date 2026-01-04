# API 接口参考文档 (v1.0.1)

后端地址: `https://hastur23.top`
认证方式: Header `Authorization: Bearer <token>`

## 1. 用户 (User)

- **POST** `/api/user/login`: 用户登录 (参数: username, password(RSA 加密))
- **GET** `/api/user/public-key`: 获取 RSA 公钥
- **GET** `/api/user/profile`: 获取当前用户信息 [Auth]
- **POST** `/api/user/change-password`: 修改密码 [Auth]

## 2. 文章 (Post)

- **GET** `/api/posts`: 获取文章列表 (分页, 筛选: category_id, tag_id, keyword)
- **GET** `/api/posts/:slug`: 获取文章详情 (通过 Slug)
- **POST** `/api/posts`: 创建文章 [Auth]
- **PUT** `/api/posts/:id`: 更新文章 [Auth]
- **DELETE** `/api/posts/:id`: 删除文章 [Auth]

## 3. 分类 (Category)

- **GET** `/api/categories`: 获取分类列表
- **POST** `/api/categories`: 创建分类 [Auth]
- **PUT** `/api/categories/:id`: 更新分类 [Auth]
- **DELETE** `/api/categories/:id`: 删除分类 [Auth]

## 4. 标签 (Tag)

- **GET** `/api/tags`: 获取标签列表
- **POST** `/api/tags`: 创建标签 [Auth]
- **PUT** `/api/tags/:id`: 更新标签 [Auth]
- **DELETE** `/api/tags/:id`: 删除标签 [Auth]

## 5. 友链 (Link)

- **GET** `/api/links`: 获取友链列表
- **POST** `/api/links`: 创建/审核友链 [Auth]
- **PUT** `/api/links/:id`: 更新友链 [Auth]
- **DELETE** `/api/links/:id`: 删除友链 [Auth]

## 6. 站点配置 (Config)

- **GET** `/api/config`: 获取站点配置 (Title, Desc, etc.)
- **PUT** `/api/config`: 更新站点配置 [Auth]

## 7. 系统

- **GET** `/api/health`: 健康检查

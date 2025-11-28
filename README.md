# **个人博客系统后端**

这是一个基于gin+GORM的博客项目，集成JWT认证、统一日志等功能。提供用户注册/登录、用户管理、文章发布/删除、评论互动功能。

## 技术栈

- **核心语言**：Go 1.21+
- **Web 框架**：Gin
- **数据库**：MySQL 8.0+ 
- **ORM**：GORM
- **认证授权**：JWT
- **日志**：zerolog
- **依赖管理**：Go Modules

## 功能清单

### 1. 用户模块

- 注册：用户名+ 密码注册（密码加密存储）+邮箱
- 登录：账号密码登录，返回 JWT Token
- 用户列表：查询所有用户

### 2. 文章模块

- 文章 CRUD：发布、查询（列表 / 详情）、修改、删除

### 3. 评论模块

- 评论操作：对文章发表评论、查询评论列表

## 快速开始

### 1. 环境准备

- 安装 Go 1.19+：[官方下载地址](https://golang.org/dl/)

- 安装数据库：MySQL 8.0+ 或 PostgreSQL 14+

- 克隆项目：

  ```bash
  git@github.com:ydh2333/personalBlog.git
  ```

### 2. 配置文件

编辑 `config.yaml`，修改数据库、JWT 等核心配置：

```yaml
consoleLoggingEnabled: true
encodeLogsAsJson: true
fileLoggingEnabled: true
directory: "./logs/"
filename: "./logs"
maxSize: 1
maxBackups: 10
maxAge: 30
level: 1
db:
  host: "127.0.0.1"
  port: 3306
  username: "root"
  password: "root"
  db_name: "personal_blog"
  charset: "utf8mb4"
jwt:
  secret: my-secret-key-ydh
  expiration: 3600
```

### 3. 依赖安装

```bash
go mod download
```

### 4. 数据库迁移

自动创建数据表（基于 GORM AutoMigrate）：

```bash
 go run .\dao\migrate\main.go
```

说明：迁移脚本会自动执行创建表，仅在首次运行前执行。

### 5. 启动服务

```bash
go run cmd/server/main.go
```

服务启动后，访问 `http://127.0.0.1:8080/对应接口` 即可调用接口

## API 接口文档

| 模块 | 接口地址               | 请求方式 | 说明                   | 权限要求 |
| ---- | ---------------------- | -------- | ---------------------- | -------- |
| 用户 | `/register`            | POST     | 用户注册               | 匿名     |
| 用户 | `/login`               | POST     | 用户登录（返回 Token） | 匿名     |
| 用户 | `/userAll`             | GET      | 获取用户列表           | 登录用户 |
| 文章 | `/createPost`          | POST     | 创建文章               | 登录用户 |
| 文章 | `/getPostList`         | GET      | 获取文章列表           | 登录用户 |
| 文章 | `/updatePost`          | POST     | 修改文章               | 文章作者 |
| 文章 | `/deletePost/:id`      | DELETE   | 删除文章               | 文章作者 |
| 评论 | `/getCommList/:postId` | GET      | 获取文章评论列表       | 匿名     |
| 评论 | `/createComm`          | POST     | 发表评论               | 登录用户 |

## 项目结构

```plaintext
personalBlog/
├── config/               # 配置相关
│   ├── config.go         # 配置解析逻辑
├── dao/                  # 数据库访问层
│   └── migrate/          # 数据库迁移
│       └── main.go/      # 自动迁移数据表，仅首次执行
├── logs/                 # 日志文件（自动生成）
├── model/                # 数据模型（数据库表结构）
├── service/              # 业务逻辑层
├── util/                 # 工具包（数据库连接、日志、中间件、JWT等）
├── go.mod                # 依赖管理
├── main.mod              # 服务启动入口
└── README.md             # 项目说明（本文档）
```


# go-gin-api

## 项目概述

go-gin-api 是一个基于 Gin 框架的模块化 API 开发框架，封装了企业级应用常用的功能组件，致力于提供快速、规范的业务研发体验。框架通过设计约束开发规范，规避混乱无序的编码方式。

**核心功能特性：**
- 支持接口限流和跨域访问
- 集成 Prometheus 监控指标记录
- 提供 Swagger 接口文档自动生成
- 支持 GraphQL 查询语言
- 内置链路追踪和性能剖析
- 统一的错误码定义和日志收集
- 支持定时任务和 WebSocket 实时通讯
- 提供 Web 管理界面和代码生成工具

## 技术栈

**后端框架：**
- Gin - HTTP Web 框架
- GORM - ORM 数据库组件
- Viper - 配置管理
- Zap - 结构化日志
- Redis - 缓存和会话存储

**监控和工具：**
- Prometheus - 指标监控
- Swagger - API 文档
- GraphQL - 查询语言
- PProf - 性能剖析

**开发工具：**
- GORM 代码生成器
- 处理器代码生成器
- 代码格式化工具

**数据库：**
- MySQL - 主数据库
- Redis - 缓存和会话

## 项目架构

### 整体架构设计
项目采用分层架构设计，清晰的模块划分确保代码的可维护性和可扩展性：

```
客户端请求 → 路由层 → 中间件层 → 控制器层 → 服务层 → 数据访问层 → 数据库/缓存
```

### 核心模块关系
- **API模块** (`internal/api/`)：处理 HTTP 请求，包含业务逻辑
- **服务层** (`internal/services/`)：封装核心业务逻辑
- **数据层** (`internal/repository/`)：提供数据访问接口
- **GraphQL模块** (`internal/graph/`)：提供 GraphQL 查询能力
- **配置管理** (`configs/`)：统一管理应用配置

### 数据流向
1. 请求通过路由分发到对应的控制器
2. 控制器调用服务层处理业务逻辑
3. 服务层通过数据层访问数据库和缓存
4. 处理结果通过统一的响应格式返回

## 目录结构

```
go-gin-api/
├── cmd/                    # 命令行工具
│   ├── gormgen/           # GORM 模型代码生成器
│   ├── handlergen/        # 处理器代码生成器
│   ├── mfmt/              # 代码格式化工具
│   └── mysqlmd/           # MySQL 元数据处理
├── configs/               # 配置文件
│   ├── configs.go         # 配置结构定义
│   ├── constants.go       # 项目常量定义
│   └── *_configs.toml     # 多环境配置文件
├── internal/              # 内部模块（不对外暴露）
│   ├── api/               # API 接口模块
│   │   ├── admin/         # 管理员管理（登录、权限、用户管理）
│   │   ├── authorized/    # API 授权管理
│   │   ├── config/        # 系统配置管理
│   │   ├── cron/          # 定时任务管理
│   │   ├── menu/          # 菜单权限管理
│   │   └── tool/          # 工具类接口
│   ├── graph/             # GraphQL 相关模块
│   │   ├── generated/     # 自动生成的 GraphQL 代码
│   │   ├── model/         # GraphQL 数据模型
│   │   └── resolvers/     # GraphQL 解析器
│   ├── alert/             # 告警模块
│   ├── code/              # 错误码定义
│   └── repository/        # 数据访问层
│       ├── mysql/         # MySQL 数据访问
│       └── redis/         # Redis 数据访问
├── assets/                # 静态资源
│   ├── templates/         # HTML 模板文件
│   └── bootstrap/         # 前端资源文件
├── deployments/           # 部署配置
│   ├── loki/              # 日志收集配置
│   └── prometheus/        # 监控配置
├── docs/                  # 文档相关
└── pkg/                   # 公共工具包
```

## 核心文件说明

### 项目入口和配置文件

**main.go** - 应用主入口
- 初始化日志系统和 HTTP 服务器
- 配置优雅关闭机制
- 启动 Prometheus 监控和 Swagger 文档

**configs/configs.go** - 配置管理
- 使用 Viper 管理多环境配置
- 支持配置热更新
- 统一的配置结构定义

**configs/constants.go** - 项目常量
- 定义项目版本、端口等常量
- 配置 Redis key 前缀和超时时间

### 核心业务逻辑实现

**internal/api/admin/func_login.go** - 管理员登录逻辑
- 处理用户认证和密码验证
- 生成登录 token 并存储会话信息
- 管理用户权限和菜单访问控制

**internal/services/admin/service.go** - 管理员服务层
- 封装管理员相关的业务逻辑
- 提供用户管理、权限控制等服务
- 实现服务层和数据层的分离

### 数据模型和API接口

**internal/graph/model/generated.go** - GraphQL 数据模型
- 自动生成的 GraphQL 类型定义
- 支持用户查询和更新操作

**internal/graph/resolvers/user.go** - GraphQL 解析器
- 实现 GraphQL 查询和变更操作
- 提供用户数据查询和更新功能

### 关键组件和服务模块

**internal/pkg/core/core.go** - 核心框架
- 封装 Gin 框架，提供统一的中间件机制
- 实现请求追踪、错误处理、日志记录
- 支持限流、跨域、监控等通用功能

**internal/repository/mysql/mysql.go** - 数据访问层
- 封装 GORM 数据库操作
- 支持读写分离和连接池配置
- 提供统一的数据访问接口

**cmd/gormgen/main.go** - 代码生成器
- 自动生成 GORM 模型代码
- 支持模板化的代码生成
- 提高开发效率和代码规范性

### 路由和中间件

**internal/router/router.go** - 路由配置
- 统一的路由注册和管理
- 支持 API、GraphQL、WebSocket 等多种路由
- 集成中间件和拦截器机制

**internal/router/interceptor/interceptor.go** - 拦截器
- 提供登录验证、权限控制、签名验证等功能
- 支持 RBAC 权限管理
- 统一的认证和授权机制

## 快速开始

### 环境要求
- Go 1.16+
- MySQL 5.7+
- Redis 5.0+

### 启动步骤
1. 配置数据库连接信息
2. 运行 `go run main.go` 启动应用
3. 访问 `http://127.0.0.1:9999` 查看管理界面
4. 访问 `http://127.0.0.1:9999/swagger/index.html` 查看 API 文档

### 开发指南
项目采用模块化设计，新功能的开发应遵循现有的架构模式：
- 在 `internal/api/` 下创建新的 API 模块
- 在 `internal/services/` 下实现业务逻辑
- 使用代码生成器快速创建基础代码
- 遵循统一的错误处理和日志记录规范

该项目设计精良，模块清晰，适合作为企业级 API 项目的起点，也适合学习和研究现代 Go 语言 Web 开发的最佳实践。
        
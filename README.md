# 酒店行李智能管理系统

## 项目简介

本项目主要实现了行李的数字化管理，负责后端方面的数据库CRUD操作。项目基于 Go + Gin 框架构建，采用 MySQL 数据库持久化存储，通过 GORM ORM 框架实现了 10 张核心数据表的创建和管理，包括员工、角色、权限、行李寄存、行李明细、标签、照片、位置、客户和酒店等表。设计了完整的 RESTful API 接口，实现了行李的增删改查功能，支持创建行李寄存记录、分页查询、更新行李信息、取件校验和删除记录等操作。同时集成了 Redis 缓存系统，实现了会话管理和数据缓存，并结合腾讯云 COS 对象存储处理照片上传，为前端提供了稳定、高效的数据服务接口，实现了行李从入库到出库的完整数字化管理流程。

## 技术栈

| 类别 | 技术 | 版本 | 用途 |
|------|------|------|------|
| 后端框架 | Gin | - | Web API 开发 |
| 数据库 | MySQL | 8.0 | 持久化存储 |
| ORM | GORM | - | 数据库操作 |
| 缓存 | Redis | 7 | 会话管理、限流、缓存 |
| 认证 | JWT | - | 双令牌机制 |
| 限流 | 令牌桶算法 | - | API 访问控制 |
| 云存储 | 腾讯云 COS | - | 照片存储 |
| 容器化 | Docker + Compose | - | 部署管理 |
| 日志 | Logrus + Lumberjack | - | 日志记录与轮转 |

## 项目结构

```
hotel/
├── api/                      # API 路由层
│   └── api.go               # 路由注册和中间件配置
├── internal/                 # 内部包（不对外暴露）
│   ├── admin/               # 管理员功能模块
│   │   ├── add.go          # 添加操作（角色、权限、员工等）
│   │   ├── change.go       # 修改操作（员工角色等）
│   │   ├── delete.go       # 删除操作（角色、权限、员工等）
│   │   └── get.go          # 查询操作（员工、角色、权限等）
│   ├── employee/            # 员工功能模块
│   │   ├── employee_action/ # 员工业务操作
│   │   │   ├── add.go      # 添加行李、标签、位置等
│   │   │   ├── delete.go   # 删除行李、标签、位置等
│   │   │   ├── get.go      # 查询行李、标签、位置等
│   │   │   ├── update.go   # 更新行李、标签、位置等
│   │   │   ├── photo.go    # 照片上传下载
│   │   │   └── count.go    # 统计查询
│   │   └── employee_check/ # 员工认证
│   │       └── employee_check.go # 登录、登出、注册、刷新令牌
│   ├── middleware/          # 中间件
│   │   └── middleware.go    # 限流、JWT、权限、日志、异常恢复
│   ├── table/               # 数据库表
│   │   └── table.go        # 表结构定义和初始化
│   └── util/                # 工具类
│       ├── config.go       # 配置管理（YAML）
│       ├── jwt.go          # JWT 令牌生成和验证
│       ├── limiter.go      # 限流器（令牌桶算法）
│       ├── logger.go       # 日志系统（Logrus）
│       ├── crud.go         # 通用 CRUD 操作
│       ├── page.go         # 分页查询
│       └── util.go         # 通用工具函数
├── models/                  # 数据模型
│   └── models.go           # 数据模型定义（10张表）
├── services/                # 服务层
│   ├── services.go         # 数据库和 Redis 连接管理
│   └── tencent.go          # 腾讯云 COS 服务
├── configs/                 # 配置文件
│   ├── dev.yaml            # 开发环境配置
│   └── prod.yaml           # 生产环境配置
├── logs/                    # 日志文件目录
│   └── app.log             # 应用日志
├── uploads/                 # 上传文件目录
│   └── *.png               # 上传的图片文件
├── main.go                  # 应用入口
├── go.mod                   # Go 模块依赖
├── go.sum                   # 依赖版本锁定
├── docker-compose.yml       # Docker Compose 编排
├── Dockerfile               # Docker 镜像构建
└── apifox_import.json       # API 文档导入文件
```

## 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                        前端应用层                             │
│                      (Web / Mobile)                           │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                    API 网关层 (Gin + CORS)                     │
├─────────────────────────────────────────────────────────────┤
│  限流  │  日志  │  JWT认证  │  权限检查  │  异常恢复           │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                     业务逻辑层                                 │
├─────────────────────────────────────────────────────────────┤
│  员工管理  │  行李管理  │  照片管理  │  标签管理              │
│  位置管理  │  酒店管理  │  客户管理  │  后台管理              │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌──────────────────┬──────────────────┬───────────────────────┐
│   MySQL数据库    │    Redis缓存     │    腾讯云COS          │
│  (10张核心表)    │  (5个数据库实例)  │   (照片对象存储)       │
└──────────────────┴──────────────────┴───────────────────────┘
```

## 数据库设计

### 核心数据表

| 表名 | 用途 | 主要字段 |
|------|------|----------|
| Employee | 员工表 | 用户名、密码、角色ID、酒店ID |
| Role | 角色表 | 角色名称、描述 |
| Permission | 权限表 | 权限类型、描述 |
| LuggageStorage | 行李寄存表 | 寄存单号、客户ID、状态 |
| Luggage | 行李明细表 | 行李类型、描述、取件码 |
| Tag | 标签表 | 标签类型(RFID/二维码)、状态 |
| Photo | 照片表 | 照片URL、行李ID |
| Location | 位置表 | 位置名称、酒店ID |
| Guest | 客户表 | 客户姓名、联系方式 |
| Hotel | 酒店表 | 酒店名称、地址 |

### Redis 数据分离设计

```
DB 0: access_token   - 访问令牌存储（15分钟有效期）
DB 1: refresh_token  - 刷新令牌存储（7天有效期）
DB 2: cache          - 业务缓存
DB 3: rate_limit     - 限流计数器
DB 4: random         - 随机数/临时数据
```

## 核心功能

### 1. 身份认证与授权

- **JWT 双令牌机制**
  - Access Token: 15分钟有效期，存储在 Redis DB0
  - Refresh Token: 7天有效期，存储在 Redis DB1
  - 刷新机制: 自动续期，无感体验

- **RBAC 权限模型**
  - 基于角色的访问控制
  - 6种权限类型: 查看/创建/更新/删除行李、管理员、内部接口
  - 2种角色: 员工、管理员

### 2. 行李管理

- ✅ 创建行李寄存记录
- ✅ 查询行李（支持分页）
- ✅ 更新行李信息
- ✅ 行李取件（校验取件码）
- ✅ 删除行李记录

### 3. 标签管理

- 支持 RFID/二维码标签
- 标签绑定行李
- 标签状态管理

### 4. 照片管理

- 照片上传至腾讯云 COS
- 照片与行李关联
- 照片回传行李记录

### 5. 位置与酒店管理

- 多酒店支持
- 多位置（房间/柜子）管理
- 酒店员工绑定

## 中间件链

```
请求 → Recovery → RequestID → RateLimit 
      → JwtCheck → AuthCheck → PermissionCheck → 业务处理
```

- **Recovery**: 全局异常恢复，捕获 panic 并返回 500 错误
- **RequestID**: 生成唯一请求 ID，便于日志追踪
- **LogRequest**: 记录请求日志，包括请求方法、路径、状态码、耗时、操作员等
- **RateLimit**: 基于令牌桶算法的限流，防止 API 滥用
- **JwtCheck**: JWT 令牌验证，解析用户信息
- **AuthCheck**: 权限查询，从数据库加载用户权限
- **CheckAction**: 权限验证，检查用户是否有权访问特定接口
- **CORS**: 跨域请求处理

## 快速开始

### 环境要求

- Go 1.25.3+
- Docker & Docker Compose
- MySQL 8.0
- Redis 7

### 配置文件

编辑 `configs/dev.yaml` 文件，配置数据库、Redis、JWT、腾讯云 COS 等信息：

```yaml
server:
  mode: debug
  port: 8080

database:
  host: localhost
  port: 3306
  username: root
  password: password
  dbname: hotel

redis:
  addr: localhost:6379
  password: ""
  db: 0

jwt:
  secret_key: your-secret-key
  access_token_duration: 15m
  refresh_token_duration: 168h

tencent:
  secret_id: your-secret-id
  secret_key: your-secret-key
  bucket: your-bucket
  region: ap-guangzhou
```

### Docker 部署

一键启动所有服务：

```bash
docker-compose up -d
```

服务列表：
- MySQL (3306) - 数据库
- Redis (6379) - 缓存
- Application (8080) - 主应用

### 本地开发

```bash
# 安装依赖
go mod download

# 运行应用
go run main.go
```

## API 接口

### 公共接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /employee/login | 员工登录 |
| POST | /employee/logout | 员工登出 |
| POST | /employee/register | 员工注册 |
| POST | /employee/refresh | 刷新令牌 |

### 内部接口（需认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /internal/luggage | 创建行李寄存 |
| GET | /internal/luggage | 查询行李（分页） |
| PUT | /internal/luggage | 更新行李信息 |
| DELETE | /internal/luggage | 删除行李记录 |
| POST | /internal/luggage/pickup | 行李取件 |

### 管理接口（需管理员权限）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /tool/employee | 查询员工列表 |
| POST | /tool/employee | 添加员工 |
| PUT | /tool/employee | 更新员工信息 |
| DELETE | /tool/employee | 删除员工 |
| GET | /tool/role | 查询角色列表 |
| POST | /tool/role | 添加角色 |
| GET | /tool/permission | 查询权限列表 |
| POST | /tool/permission | 添加权限 |

## 技术亮点

### 1. 智能限流机制
- 基于令牌桶算法
- 支持多位置差异化限流
- 默认: 容量 100，填充速率 10ms

### 2. 完善的日志系统
- 结构化日志（Logrus）
- 支持控制台 + 文件双输出
- 自动日志轮转（最大 100MB，保留 3 个，28 天）

### 3. 安全防护
- 密码 bcrypt 加密
- JWT 令牌黑名单机制
- CORS 跨域配置
- 全局异常捕获

### 4. 容器化部署
- Docker Compose 一键启动
- 健康检查机制
- 数据持久化卷挂载

### 5. 配置化管理
- YAML 配置文件
- 环境区分（dev/test/prod）
- 配置热加载

## 开发思路

采用分层架构设计，遵循"关注点分离"原则：

1. **整体架构设计**: API 层、业务逻辑层、服务层、数据层
2. **核心技术选型**: Gin + MySQL + Redis + JWT + 腾讯云 COS
3. **开发流程**:
   - 配置管理
   - 基础设施搭建
   - 数据库设计
   - 中间件开发
   - 业务功能实现
   - 接口设计
4. **安全设计**: 密码加密、JWT 双令牌、令牌黑名单、CORS 配置、限流保护
5. **性能优化**: Redis 缓存、连接池、日志轮转、容器化部署
6. **可扩展性设计**: 中间件模式、RBAC 权限、配置化、模块化
7. **部署策略**: Docker 容器化、Docker Compose 编排、健康检查机制、数据卷持久化

## 项目亮点与收获

### 技术收获
- ✅ 掌握 Go + Gin 框架开发
- ✅ 实现 JWT 双令牌认证机制
- ✅ 设计 RBAC 权限系统
- ✅ 集成腾讯云 COS 对象存储
- ✅ 实现分布式限流算法
- ✅ 掌握 Docker 容器化部署

### 工程实践
- ✅ 代码分层清晰（API/Service/Model）
- ✅ 中间件复用性强
- ✅ 配置文件统一管理
- ✅ 日志系统完善
- ✅ 错误处理规范

## 未来改进方向

1. **前端开发**: 开发 Web/移动端管理界面
2. **功能扩展**:
   - 行李历史追踪
   - 异常报警功能
   - 数据统计报表
   - 多语言支持
3. **性能优化**:
   - 数据库读写分离
   - Redis 集群部署
   - API 响应缓存
4. **安全增强**:
   - HTTPS 部署
   - 操作审计日志
   - 敏感数据加密

## 许可证

MIT License

## 联系方式

如有问题或建议，欢迎提出 Issue 或 Pull Request。

# Design Server 后端服务

## 项目简介

Design Server 是一个基于 Gin + Vue 的全栈开发基础平台，提供完整的系统管理功能和门户网站功能。项目采用现代化的技术栈，支持快速开发和部署。

## 技术栈

### 后端技术
- **框架**: Gin (Go Web框架)
- **数据库**: MySQL + GORM
- **缓存**: Redis
- **认证**: JWT
- **权限**: Casbin
- **日志**: Zap
- **配置**: Viper
- **文档**: Swagger
- **定时任务**: Cron
- **WebSocket**: Gorilla WebSocket

### 核心依赖
- `github.com/gin-gonic/gin` - Web框架
- `gorm.io/gorm` - ORM框架
- `github.com/casbin/casbin/v2` - 权限管理
- `go.uber.org/zap` - 日志库
- `github.com/spf13/viper` - 配置管理
- `github.com/swaggo/swag` - API文档

## 项目结构

```
backend/
├── api/                    # API接口层
│   └── v1/                # API版本
│       ├── system/        # 系统管理接口
│       ├── portal/        # 门户网站接口
│       └── example/       # 示例接口
├── config/                # 配置文件
│   ├── config.go          # 主配置
│   ├── system.go          # 系统配置
│   ├── gorm_mysql.go      # 数据库配置
│   ├── jwt.go             # JWT配置
│   ├── cors.go            # CORS配置
│   ├── captcha.go         # 验证码配置
│   ├── oss_local.go       # 文件存储配置
│   └── zap.go             # 日志配置
├── core/                  # 核心功能
├── docs/                  # API文档
├── global/                # 全局变量
├── initialize/            # 初始化模块
├── middleware/            # 中间件
├── model/                 # 数据模型
│   ├── system/           # 系统模型
│   ├── portal/           # 门户模型
│   ├── example/          # 示例模型
│   └── common/           # 公共模型
├── plugin/               # 插件
├── router/               # 路由配置
│   ├── system/           # 系统路由
│   ├── portal/           # 门户路由
│   └── example/          # 示例路由
├── service/              # 业务逻辑层
│   ├── system/           # 系统服务
│   ├── portal/           # 门户服务
│   └── example/          # 示例服务
├── source/               # 资源文件
├── task/                 # 定时任务
├── uploads/              # 上传文件目录
├── utils/                # 工具函数
├── packfile/             # 打包文件
├── main.go               # 程序入口
├── go.mod                # Go模块文件
├── go.sum                # 依赖校验文件
└── config.yaml           # 配置文件
```

## 功能模块

### 系统管理模块 (system)
- **用户管理**: 用户CRUD、角色分配、状态管理
- **权限管理**: 基于Casbin的RBAC权限控制
- **菜单管理**: 动态菜单配置
- **API管理**: 接口权限配置
- **字典管理**: 系统字典维护
- **操作日志**: 用户操作记录
- **JWT黑名单**: Token失效管理
- **验证码**: 图形验证码生成
- **仪表板**: 系统统计信息

### 门户网站模块 (portal)
- **文章管理**: 文章CRUD、分类、标签
- **分类管理**: 文章分类维护
- **标签管理**: 文章标签维护
- **主题管理**: 网站主题配置
- **消息管理**: 系统消息通知

## 安装和运行

### 环境要求
- Go 1.23+
- MySQL 8.0+
- Redis 6.0+

### 安装步骤

1. **克隆项目**
```bash
git clone <repository-url>
cd backend
```

2. **安装依赖**
```bash
go mod tidy
go mod download
```

3. **配置数据库**
- 创建MySQL数据库
- 修改 `config.yaml` 中的数据库连接信息

4. **运行项目**
```bash
go run main.go
```

## API文档

启动服务后，访问以下地址查看API文档：
- Swagger UI: `http://localhost:8888/swagger/index.html`
- API文档: `http://localhost:8888/swagger/doc.json`

## 开发指南

### 添加新模块

1. **创建数据模型** (`model/`)
2. **创建业务逻辑** (`service/`)
3. **创建API接口** (`api/v1/`)
4. **配置路由** (`router/`)
5. **更新初始化** (`initialize/`)

### 代码规范

- 遵循Go语言官方代码规范
- 使用 `gofmt` 格式化代码
- 添加必要的注释和文档
- 编写单元测试

## 部署

### Docker部署

```bash
# 构建镜像
docker build -t design-server .

# 运行容器
docker run -d -p 8888:8888 design-server
```

### 生产环境配置

1. 修改 `config.yaml` 中的环境配置
2. 配置生产环境数据库
3. 设置日志级别
4. 配置反向代理

## 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交代码
4. 创建 Pull Request

## 许可证

本项目采用 MIT 许可证。

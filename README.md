# LingJian 低代码平台

LingJian是一个开源的低代码平台，支持多应用管理、可视化配置、权限控制等功能。

## 功能特点

- 多应用管理：支持创建和管理多个应用，共享用户系统
- RBAC权限控制：细粒度的权限管理，支持角色和权限的动态配置
- 可视化配置：支持数据表、维度、数据模型、表单等的可视化配置
- 定时任务：支持配置SQL和HTTP类型的定时任务
- 版本管理：配置信息支持版本控制，方便追踪和回滚
- 模板商城：支持应用模板的分享和复用

## 技术架构

### 后端技术栈

- 语言：Go
- Web框架：Gin
- 数据库：MySQL 8.0
- 缓存：Redis
- 消息队列：RabbitMQ
- API文档：Swagger
- 配置管理：Viper

### 前端技术栈

- 框架：React
- UI组件：Ant Design
- 状态管理：Redux
- 拖拽组件：dnd-kit

## 系统架构

```
├── api                 # API接口层
│   └── v1             # V1版本接口
│       ├── auth       # 认证相关
│       ├── config     # 配置管理
│       ├── element    # 元素管理
│       ├── rbac       # 权限管理
│       └── task       # 任务管理
├── cmd                # 主程序入口
│   ├── server        # API服务器
│   └── worker        # 任务处理器
├── config            # 配置文件
├── docs              # API文档
├── internal          # 内部包
│   ├── middleware    # 中间件
│   ├── model        # 数据模型
│   ├── service      # 业务逻辑
│   └── test         # 测试用例
└── pkg              # 公共包
    ├── queue        # 消息队列
    ├── redis        # Redis工具
    ├── store        # 存储接口
    └── utils        # 工具函数
```
```

## 数据库设计

### 核心表结构

- users：用户表
- roles：角色表
- permissions：权限表
- applications：应用表
- config_tables：数据表配置
- config_dimensions：维度配置
- config_data_models：数据模型配置
- config_forms：表单配置
- config_menus：菜单配置
- scheduled_tasks：定时任务表

详细的数据库设计请参考 `internal/model/schema_*.sql` 文件。

## 快速开始

### 环境要求

- Go 1.20+
- MySQL 8.0+
- Redis 6.0+
- RabbitMQ 3.8+

### 安装步骤

1. 克隆代码
```bash
git clone https://github.com/iiwish/lingjian.git
cd lingjian
```

2. 安装依赖
```bash
make deps
```

3. 初始化数据库
```bash
# 修改 Makefile 中的数据库配置
make init-db
```

4. 修改配置文件
```bash
cp config/config.yaml.example config/config.yaml
# 编辑 config.yaml，配置数据库、Redis、RabbitMQ等信息
```

5. 启动服务
```bash
# 启动前需要先启动redis、rabbitmq、mysql服务
# 开发模式启动服务器（支持热重载）
make dev-server

# 或者正常启动
make run-server  # 启动API服务器
make run-worker  # 启动任务处理器
```

### 默认账号

在执行 `make init-db` 初始化数据库时，系统会自动创建以下默认账号：

- 超级管理员
  - 用户名：admin
  - 密码：admin1324
  - 权限：系统所有权限

请在首次登录后及时修改默认密码。

### 开发命令

```bash
# 构建项目
make build

# 运行测试
make test

# 代码检查
make lint

# 格式化代码
make fmt

# 生成API文档
make swagger

# 查看所有可用命令
make help
```


## API文档

### 认证相关

#### 登录
```
POST /api/v1/auth/login
Content-Type: application/json

{
    "username": "admin",
    "password": "password"
}
```

#### 刷新Token
```
POST /api/v1/auth/refresh
X-Refresh-Token: {refresh_token}
```

### 应用管理

#### 创建应用
```
POST /api/v1/applications
Content-Type: application/json

{
    "name": "测试应用",
    "code": "test_app",
    "description": "这是一个测试应用"
}
```

### 配置管理

#### 创建数据表配置
```
POST /api/v1/config/tables
Content-Type: application/json

{
    "application_id": 1,
    "name": "用户信息表",
    "code": "user_info",
    "mysql_table_name": "user_info",
    "fields": [
        {
            "name": "id",
            "type": "int",
            "required": true
        },
        {
            "name": "name",
            "type": "varchar",
            "length": 50,
            "required": true
        }
    ]
}
```

更多API文档请参考 Swagger 文档。

## 开发指南

### 目录结构说明

- `api/`: API接口定义和实现
- `cmd/`: 主程序入口
- `config/`: 配置文件
- `docs/`: 项目文档
- `internal/`: 内部包
  - `middleware/`: 中间件
  - `model/`: 数据模型
  - `service/`: 业务逻辑
- `pkg/`: 公共包
  - `queue/`: 消息队列
  - `utils/`: 工具函数

### 开发规范

1. 代码风格
   - 使用gofmt格式化代码
   - 遵循Go语言编码规范
   - 使用英文注释

2. 错误处理
   - 使用统一的错误响应格式
   - 记录关键错误日志
   - 避免panic

3. 数据库操作
   - 使用事务确保数据一致性
   - 合理使用索引
   - 避免大事务

4. API设计
   - 遵循RESTful规范
   - 版本控制
   - 统一的响应格式

### 测试

1. 单元测试
```bash
go test ./...
```

2. API测试
```bash
# 使用 Postman 导入 docs/postman_collection.json
```

## 常见问题

1. 启动失败
   - 检查配置文件是否正确
   - 检查数据库连接是否正常
   - 检查端口是否被占用

2. 权限问题
   - 确认用户角色是否配置正确
   - 检查权限表是否正确配置
   - 查看日志了解具体错误

## 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交代码
4. 创建 Pull Request

## 许可证

本项目采用 Apache-2.0 许可证 - 查看 [LICENSE](LICENSE) 文件了解详细信息。

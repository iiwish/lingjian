# LingJian 低代码平台

LingJian是一个基于Go语言开发的低代码平台，支持多租户、多应用管理，提供灵活的配置和强大的扩展能力。

## 技术栈

- 后端框架：Gin
- 数据库：MySQL 8.0
- 缓存：Redis
- 消息队列：RabbitMQ
- API文档：Swagger
- 配置管理：Viper

## 主要特性

- 多租户支持
  - 独立数据库模式
  - 租户级别配置
  - 数据隔离

- 多应用管理
  - 应用模板
  - 灵活配置
  - 权限隔离

- RBAC权限控制
  - 细粒度权限管理
  - 动态权限加载
  - 多级角色

- 用户认证
  - JWT认证
  - 双Token机制
  - OAuth2.0支持

## 项目结构

```
.
├── api/                    # API接口定义
├── cmd/                    # 主程序入口
│   └── server/            # 服务器启动
├── config/                # 配置文件
├── internal/              # 内部包
│   ├── model/            # 数据模型
│   ├── middleware/       # 中间件
│   ├── handler/         # HTTP处理器
│   ├── service/         # 业务逻辑
│   ├── repository/      # 数据访问层
│   └── server/          # 服务器初始化
├── pkg/                   # 公共包
│   ├── config/          # 配置管理
│   ├── database/        # 数据库连接
│   └── redis/           # Redis连接
└── scripts/              # SQL脚本
```

## 快速开始

1. 克隆项目
```bash
git clone https://github.com/iiwish/lingjian.git
```

2. 安装依赖
```bash
go mod download
```

3. 配置数据库
```bash
# 创建数据库和表
mysql -u root -p < scripts/schema.sql

# 初始化数据
mysql -u root -p < scripts/init_data.sql
```

4. 配置环境
```bash
# 复制配置文件
cp config/config.yaml.example config/config.yaml

# 修改配置
vim config/config.yaml
```

5. 运行项目
```bash
go run cmd/server/main.go
```

## 开发计划

- [x] 项目基础框架
- [x] 数据库设计
- [ ] 用户认证模块
- [ ] RBAC权限控制
- [ ] 多租户支持
- [ ] 多应用管理
- [ ] API文档
- [ ] 单元测试
- [ ] 部署文档

## 贡献指南

欢迎提交Issue和Pull Request。

## 许可证

MIT License

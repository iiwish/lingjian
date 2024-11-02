#!/bin/bash

# 显示颜色
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}开始设置开发环境...${NC}"

# 检查MySQL服务是否运行
if ! mysqladmin ping -h localhost --silent; then
    echo -e "${RED}错误: MySQL服务未运行${NC}"
    echo "请先启动MySQL服务"
    echo "在macOS上: brew services start mysql"
    echo "在Linux上: sudo service mysql start"
    exit 1
fi

# 检查Redis服务是否运行
if ! redis-cli ping >/dev/null 2>&1; then
    echo -e "${RED}错误: Redis服务未运行${NC}"
    echo "请先启动Redis服务"
    echo "在macOS上: brew services start redis"
    echo "在Linux上: sudo service redis start"
    exit 1
fi

# 创建配置文件
if [ ! -f "config/config.yaml" ]; then
    echo -e "${GREEN}创建配置文件...${NC}"
    cp config/config.yaml.example config/config.yaml
    echo "配置文件已创建，请根据需要修改 config/config.yaml"
fi

# 创建数据库和表
echo -e "${GREEN}创建数据库和表...${NC}"
mysql -u root -p < scripts/schema.sql
if [ $? -ne 0 ]; then
    echo -e "${RED}创建数据库和表失败${NC}"
    exit 1
fi

# 初始化数据
echo -e "${GREEN}初始化数据...${NC}"
mysql -u root -p < scripts/init_data.sql
if [ $? -ne 0 ]; then
    echo -e "${RED}初始化数据失败${NC}"
    exit 1
fi

# 创建日志目录
echo -e "${GREEN}创建日志目录...${NC}"
mkdir -p logs
mkdir -p uploads

echo -e "${GREEN}开发环境设置完成！${NC}"
echo "现在您可以："
echo "1. 修改 config/config.yaml 中的配置"
echo "2. 运行 'go run cmd/server/main.go' 启动服务器"
echo ""
echo "默认超级管理员账号："
echo "用户名: admin"
echo "密码: admin123"

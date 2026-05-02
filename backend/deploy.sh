#!/bin/bash

# 农产品溯源系统部署脚本

set -e

echo "========================================="
echo "农产品溯源系统 - 生产环境部署脚本"
echo "========================================="

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_ROOT"

echo -e "${GREEN}[1/6]${NC} 检查Python环境..."
if ! command -v python3 &> /dev/null; then
    echo -e "${RED}错误: 未找到Python3${NC}"
    exit 1
fi
echo "Python版本: $(python3 --version)"

echo -e "${GREEN}[2/6]${NC} 检查Node.js环境..."
if ! command -v node &> /dev/null; then
    echo -e "${RED}错误: 未找到Node.js${NC}"
    exit 1
fi
echo "Node.js版本: $(node --version)"
echo "npm版本: $(npm --version)"

echo -e "${GREEN}[3/6]${NC} 安装后端依赖..."
cd backend
pip3 install -r requirements.txt 2>/dev/null || pip install -r requirements.txt || echo "依赖可能已安装"

echo -e "${GREEN}[4/6]${NC} 执行数据库迁移..."
python3 manage.py migrate

echo -e "${GREEN}[5/6]${NC} 安装前端依赖并构建..."
cd ../frontend
npm install
npm run build

echo -e "${GREEN}[6/6]${NC} 收集静态文件..."
cd ../backend
python3 manage.py collectstatic --noinput

echo ""
echo -e "${GREEN}========================================="
echo -e "部署完成！"
echo "=========================================${NC}"
echo ""
echo "启动服务："
echo "  后端: cd backend && gunicorn backend.wsgi:application --bind 0.0.0.0:8000"
echo "  前端: cd frontend && npx serve -s dist -l 3000"
echo ""
echo "或使用Nginx配置反向代理（参见README.md）"

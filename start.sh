#!/bin/bash

# 农产品溯源系统启动脚本

echo "=========================================="
echo "农产品溯源系统"
echo "=========================================="

# 检查Python
if ! command -v python3 &> /dev/null; then
    echo "错误: 未找到Python3，请先安装Python 3.10+"
    exit 1
fi

# 检查Node.js
if ! command -v node &> /dev/null; then
    echo "错误: 未找到Node.js，请先安装Node.js"
    exit 1
fi

echo ""
echo "[1/5] 安装后端依赖..."
cd backend
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
cd ..

echo ""
echo "[2/5] 数据库迁移..."
cd backend
source venv/bin/activate
python manage.py makemigrations
python manage.py migrate
python manage.py createsuperuser --noinput --username admin --email admin@example.com || true
cd ..

echo ""
echo "[3/5] 安装前端依赖..."
cd frontend/trace-web
npm install
cd ../trace-admin
npm install
cd ../../

echo ""
echo "=========================================="
echo "安装完成！"
echo "=========================================="
echo ""
echo "启动方式："
echo ""
echo "终端1 - 启动后端服务："
echo "  cd backend"
echo "  source venv/bin/activate"
echo "  python manage.py runserver"
echo ""
echo "终端2 - 启动消费者端(开发模式)："
echo "  cd frontend/trace-web"
echo "  npm run dev"
echo ""
echo "终端3 - 启动管理端(开发模式)："
echo "  cd frontend/trace-admin"
echo "  npm run dev"
echo ""
echo "访问地址："
echo "  消费者端: http://localhost:5173"
echo "  管理端: http://localhost:3000"
echo "  后端API: http://localhost:8000"
echo ""
echo "管理后台账号: admin / admin123"
echo "=========================================="

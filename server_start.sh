#!/bin/bash

# 服务器启动脚本 - 农产品二维码追溯系统

echo "=== 农产品二维码追溯系统服务器启动脚本 ==="

# 检查项目目录是否存在
if [ ! -d "/opt/qrTraceability" ]; then
    echo "错误: 项目目录不存在"
    echo "请先运行 ubuntu_setup.sh 脚本进行初始化"
    exit 1
fi

cd /opt/qrTraceability

# 检查是否需要更新代码
echo "检查代码更新..."
git pull

# 重新构建前端项目
echo "重新构建前端项目..."
cd client
npm install
npm run build

# 重新编译后端项目
echo "重新编译后端项目..."
cd ../server

# 优化 Go 依赖管理
go mod tidy  # 清理和更新依赖
go mod verify  # 验证依赖完整性

# 优化编译参数（静态编译，减少依赖）
go build -ldflags="-s -w" -o server cmd/main.go

# 检查编译结果
if [ ! -f "server" ]; then
    echo "错误: 后端编译失败"
    exit 1
fi

chmod +x server

# 重启服务
echo "重启服务..."
sudo systemctl restart qrbackend
sudo systemctl restart nginx

# 检查服务状态
echo "检查服务状态..."
sudo systemctl status qrbackend --no-pager
sudo systemctl status nginx --no-pager

echo "服务启动完成！"
echo "前端地址: http://$(hostname -I | awk '{print $1}')"
echo "后端API地址: http://$(hostname -I | awk '{print $1}'):8080"

echo "=== 启动脚本执行完成 ==="

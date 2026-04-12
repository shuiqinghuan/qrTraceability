#!/bin/bash

# 服务器启动脚本

echo "=== 农产品二维码追溯系统服务器启动脚本 ==="

# 检查是否安装了Docker和Docker Compose
if ! command -v docker &> /dev/null; then
    echo "错误: Docker 未安装"
    echo "正在安装 Docker..."
    sudo apt update
    sudo apt install -y apt-transport-https ca-certificates curl gnupg-agent software-properties-common
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
    sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
    sudo apt update
    sudo apt install -y docker-ce docker-ce-cli containerd.io
    sudo usermod -aG docker $USER
    echo "Docker 安装完成，请重新登录后运行此脚本"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "错误: Docker Compose 未安装"
    echo "正在安装 Docker Compose..."
    sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    echo "Docker Compose 安装完成"
fi

# 检查是否已克隆项目
if [ ! -d "qrTraceability" ]; then
    echo "正在从 GitHub 克隆项目..."
    git clone https://github.com/shuiqinghuan/qrTraceability.git
    cd qrTraceability
else
    echo "项目已存在，正在更新..."
    cd qrTraceability
    git pull
fi

# 构建和启动服务
echo "正在构建和启动服务..."
docker-compose up -d --build

echo "服务启动完成！"
echo "前端地址: http://$(hostname -I | awk '{print $1}')"
echo "后端API地址: http://$(hostname -I | awk '{print $1}'):8080"
echo "数据库地址: http://$(hostname -I | awk '{print $1}'):5432"

echo "=== 启动脚本执行完成 ==="

#!/bin/bash

# 安装依赖
if [ ! -d "node_modules" ]; then
    echo "Installing dependencies..."
    npm install
fi

# 构建前端应用
echo "Building frontend..."
npm run build

# 启动静态文件服务器
echo "Starting frontend server..."
npx serve -s dist -l 3000
#!/bin/bash

# 设置环境变量
export DB_HOST=${DB_HOST:-localhost}
export DB_PORT=${DB_PORT:-5432}
export DB_USER=${DB_USER:-postgres}
export DB_PASSWORD=${DB_PASSWORD:-postgres}
export DB_NAME=${DB_NAME:-plantation}
export JWT_SECRET=${JWT_SECRET:-your-secret-key}
export SERVER_PORT=${SERVER_PORT:-8080}

# 检查是否存在可执行文件
if [ ! -f "./server" ]; then
    echo "Building server..."
    go build -o server ./cmd/main.go
fi

# 启动服务器
echo "Starting server on port $SERVER_PORT..."
./server
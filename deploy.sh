#!/bin/bash

# 部署脚本 - 农产品二维码追溯系统

# 显示帮助信息
show_help() {
    echo "部署脚本 - 农产品二维码追溯系统"
    echo ""
    echo "用法: ./deploy.sh [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help                显示帮助信息"
    echo "  -i, --ip <IP地址>         设置服务器IP地址"
    echo "  -p, --port <端口>          设置服务器端口（默认8080）"
    echo "  -d, --db-host <数据库主机>  设置数据库主机（默认localhost）"
    echo "  -u, --db-user <数据库用户>  设置数据库用户（默认qruser）"
    echo "  -w, --db-pass <数据库密码>  设置数据库密码（默认6180680）"
    echo "  -n, --db-name <数据库名>    设置数据库名（默认qrtraceability）"
    echo "  -j, --jwt-secret <密钥>    设置JWT密钥"
    echo "  -b, --build                构建前端项目"
    echo "  -s, --start                启动服务"
    echo "  -r, --restart              重启服务"
    echo ""
    echo "示例: ./deploy.sh -i 192.168.1.100 -p 8080 -b -s"
}

# 默认值
SERVER_IP=""
SERVER_PORT="8080"
DB_HOST="localhost"
DB_USER="qruser"
DB_PASS="6180680"
DB_NAME="qrtraceability"
JWT_SECRET="your_jwt_secret_key_for_production_change_this"
BUILD=false
START=false
RESTART=false

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -i|--ip)
            SERVER_IP="$2"
            shift 2
            ;;
        -p|--port)
            SERVER_PORT="$2"
            shift 2
            ;;
        -d|--db-host)
            DB_HOST="$2"
            shift 2
            ;;
        -u|--db-user)
            DB_USER="$2"
            shift 2
            ;;
        -w|--db-pass)
            DB_PASS="$2"
            shift 2
            ;;
        -n|--db-name)
            DB_NAME="$2"
            shift 2
            ;;
        -j|--jwt-secret)
            JWT_SECRET="$2"
            shift 2
            ;;
        -b|--build)
            BUILD=true
            shift
            ;;
        -s|--start)
            START=true
            shift
            ;;
        -r|--restart)
            RESTART=true
            shift
            ;;
        *)
            echo "未知参数: $1"
            show_help
            exit 1
            ;;
    esac
done

# 检查必要参数
if [ -z "$SERVER_IP" ]; then
    echo "错误: 必须指定服务器IP地址"
    show_help
    exit 1
fi

echo "=== 农产品二维码追溯系统部署脚本 ==="
echo "服务器IP: $SERVER_IP"
echo "服务器端口: $SERVER_PORT"
echo "数据库主机: $DB_HOST"
echo "数据库用户: $DB_USER"
echo "数据库名: $DB_NAME"

# 构建前端项目
if [ "$BUILD" = true ]; then
    echo ""
    echo "=== 构建前端项目 ==="
    cd client
    
    # 安装依赖
    echo "安装前端依赖..."
    npm install
    
    # 构建项目
    echo "构建前端项目..."
    npm run build
    
    if [ $? -ne 0 ]; then
        echo "错误: 前端构建失败"
        exit 1
    fi
    
    echo "前端构建成功"
    cd ..
fi

# 编译后端项目
echo ""
echo "=== 编译后端项目 ==="
cd server

# 清理和更新依赖
echo "更新Go依赖..."
go mod tidy
go mod verify

# 编译后端
echo "编译后端项目..."
go build -ldflags="-s -w" -o server cmd/main.go

if [ $? -ne 0 ]; then
    echo "错误: 后端编译失败"
    exit 1
fi

chmod +x server
echo "后端编译成功"
cd ..

# 创建环境变量文件
echo ""
echo "=== 创建环境变量配置 ==="
cat > .env << EOF
# 服务器配置
SERVER_PORT=$SERVER_PORT

# 数据库配置
DB_HOST=$DB_HOST
DB_PORT=5432
DB_USER=$DB_USER
DB_PASSWORD=$DB_PASS
DB_NAME=$DB_NAME

# JWT配置
JWT_SECRET=$JWT_SECRET

# 前端配置
FRONTEND_URL=http://$SERVER_IP
API_URL=http://$SERVER_IP:$SERVER_PORT/api
EOF

echo "环境变量配置创建成功"

# 启动服务
if [ "$START" = true ] || [ "$RESTART" = true ]; then
    echo ""
    echo "=== 启动服务 ==="
    
    # 停止现有服务（如果运行）
    if [ "$RESTART" = true ]; then
        echo "停止现有服务..."
        pkill -f "server" 2>/dev/null
        sleep 2
    fi
    
    # 启动后端服务
    echo "启动后端服务..."
    cd server
    
    # 使用nohup运行服务
    nohup ./server > server.log 2>&1 &
    
    # 检查服务是否启动成功
    sleep 3
    if pgrep -f "server" > /dev/null; then
        echo "后端服务启动成功"
    else
        echo "错误: 后端服务启动失败"
        cat server.log
        exit 1
    fi
    
    cd ..
    
    # 配置Nginx（如果存在）
    if command -v nginx &> /dev/null; then
        echo "配置Nginx..."
        cat > /etc/nginx/sites-available/qr-traceability << EOF
server {
    listen 80;
    server_name $SERVER_IP;
    
    root /opt/qrTraceability/client/dist;
    index index.html;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    location /api {
        proxy_pass http://localhost:$SERVER_PORT;
        proxy_set_header Host host;
        proxy_set_header X-Real-IP remote_addr;
        proxy_set_header X-Forwarded-For proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto scheme;
    }
}
EOF
        
        # 创建符号链接
        if [ ! -f /etc/nginx/sites-enabled/qr-traceability ]; then
            ln -s /etc/nginx/sites-available/qr-traceability /etc/nginx/sites-enabled/
        fi
        
        # 测试Nginx配置
        nginx -t
        
        # 重启Nginx
        systemctl restart nginx
        echo "Nginx配置成功"
    else
        echo "Nginx未安装，跳过Nginx配置"
    fi
fi

echo ""
echo "=== 部署完成 ==="
echo "前端地址: http://$SERVER_IP"
echo "后端API地址: http://$SERVER_IP:$SERVER_PORT/api"
echo "后台管理地址: http://$SERVER_IP/admin/login"
echo "默认管理员账号: lhseed"
echo "默认管理员密码: 123456"
echo ""
echo "部署脚本执行完成！"

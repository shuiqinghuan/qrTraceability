#!/bin/bash

# Ubuntu 服务器配置脚本 - 农产品二维码追溯系统

echo "=== 农产品二维码追溯系统 - Ubuntu 服务器配置脚本 ==="

echo "\n1. 更新系统包..."
sudo apt update && sudo apt upgrade -y

echo "\n2. 安装必要的依赖包..."
sudo apt install -y git curl wget unzip build-essential nginx

echo "\n3. 安装 Go (使用 apt 包管理器)..."
sudo apt install -y golang-go

echo "\n4. 安装 Node.js 20.x (使用 apt 包管理器)..."
sudo apt install -y ca-certificates curl gnupg
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg
NODE_MAJOR=20
echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$NODE_MAJOR.x nodistro main" | sudo tee /etc/apt/sources.list.d/nodesource.list
sudo apt update
sudo apt install -y nodejs

echo "\n5. 安装 PostgreSQL (使用 apt 包管理器)..."
sudo apt install -y postgresql postgresql-contrib

# 启动 PostgreSQL 服务
sudo systemctl start postgresql
sudo systemctl enable postgresql

echo "\n6. 配置数据库..."
sudo -u postgres psql -c "CREATE USER qruser WITH PASSWORD 'qrpassword';"
sudo -u postgres psql -c "CREATE DATABASE qrtraceability OWNER qruser;"

echo "\n7. 创建项目目录..."
sudo mkdir -p /opt/qrTraceability
sudo chown $USER:$USER /opt/qrTraceability
cd /opt/qrTraceability

echo "\n8. 从 GitHub 拉取项目代码..."
git clone https://github.com/shuiqinghuan/qrTraceability.git .

echo "\n9. 构建前端项目..."
cd client
npm install
npm run build

echo "\n10. 编译后端项目..."
cd ../server

# 优化 Go 依赖管理
go mod tidy  # 清理和更新依赖
go mod verify  # 验证依赖完整性

# 优化编译参数（静态编译，减少依赖）
go build -ldflags="-s -w" -o server cmd/main.go

# 检查编译结果
if [ -f "server" ]; then
    echo "后端编译成功"
    chmod +x server
else
    echo "后端编译失败"
    exit 1
fi

echo "\n11. 配置环境变量..."
echo "DB_HOST=localhost\nDB_PORT=5432\nDB_USER=qruser\nDB_PASSWORD=qrpassword\nDB_NAME=qrtraceability\nJWT_SECRET=your_jwt_secret\nSERVER_PORT=8080" > .env

echo "\n12. 执行数据库模式脚本..."
psql -U qruser -d qrtraceability -f internal/db/schema.sql

echo "\n13. 创建 Systemd 服务..."
sudo bash -c 'cat > /etc/systemd/system/qrbackend.service << EOF
[Unit]
Description=QR Traceability Backend Service
After=network.target postgresql.service

[Service]
WorkingDirectory=/opt/qrTraceability/server
EnvironmentFile=/opt/qrTraceability/server/.env
ExecStart=/opt/qrTraceability/server/server
Restart=always
RestartSec=5
User=ubuntu

[Install]
WantedBy=multi-user.target
EOF'

echo "\n14. 配置 Nginx 反向代理..."
sudo bash -c 'cat > /etc/nginx/sites-available/qrtraceability << EOF
server {
    listen 80;
    server_name $(hostname -I | awk "{print $1}");

    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location / {
        root /opt/qrTraceability/client/dist;
        index index.html;
        try_files $uri $uri/ /index.html;
    }
}
EOF'

echo "\n15. 启动服务..."
sudo systemctl daemon-reload
sudo systemctl start qrbackend
sudo systemctl enable qrbackend
sudo ln -s /etc/nginx/sites-available/qrtraceability /etc/nginx/sites-enabled/
sudo systemctl restart nginx

echo "\n16. 配置防火墙..."
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 8080/tcp
sudo ufw allow 5432/tcp
sudo ufw --force enable

# 创建 IP 地址替换脚本（钩子）
echo '#!/bin/bash

# 服务器 IP 地址替换脚本
SERVER_IP=$1

if [ -z "$SERVER_IP" ]; then
    echo "请提供服务器 IP 地址"
    echo "用法: $0 <服务器IP地址>"
    exit 1
fi

# 替换 Nginx 配置中的 server_name
sudo sed -i "s/server_name .*/server_name $SERVER_IP;/" /etc/nginx/sites-available/qrtraceability

# 重启 Nginx
sudo systemctl restart nginx

echo "服务器 IP 地址已更新为: $SERVER_IP"
echo "Nginx 已重启"
' > /opt/qrTraceability/update_server_ip.sh

chmod +x /opt/qrTraceability/update_server_ip.sh

echo "\n=== 配置完成！==="
echo "\n服务访问地址："
echo "前端：http://$(hostname -I | awk '{print $1}')"
echo "后端API：http://$(hostname -I | awk '{print $1}'):8080"
echo "\n查看服务状态："
echo "sudo systemctl status qrbackend"
echo "sudo systemctl status nginx"
echo "\n查看服务日志："
echo "sudo journalctl -u qrbackend"
echo "sudo journalctl -u nginx"
echo "\n更新服务器 IP 地址："
echo "sudo /opt/qrTraceability/update_server_ip.sh <新IP地址>"
echo "\n=== 配置脚本执行完成 ==="

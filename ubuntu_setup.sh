#!/bin/bash

# Ubuntu 服务器配置脚本 - 农产品二维码追溯系统

echo "=== 农产品二维码追溯系统 - Ubuntu 服务器配置脚本 ==="

echo "\n1. 更新系统包..."
sudo apt update && sudo apt upgrade -y

echo "\n2. 安装必要的依赖包..."
sudo apt install -y git curl wget unzip build-essential nginx

echo "\n3. 安装 Node.js 18.x..."
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt install -y nodejs

echo "\n4. 安装 Go 1.25.1..."
wget https://go.dev/dl/go1.25.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.25.1.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

echo "\n5. 安装 PostgreSQL..."
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
go build -o server cmd/main.go

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
echo "\n=== 配置脚本执行完成 ==="

# 农产品二维码追溯系统

## 项目介绍

农产品二维码追溯系统是一个基于React和Go的全栈应用，用于追踪和管理农产品从种子到成品的整个生命周期。系统通过二维码技术，实现了农产品信息的快速查询和追溯，提高了农产品的透明度和可信度。

### 主要功能

- **种子信息管理**：记录和管理种子的基本信息
- **播种定植管理**：追踪种植过程和时间节点
- **生长媒体管理**：上传和管理生长过程中的图片和视频
- **产品品质管理**：记录和分析产品的品质指标
- **用户交互功能**：收藏和点赞功能
- **二维码生成**：为每个种植批次生成唯一的二维码

## 技术栈

- **前端**：Vue.js 3 + Vite 5.0.0 + Vue Router 4.2.0
- **后端**：Go 1.25.1 + Gin 1.9.1 + PostgreSQL + JWT
- **数据库**：PostgreSQL 14.x 或更高版本
- **其他**：QR码生成、Nginx反向代理、Systemd服务管理

## 环境要求

- Ubuntu 24.04.4 LTS
- Go 1.25.1 或更高版本
- Node.js 20.x 或更高版本
- PostgreSQL 14.x 或更高版本
- Git
- Nginx
- Systemd


### 更新服务器 IP 地址
如果需要更新服务器 IP 地址，可使用：
   ```bash
   sudo /opt/qrTraceability/update_server_ip.sh <新IP地址>
   ```

## 手动安装步骤

### 1. 服务器环境准备

```bash
# 更新系统包
sudo apt update && sudo apt upgrade -y

# 安装必要的依赖包
sudo apt install -y git curl wget unzip build-essential nginx

# 安装 Go (使用 apt 包管理器)
sudo apt install -y golang-go

# 安装 Node.js 20.x (使用 apt 包管理器)
sudo apt install -y ca-certificates curl gnupg
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg
NODE_MAJOR=20
echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$NODE_MAJOR.x nodistro main" | sudo tee /etc/apt/sources.list.d/nodesource.list
sudo apt update
sudo apt install -y nodejs

# 安装 PostgreSQL (使用 apt 包管理器)
sudo apt install -y postgresql postgresql-contrib

# 配置防火墙
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 8080/tcp
sudo ufw allow 5432/tcp
sudo ufw enable
```

### 2. 数据库配置

```bash
# 启动 PostgreSQL 服务
sudo systemctl start postgresql
sudo systemctl enable postgresql

# 创建数据库用户和数据库
sudo -u postgres psql
CREATE USER qruser WITH PASSWORD 'qrpassword';
CREATE DATABASE qrtraceability OWNER qruser;
\q

# 执行数据库模式脚本
psql -U qruser -d qrtraceability -f /opt/qrTraceability/server/internal/db/schema.sql
```

### 3. 项目拉取和构建

```bash
# 从 GitHub 拉取项目代码
sudo mkdir -p /opt/qrTraceability
sudo chown $USER:$USER /opt/qrTraceability
cd /opt/qrTraceability
git clone https://github.com/shuiqinghuan/qrTraceability.git .

# 构建前端项目
cd client
npm install
npm run build

# 编译后端项目
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
```

### 4. 服务配置和启动

#### 4.1 配置环境变量

```bash
# 在 server 目录下创建 .env 文件
echo "DB_HOST=localhost\nDB_PORT=5432\nDB_USER=qruser\nDB_PASSWORD=qrpassword\nDB_NAME=qrtraceability\nJWT_SECRET=your_jwt_secret\nSERVER_PORT=8080" > .env
```

#### 4.2 创建 Systemd 服务

```bash
# 创建 systemd 服务文件
sudo nano /etc/systemd/system/qrbackend.service
```

添加以下内容：

```ini
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
```

#### 4.3 配置 Nginx 反向代理

```bash
# 创建 Nginx 配置文件
sudo nano /etc/nginx/sites-available/qrtraceability
```

添加以下内容：

```nginx
server {
    listen 80;
    server_name example.com;

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
```

#### 4.4 启动服务

```bash
# 启动后端服务
sudo systemctl daemon-reload
sudo systemctl start qrbackend
sudo systemctl enable qrbackend

# 启动 Nginx 服务
sudo ln -s /etc/nginx/sites-available/qrtraceability /etc/nginx/sites-enabled/
sudo systemctl restart nginx
```

## 服务管理

### 查看服务状态

```bash
# 查看后端服务状态
sudo systemctl status qrbackend

# 查看 Nginx 服务状态
sudo systemctl status nginx
```

### 重启服务

```bash
# 重启后端服务
sudo systemctl restart qrbackend

# 重启 Nginx 服务
sudo systemctl restart nginx
```

### 查看服务日志

```bash
# 查看后端服务日志
sudo journalctl -u qrbackend

# 查看 Nginx 服务日志
sudo journalctl -u nginx
```

## 故障排查

### 1. 数据库连接问题

- 检查 PostgreSQL 服务是否运行：`sudo systemctl status postgresql`
- 检查数据库连接配置是否正确
- 检查防火墙是否允许 5432 端口

### 2. 后端服务启动失败

- 查看服务日志：`sudo journalctl -u qrbackend`
- 检查环境变量配置是否正确
- 检查数据库连接是否正常

### 3. 前端页面无法访问

- 检查 Nginx 服务是否运行：`sudo systemctl status nginx`
- 检查 Nginx 配置是否正确
- 检查防火墙是否允许 80 端口
- 检查前端构建是否成功

### 4. API 接口无法访问

- 检查后端服务是否运行：`sudo systemctl status qrbackend`
- 检查防火墙是否允许 8080 端口
- 检查 API 地址是否正确

## 目录结构

```
qrTraceability/
├── client/               # 前端代码
│   ├── public/           # 静态资源
│   ├── src/              # 源代码
│   │   ├── components/   # 组件
│   │   ├── services/     # API 服务
│   │   ├── App.jsx       # 主应用
│   │   └── main.jsx      # 入口文件
│   ├── package.json      # 前端依赖
│   └── vite.config.js    # Vite 配置
├── server/               # 后端代码
│   ├── cmd/              # 命令行入口
│   │   └── main.go       # 主程序
│   ├── internal/         # 内部包
│   │   ├── api/          # API 接口
│   │   ├── auth/         # 认证
│   │   ├── db/           # 数据库
│   │   ├── middleware/   # 中间件
│   │   └── utils/        # 工具函数
│   ├── go.mod            # Go 依赖
│   └── server            # 编译后的可执行文件
├── .env                  # 环境变量配置
└── README.md             # 项目说明
```

## 注意事项

1. 请根据实际服务器配置修改环境变量和配置文件
2. 生产环境中请使用强密码和HTTPS
3. 定期备份数据库和重要文件
4. 保持系统和依赖包的更新

## 许可证

本项目采用 MIT 许可证。

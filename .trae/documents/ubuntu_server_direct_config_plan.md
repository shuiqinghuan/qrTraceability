# Ubuntu服务器直接配置计划

## 项目分析

### 项目结构

- **前端（client）**：React 19.2.4 + Vite 8.0.4 + React Router DOM 7.14.0
- **后端（server）**：Go 1.25.1 + Gin 1.9.1 + PostgreSQL + JWT + QR码生成
- **数据库**：PostgreSQL 数据库，包含用户、种子信息、播种定植、生长媒体、产品品质等表

### 服务器环境要求

- Ubuntu 24.04.4 LTS
- Go 1.25.1 或更高版本
- Node.js 20.x 或更高版本
- PostgreSQL 14.x 或更高版本
- Git
- Nginx（用于反向代理）
- Systemd（用于服务管理）

## 配置步骤

### 1. 服务器环境准备

1. 更新系统包
2. 安装必要的依赖包
3. 安装 Go、Node.js、PostgreSQL
4. 配置防火墙

### 2. 数据库配置

1. 初始化 PostgreSQL 数据库
2. 创建数据库用户和数据库
3. 执行数据库模式脚本

### 3. 项目拉取和构建

1. 从 GitHub 拉取项目代码
2. 构建前端项目
3. 编译后端项目

### 4. 服务配置和启动

1. 配置环境变量
2. 设置后端服务为 systemd 服务
3. 配置 Nginx 反向代理
4. 启动服务

### 5. README.md 生成

1. 编写详细的安装和使用说明
2. 包含环境要求、安装步骤、配置方法、启动命令等

## 详细步骤

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

```bash
# 启动服务
sudo systemctl daemon-reload
sudo systemctl start qrbackend
sudo systemctl enable qrbackend
sudo ln -s /etc/nginx/sites-available/qrtraceability /etc/nginx/sites-enabled/
sudo systemctl restart nginx
```

```bash
# 启动服务
sudo systemctl daemon-reload
sudo systemctl start qrbackend
sudo systemctl enable qrbackend
sudo ln -s /etc/nginx/sites-available/qrtraceability /etc/nginx/sites-enabled/
sudo systemctl restart nginx
```

### 5. README.md 生成

创建详细的 README.md 文件，包含：

- 项目介绍
- 环境要求
- 安装步骤
- 配置方法
- 启动命令
- 故障排查

## 风险处理

1. **依赖版本问题**：确保使用指定版本的依赖
2. **数据库连接问题**：检查数据库配置和网络连接
3. **权限问题**：确保文件和目录权限正确
4. **端口冲突**：确保使用的端口未被占用
5. **防火墙问题**：确保防火墙规则正确配置
6. **服务启动问题**：检查 systemd 服务状态和日志

## 结论

本计划提供了在 Ubuntu 服务器上直接配置项目的详细步骤，包括环境准备、数据库配置、项目构建和服务启动。通过生成的 README.md 文件，可以方便地在服务器上部署和管理项目。

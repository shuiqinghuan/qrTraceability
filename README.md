# 农产品溯源系统

基于 Vue 3 + Django REST Framework 的农产品溯源信息展示系统。消费者可通过本系统追溯农产品的品种信息、种植过程、采收质量等全流程数据，提升产品透明度和信任度。

## 功能特性

- **产品基本信息展示**：品种名称、品种编码、定植地点、定植时间
- **多媒体信息展示**：产品图片轮播、视频播放功能
- **采收质量信息展示**：采收时间、糖度、重量、口感描述、适应人群、品质小结
- **响应式设计**：适配桌面端、平板端、移动端
- **RESTful API**：标准化的后端接口设计

## 技术栈

### 前端

- Vue 3.4+（组合式 API）
- Vite 5（构建工具）
- Vue Router 5（路由管理）
- Axios（HTTP 请求库）

### 后端

- Django 4.2+（Web 框架）
- Django REST Framework（API 开发）
- SQLite（默认数据库）
- Gunicorn（WSGI 服务器）

## 项目结构

```
agricultural-traceability-system/
├── frontend/                      # 前端项目
│   ├── src/
│   │   ├── api/                  # API 请求封装
│   │   │   ├── index.js          # Axios 配置
│   │   │   └── product.js        # 产品接口
│   │   ├── components/           # Vue 组件
│   │   │   ├── ProductInfo.vue   # 产品基本信息
│   │   │   ├── MediaGallery.vue  # 多媒体展示
│   │   │   └── HarvestQuality.vue # 采收质量
│   │   ├── views/
│   │   │   └── ProductTrace.vue  # 产品溯源页面
│   │   ├── router/
│   │   │   └── index.js          # 路由配置
│   │   ├── styles/
│   │   │   └── main.css          # 全局样式
│   │   ├── App.vue
│   │   └── main.js
│   ├── public/                   # 静态资源
│   ├── index.html
│   ├── package.json
│   └── vite.config.js            # Vite 配置
├── backend/                      # 后端项目
│   ├── products/                 # 产品应用
│   │   ├── models.py             # 数据模型
│   │   ├── views.py              # API 视图
│   │   ├── serializers.py         # 序列化器
│   │   ├── urls.py               # 路由配置
│   │   └── admin.py              # Admin 配置
│   ├── manage.py
│   ├── settings.py               # 项目配置
│   ├── wsgi.py                  # WSGI 配置
│   ├── requirements.txt          # Python 依赖
│   ├── deploy.sh                # 部署脚本
│   └── .env.example             # 环境变量示例
├── prd.md                       # 产品需求文档
├── database-design.md           # 数据库设计文档
└── API接口文档.md               # API 接口文档
```

## 快速开始

### 环境要求

- Python 3.10+
- Node.js 18+
- npm 9+

### 1. 克隆项目

```bash
git clone <repository-url>
cd agricultural-traceability-system
```

### 2. 后端设置

```bash
cd backend

# 创建虚拟环境（推荐）
python3 -m venv venv
source venv/bin/activate  # Linux/Mac
# venv\Scripts\activate    # Windows

# 安装依赖
pip install -r requirements.txt

# 执行数据库迁移
python manage.py migrate

# 创建测试数据
python manage.py create_test_data

# 创建管理员账户（可选）
python manage.py createsuperuser

# 启动开发服务器
python manage.py runserver 0.0.0.0:8000
```

### 3. 前端设置

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

### 4. 访问应用

- 前端页面：http://localhost:5173
- 后端 API：http://localhost:8000/api
- 管理后台：http://localhost:8000/admin

## 生产环境部署

### 方式一：使用 Nginx + Gunicorn（推荐）

#### 1. 构建前端

```bash
cd frontend
npm install
npm run build
```

构建产物位于 `frontend/dist` 目录。

#### 2. 配置 Nginx

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 前端静态文件
    location / {
        root /path/to/frontend/dist;
        try_files $uri $uri/ /index.html;
    }

    # API 代理
    location /api/ {
        proxy_pass http://127.0.0.1:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 静态文件
    location /static/ {
        alias /path/to/backend/staticfiles/;
    }
}
```

#### 3. 使用 Gunicorn 运行后端

```bash
cd backend
gunicorn backend.wsgi:application --bind 0.0.0.0:8000 --workers 3
```

#### 4. 配置系统服务（systemd）

创建服务文件 `/etc/systemd/system/traceability.service`：

```ini
[Unit]
Description=Agricultural Traceability System
After=network.target

[Service]
User=www-data
Group=www-data
WorkingDirectory=/path/to/backend
ExecStart=/path/to/backend/venv/bin/gunicorn backend.wsgi:application --bind 127.0.0.1:8000 --workers 3
Restart=always

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable traceability
sudo systemctl start traceability
```

### 方式二：使用 Docker（可选）

创建 `Dockerfile`：

```dockerfile
FROM python:3.11-slim

WORKDIR /app

COPY backend/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY backend/ .

RUN python manage.py migrate
RUN python manage.py collectstatic --noinput

EXPOSE 8000

CMD ["gunicorn", "backend.wsgi:application", "--bind", "0.0.0.0:8000"]
```

构建并运行：

```bash
docker build -t traceability .
docker run -d -p 8000:8000 traceability
```

### 环境变量配置

在生产环境中，创建 `backend/.env` 文件：

```bash
DJANGO_SECRET_KEY=your-super-secret-key-here
DEBUG=False
ALLOWED_HOSTS=your-domain.com,www.your-domain.com
CORS_ALLOWED_ORIGINS=https://your-frontend-domain.com
```

## API 接口

### 基础信息

- 基础 URL：`http://localhost:8000/api`
- 数据格式：JSON
- 认证方式：暂无（后续可扩展）

### 主要接口

| 接口地址 | 方法 | 说明 |
|---------|------|------|
| `/api/products/` | GET | 获取产品列表 |
| `/api/products/{code}/` | GET | 获取产品详情（按编码） |
| `/api/products/` | POST | 创建产品 |
| `/api/products/{id}/media/` | GET/POST | 多媒体列表/添加 |
| `/api/products/{id}/harvest/` | GET/POST | 采收质量信息 |

详细接口文档请参阅 [API接口文档.md](./API接口文档.md)。

## 数据库设计

系统使用 SQLite 数据库，包含以下主要表：

- **product**：产品基本信息表
- **media**：多媒体信息表（图片/视频）
- **harvest_quality**：采收质量信息表

详细设计请参阅 [database-design.md](./database-design.md)。

## 配置说明

### 前端环境变量

创建 `frontend/.env` 文件：

```bash
VITE_API_BASE_URL=http://localhost:8000/api
```

### 后端环境变量

```bash
DJANGO_SECRET_KEY=your-secret-key
DEBUG=False
ALLOWED_HOSTS=localhost,127.0.0.1
```

## 常见问题

### 1. 前端无法访问后端 API

检查后端是否正常启动，并确认 CORS 配置正确。开发环境下前端已配置 API 代理。

### 2. 数据库迁移失败

确保在正确的目录下执行 `python manage.py migrate`，并检查数据库文件权限。

### 3. 静态文件无法加载

执行 `python manage.py collectstatic`，并确保 Nginx 配置了正确的静态文件路径。

### 4. 前端构建失败

检查 Node.js 版本，确保为 18+。删除 `node_modules` 和 `package-lock.json` 后重新安装依赖。

## 开发指南

### 添加新产品功能

1. 在 `backend/products/models.py` 中添加模型
2. 执行 `python manage.py makemigrations`
3. 执行 `python manage.py migrate`
4. 在 `backend/products/serializers.py` 中添加序列化器
5. 在 `backend/products/views.py` 中添加视图
6. 在 `backend/products/urls.py` 中添加路由
7. 更新前端组件和 API 调用

### 修改现有功能

1. 修改对应的模型/序列化器/视图
2. 更新前端组件
3. 测试功能正常

## 许可证

本项目仅供学习交流使用。

## 联系方式

如有问题，请提交 Issue 或联系开发者。

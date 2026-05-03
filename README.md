# 农产品溯源系统

基于 Vue 3 + Django REST Framework 的农产品溯源信息展示系统。消费者可通过扫描二维码或访问链接追溯农产品的品种信息、种植过程、采收质量等全流程数据，提升产品透明度和信任度。

**线上地址：** http://47.104.189.148/

## 功能特性

- **产品列表分页展示**：首页展示所有品种，支持分页浏览，点击进入详情页
- **产品基本信息展示**：品种名称、品种编码、定植地点、定植时间
- **多媒体信息展示**：产品图片轮播、支持文件上传与URL两种方式
- **采收质量信息展示**：采收时间、糖度(Brix)、单果重量、口感描述、适应人群、品质小结
- **响应式设计**：适配桌面端、平板端、移动端
- **RESTful API**：标准化的后端接口设计，支持分页查询

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

### 部署

- Docker + Docker Compose 容器化部署
- Nginx 反向代理与静态资源服务

## 项目结构

```
qrTraceability/
├── frontend/                          # 前端项目
│   ├── Dockerfile                     # 前端容器构建文件
│   ├── src/
│   │   ├── api/                       # API 请求封装
│   │   │   ├── index.js               # Axios 实例配置
│   │   │   └── product.js             # 产品相关接口
│   │   ├── components/                # Vue 组件
│   │   │   ├── ProductInfo.vue        # 产品基本信息卡片
│   │   │   ├── MediaGallery.vue       # 多媒体图片轮播
│   │   │   ├── MediaUpload.vue        # 多媒体文件上传
│   │   │   └── HarvestQuality.vue     # 采收质量信息卡片
│   │   ├── views/
│   │   │   ├── Home.vue               # 首页（产品列表分页）
│   │   │   ├── ProductTrace.vue       # 产品溯源详情页
│   │   │   └── About.vue              # 关于页面
│   │   ├── router/
│   │   │   └── index.js               # 路由配置
│   │   ├── styles/
│   │   │   └── main.css               # 全局样式与CSS变量
│   │   ├── App.vue                    # 根组件（router-view）
│   │   └── main.js                    # 入口文件
│   ├── public/
│   │   └── upload.html                # 独立上传页面
│   ├── package.json
│   └── vite.config.js
├── backend/                           # 后端项目
│   ├── settings.py                    # Django 配置
│   ├── urls.py                        # URL 路由
│   ├── wsgi.py                        # WSGI 配置
│   └── requirements.txt               # Python 依赖
├── products/                          # 产品应用（Django App）
│   ├── models.py                      # 数据模型
│   ├── views.py                       # API 视图
│   ├── serializers.py                 # 序列化器
│   ├── urls.py                        # 应用路由
│   ├── admin.py                       # Admin 后台配置
│   └── management/
│       └── commands/
│           └── create_test_data.py    # 测试数据初始化命令
├── Dockerfile                         # 后端容器构建文件
├── docker-compose.yml                 # Docker Compose 编排文件
├── nginx.conf                         # Nginx 反向代理配置
├── manage.py                          # Django 管理脚本
├── db.sqlite3                         # SQLite 数据库文件
├── prd.md                             # 产品需求文档
├── database-design.md                 # 数据库设计文档
└── API接口文档.md                      # API 接口文档
```

## 快速开始

### 环境要求

- Docker 20+
- Docker Compose 2+

### 使用 Docker Compose 部署（推荐）

```bash
# 克隆项目
git clone <repository-url>
cd qrTraceability

# 构建并启动所有服务
docker compose up -d --build

# 查看运行状态
docker compose ps

# 查看日志
docker compose logs -f
```

启动完成后访问：
- 前端页面：http://localhost
- 后端 API：http://localhost:8000/api
- 管理后台：http://localhost:8000/admin

### 本地开发

#### 后端

```bash
# 创建虚拟环境
python3 -m venv venv
source venv/bin/activate

# 安装依赖
pip install -r backend/requirements.txt

# 执行数据库迁移
python manage.py migrate

# 创建测试数据
python manage.py create_test_data

# 启动开发服务器
python manage.py runserver 0.0.0.0:8000
```

#### 前端

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端开发服务器默认运行在 http://localhost:5173，API 请求通过 Vite 代理转发到后端。

## Docker 部署架构

```
┌─────────────────────────────────────────────────┐
│                   Docker Compose                 │
│                                                  │
│  ┌──────────────┐       ┌─────────────────────┐  │
│  │   frontend   │       │      backend        │  │
│  │  (Nginx:80)  │──────>│  (Gunicorn:8000)    │  │
│  │  Vue SPA     │ /api  │  Django REST API    │  │
│  └──────────────┘       └─────────────────────┘  │
│         │                        │               │
│         │ /media                 │ db.sqlite3    │
│         ▼                        ▼               │
│  ┌──────────────┐       ┌─────────────────────┐  │
│  │  media 卷    │       │   db.sqlite3 卷     │  │
│  └──────────────┘       └─────────────────────┘  │
└─────────────────────────────────────────────────┘
```

- **frontend 容器**：构建 Vue 项目，通过 Nginx 提供静态资源，反向代理 `/api` 到后端
- **backend 容器**：运行 Django 应用，执行数据库迁移、创建测试数据、启动 Gunicorn
- **共享卷**：`media` 目录和 `db.sqlite3` 通过 volume 在两个容器间共享

## API 接口

### 主要接口

| 接口地址 | 方法 | 说明 |
|---------|------|------|
| `/api/products/` | GET | 获取产品列表（支持分页） |
| `/api/products/` | POST | 创建产品 |
| `/api/products/{code}/` | GET | 获取产品详情（按编码或ID） |
| `/api/products/{id}/` | PUT | 更新产品信息 |
| `/api/products/{id}/` | DELETE | 删除产品 |
| `/api/products/{id}/media/` | GET | 获取产品多媒体列表 |
| `/api/products/{id}/media/` | POST | 上传多媒体文件 |
| `/api/media/{id}/` | DELETE | 删除多媒体 |
| `/api/products/{id}/harvest/` | GET | 获取采收质量信息 |
| `/api/products/{id}/harvest/` | POST | 创建/更新采收质量信息 |

### 通用响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": { ... }
}
```

详细接口文档请参阅 [API接口文档.md](./API接口文档.md)。

## 数据库设计

系统使用 SQLite 数据库，包含三张核心表：

| 表名 | 中文名 | 说明 | 关系 |
|------|--------|------|------|
| product | 产品信息表 | 存储品种名称、编码、定植信息 | 主表 |
| media | 多媒体信息表 | 存储产品图片和视频（支持文件上传和URL） | 多对一关联 product |
| harvest_quality | 采收质量信息表 | 存储糖度、重量、口感等质量数据 | 一对一关联 product |

详细设计请参阅 [database-design.md](./database-design.md)。

## 路由说明

| 路径 | 页面 | 说明 |
|------|------|------|
| `/` | 首页 | 产品列表分页展示，点击产品卡片进入详情 |
| `/product/{code}` | 详情页 | 展示产品溯源信息，包含返回列表按钮 |

## 常见问题

### 1. 前端无法访问后端 API

检查后端容器是否正常运行：`docker compose ps`，查看后端日志：`docker compose logs backend`。

### 2. 数据库迁移失败

进入后端容器执行迁移：`docker compose exec backend python manage.py migrate`。

### 3. 上传文件无法访问

确认 `media` 目录已正确挂载，Nginx 配置了 `/media/` 路径的 alias。

### 4. 页面刷新后 404

确认 Nginx 配置了 `try_files $uri $uri/ /index.html;` 以支持 Vue Router 的 history 模式。

## 开发计划

| 阶段 | 任务 | 状态 |
|------|------|------|
| 第一阶段 | PRD 文档编写 | 已完成 |
| 第二阶段 | Vue 前端开发 | 已完成 |
| 第三阶段 | 数据库设计 | 已完成 |
| 第四阶段 | API 接口文档 | 已完成 |
| 第五阶段 | Django 后端开发 | 已完成 |
| 第六阶段 | 前后端对接 | 已完成 |
| 第七阶段 | Docker 容器化部署 | 已完成 |
| 第八阶段 | 测试与优化 | 进行中 |

## 许可证

本项目仅供学习交流使用。

# 农产品二维码追溯系统

这是一个完整的农产品全生命周期追溯系统，通过二维码技术实现农产品从种植到销售的全程可追溯。

## 系统特点

- **产品信息管理**：完整的产品基本信息、批次管理和媒体文件展示
- **采收质量追踪**：记录采收时间、糖度、重量等质量指标
- **用户交互反馈**：支持点赞、分享、收藏功能
- **安全保障**：Redis限流机制，防止恶意刷量
- **容器化部署**：使用Docker和Docker Compose实现一键部署

## 技术栈

### 后端
- **Go 1.26** ：高性能服务端语言 (最新稳定版本)
- **Gin** ：轻量级Web框架
- **GORM** ：ORM库
- **PostgreSQL** ：关系型数据库
- **Redis** ：缓存和限流

### 前端
- **Vue 3** ：现代化前端框架
- **Vant UI** ：移动端组件库
- **Axios** ：HTTP客户端
- **Vite** ：构建工具

### 基础设施
- **Docker** ：容器化
- **Docker Compose** ：编排工具

## 快速开始

### 前置要求

- Docker 和 Docker Compose

### 安装步骤

1. 克隆项目
```bash
git clone https://github.com/shuiqinghuan/qrTraceability
cd qrTraceability
```

2. 启动服务
```bash
docker-compose up -d
```

3. 访问应用

- 本地开发：
  - 前端：http://localhost:3000
  - 后端API：http://localhost:8000
- 云服务器部署：
  - 前端：http://139.155.97.74:3000
  - 后端API：http://139.155.97.74:8000

### 停止服务

```bash
docker-compose down
```

## 项目结构

```
.
├── backend/                 # 后端项目
│   ├── cmd/
│   │   └── api/            # API服务入口
│   ├── internal/
│   │   ├── config/         # 配置管理
│   │   ├── handlers/       # API处理器
│   │   ├── models/         # 数据模型
│   │   ├── repository/     # 数据访问层
│   │   ├── service/        # 业务逻辑层
│   │   └── utils/          # 工具函数
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
├── frontend/               # 前端项目
│   ├── src/
│   │   ├── App.vue        # 主应用组件
│   │   ├── main.js        # 入口文件
│   │   └── style.css      # 全局样式
│   ├── Dockerfile
│   ├── package.json
│   └── vite.config.js
├── docker-compose.yml      # Docker Compose配置
└── README.md              # 项目文档
```

## 核心功能

### 1. 产品管理
- 创建和编辑产品信息
- 产品编码管理
- 产品描述和属性

### 2. 批次管理
- 创建产品批次
- 批次编号和唯一标识
- 定植地点和时间记录
- 媒体文件上传（图片、视频）

### 3. 采收质量
- 采收时间范围记录
- 糖度、重量等指标
- 口感和适应人群描述
- 品质小结

### 4. 用户交互
- 点赞功能
- 分享功能（支持Web Share API）
- 收藏功能
- 实时统计更新

## API文档

### 产品API

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /api/products | 创建产品 |
| GET | /api/products/:id | 获取产品详情 |
| GET | /api/products | 获取产品列表 |
| PUT | /api/products/:id | 更新产品 |
| DELETE | /api/products/:id | 删除产品 |

### 批次API

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /api/batches | 创建批次 |
| GET | /api/batches/:id | 获取批次详情 |
| GET | /api/batches/unique/:unique_id | 通过唯一标识获取批次 |
| GET | /api/batches/product/:product_id | 获取产品的批次列表 |
| PUT | /api/batches/:id | 更新批次 |
| DELETE | /api/batches/:id | 删除批次 |

### 交互API

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /api/interactions | 记录用户交互 |
| GET | /api/interactions/batch/:batch_id | 获取批次交互统计 |

## 数据模型

### Product（产品）
```go
type Product struct {
    ID          uint   `gorm:"primaryKey"`
    Code        string `gorm:"uniqueIndex;not null"`
    Name        string `gorm:"not null"`
    Description string
}
```

### ProductBatch（产品批次）
```go
type ProductBatch struct {
    ID               uint      `gorm:"primaryKey"`
    ProductID        uint      `gorm:"not null"`
    BatchNumber      string    `gorm:"not null"`
    UniqueID         string    `gorm:"uniqueIndex;not null"`
    PlantingLocation string
    PlantingDate     time.Time
}
```

### UserInteraction（用户交互）
```go
type UserInteraction struct {
    ID         uint   `gorm:"primaryKey"`
    BatchID    uint   `gorm:"not null"`
    IP         string `gorm:"not null"`
    ActionType string `gorm:"not null"` // like, share, collect
}
```

## 开发说明

### 后端开发

1. 进入backend目录
```bash
cd backend
```

2. 安装依赖
```bash
go mod download
```

3. 运行服务
```bash
go run cmd/api/main.go
```

### 前端开发

1. 进入frontend目录
```bash
cd frontend
```

2. 安装依赖
```bash
npm install
```

3. 启动开发服务器
```bash
npm run dev
```

4. 构建生产版本
```bash
npm run build
```

## 配置说明

后端配置通过环境变量控制：

- `SERVER_PORT`：服务端口，默认8000
- `DB_HOST`：数据库主机，默认postgres
- `DB_PORT`：数据库端口，默认5432
- `DB_USER`：数据库用户，默认postgres
- `DB_PASSWORD`：数据库密码，默认password
- `DB_NAME`：数据库名称，默认qr_traceability
- `DB_SSLMODE`：数据库SSL模式，默认disable
- `REDIS_HOST`：Redis主机，默认redis
- `REDIS_PORT`：Redis端口，默认6379
- `REDIS_PASSWORD`：Redis密码，默认为空
- `REDIS_DB`：Redis数据库，默认0

## 许可证

MIT License

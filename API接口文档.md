# 农产品溯源系统 API接口文档

## 1. 接口概述

### 1.1 基础信息
- 基础URL：`http://47.104.189.148/api`
- 数据格式：JSON
- 编码：UTF-8
- 接口风格：RESTful

### 1.2 通用响应格式

**成功响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": { ... }
}
```

**错误响应：**
```json
{
  "code": 400,
  "message": "错误描述",
  "data": null
}
```

### 1.3 通用状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 2. 产品信息接口

### 2.1 获取产品详情

**接口描述：** 根据产品编码获取产品完整信息，包括基本信息、多媒体信息、采收质量信息。

**请求地址：** `GET /api/products/{code}/`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| code | String | 是 | 产品编码（URL路径参数） |

**请求示例：**
```
GET /api/products/4395/
```

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "枣甜5号",
    "code": "4395",
    "plantingLocation": "山东省济南市历城区农业示范园",
    "plantingDate": "2024-03-15",
    "images": [
      "https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=新鲜红枣果实特写&image_size=landscape_16_9",
      "https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=红枣种植园风景&image_size=landscape_16_9",
      "https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=新鲜水果采摘场景&image_size=landscape_16_9"
    ],
    "videos": [],
    "harvest": {
      "startDate": "2024-07-01",
      "endDate": "2024-07-15",
      "sugarContent": 15.5,
      "weight": 280.5,
      "taste": "肉质细腻，汁多味甜，口感爽脆",
      "suitableCrowd": "老少皆宜，特别适合血糖稳定人群",
      "qualitySummary": "果实饱满，色泽鲜艳，糖度适中，品质优良"
    },
    "createdAt": "2024-01-01T10:00:00Z",
    "updatedAt": "2024-01-01T10:00:00Z"
  }
}
```

### 2.2 获取产品列表

**接口描述：** 获取所有产品列表。

**请求地址：** `GET /api/products/`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | Integer | 否 | 页码，默认1 |
| pageSize | Integer | 否 | 每页数量，默认10 |

**请求示例：**
```
GET /api/products/?page=1&pageSize=10
```

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total": 2,
    "page": 1,
    "pageSize": 10,
    "list": [
      {
        "id": 3,
        "name": "甜美4K",
        "code": "1206",
        "plantingLocation": "陕西省渭南市蒲城县龙池镇",
        "plantingDate": "2026-02-14"
      },
      {
        "id": 1,
        "name": "枣甜5号",
        "code": "4395",
        "plantingLocation": "四川省乐山市夹江县甘江镇",
        "plantingDate": "2026-01-18"
      }
    ]
  }
}
```

### 2.3 创建产品

**接口描述：** 创建新的产品信息。

**请求地址：** `POST /api/products/`

**请求头：**
```
Content-Type: application/json
```

**请求体：**
```json
{
  "name": "枣甜5号",
  "code": "4395",
  "plantingLocation": "山东省济南市历城区农业示范园",
  "plantingDate": "2024-03-15"
}
```

**响应示例：**
```json
{
  "code": 201,
  "message": "创建成功",
  "data": {
    "id": 1,
    "name": "枣甜5号",
    "code": "4395",
    "plantingLocation": "山东省济南市历城区农业示范园",
    "plantingDate": "2024-03-15",
    "createdAt": "2024-01-01T10:00:00Z",
    "updatedAt": "2024-01-01T10:00:00Z"
  }
}
```

### 2.4 更新产品

**接口描述：** 更新产品基本信息。

**请求地址：** `PUT /api/products/{id}/`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | Integer | 是 | 产品ID（URL路径参数） |

**请求体：**
```json
{
  "name": "枣甜5号",
  "code": "4395",
  "plantingLocation": "山东省济南市历城区农业示范园",
  "plantingDate": "2024-03-15"
}
```

**响应示例：**
```json
{
  "code": 200,
  "message": "更新成功",
  "data": {
    "id": 1,
    "name": "枣甜5号",
    "code": "4395",
    "plantingLocation": "山东省济南市历城区农业示范园",
    "plantingDate": "2024-03-15",
    "createdAt": "2024-01-01T10:00:00Z",
    "updatedAt": "2024-01-01T11:00:00Z"
  }
}
```

### 2.5 删除产品

**接口描述：** 删除产品及其关联数据。

**请求地址：** `DELETE /api/products/{id}/`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | Integer | 是 | 产品ID（URL路径参数） |

**响应示例：**
```json
{
  "code": 200,
  "message": "删除成功",
  "data": null
}
```

## 3. 多媒体信息接口

### 3.1 获取产品多媒体列表

**接口描述：** 获取指定产品的所有多媒体信息。

**请求地址：** `GET /api/products/{productId}/media/`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| productId | Integer | 是 | 产品ID（URL路径参数） |
| type | String | 否 | 媒体类型过滤：image/video |

**请求示例：**
```
GET /api/products/1/media/?type=image
```

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "productId": 1,
      "mediaType": "image",
      "url": "https://example.com/image1.jpg",
      "title": "产品图片1",
      "description": "产品主图",
      "sortOrder": 1,
      "createdAt": "2024-01-01T10:00:00Z"
    }
  ]
}
```

### 3.2 添加多媒体

**接口描述：** 为产品添加图片或视频，支持通过URL或文件上传两种方式。

**请求地址：** `POST /api/products/{productId}/media/`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| productId | Integer | 是 | 产品ID（URL路径参数） |

**方式一：通过URL添加**

请求头：`Content-Type: application/json`

```json
{
  "mediaType": "image",
  "url": "https://example.com/image1.jpg",
  "title": "产品图片1",
  "description": "产品主图",
  "sortOrder": 1
}
```

**方式二：通过文件上传**

请求头：`Content-Type: multipart/form-data`

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| mediaType | String | 是 | 媒体类型：image/video |
| file | File | 是 | 上传的媒体文件 |
| title | String | 否 | 媒体标题 |
| description | String | 否 | 媒体描述 |
| sortOrder | Integer | 否 | 排序顺序，默认0 |

**响应示例：**
```json
{
  "code": 201,
  "message": "上传成功",
  "data": {
    "id": 1,
    "productId": 1,
    "mediaType": "image",
    "url": "/media/products/2024/01/01/image.jpg",
    "title": "产品图片1",
    "description": "产品主图",
    "sortOrder": 1,
    "createdAt": "2024-01-01T10:00:00Z"
  }
}
```

### 3.3 删除多媒体

**接口描述：** 删除指定的多媒体信息。

**请求地址：** `DELETE /api/media/{id}/`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | Integer | 是 | 多媒体ID（URL路径参数） |

**响应示例：**
```json
{
  "code": 200,
  "message": "删除成功",
  "data": null
}
```

## 4. 采收质量信息接口

### 4.1 获取采收质量信息

**接口描述：** 获取指定产品的采收质量信息。

**请求地址：** `GET /api/products/{productId}/harvest/`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| productId | Integer | 是 | 产品ID（URL路径参数） |

**请求示例：**
```
GET /api/products/1/harvest/
```

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "productId": 1,
    "startDate": "2024-07-01",
    "endDate": "2024-07-15",
    "sugarContent": 15.5,
    "weight": 280.5,
    "taste": "肉质细腻，汁多味甜，口感爽脆",
    "suitableCrowd": "老少皆宜，特别适合血糖稳定人群",
    "qualitySummary": "果实饱满，色泽鲜艳，糖度适中，品质优良",
    "createdAt": "2024-01-01T10:00:00Z",
    "updatedAt": "2024-01-01T10:00:00Z"
  }
}
```

### 4.2 创建/更新采收质量信息

**接口描述：** 创建或更新产品的采收质量信息（存在则更新，不存在则创建）。

**请求地址：** `POST /api/products/{productId}/harvest/`

**请求参数：**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| productId | Integer | 是 | 产品ID（URL路径参数） |

**请求体：**
```json
{
  "startDate": "2024-07-01",
  "endDate": "2024-07-15",
  "sugarContent": 15.5,
  "weight": 280.5,
  "taste": "肉质细腻，汁多味甜，口感爽脆",
  "suitableCrowd": "老少皆宜，特别适合血糖稳定人群",
  "qualitySummary": "果实饱满，色泽鲜艳，糖度适中，品质优良"
}
```

**响应示例：**
```json
{
  "code": 200,
  "message": "保存成功",
  "data": {
    "id": 1,
    "productId": 1,
    "startDate": "2024-07-01",
    "endDate": "2024-07-15",
    "sugarContent": 15.5,
    "weight": 280.5,
    "taste": "肉质细腻，汁多味甜，口感爽脆",
    "suitableCrowd": "老少皆宜，特别适合血糖稳定人群",
    "qualitySummary": "果实饱满，色泽鲜艳，糖度适中，品质优良",
    "createdAt": "2024-01-01T10:00:00Z",
    "updatedAt": "2024-01-01T11:00:00Z"
  }
}
```

## 5. 接口调用流程

### 5.1 前端页面加载流程

```
1. 用户访问产品溯源页面
   ↓
2. 前端调用 GET /api/products/{code}/
   ↓
3. 后端返回完整产品信息（包含多媒体和采收质量）
   ↓
4. 前端渲染页面
```

### 5.2 前端API调用示例

```javascript
// src/api/product.js
import request from './index'

// 获取产品详情
export function getProductByCode(code) {
  return request.get(`/products/${code}/`)
}

// 获取产品列表
export function getProductList(params) {
  return request.get('/products/', { params })
}

// 创建产品
export function createProduct(data) {
  return request.post('/products/', data)
}

// 更新产品
export function updateProduct(id, data) {
  return request.put(`/products/${id}/`, data)
}

// 删除产品
export function deleteProduct(id) {
  return request.delete(`/products/${id}/`)
}

// 获取多媒体列表
export function getMediaList(productId, type) {
  return request.get(`/products/${productId}/media/`, { params: { type } })
}

// 添加多媒体
export function addMedia(productId, data) {
  return request.post(`/products/${productId}/media/`, data)
}

// 删除多媒体
export function deleteMedia(id) {
  return request.delete(`/media/${id}/`)
}

// 获取采收质量信息
export function getHarvestQuality(productId) {
  return request.get(`/products/${productId}/harvest/`)
}

// 保存采收质量信息
export function saveHarvestQuality(productId, data) {
  return request.post(`/products/${productId}/harvest/`, data)
}
```

## 6. 接口鉴权

当前版本暂不实现用户鉴权，后续可扩展添加以下鉴权方式：
- Token认证
- JWT认证
- Session认证

## 7. 跨域配置

后端需要配置CORS允许前端跨域访问：

```python
# Django settings.py
CORS_ALLOWED_ORIGINS = [
    "http://localhost:5173",
    "http://localhost:5174",
]

CORS_ALLOW_ALL_ORIGINS = True  # 开发环境可使用
```

## 8. 更新日志

| 日期 | 版本 | 更新内容 | 更新人 |
|------|------|----------|--------|
| 2024-01-01 | v1.0 | 初始接口设计 | 系统 |
| 2026-05-03 | v1.1 | 补充文件上传接口说明，更新示例数据 | 系统 |

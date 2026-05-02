# 农产品溯源系统 数据库设计文档

## 1. 数据库概述

### 1.1 数据库类型
SQLite（Django自带数据库）

### 1.2 数据库命名规范
- 表名：小写字母，使用下划线分隔（snake_case）
- 字段名：小写字母，使用下划线分隔（snake_case）
- 主键：默认使用 `id` 作为自增主键
- 外键：关联表名_id 格式

## 2. 数据表设计

### 2.1 产品基本信息表 (product)

存储农产品的基本品种信息。

| 字段名 | 字段类型 | 是否必填 | 默认值 | 中文说明 |
|--------|----------|----------|--------|----------|
| id | INTEGER | 是 | 自增 | 主键ID |
| name | VARCHAR(100) | 是 | - | 品种名称 |
| code | VARCHAR(50) | 是 | - | 品种编码（唯一标识） |
| planting_location | VARCHAR(255) | 是 | - | 定植地点 |
| planting_date | DATE | 是 | - | 定植时间 |
| created_at | DATETIME | 是 | 当前时间 | 创建时间 |
| updated_at | DATETIME | 是 | 当前时间 | 更新时间 |

**索引设计：**
- PRIMARY KEY: id
- UNIQUE INDEX: code

**Django模型定义：**
```python
class Product(models.Model):
    name = models.CharField('品种名称', max_length=100)
    code = models.CharField('品种编码', max_length=50, unique=True)
    planting_location = models.CharField('定植地点', max_length=255)
    planting_date = models.DateField('定植时间')
    created_at = models.DateTimeField('创建时间', auto_now_add=True)
    updated_at = models.DateTimeField('更新时间', auto_now=True)

    class Meta:
        db_table = 'product'
        verbose_name = '产品信息'
        verbose_name_plural = '产品信息'

    def __str__(self):
        return f'{self.name}({self.code})'
```

### 2.2 多媒体信息表 (media)

存储产品的图片和视频信息。

| 字段名 | 字段类型 | 是否必填 | 默认值 | 中文说明 |
|--------|----------|----------|--------|----------|
| id | INTEGER | 是 | 自增 | 主键ID |
| product_id | INTEGER | 是 | - | 关联产品ID（外键） |
| media_type | VARCHAR(20) | 是 | - | 媒体类型（image/video） |
| url | VARCHAR(500) | 是 | - | 媒体文件URL |
| title | VARCHAR(100) | 否 | - | 媒体标题 |
| description | TEXT | 否 | - | 媒体描述 |
| sort_order | INTEGER | 否 | 0 | 排序顺序 |
| created_at | DATETIME | 是 | 当前时间 | 创建时间 |

**索引设计：**
- PRIMARY KEY: id
- INDEX: product_id
- FOREIGN KEY: product_id REFERENCES product(id)

**Django模型定义：**
```python
class Media(models.Model):
    MEDIA_TYPE_CHOICES = [
        ('image', '图片'),
        ('video', '视频'),
    ]

    product = models.ForeignKey(Product, on_delete=models.CASCADE, related_name='media', verbose_name='关联产品')
    media_type = models.CharField('媒体类型', max_length=20, choices=MEDIA_TYPE_CHOICES)
    url = models.CharField('媒体URL', max_length=500)
    title = models.CharField('媒体标题', max_length=100, blank=True, null=True)
    description = models.TextField('媒体描述', blank=True, null=True)
    sort_order = models.IntegerField('排序顺序', default=0)
    created_at = models.DateTimeField('创建时间', auto_now_add=True)

    class Meta:
        db_table = 'media'
        verbose_name = '多媒体信息'
        verbose_name_plural = '多媒体信息'
        ordering = ['sort_order', 'id']

    def __str__(self):
        return f'{self.get_media_type_display()} - {self.product.name}'
```

### 2.3 采收质量信息表 (harvest_quality)

存储产品的采收质量检测信息。

| 字段名 | 字段类型 | 是否必填 | 默认值 | 中文说明 |
|--------|----------|----------|--------|----------|
| id | INTEGER | 是 | 自增 | 主键ID |
| product_id | INTEGER | 是 | - | 关联产品ID（外键，一对一） |
| harvest_start_date | DATE | 是 | - | 采收起始时间 |
| harvest_end_date | DATE | 是 | - | 采收终止时间 |
| sugar_content | DECIMAL(5,2) | 是 | - | 糖度（Brix） |
| weight | DECIMAL(8,2) | 是 | - | 单果重量（克） |
| taste | VARCHAR(255) | 是 | - | 口感描述 |
| suitable_crowd | VARCHAR(255) | 是 | - | 适应人群 |
| quality_summary | TEXT | 是 | - | 品质小结 |
| created_at | DATETIME | 是 | 当前时间 | 创建时间 |
| updated_at | DATETIME | 是 | 当前时间 | 更新时间 |

**索引设计：**
- PRIMARY KEY: id
- UNIQUE INDEX: product_id
- FOREIGN KEY: product_id REFERENCES product(id)

**Django模型定义：**
```python
class HarvestQuality(models.Model):
    product = models.OneToOneField(Product, on_delete=models.CASCADE, related_name='harvest_quality', verbose_name='关联产品')
    harvest_start_date = models.DateField('采收起始时间')
    harvest_end_date = models.DateField('采收终止时间')
    sugar_content = models.DecimalField('糖度(Brix)', max_digits=5, decimal_places=2)
    weight = models.DecimalField('单果重量(克)', max_digits=8, decimal_places=2)
    taste = models.CharField('口感描述', max_length=255)
    suitable_crowd = models.CharField('适应人群', max_length=255)
    quality_summary = models.TextField('品质小结')
    created_at = models.DateTimeField('创建时间', auto_now_add=True)
    updated_at = models.DateTimeField('更新时间', auto_now=True)

    class Meta:
        db_table = 'harvest_quality'
        verbose_name = '采收质量信息'
        verbose_name_plural = '采收质量信息'

    def __str__(self):
        return f'{self.product.name} - 采收质量'
```

## 3. 数据库关系图

```
┌─────────────────────┐
│      product        │
├─────────────────────┤
│ id (PK)             │
│ name                │
│ code (UNIQUE)       │
│ planting_location   │
│ planting_date       │
│ created_at          │
│ updated_at          │
└──────────┬──────────┘
           │
           │ 1:N
           ▼
┌─────────────────────┐
│       media         │
├─────────────────────┤
│ id (PK)             │
│ product_id (FK)     │
│ media_type          │
│ url                 │
│ title               │
│ description         │
│ sort_order          │
│ created_at          │
└─────────────────────┘

           │
           │ 1:1
           ▼
┌─────────────────────┐
│  harvest_quality    │
├─────────────────────┤
│ id (PK)             │
│ product_id (FK/UNQ) │
│ harvest_start_date  │
│ harvest_end_date    │
│ sugar_content       │
│ weight              │
│ taste               │
│ suitable_crowd      │
│ quality_summary     │
│ created_at          │
│ updated_at          │
└─────────────────────┘
```

## 4. 初始测试数据

### 4.1 产品信息测试数据
```sql
INSERT INTO product (name, code, planting_location, planting_date, created_at, updated_at)
VALUES ('枣甜5号', '4395', '山东省济南市历城区农业示范园', '2024-03-15', datetime('now'), datetime('now'));
```

### 4.2 多媒体信息测试数据
```sql
INSERT INTO media (product_id, media_type, url, title, sort_order, created_at)
VALUES 
(1, 'image', 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=新鲜红枣果实特写&image_size=landscape_16_9', '产品图片1', 1, datetime('now')),
(1, 'image', 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=红枣种植园风景&image_size=landscape_16_9', '产品图片2', 2, datetime('now')),
(1, 'image', 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=新鲜水果采摘场景&image_size=landscape_16_9', '产品图片3', 3, datetime('now'));
```

### 4.3 采收质量测试数据
```sql
INSERT INTO harvest_quality (product_id, harvest_start_date, harvest_end_date, sugar_content, weight, taste, suitable_crowd, quality_summary, created_at, updated_at)
VALUES (1, '2024-07-01', '2024-07-15', 15.5, 280.5, '肉质细腻，汁多味甜，口感爽脆', '老少皆宜，特别适合血糖稳定人群', '果实饱满，色泽鲜艳，糖度适中，品质优良', datetime('now'), datetime('now'));
```

## 5. 数据库操作说明

### 5.1 Django迁移命令
```bash
# 创建迁移文件
python manage.py makemigrations

# 执行迁移
python manage.py migrate
```

### 5.2 Django Admin注册
在应用的 `admin.py` 中注册模型：
```python
from django.contrib import admin
from .models import Product, Media, HarvestQuality

@admin.register(Product)
class ProductAdmin(admin.ModelAdmin):
    list_display = ['id', 'name', 'code', 'planting_location', 'planting_date']
    search_fields = ['name', 'code']

@admin.register(Media)
class MediaAdmin(admin.ModelAdmin):
    list_display = ['id', 'product', 'media_type', 'title', 'sort_order']
    list_filter = ['media_type']

@admin.register(HarvestQuality)
class HarvestQualityAdmin(admin.ModelAdmin):
    list_display = ['id', 'product', 'harvest_start_date', 'harvest_end_date', 'sugar_content', 'weight']
```

## 6. 数据字典汇总

| 表名 | 中文名 | 说明 | 关系 |
|------|--------|------|------|
| product | 产品信息表 | 存储产品基本信息 | 主表 |
| media | 多媒体信息表 | 存储产品图片和视频 | 多对一关联product |
| harvest_quality | 采收质量信息表 | 存储采收质量数据 | 一对一关联product |

## 7. 更新日志

| 日期 | 版本 | 更新内容 | 更新人 |
|------|------|----------|--------|
| 2024-01-01 | v1.0 | 初始数据库设计 | 系统 |

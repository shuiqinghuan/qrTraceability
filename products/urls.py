"""
产品URL路由模块

定义产品相关的API路由映射：
- /products/                产品列表（GET）和创建产品（POST）
- /products/<code_or_id>/   产品详情（GET）、更新（PUT）、删除（DELETE），支持编码或ID
- /products/<id>/media/     产品多媒体资源列表（GET）和上传（POST）
- /media/<pk>/              单个多媒体资源操作（DELETE）
- /products/<id>/harvest/   采收质量信息（GET）和创建/更新（POST）
"""

from django.urls import path
from .views import (
    ProductListView,
    ProductDetailView,
    ProductMediaListView,
    MediaDetailView,
    HarvestQualityView,
)

urlpatterns = [
    # 产品列表与创建
    path('products/', ProductListView.as_view(), name='product-list'),
    # 产品详情、更新、删除（支持品种编码或数字ID）
    path('products/<str:code_or_id>/', ProductDetailView.as_view(), name='product-detail'),
    # 产品的多媒体资源管理
    path('products/<int:productId>/media/', ProductMediaListView.as_view(), name='product-media-list'),
    # 单个多媒体资源删除
    path('media/<int:pk>/', MediaDetailView.as_view(), name='media-detail'),
    # 产品的采收质量信息
    path('products/<int:productId>/harvest/', HarvestQualityView.as_view(), name='harvest-quality'),
]

"""
URL路由配置模块。

定义溯源应用的API路由，映射URL路径到对应的视图处理函数。
路由由上层项目统一挂载，此处仅定义应用内部路径。
"""

from django.urls import path
from . import views

urlpatterns = [
    # 产品列表（GET）和创建（POST）
    path('products/', views.ProductListCreateView.as_view(), name='product-list-create'),
    # 单个产品的详情（GET）、更新（PUT/PATCH）和删除（DELETE）
    path('products/<int:pk>/', views.ProductDetailView.as_view(), name='product-detail'),
]

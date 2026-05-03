"""
农产品溯源系统 - 主URL路由配置。

挂载Django Admin后台管理和API接口路由。
"""

from django.contrib import admin
from django.urls import path, include

urlpatterns = [
    # Django自带的管理后台
    path("admin/", admin.site.urls),
    # 产品相关API接口，统一以 /api/ 为前缀
    path("api/", include('products.urls')),
]

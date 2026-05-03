"""
农产品溯源系统 - 开发环境URL路由配置。

挂载Django Admin后台管理和溯源应用API接口。
"""

from django.contrib import admin
from django.urls import path, include

urlpatterns = [
    # Django管理后台
    path('admin/', admin.site.urls),
    # 溯源应用API接口
    path('api/', include('apps.traceability.urls')),
]

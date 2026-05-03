"""
农产品溯源系统 - 开发环境ASGI入口模块。
"""

import os
from django.core.asgi import get_asgi_application

# 设置默认配置模块为开发环境配置
os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'config.settings')

application = get_asgi_application()

"""
农产品溯源系统 - ASGI入口模块。

暴露ASGI可调用对象，用于支持异步服务器（如Daphne、Uvicorn）部署。
"""

import os

from django.core.asgi import get_asgi_application

os.environ.setdefault("DJANGO_SETTINGS_MODULE", "backend.settings")

application = get_asgi_application()

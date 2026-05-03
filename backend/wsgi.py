"""
农产品溯源系统 - WSGI入口模块。

用于Gunicorn、uWSGI等WSGI服务器加载Django应用。
"""

import os

from django.core.wsgi import get_wsgi_application

os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'backend.settings')

application = get_wsgi_application()

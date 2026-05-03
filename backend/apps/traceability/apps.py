"""
应用配置模块。

Django应用的初始化配置，定义应用名称和在后台管理中显示的中文名称。
"""

from django.apps import AppConfig


class TraceabilityConfig(AppConfig):
    """农产品溯源应用配置。"""

    default_auto_field = 'django.db.models.BigAutoField'
    name = 'apps.traceability'
    verbose_name = '农产品溯源'

"""
产品应用配置模块

定义products应用的Django配置，指定应用名称。
"""

from django.apps import AppConfig


class ProductsConfig(AppConfig):
    """产品应用配置类"""
    name = "products"

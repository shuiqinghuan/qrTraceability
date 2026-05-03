"""
后台管理配置模块。

注册各模型到Django Admin后台，配置列表展示字段、搜索和筛选功能，
方便管理员通过Web界面管理溯源数据。
"""

from django.contrib import admin
from .models import Product, Media, HarvestQuality


@admin.register(Product)
class ProductAdmin(admin.ModelAdmin):
    """产品信息后台管理，支持按名称和编码搜索。"""

    list_display = ['id', 'name', 'code', 'planting_location', 'planting_date']
    search_fields = ['name', 'code']


@admin.register(Media)
class MediaAdmin(admin.ModelAdmin):
    """多媒体资源后台管理，支持按媒体类型筛选。"""

    list_display = ['id', 'product', 'media_type', 'title', 'sort_order']
    list_filter = ['media_type']


@admin.register(HarvestQuality)
class HarvestQualityAdmin(admin.ModelAdmin):
    """采收质量信息后台管理。"""

    list_display = ['id', 'product', 'harvest_start_date', 'harvest_end_date', 'sugar_content', 'weight']

"""
产品后台管理模块

配置Django Admin后台的产品管理界面，包括：
- ProductAdmin: 产品信息管理，支持按名称和编码搜索
- MediaAdmin: 多媒体资源管理，支持按类型筛选和媒体预览
- HarvestQualityAdmin: 采收质量信息管理
"""

from django.contrib import admin
from django.utils.html import format_html
from .models import Product, Media, HarvestQuality


@admin.register(Product)
class ProductAdmin(admin.ModelAdmin):
    """产品后台管理配置，支持搜索和列表展示"""
    list_display = ['id', 'name', 'code', 'planting_location', 'planting_date']
    search_fields = ['name', 'code']


@admin.register(Media)
class MediaAdmin(admin.ModelAdmin):
    """多媒体资源后台管理配置，支持类型筛选和在线预览"""
    list_display = ['id', 'product', 'media_type', 'title', 'media_preview', 'sort_order']
    list_filter = ['media_type']
    readonly_fields = ['media_preview']

    def media_preview(self, obj):
        """在后台列表中生成媒体预览，图片显示缩略图，视频显示播放器"""
        if obj.file:
            if obj.media_type == 'image':
                return format_html('<img src="{}" style="max-height: 50px; max-width: 100px;" />', obj.file.url)
            elif obj.media_type == 'video':
                return format_html('<video src="{}" style="max-height: 50px; max-width: 100px;" controls></video>', obj.file.url)
        return obj.url if obj.url else '-'
    media_preview.short_description = '预览'


@admin.register(HarvestQuality)
class HarvestQualityAdmin(admin.ModelAdmin):
    """采收质量后台管理配置"""
    list_display = ['id', 'product', 'harvest_start_date', 'harvest_end_date', 'sugar_content', 'weight']

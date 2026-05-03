from django.contrib import admin
from django.utils.html import format_html
from .models import Product, Media, HarvestQuality


@admin.register(Product)
class ProductAdmin(admin.ModelAdmin):
    list_display = ['id', 'name', 'code', 'planting_location', 'planting_date']
    search_fields = ['name', 'code']


@admin.register(Media)
class MediaAdmin(admin.ModelAdmin):
    list_display = ['id', 'product', 'media_type', 'title', 'media_preview', 'sort_order']
    list_filter = ['media_type']
    readonly_fields = ['media_preview']

    def media_preview(self, obj):
        if obj.file:
            if obj.media_type == 'image':
                return format_html('<img src="{}" style="max-height: 50px; max-width: 100px;" />', obj.file.url)
            elif obj.media_type == 'video':
                return format_html('<video src="{}" style="max-height: 50px; max-width: 100px;" controls></video>', obj.file.url)
        return obj.url if obj.url else '-'
    media_preview.short_description = '预览'


@admin.register(HarvestQuality)
class HarvestQualityAdmin(admin.ModelAdmin):
    list_display = ['id', 'product', 'harvest_start_date', 'harvest_end_date', 'sugar_content', 'weight']

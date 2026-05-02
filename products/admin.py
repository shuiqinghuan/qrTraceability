from django.contrib import admin
from .models import Product, Media, HarvestQuality


@admin.register(Product)
class ProductAdmin(admin.ModelAdmin):
    list_display = ['id', 'name', 'code', 'planting_location', 'planting_date']
    search_fields = ['name', 'code']


@admin.register(Media)
class MediaAdmin(admin.ModelAdmin):
    list_display = ['id', 'product', 'media_type', 'title', 'sort_order']
    list_filter = ['media_type']


@admin.register(HarvestQuality)
class HarvestQualityAdmin(admin.ModelAdmin):
    list_display = ['id', 'product', 'harvest_start_date', 'harvest_end_date', 'sugar_content', 'weight']

from django.contrib import admin
from .models import Product, MediaFile, AdminUser


@admin.register(Product)
class ProductAdmin(admin.ModelAdmin):
    list_display = ['name', 'code', 'planting_location', 'planting_date', 'created_at']
    list_filter = ['planting_date', 'quality']
    search_fields = ['name', 'code', 'planting_location']
    ordering = ['-created_at']


@admin.register(MediaFile)
class MediaFileAdmin(admin.ModelAdmin):
    list_display = ['filename', 'product', 'media_type', 'file_size', 'uploaded_at']
    list_filter = ['media_type', 'uploaded_at']
    search_fields = ['filename', 'product__name']


@admin.register(AdminUser)
class AdminUserAdmin(admin.ModelAdmin):
    list_display = ['user', 'phone']
    search_fields = ['user__username', 'phone']

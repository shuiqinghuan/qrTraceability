"""
序列化器模块。

定义各模型的DRF序列化器，用于API数据的验证、序列化与反序列化。
ProductSerializer中嵌套了多媒体和采收质量的序列化器，实现关联数据的一次性返回。
"""

from rest_framework import serializers
from .models import Product, Media, HarvestQuality


class MediaSerializer(serializers.ModelSerializer):
    """多媒体资源序列化器。"""
    class Meta:
        model = Media
        fields = '__all__'


class HarvestQualitySerializer(serializers.ModelSerializer):
    """采收质量信息序列化器。"""
    class Meta:
        model = HarvestQuality
        fields = '__all__'


class ProductSerializer(serializers.ModelSerializer):
    """产品信息序列化器，嵌套返回关联的多媒体列表和采收质量信息。"""

    # 嵌套只读字段：一对多的多媒体资源（列表形式）
    media = MediaSerializer(many=True, read_only=True)
    # 嵌套只读字段：一对一的采收质量信息
    harvest_quality = HarvestQualitySerializer(read_only=True)

    class Meta:
        model = Product
        fields = '__all__'

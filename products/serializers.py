"""
产品序列化器模块

将产品相关数据模型转换为JSON格式的API响应，主要包括：
- ProductSerializer: 产品完整信息序列化（含读写）
- ProductListSerializer: 产品列表精简序列化（不含时间戳）
- MediaSerializer: 多媒体资源序列化（支持文件上传和URL两种方式）
- HarvestQualitySerializer: 采收质量信息序列化（字段名转驼峰）
- ProductDetailSerializer: 产品详情聚合序列化（关联图片、视频、采收质量）
"""

from rest_framework import serializers
from .models import Product, Media, HarvestQuality


class ProductSerializer(serializers.ModelSerializer):
    """产品序列化器，将蛇形命名的数据库字段转为驼峰命名供前端使用"""
    plantingLocation = serializers.CharField(source='planting_location')
    plantingDate = serializers.DateField(source='planting_date')
    createdAt = serializers.DateTimeField(source='created_at', read_only=True)
    updatedAt = serializers.DateTimeField(source='updated_at', read_only=True)

    class Meta:
        model = Product
        fields = ['id', 'name', 'code', 'plantingLocation', 'plantingDate', 'createdAt', 'updatedAt']


class ProductListSerializer(serializers.ModelSerializer):
    """产品列表序列化器，仅包含列表展示所需的精简字段"""
    plantingLocation = serializers.CharField(source='planting_location')
    plantingDate = serializers.DateField(source='planting_date')

    class Meta:
        model = Product
        fields = ['id', 'name', 'code', 'plantingLocation', 'plantingDate']


class MediaSerializer(serializers.ModelSerializer):
    """多媒体资源序列化器，支持文件上传和外部URL两种媒体来源"""
    productId = serializers.IntegerField(source='product_id', read_only=True)
    mediaType = serializers.CharField(source='media_type')
    sortOrder = serializers.IntegerField(source='sort_order', required=False, default=0)
    url = serializers.SerializerMethodField()

    class Meta:
        model = Media
        fields = ['id', 'productId', 'mediaType', 'url', 'file', 'title', 'description', 'sortOrder', 'created_at']
        extra_kwargs = {
            'file': {'write_only': False, 'required': False}
        }

    def get_url(self, obj):
        """自定义URL字段，优先返回文件地址，其次返回外部链接"""
        return obj.get_media_url()


class HarvestQualitySerializer(serializers.ModelSerializer):
    """采收质量序列化器，将数据库字段名转为驼峰命名"""
    productId = serializers.IntegerField(source='product_id', read_only=True)
    startDate = serializers.DateField(source='harvest_start_date')
    endDate = serializers.DateField(source='harvest_end_date')
    sugarContent = serializers.DecimalField(source='sugar_content', max_digits=5, decimal_places=2)
    suitableCrowd = serializers.CharField(source='suitable_crowd')
    qualitySummary = serializers.CharField(source='quality_summary')

    class Meta:
        model = HarvestQuality
        fields = ['id', 'productId', 'startDate', 'endDate', 'sugarContent', 'weight', 'taste', 'suitableCrowd', 'qualitySummary', 'created_at', 'updated_at']


class ProductDetailSerializer(serializers.ModelSerializer):
    """产品详情序列化器，聚合产品的图片、视频和采收质量信息"""
    images = serializers.SerializerMethodField()
    videos = serializers.SerializerMethodField()
    harvest = serializers.SerializerMethodField()
    plantingLocation = serializers.CharField(source='planting_location')
    plantingDate = serializers.DateField(source='planting_date')
    createdAt = serializers.DateTimeField(source='created_at')
    updatedAt = serializers.DateTimeField(source='updated_at')

    class Meta:
        model = Product
        fields = ['id', 'name', 'code', 'plantingLocation', 'plantingDate', 'images', 'videos', 'harvest', 'createdAt', 'updatedAt']

    def get_images(self, obj):
        """获取产品关联的所有图片URL列表"""
        media = obj.media.filter(media_type='image')
        return [m.get_media_url() for m in media]

    def get_videos(self, obj):
        """获取产品关联的所有视频URL列表"""
        media = obj.media.filter(media_type='video')
        return [m.get_media_url() for m in media]

    def get_harvest(self, obj):
        """获取产品的采收质量信息，若不存在则返回None"""
        try:
            harvest = obj.harvest_quality
            return {
                'startDate': harvest.harvest_start_date,
                'endDate': harvest.harvest_end_date,
                'sugarContent': float(harvest.sugar_content),
                'weight': float(harvest.weight),
                'taste': harvest.taste,
                'suitableCrowd': harvest.suitable_crowd,
                'qualitySummary': harvest.quality_summary,
            }
        except HarvestQuality.DoesNotExist:
            return None

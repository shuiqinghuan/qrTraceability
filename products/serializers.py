from rest_framework import serializers
from .models import Product, Media, HarvestQuality


class ProductSerializer(serializers.ModelSerializer):
    plantingLocation = serializers.CharField(source='planting_location')
    plantingDate = serializers.DateField(source='planting_date')
    createdAt = serializers.DateTimeField(source='created_at', read_only=True)
    updatedAt = serializers.DateTimeField(source='updated_at', read_only=True)

    class Meta:
        model = Product
        fields = ['id', 'name', 'code', 'plantingLocation', 'plantingDate', 'createdAt', 'updatedAt']


class ProductListSerializer(serializers.ModelSerializer):
    plantingLocation = serializers.CharField(source='planting_location')
    plantingDate = serializers.DateField(source='planting_date')

    class Meta:
        model = Product
        fields = ['id', 'name', 'code', 'plantingLocation', 'plantingDate']


class MediaSerializer(serializers.ModelSerializer):
    productId = serializers.IntegerField(source='product_id', read_only=True)
    mediaType = serializers.CharField(source='media_type')
    sortOrder = serializers.IntegerField(source='sort_order', required=False, default=0)

    class Meta:
        model = Media
        fields = ['id', 'productId', 'mediaType', 'url', 'title', 'description', 'sortOrder', 'created_at']


class HarvestQualitySerializer(serializers.ModelSerializer):
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
        media = obj.media.filter(media_type='image')
        return [m.url for m in media]

    def get_videos(self, obj):
        media = obj.media.filter(media_type='video')
        return [m.url for m in media]

    def get_harvest(self, obj):
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

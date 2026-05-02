from rest_framework import serializers
from .models import Product, MediaFile


class MediaFileSerializer(serializers.ModelSerializer):
    url = serializers.SerializerMethodField()
    
    class Meta:
        model = MediaFile
        fields = ['id', 'url', 'media_type', 'filename', 'file_size', 'uploaded_at']
    
    def get_url(self, obj):
        request = self.context.get('request')
        if obj.file and request:
            return request.build_absolute_uri(obj.file.url)
        return None


class ProductListSerializer(serializers.ModelSerializer):
    class Meta:
        model = Product
        fields = ['id', 'name', 'code', 'planting_location', 'planting_date', 'images', 'video']


class ProductDetailSerializer(serializers.ModelSerializer):
    media_files = MediaFileSerializer(many=True, read_only=True)
    
    class Meta:
        model = Product
        fields = [
            'id', 'name', 'code', 'planting_location', 'planting_date',
            'images', 'video', 'media_files',
            'harvest_start_date', 'harvest_end_date', 'sugar_content',
            'weight', 'taste', 'quality', 'quality_summary', 'suitable_for',
            'created_at', 'updated_at'
        ]


class ProductCreateSerializer(serializers.ModelSerializer):
    class Meta:
        model = Product
        fields = [
            'name', 'code', 'planting_location', 'planting_date',
            'images', 'video', 'harvest_start_date', 'harvest_end_date',
            'sugar_content', 'weight', 'taste', 'quality',
            'quality_summary', 'suitable_for'
        ]
    
    def validate_code(self, value):
        if Product.objects.filter(code=value).exists():
            raise serializers.ValidationError("该品种编码已存在")
        return value

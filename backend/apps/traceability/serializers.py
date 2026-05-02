from rest_framework import serializers
from .models import Product, Media, HarvestQuality


class MediaSerializer(serializers.ModelSerializer):
    class Meta:
        model = Media
        fields = '__all__'


class HarvestQualitySerializer(serializers.ModelSerializer):
    class Meta:
        model = HarvestQuality
        fields = '__all__'


class ProductSerializer(serializers.ModelSerializer):
    media = MediaSerializer(many=True, read_only=True)
    harvest_quality = HarvestQualitySerializer(read_only=True)

    class Meta:
        model = Product
        fields = '__all__'

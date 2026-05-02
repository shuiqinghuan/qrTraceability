from rest_framework import viewsets, status
from rest_framework.decorators import action
from rest_framework.response import Response
from rest_framework.permissions import AllowAny, IsAuthenticated
from django_filters.rest_framework import DjangoFilterBackend

from .models import Product, MediaFile
from .serializers import (
    ProductListSerializer, ProductDetailSerializer,
    ProductCreateSerializer, MediaFileSerializer
)


class ProductViewSet(viewsets.ModelViewSet):
    queryset = Product.objects.all()
    filter_backends = [DjangoFilterBackend]
    filterset_fields = ['code', 'name']
    
    def get_serializer_class(self):
        if self.action == 'list':
            return ProductListSerializer
        elif self.action in ['create', 'update', 'partial_update']:
            return ProductCreateSerializer
        return ProductDetailSerializer
    
    def get_permissions(self):
        if self.action in ['list', 'retrieve', 'public_by_code']:
            return [AllowAny()]
        return [IsAuthenticated()]
    
    @action(detail=False, methods=['get'], url_path='public/code/(?P<code>[^/.]+)')
    def public_by_code(self, request, code=None):
        try:
            product = Product.objects.get(code=code)
            serializer = ProductDetailSerializer(product, context={'request': request})
            return Response(serializer.data)
        except Product.DoesNotExist:
            return Response(
                {'error': '产品不存在'},
                status=status.HTTP_404_NOT_FOUND
            )


class MediaFileViewSet(viewsets.ModelViewSet):
    queryset = MediaFile.objects.all()
    serializer_class = MediaFileSerializer
    permission_classes = [IsAuthenticated]
    
    def perform_create(self, serializer):
        serializer.save()
    
    def create(self, request, *args, **kwargs):
        file = request.FILES.get('file')
        if not file:
            return Response(
                {'error': '未上传文件'},
                status=status.HTTP_400_BAD_REQUEST
            )
        
        media_type = 'image' if file.content_type.startswith('image') else 'video'
        
        product_id = request.data.get('product_id')
        if not product_id:
            return Response(
                {'error': '缺少产品ID'},
                status=status.HTTP_400_BAD_REQUEST
            )
        
        try:
            product = Product.objects.get(id=product_id)
        except Product.DoesNotExist:
            return Response(
                {'error': '产品不存在'},
                status=status.HTTP_404_NOT_FOUND
            )
        
        media = MediaFile.objects.create(
            product=product,
            file=file,
            media_type=media_type,
            filename=file.name,
            file_size=file.size
        )
        
        serializer = self.get_serializer(media)
        return Response(serializer.data, status=status.HTTP_201_CREATED)

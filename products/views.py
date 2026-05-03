from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework import status
from rest_framework.pagination import PageNumberPagination
from rest_framework.parsers import MultiPartParser, FormParser
from django.shortcuts import get_object_or_404
from .models import Product, Media, HarvestQuality
from .serializers import (
    ProductSerializer,
    ProductListSerializer,
    ProductDetailSerializer,
    MediaSerializer,
    HarvestQualitySerializer,
)


def api_response(data=None, code=200, message='success'):
    return Response({
        'code': code,
        'message': message,
        'data': data
    }, status=code if code < 400 else 400)


class ProductPagination(PageNumberPagination):
    page_size = 10
    page_size_query_param = 'pageSize'
    page_query_param = 'page'


class ProductListView(APIView):
    def get(self, request):
        paginator = ProductPagination()
        products = Product.objects.all().order_by('-created_at')
        result_page = paginator.paginate_queryset(products, request)
        serializer = ProductListSerializer(result_page, many=True)
        return api_response({
            'total': products.count(),
            'page': paginator.page.number if hasattr(paginator, 'page') else 1,
            'pageSize': paginator.get_page_size(request),
            'list': serializer.data
        })

    def post(self, request):
        serializer = ProductSerializer(data=request.data)
        if serializer.is_valid():
            serializer.save()
            return api_response(serializer.data, code=201, message='创建成功')
        return api_response(code=400, message=list(serializer.errors.values())[0][0] if serializer.errors else '参数错误')


class ProductDetailView(APIView):
    def get(self, request, code_or_id):
        product = None
        try:
            product = Product.objects.get(code=code_or_id)
        except Product.DoesNotExist:
            if code_or_id.isdigit():
                try:
                    product = Product.objects.get(pk=int(code_or_id))
                except Product.DoesNotExist:
                    pass
        
        if product is None:
            return api_response(code=404, message='产品不存在')
        
        serializer = ProductDetailSerializer(product)
        return api_response(serializer.data)

    def put(self, request, code_or_id):
        if not code_or_id.isdigit():
            return api_response(code=400, message='更新产品需要使用产品ID')
        
        try:
            product = Product.objects.get(pk=int(code_or_id))
        except Product.DoesNotExist:
            return api_response(code=404, message='产品不存在')
        
        serializer = ProductSerializer(product, data=request.data)
        if serializer.is_valid():
            serializer.save()
            return api_response(serializer.data, message='更新成功')
        return api_response(code=400, message=list(serializer.errors.values())[0][0] if serializer.errors else '参数错误')

    def delete(self, request, code_or_id):
        if not code_or_id.isdigit():
            return api_response(code=400, message='删除产品需要使用产品ID')
        
        try:
            product = Product.objects.get(pk=int(code_or_id))
        except Product.DoesNotExist:
            return api_response(code=404, message='产品不存在')
        
        product.delete()
        return api_response(message='删除成功')


class ProductMediaListView(APIView):
    parser_classes = [MultiPartParser, FormParser]

    def get(self, request, productId):
        try:
            product = Product.objects.get(pk=productId)
        except Product.DoesNotExist:
            return api_response(code=404, message='产品不存在')

        media_type = request.query_params.get('type')
        media = product.media.all()
        if media_type:
            media = media.filter(media_type=media_type)

        serializer = MediaSerializer(media, many=True, context={'request': request})
        return api_response(serializer.data)

    def post(self, request, productId):
        try:
            product = Product.objects.get(pk=productId)
        except Product.DoesNotExist:
            return api_response(code=404, message='产品不存在')

        data = request.data.copy()
        serializer = MediaSerializer(data=data, context={'request': request})
        if serializer.is_valid():
            serializer.save(product=product)
            return api_response(serializer.data, code=201, message='上传成功')
        return api_response(code=400, message=list(serializer.errors.values())[0][0] if serializer.errors else '参数错误')


class MediaDetailView(APIView):
    def delete(self, request, pk):
        try:
            media = Media.objects.get(pk=pk)
        except Media.DoesNotExist:
            return api_response(code=404, message='多媒体不存在')
        
        media.delete()
        return api_response(message='删除成功')


class HarvestQualityView(APIView):
    def get(self, request, productId):
        try:
            product = Product.objects.get(pk=productId)
        except Product.DoesNotExist:
            return api_response(code=404, message='产品不存在')
        
        try:
            harvest = product.harvest_quality
        except HarvestQuality.DoesNotExist:
            return api_response(code=404, message='采收质量信息不存在')
        
        serializer = HarvestQualitySerializer(harvest)
        return api_response(serializer.data)

    def post(self, request, productId):
        try:
            product = Product.objects.get(pk=productId)
        except Product.DoesNotExist:
            return api_response(code=404, message='产品不存在')
        
        try:
            harvest = product.harvest_quality
            serializer = HarvestQualitySerializer(harvest, data=request.data)
        except HarvestQuality.DoesNotExist:
            serializer = HarvestQualitySerializer(data=request.data)
        
        if serializer.is_valid():
            serializer.save(product=product)
            return api_response(serializer.data, message='保存成功')
        return api_response(code=400, message=list(serializer.errors.values())[0][0] if serializer.errors else '参数错误')

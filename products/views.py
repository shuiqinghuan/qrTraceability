"""
产品API视图模块

提供农产品溯源系统的RESTful API接口，包括：
- 产品的增删改查（列表、详情、创建、更新、删除）
- 产品多媒体资源的管理（上传、查看、删除）
- 采收质量信息的查看与更新
"""

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
    """统一API响应格式，返回包含code、message、data的标准结构"""
    return Response({
        'code': code,
        'message': message,
        'data': data
    }, status=code if code < 400 else 400)


class ProductPagination(PageNumberPagination):
    """产品列表分页器，默认每页10条记录"""
    page_size = 10
    page_size_query_param = 'pageSize'
    page_query_param = 'page'


class ProductListView(APIView):
    """产品列表视图，支持分页查询和新建产品"""

    def get(self, request):
        """获取产品列表，按创建时间倒序排列"""
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
        """创建新产品"""
        serializer = ProductSerializer(data=request.data)
        if serializer.is_valid():
            serializer.save()
            return api_response(serializer.data, code=201, message='创建成功')
        return api_response(code=400, message=list(serializer.errors.values())[0][0] if serializer.errors else '参数错误')


class ProductDetailView(APIView):
    """产品详情视图，支持通过品种编码或产品ID查询、更新和删除"""

    def get(self, request, code_or_id):
        """获取产品详情，支持按品种编码或产品ID查找（优先按编码查找）"""
        product = None
        try:
            product = Product.objects.get(code=code_or_id)
        except Product.DoesNotExist:
            # 编码未匹配时，尝试按数字ID查找
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
        """更新产品信息，仅支持通过数字ID更新"""
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
        """删除产品，仅支持通过数字ID删除"""
        if not code_or_id.isdigit():
            return api_response(code=400, message='删除产品需要使用产品ID')

        try:
            product = Product.objects.get(pk=int(code_or_id))
        except Product.DoesNotExist:
            return api_response(code=404, message='产品不存在')

        product.delete()
        return api_response(message='删除成功')


class ProductMediaListView(APIView):
    """产品多媒体资源视图，支持查看和上传产品的图片/视频"""
    # 支持multipart表单数据和文件上传
    parser_classes = [MultiPartParser, FormParser]

    def get(self, request, productId):
        """获取产品的多媒体列表，支持通过type查询参数筛选媒体类型"""
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
        """为产品上传多媒体文件"""
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
    """单个多媒体资源视图，支持删除操作"""

    def delete(self, request, pk):
        """删除指定的多媒体资源"""
        try:
            media = Media.objects.get(pk=pk)
        except Media.DoesNotExist:
            return api_response(code=404, message='多媒体不存在')

        media.delete()
        return api_response(message='删除成功')


class HarvestQualityView(APIView):
    """采收质量信息视图，支持查看和创建/更新产品的采收质量数据"""

    def get(self, request, productId):
        """获取指定产品的采收质量信息"""
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
        """创建或更新采收质量信息，若已存在则更新，否则新建"""
        try:
            product = Product.objects.get(pk=productId)
        except Product.DoesNotExist:
            return api_response(code=404, message='产品不存在')

        try:
            harvest = product.harvest_quality
            # 已有记录，执行更新操作
            serializer = HarvestQualitySerializer(harvest, data=request.data)
        except HarvestQuality.DoesNotExist:
            # 无已有记录，执行新建操作
            serializer = HarvestQualitySerializer(data=request.data)

        if serializer.is_valid():
            serializer.save(product=product)
            return api_response(serializer.data, message='保存成功')
        return api_response(code=400, message=list(serializer.errors.values())[0][0] if serializer.errors else '参数错误')

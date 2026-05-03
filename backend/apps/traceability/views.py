"""
API视图模块。

提供产品信息的RESTful接口，支持列表查询、创建、详情查看、更新和删除操作。
"""

from rest_framework import generics
from .models import Product
from .serializers import ProductSerializer


class ProductListCreateView(generics.ListCreateAPIView):
    """产品列表与创建接口。GET返回产品列表，POST创建新产品。"""

    queryset = Product.objects.all()
    serializer_class = ProductSerializer


class ProductDetailView(generics.RetrieveUpdateDestroyAPIView):
    """产品详情接口。支持查看（GET）、更新（PUT/PATCH）和删除（DELETE）单个产品。"""

    queryset = Product.objects.all()
    serializer_class = ProductSerializer

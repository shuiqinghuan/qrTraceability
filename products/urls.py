from django.urls import path
from .views import (
    ProductListView,
    ProductDetailView,
    ProductMediaListView,
    MediaDetailView,
    HarvestQualityView,
)

urlpatterns = [
    path('products/', ProductListView.as_view(), name='product-list'),
    path('products/<str:code_or_id>/', ProductDetailView.as_view(), name='product-detail'),
    path('products/<int:productId>/media/', ProductMediaListView.as_view(), name='product-media-list'),
    path('media/<int:pk>/', MediaDetailView.as_view(), name='media-detail'),
    path('products/<int:productId>/harvest/', HarvestQualityView.as_view(), name='harvest-quality'),
]

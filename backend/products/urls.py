from django.urls import path, include
from rest_framework.routers import DefaultRouter
from .views import ProductViewSet, MediaFileViewSet

router = DefaultRouter()
router.register(r'products', ProductViewSet, basename='product')
router.register(r'media', MediaFileViewSet, basename='media')

urlpatterns = [
    path('', include(router.urls)),
]

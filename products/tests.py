"""
产品应用测试模块

覆盖模型、序列化器和API接口的功能测试，包括：
- ModelTestCase: 产品、多媒体、采收质量模型的创建与关系测试
- ProductListViewTest: 产品列表和创建接口测试
- ProductDetailViewTest: 产品详情、更新、删除接口测试
- ProductMediaListViewTest: 多媒体资源列表和上传接口测试
- MediaDetailViewTest: 多媒体资源删除接口测试
- HarvestQualityViewTest: 采收质量信息查看和创建/更新接口测试
"""

import datetime
from decimal import Decimal
from django.test import TestCase, override_settings
from django.core.files.uploadedfile import SimpleUploadedFile
from rest_framework.test import APIClient
from .models import Product, Media, HarvestQuality


@override_settings(STATICFILES_DIRS=[])
class ModelTestCase(TestCase):
    """数据模型的单元测试"""

    def _create_product(self, name='枣甜5号', code='4395'):
        return Product.objects.create(
            name=name,
            code=code,
            planting_location='山东省济南市历城区农业示范园',
            planting_date='2024-03-15',
        )

    def _create_harvest(self, product):
        return HarvestQuality.objects.create(
            product=product,
            harvest_start_date='2024-07-01',
            harvest_end_date='2024-07-15',
            sugar_content=Decimal('15.50'),
            weight=Decimal('280.50'),
            taste='肉质细腻，汁多味甜',
            suitable_crowd='老少皆宜',
            quality_summary='品质优良',
        )

    def test_product_create_and_str(self):
        """测试产品创建和__str__方法"""
        product = self._create_product()
        self.assertEqual(str(product), '枣甜5号(4395)')
        self.assertEqual(Product.objects.count(), 1)

    def test_product_unique_code(self):
        """测试品种编码唯一约束"""
        self._create_product()
        with self.assertRaises(Exception):
            self._create_product(code='4395')

    def test_media_create(self):
        """测试多媒体资源创建和排序"""
        product = self._create_product()
        m1 = Media.objects.create(
            product=product, media_type='image',
            url='https://example.com/1.jpg', sort_order=1
        )
        m2 = Media.objects.create(
            product=product, media_type='video',
            url='https://example.com/1.mp4', sort_order=0
        )
        media_list = list(product.media.all())
        self.assertEqual(media_list[0], m2)  # sort_order=0 排在前面
        self.assertEqual(media_list[1], m1)

    def test_media_get_url_file_priority(self):
        """测试媒体URL优先返回上传文件地址"""
        product = self._create_product()
        media = Media.objects.create(
            product=product, media_type='image',
            url='https://example.com/fallback.jpg',
            file='products/2024/01/01/uploaded.jpg',
        )
        self.assertEqual(media.get_media_url(), '/media/products/2024/01/01/uploaded.jpg')

    def test_media_get_url_external(self):
        """测试无上传文件时返回外部URL"""
        product = self._create_product()
        media = Media.objects.create(
            product=product, media_type='image',
            url='https://example.com/img.jpg',
        )
        self.assertEqual(media.get_media_url(), 'https://example.com/img.jpg')

    def test_harvest_quality_one_to_one(self):
        """测试采收质量与产品的一对一关系"""
        product = self._create_product()
        self._create_harvest(product)
        self.assertEqual(product.harvest_quality.sugar_content, Decimal('15.50'))

    def test_product_cascade_delete(self):
        """测试删除产品时级联删除关联数据"""
        product = self._create_product()
        Media.objects.create(product=product, media_type='image', url='test.jpg')
        self._create_harvest(product)
        product.delete()
        self.assertEqual(Media.objects.count(), 0)
        self.assertEqual(HarvestQuality.objects.count(), 0)


@override_settings(STATICFILES_DIRS=[])
class ProductListViewTest(TestCase):
    """产品列表与创建接口测试"""

    def setUp(self):
        self.client = APIClient()

    def test_empty_list(self):
        """测试空列表返回"""
        resp = self.client.get('/api/products/')
        self.assertEqual(resp.status_code, 200)
        body = resp.json()
        self.assertEqual(body['data']['total'], 0)
        self.assertEqual(body['data']['list'], [])

    def test_list_with_pagination(self):
        """测试分页功能"""
        for i in range(15):
            Product.objects.create(
                name=f'品种{i}', code=f'CODE{i}',
                planting_location='测试地点', planting_date='2024-01-01',
            )
        resp = self.client.get('/api/products/')
        body = resp.json()
        self.assertEqual(body['data']['total'], 15)
        self.assertEqual(len(body['data']['list']), 10)

        resp2 = self.client.get('/api/products/?page=2')
        self.assertEqual(len(resp2.json()['data']['list']), 5)

    def test_create_product(self):
        """测试创建产品"""
        data = {
            'name': '枣甜5号', 'code': '4395',
            'plantingLocation': '山东省济南市',
            'plantingDate': '2024-03-15',
        }
        resp = self.client.post('/api/products/', data, format='json')
        self.assertEqual(resp.status_code, 200)
        self.assertEqual(resp.json()['data']['name'], '枣甜5号')

    def test_create_product_missing_field(self):
        """测试创建产品缺少必填字段"""
        resp = self.client.post('/api/products/', {'name': '只有名称'}, format='json')
        self.assertEqual(resp.json()['code'], 400)


@override_settings(STATICFILES_DIRS=[])
class ProductDetailViewTest(TestCase):
    """产品详情、更新、删除接口测试"""

    def setUp(self):
        self.client = APIClient()
        self.product = Product.objects.create(
            name='枣甜5号', code='4395',
            planting_location='山东省济南市历城区农业示范园',
            planting_date='2024-03-15',
        )

    def test_get_by_code(self):
        """测试通过品种编码查询产品详情"""
        resp = self.client.get('/api/products/4395/')
        self.assertEqual(resp.status_code, 200)
        self.assertEqual(resp.json()['data']['name'], '枣甜5号')

    def test_get_by_id(self):
        """测试通过数字ID查询产品详情"""
        resp = self.client.get(f'/api/products/{self.product.pk}/')
        self.assertEqual(resp.status_code, 200)

    def test_get_not_found(self):
        """测试查询不存在的产品"""
        resp = self.client.get('/api/products/NONEXIST/')
        self.assertEqual(resp.json()['code'], 404)

    def test_update_product(self):
        """测试更新产品信息"""
        resp = self.client.put(
            f'/api/products/{self.product.pk}/',
            {'name': '新名称', 'code': 'NEW1',
             'plantingLocation': '新地点', 'plantingDate': '2025-01-01'},
            format='json',
        )
        self.assertEqual(resp.json()['message'], '更新成功')
        self.product.refresh_from_db()
        self.assertEqual(self.product.name, '新名称')

    def test_update_with_code_returns_error(self):
        """测试通过品种编码更新产品时返回错误"""
        resp = self.client.put(
            '/api/products/4395/', {'name': 'x'}, format='json'
        )
        self.assertEqual(resp.json()['code'], 400)

    def test_delete_product(self):
        """测试删除产品"""
        resp = self.client.delete(f'/api/products/{self.product.pk}/')
        self.assertEqual(resp.json()['message'], '删除成功')
        self.assertFalse(Product.objects.filter(pk=self.product.pk).exists())

    def test_delete_with_code_returns_error(self):
        """测试通过品种编码删除产品时返回错误"""
        resp = self.client.delete('/api/products/4395/')
        self.assertEqual(resp.json()['code'], 400)

    def test_detail_includes_harvest_and_media(self):
        """测试详情接口包含采收质量和多媒体数据"""
        Media.objects.create(
            product=self.product, media_type='image', url='https://example.com/1.jpg'
        )
        HarvestQuality.objects.create(
            product=self.product,
            harvest_start_date='2024-07-01', harvest_end_date='2024-07-15',
            sugar_content=Decimal('15.50'), weight=Decimal('280.50'),
            taste='甜', suitable_crowd='所有人群', quality_summary='优',
        )
        resp = self.client.get('/api/products/4395/')
        data = resp.json()['data']
        self.assertEqual(len(data['images']), 1)
        self.assertEqual(data['harvest']['sugarContent'], 15.5)


@override_settings(STATICFILES_DIRS=[])
class ProductMediaListViewTest(TestCase):
    """产品多媒体资源列表和上传接口测试"""

    def setUp(self):
        self.client = APIClient()
        self.product = Product.objects.create(
            name='测试品种', code='T001',
            planting_location='测试地点', planting_date='2024-01-01',
        )

    def test_list_media_empty(self):
        """测试空多媒体列表"""
        resp = self.client.get(f'/api/products/{self.product.pk}/media/')
        self.assertEqual(resp.json()['data'], [])

    def test_list_media_filter_by_type(self):
        """测试按类型筛选多媒体"""
        Media.objects.create(product=self.product, media_type='image', url='img.jpg')
        Media.objects.create(product=self.product, media_type='video', url='vid.mp4')
        resp = self.client.get(f'/api/products/{self.product.pk}/media/?type=image')
        self.assertEqual(len(resp.json()['data']), 1)

    def test_upload_media_with_url(self):
        """测试通过URL方式上传多媒体"""
        resp = self.client.post(
            f'/api/products/{self.product.pk}/media/',
            {'mediaType': 'image', 'url': 'https://example.com/img.jpg'},
            format='json',
        )
        self.assertEqual(resp.status_code, 200)

    def test_media_product_not_found(self):
        """测试对不存在的产品上传多媒体"""
        resp = self.client.get('/api/products/99999/media/')
        self.assertEqual(resp.json()['code'], 404)


@override_settings(STATICFILES_DIRS=[])
class MediaDetailViewTest(TestCase):
    """多媒体资源删除接口测试"""

    def setUp(self):
        self.client = APIClient()
        product = Product.objects.create(
            name='测试品种', code='T002',
            planting_location='测试地点', planting_date='2024-01-01',
        )
        self.media = Media.objects.create(
            product=product, media_type='image', url='test.jpg'
        )

    def test_delete_media(self):
        """测试删除多媒体资源"""
        resp = self.client.delete(f'/api/media/{self.media.pk}/')
        self.assertEqual(resp.json()['message'], '删除成功')
        self.assertFalse(Media.objects.filter(pk=self.media.pk).exists())

    def test_delete_media_not_found(self):
        """测试删除不存在的多媒体"""
        resp = self.client.delete('/api/media/99999/')
        self.assertEqual(resp.json()['code'], 404)


@override_settings(STATICFILES_DIRS=[])
class HarvestQualityViewTest(TestCase):
    """采收质量信息查看和创建/更新接口测试"""

    def setUp(self):
        self.client = APIClient()
        self.product = Product.objects.create(
            name='测试品种', code='T003',
            planting_location='测试地点', planting_date='2024-01-01',
        )
        self.valid_data = {
            'startDate': '2024-07-01', 'endDate': '2024-07-15',
            'sugarContent': '15.50', 'weight': '280.50',
            'taste': '肉质细腻，汁多味甜',
            'suitableCrowd': '老少皆宜',
            'qualitySummary': '品质优良',
        }

    def test_get_harvest_not_exists(self):
        """测试获取不存在的采收质量信息"""
        resp = self.client.get(f'/api/products/{self.product.pk}/harvest/')
        self.assertEqual(resp.json()['code'], 404)

    def test_create_harvest(self):
        """测试创建采收质量信息"""
        resp = self.client.post(
            f'/api/products/{self.product.pk}/harvest/',
            self.valid_data, format='json',
        )
        self.assertEqual(resp.json()['message'], '保存成功')
        self.assertEqual(self.product.harvest_quality.sugar_content, Decimal('15.50'))

    def test_update_harvest(self):
        """测试更新已有的采收质量信息"""
        self.client.post(
            f'/api/products/{self.product.pk}/harvest/',
            self.valid_data, format='json',
        )
        updated = {**self.valid_data, 'sugarContent': '18.00'}
        self.client.post(
            f'/api/products/{self.product.pk}/harvest/',
            updated, format='json',
        )
        self.product.harvest_quality.refresh_from_db()
        self.assertEqual(self.product.harvest_quality.sugar_content, Decimal('18.00'))

    def test_harvest_product_not_found(self):
        """测试对不存在的产品创建采收质量"""
        resp = self.client.post('/api/products/99999/harvest/', self.valid_data, format='json')
        self.assertEqual(resp.json()['code'], 404)

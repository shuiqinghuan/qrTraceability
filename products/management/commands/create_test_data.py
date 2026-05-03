"""
创建测试数据管理命令模块

用于快速生成产品的测试数据，包括产品信息、多媒体资源和采收质量信息。
使用方式：python manage.py create_test_data
"""

from django.core.management.base import BaseCommand
from products.models import Product, Media, HarvestQuality
from decimal import Decimal
from datetime import date


class Command(BaseCommand):
    """创建初始测试数据的Django管理命令"""
    help = '创建初始测试数据'

    def handle(self, *args, **options):
        """执行命令：依次创建产品、多媒体资源和采收质量测试数据"""
        # 使用get_or_create避免重复创建，以品种编码为唯一标识
        product, created = Product.objects.get_or_create(
            code='4395',
            defaults={
                'name': '枣甜5号',
                'planting_location': '山东省济南市历城区农业示范园',
                'planting_date': date(2024, 3, 15),
            }
        )

        if created:
            self.stdout.write(self.style.SUCCESS(f'创建产品: {product.name}'))
        else:
            self.stdout.write(self.style.WARNING(f'产品已存在: {product.name}'))

        # 定义产品关联的多媒体资源数据
        media_data = [
            {
                'media_type': 'image',
                'url': 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=新鲜红枣果实特写&image_size=landscape_16_9',
                'title': '产品图片1',
                'sort_order': 1,
            },
            {
                'media_type': 'image',
                'url': 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=红枣种植园风景&image_size=landscape_16_9',
                'title': '产品图片2',
                'sort_order': 2,
            },
            {
                'media_type': 'image',
                'url': 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=新鲜水果采摘场景&image_size=landscape_16_9',
                'title': '产品图片3',
                'sort_order': 3,
            },
        ]

        # 逐条创建多媒体资源，以产品和URL为唯一标识避免重复
        for media_item in media_data:
            media, created = Media.objects.get_or_create(
                product=product,
                url=media_item['url'],
                defaults=media_item
            )
            if created:
                self.stdout.write(self.style.SUCCESS(f'创建多媒体: {media.title}'))

        # 创建采收质量信息，每个产品只能有一条记录
        harvest, created = HarvestQuality.objects.get_or_create(
            product=product,
            defaults={
                'harvest_start_date': date(2024, 7, 1),
                'harvest_end_date': date(2024, 7, 15),
                'sugar_content': Decimal('10.50'),
                'weight': Decimal('18.50'),
                'taste': '肉质细腻，汁多味甜，口感爽脆',
                'suitable_crowd': '老少皆宜，特别适合血糖稳定人群',
                'quality_summary': '果实饱满，色泽鲜艳，糖度适中，品质优良',
            }
        )

        if created:
            self.stdout.write(self.style.SUCCESS(f'创建采收质量信息: {product.name}'))
        else:
            self.stdout.write(self.style.WARNING(f'采收质量信息已存在: {product.name}'))

        self.stdout.write(self.style.SUCCESS('测试数据创建完成!'))

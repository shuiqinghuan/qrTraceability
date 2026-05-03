from django.core.management.base import BaseCommand
from products.models import Product, Media, HarvestQuality
from decimal import Decimal
from datetime import date


class Command(BaseCommand):
    help = '创建初始测试数据'

    def handle(self, *args, **options):
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

        for media_item in media_data:
            media, created = Media.objects.get_or_create(
                product=product,
                url=media_item['url'],
                defaults=media_item
            )
            if created:
                self.stdout.write(self.style.SUCCESS(f'创建多媒体: {media.title}'))

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

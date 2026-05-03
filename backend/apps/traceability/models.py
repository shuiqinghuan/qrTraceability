"""
溯源数据模型模块。

定义农产品溯源系统的核心数据表，包括产品信息、多媒体资源和采收质量信息。
各模型之间通过外键/一对一关系关联，构成完整的溯源数据链路。
"""

from django.db import models


class Product(models.Model):
    """产品信息模型，记录农产品的品种、定植地点等基础信息。"""
    name = models.CharField('品种名称', max_length=100)
    code = models.CharField('品种编码', max_length=50, unique=True)
    planting_location = models.CharField('定植地点', max_length=255)
    planting_date = models.DateField('定植时间')
    created_at = models.DateTimeField('创建时间', auto_now_add=True)
    updated_at = models.DateTimeField('更新时间', auto_now=True)

    class Meta:
        db_table = 'product'
        verbose_name = '产品信息'
        verbose_name_plural = '产品信息'

    def __str__(self):
        return f'{self.name}({self.code})'


class Media(models.Model):
    """多媒体资源模型，存储产品相关的图片和视频信息。"""

    # 媒体类型枚举：图片或视频
    MEDIA_TYPE_CHOICES = [
        ('image', '图片'),
        ('video', '视频'),
    ]

    # 级联删除：当关联产品被删除时，对应的多媒体资源也会被删除
    product = models.ForeignKey(Product, on_delete=models.CASCADE, related_name='media', verbose_name='关联产品')
    media_type = models.CharField('媒体类型', max_length=20, choices=MEDIA_TYPE_CHOICES)
    url = models.CharField('媒体URL', max_length=500)
    title = models.CharField('媒体标题', max_length=100, blank=True, null=True)
    description = models.TextField('媒体描述', blank=True, null=True)
    sort_order = models.IntegerField('排序顺序', default=0)
    created_at = models.DateTimeField('创建时间', auto_now_add=True)

    class Meta:
        db_table = 'media'
        verbose_name = '多媒体信息'
        verbose_name_plural = '多媒体信息'
        # 先按排序字段升序，再按主键排序，确保展示顺序稳定
        ordering = ['sort_order', 'id']

    def __str__(self):
        return f'{self.get_media_type_display()} - {self.product.name}'


class HarvestQuality(models.Model):
    """采收质量信息模型，记录产品的采收周期、品质检测数据等。"""

    # 一对一关系：每个产品仅对应一条采收质量记录
    product = models.OneToOneField(Product, on_delete=models.CASCADE, related_name='harvest_quality', verbose_name='关联产品')
    harvest_start_date = models.DateField('采收起始时间')
    harvest_end_date = models.DateField('采收终止时间')
    sugar_content = models.DecimalField('糖度(Brix)', max_digits=5, decimal_places=2)
    weight = models.DecimalField('单果重量(克)', max_digits=8, decimal_places=2)
    taste = models.CharField('口感描述', max_length=255)
    suitable_crowd = models.CharField('适应人群', max_length=255)
    quality_summary = models.TextField('品质小结')
    created_at = models.DateTimeField('创建时间', auto_now_add=True)
    updated_at = models.DateTimeField('更新时间', auto_now=True)

    class Meta:
        db_table = 'harvest_quality'
        verbose_name = '采收质量信息'
        verbose_name_plural = '采收质量信息'

    def __str__(self):
        return f'{self.product.name} - 采收质量'

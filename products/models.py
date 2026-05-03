"""
产品数据模型模块

定义农产品溯源系统的核心数据模型，包括：
- Product: 产品基本信息（品种名称、编码、定植地点等）
- Media: 产品关联的多媒体资源（图片、视频）
- HarvestQuality: 产品的采收质量信息（糖度、重量、口感等）
"""

from django.db import models


class Product(models.Model):
    """产品模型，存储农产品的基本信息"""
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
    """多媒体模型，存储产品关联的图片和视频资源"""
    MEDIA_TYPE_CHOICES = [
        ('image', '图片'),
        ('video', '视频'),
    ]

    product = models.ForeignKey(Product, on_delete=models.CASCADE, related_name='media', verbose_name='关联产品')
    media_type = models.CharField('媒体类型', max_length=20, choices=MEDIA_TYPE_CHOICES)
    url = models.CharField('媒体URL', max_length=500, blank=True, null=True)
    file = models.FileField('媒体文件', upload_to='products/%Y/%m/%d/', blank=True, null=True)
    title = models.CharField('媒体标题', max_length=100, blank=True, null=True)
    description = models.TextField('媒体描述', blank=True, null=True)
    sort_order = models.IntegerField('排序顺序', default=0)
    created_at = models.DateTimeField('创建时间', auto_now_add=True)

    class Meta:
        db_table = 'media'
        verbose_name = '多媒体信息'
        verbose_name_plural = '多媒体信息'
        ordering = ['sort_order', 'id']

    def __str__(self):
        return f'{self.get_media_type_display()} - {self.product.name}'

    def get_media_url(self):
        """获取媒体资源的访问地址，优先返回上传文件的URL，其次返回外部链接"""
        if self.file:
            return self.file.url
        return self.url


class HarvestQuality(models.Model):
    """采收质量模型，记录产品的品质检测数据，与产品为一对一关系"""
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

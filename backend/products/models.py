from django.db import models
from django.contrib.auth.models import User


class Product(models.Model):
    name = models.CharField('品种名称', max_length=100)
    code = models.CharField('品种编码', max_length=50, unique=True)
    planting_location = models.CharField('定植地点', max_length=200)
    planting_date = models.DateField('定植时间')
    images = models.JSONField('产品图片', default=list, blank=True)
    video = models.URLField('视频链接', blank=True, null=True)
    
    harvest_start_date = models.DateField('采收起始时间', blank=True, null=True)
    harvest_end_date = models.DateField('采收终止时间', blank=True, null=True)
    sugar_content = models.DecimalField('糖度', max_digits=4, decimal_places=2, null=True, blank=True)
    weight = models.DecimalField('单果重量(克)', max_digits=6, decimal_places=2, null=True, blank=True)
    taste = models.CharField('口感描述', max_length=200, blank=True)
    quality = models.CharField('品质等级', max_length=50, blank=True)
    
    quality_summary = models.TextField('品质小结', blank=True)
    suitable_for = models.JSONField('适应人群', default=list, blank=True)
    
    created_at = models.DateTimeField('创建时间', auto_now_add=True)
    updated_at = models.DateTimeField('更新时间', auto_now=True)
    
    class Meta:
        db_table = 'products'
        ordering = ['-created_at']
        verbose_name = '产品'
        verbose_name_plural = '产品列表'
    
    def __str__(self):
        return f"{self.name} ({self.code})"


class MediaFile(models.Model):
    MEDIA_TYPES = [
        ('image', '图片'),
        ('video', '视频'),
    ]
    
    product = models.ForeignKey(Product, on_delete=models.CASCADE, related_name='media_files')
    file = models.FileField('文件', upload_to='media/%Y/%m/')
    media_type = models.CharField('媒体类型', max_length=10, choices=MEDIA_TYPES)
    filename = models.CharField('文件名', max_length=255)
    file_size = models.IntegerField('文件大小(字节)', default=0)
    
    uploaded_at = models.DateTimeField('上传时间', auto_now_add=True)
    
    class Meta:
        db_table = 'media_files'
        ordering = ['-uploaded_at']
        verbose_name = '媒体文件'
        verbose_name_plural = '媒体文件列表'
    
    def __str__(self):
        return f"{self.product.name} - {self.filename}"


class AdminUser(models.Model):
    user = models.OneToOneField(User, on_delete=models.CASCADE)
    phone = models.CharField('手机号', max_length=20, blank=True)
    avatar = models.URLField('头像', blank=True)
    
    class Meta:
        db_table = 'admin_users'
        verbose_name = '管理员'
        verbose_name_plural = '管理员列表'

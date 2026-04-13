<template>
  <section class="product-section quality-section">
    <div class="section-header">
      <h2 class="section-title">
        <i class="icon-quality"></i> 品质信息
      </h2>
    </div>
    <div class="section-content">
      <div class="quality-grid">
        <div class="quality-item">
          <label class="quality-label">采收起始时间</label>
          <div class="quality-value">
            {{ quality.harvest_start_date ? formatDate(quality.harvest_start_date) : '未设置' }}
          </div>
        </div>
        <div class="quality-item">
          <label class="quality-label">采收终止时间</label>
          <div class="quality-value">
            {{ quality.harvest_end_date ? formatDate(quality.harvest_end_date) : '未设置' }}
          </div>
        </div>
        <div class="quality-item">
          <label class="quality-label">糖度</label>
          <div class="quality-value">
            {{ quality.sugar_content ? quality.sugar_content + '°Bx' : '未检测' }}
          </div>
        </div>
        <div class="quality-item">
          <label class="quality-label">重量</label>
          <div class="quality-value">
            {{ quality.weight ? quality.weight + 'g' : '未称重' }}
          </div>
        </div>
        <div class="quality-item full-width">
          <label class="quality-label">口感描述</label>
          <div class="quality-description">
            {{ quality.taste_description || '暂无口感描述' }}
          </div>
        </div>
        <div class="quality-item full-width">
          <label class="quality-label">适应人群</label>
          <div class="quality-description">
            {{ quality.suitable_for || '暂无适应人群描述' }}
          </div>
        </div>
        <div class="quality-item full-width">
          <label class="quality-label">品质小结</label>
          <div class="quality-summary">
            {{ quality.quality_summary || '暂无品质小结' }}
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  quality: {
    type: Object,
    required: true,
    default: () => ({
      harvest_start_date: null,
      harvest_end_date: null,
      sugar_content: 0,
      weight: 0,
      taste_description: '',
      suitable_for: '',
      quality_summary: ''
    })
  }
})

const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}
</script>

<style scoped>
.quality-section {
  margin-bottom: 24px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  overflow: hidden;
}

.section-header {
  background: #f5f5f5;
  padding: 16px 20px;
  border-bottom: 1px solid #e0e0e0;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.section-content {
  padding: 20px;
}

.quality-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
}

.quality-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.quality-item.full-width {
  grid-column: 1 / -1;
}

.quality-label {
  font-size: 14px;
  color: #666;
  font-weight: 500;
}

.quality-value {
  font-size: 16px;
  color: #333;
  font-weight: 500;
  padding: 8px 0;
}

.quality-description,
.quality-summary {
  font-size: 14px;
  color: #555;
  line-height: 1.5;
  padding: 12px;
  background: #f9f9f9;
  border-radius: 4px;
}

.quality-description {
  border-left: 4px solid #FF9800;
}

.quality-summary {
  border-left: 4px solid #9C27B0;
}

@media (max-width: 768px) {
  .quality-grid {
    grid-template-columns: 1fr;
  }
  
  .section-content {
    padding: 16px;
  }
}

@media (max-width: 480px) {
  .quality-item {
    margin-bottom: 16px;
  }
}
</style>
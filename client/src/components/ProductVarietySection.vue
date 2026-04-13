<template>
  <section class="product-section variety-section">
    <div class="section-header">
      <h2 class="section-title">
        <i class="icon-seed"></i> 品种信息
      </h2>
    </div>
    <div class="section-content">
      <div class="variety-grid">
        <div class="variety-item">
          <label class="variety-label">品种名称</label>
          <div class="variety-value">{{ seedInfo.name }}</div>
        </div>
        <div class="variety-item">
          <label class="variety-label">品种类型</label>
          <div class="variety-value">{{ seedInfo.variety }}</div>
        </div>
        <div class="variety-item">
          <label class="variety-label">品种编码</label>
          <div class="variety-value">{{ seedInfo.variety_code || '未设置' }}</div>
        </div>
        <div class="variety-item">
          <label class="variety-label">品种描述</label>
          <div class="variety-description">{{ seedInfo.description || '暂无描述' }}</div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
defineProps({
  seedInfo: {
    type: Object,
    required: true,
    default: () => ({
      name: '',
      variety: '',
      variety_code: '',
      description: ''
    })
  }
})
</script>

<style scoped>
.variety-section {
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

.variety-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
}

.variety-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.variety-label {
  font-size: 14px;
  color: #666;
  font-weight: 500;
}

.variety-value {
  font-size: 16px;
  color: #333;
  font-weight: 500;
  padding: 8px 0;
}

.variety-description {
  font-size: 14px;
  color: #555;
  line-height: 1.5;
  padding: 12px;
  background: #f9f9f9;
  border-radius: 4px;
  border-left: 4px solid #4CAF50;
}

@media (max-width: 768px) {
  .variety-grid {
    grid-template-columns: 1fr;
  }
  
  .section-content {
    padding: 16px;
  }
}
</style>
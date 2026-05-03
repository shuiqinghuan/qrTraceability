<template>
  <div class="product-info card">
    <h2 class="module-title">
      <span class="icon">🌱</span>
      产品基本信息
    </h2>
    <div class="info-content">
      <div class="info-item product-name">
        <span class="label">品种名称</span>
        <span class="value highlight">{{ product.name }}</span>
        <span class="code-badge">编码: {{ product.code }}</span>
      </div>

      <div class="info-item">
        <span class="icon-label">
          <span class="icon">📍</span>
          <span class="label">定植地点</span>
        </span>
        <span class="value">{{ product.plantingLocation }}</span>
      </div>

      <div class="info-item">
        <span class="icon-label">
          <span class="icon">📅</span>
          <span class="label">定植时间</span>
        </span>
        <span class="value">{{ product.plantingDate }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
/**
 * 产品基本信息展示组件
 * 展示农产品的品种名称、编码、定植地点和定植时间
 */

// 接收父组件传入的产品基础信息
defineProps({
  product: {
    type: Object,
    required: true,
    default: () => ({
      name: '',
      code: '',
      plantingLocation: '',
      plantingDate: ''
    })
  }
})
</script>

<style scoped>
/* 组件根容器 */
.product-info {
  height: 100%;
}

.module-title {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--primary-color);
  border-bottom: 2px solid var(--secondary-color);
  padding-bottom: 12px;
  margin-bottom: 20px;
}

.module-title .icon {
  font-size: 24px;
}

/* 信息内容区域 */
.info-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

/* 产品名称高亮卡片 */
.info-item.product-name {
  background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
  padding: 16px;
  border-radius: 8px;
  color: white;
}

.product-name .label {
  font-size: 12px;
  opacity: 0.9;
}

.product-name .value {
  font-size: 24px;
  font-weight: bold;
}

/* 产品编码徽章 */
.product-name .code-badge {
  display: inline-block;
  background: rgba(255, 255, 255, 0.2);
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  margin-top: 8px;
  width: fit-content;
}

.icon-label {
  display: flex;
  align-items: center;
  gap: 6px;
}

.icon-label .icon {
  font-size: 16px;
}

.label {
  font-size: 12px;
  color: var(--secondary-text);
}

.value {
  font-size: 14px;
  color: var(--text-color);
}

.info-item:not(.product-name) .value {
  padding-left: 22px;
}
</style>

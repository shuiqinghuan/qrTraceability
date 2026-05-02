<script setup lang="ts">
import type { Product } from '@/types'

defineProps<{
  product: Product
}>()

const formatDate = (dateStr: string) => {
  if (!dateStr) return '--'
  const date = new Date(dateStr)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}
</script>

<template>
  <div class="product-card">
    <div class="card-header">
      <div class="product-image" v-if="product.images?.[0]">
        <img :src="product.images[0]" :alt="product.name" />
      </div>
      <div class="product-image placeholder" v-else>
        <span class="i-ph-image text-4xl text-gray-300"></span>
      </div>
      <div class="product-info">
        <h3 class="product-name">{{ product.name }}</h3>
        <p class="product-code">品种编码: {{ product.code }}</p>
      </div>
    </div>
    
    <div class="card-body">
      <div class="info-row">
        <span class="info-icon">
          <span class="i-ph-map-pin"></span>
        </span>
        <div class="info-content">
          <span class="info-label">定植地点</span>
          <span class="info-value">{{ product.planting_location }}</span>
        </div>
      </div>
      
      <div class="info-row">
        <span class="info-icon">
          <span class="i-ph-calendar"></span>
        </span>
        <div class="info-content">
          <span class="info-label">定植时间</span>
          <span class="info-value">{{ formatDate(product.planting_date) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.product-card {
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(45, 90, 39, 0.08);
  overflow: hidden;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: linear-gradient(135deg, rgba(45, 90, 39, 0.05) 0%, rgba(232, 168, 56, 0.05) 100%);
}

.product-image {
  width: 80px;
  height: 80px;
  border-radius: 12px;
  overflow: hidden;
  flex-shrink: 0;
}

.product-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.product-image.placeholder {
  background: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
}

.product-info {
  flex: 1;
  min-width: 0;
}

.product-name {
  font-size: 20px;
  font-weight: 700;
  color: #333;
  margin-bottom: 4px;
}

.product-code {
  font-size: 12px;
  color: #999;
  background: rgba(45, 90, 39, 0.1);
  padding: 4px 10px;
  border-radius: 12px;
  display: inline-block;
}

.card-body {
  padding: 16px 20px 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.info-row {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.info-icon {
  width: 32px;
  height: 32px;
  background: rgba(45, 90, 39, 0.1);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #2D5A27;
  font-size: 16px;
  flex-shrink: 0;
}

.info-content {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.info-label {
  font-size: 12px;
  color: #999;
}

.info-value {
  font-size: 14px;
  color: #333;
}
</style>

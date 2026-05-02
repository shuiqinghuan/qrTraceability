<script setup lang="ts">
import { computed } from 'vue'
import { useProductStore } from '@/stores/product'
import { storeToRefs } from 'pinia'
import ProductCard from '@/components/ProductCard.vue'
import QualitySummary from '@/components/QualitySummary.vue'

const productStore = useProductStore()
const { currentProduct: product, loading, error } = storeToRefs(productStore)

const formattedDate = computed(() => {
  if (!product.value?.planting_date) return ''
  return product.value.planting_date
})
</script>

<template>
  <div class="home-view">
    <div class="container">
      <div v-if="loading" class="loading-state">
        <div class="skeleton" style="height: 200px; margin-bottom: 16px;"></div>
        <div class="skeleton" style="height: 120px; margin-bottom: 16px;"></div>
        <div class="skeleton" style="height: 100px;"></div>
      </div>

      <div v-else-if="error" class="error-state">
        <span class="i-ph-warning-circle text-5xl text-gray-400"></span>
        <p class="mt-4 text-gray-600">{{ error }}</p>
        <p class="text-sm text-gray-400 mt-2">请检查网络连接或稍后重试</p>
      </div>

      <div v-else-if="product" class="product-content">
        <section class="hero-section">
          <div class="hero-bg"></div>
          <div class="hero-content">
            <h1 class="hero-title">{{ product.name }}</h1>
            <p class="hero-code">编码: {{ product.code }}</p>
          </div>
        </section>

        <ProductCard :product="product" class="mt-4" />

        <section class="info-section mt-6">
          <h2 class="section-title">
            <span class="i-ph-map-pin text-accent"></span>
            种植信息
          </h2>
          <div class="info-grid">
            <div class="info-item">
              <span class="info-label">定植地点</span>
              <span class="info-value">{{ product.planting_location }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">定植时间</span>
              <span class="info-value">{{ formattedDate }}</span>
            </div>
          </div>
        </section>

        <section class="info-section mt-6" v-if="product.quality_summary">
          <QualitySummary 
            :summary="product.quality_summary" 
            :suitable-for="product.suitable_for" 
          />
        </section>

        <section class="quick-actions mt-6">
          <router-link to="/media" class="action-card">
            <span class="i-ph-images text-3xl text-primary"></span>
            <div class="action-info">
              <span class="action-title">产品相册</span>
              <span class="action-desc">{{ product.images?.length || 0 }} 张图片</span>
            </div>
            <span class="i-ph-caret-right text-gray-400"></span>
          </router-link>
          
          <router-link to="/quality" class="action-card">
            <span class="i-ph-chart-bar text-3xl text-accent"></span>
            <div class="action-info">
              <span class="action-title">质量报告</span>
              <span class="action-desc">查看详细数据</span>
            </div>
            <span class="i-ph-caret-right text-gray-400"></span>
          </router-link>
        </section>
      </div>

      <div v-else class="empty-state">
        <span class="i-ph-plant text-6xl text-gray-300"></span>
        <p class="mt-4 text-gray-500">暂无产品信息</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.home-view {
  padding-top: 24px;
}

.hero-section {
  position: relative;
  border-radius: 20px;
  overflow: hidden;
  margin-bottom: 16px;
}

.hero-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, #2D5A27 0%, #4A8A44 100%);
}

.hero-content {
  position: relative;
  padding: 40px 24px;
  color: white;
  text-align: center;
}

.hero-title {
  font-size: 32px;
  font-weight: 700;
  margin-bottom: 8px;
}

.hero-code {
  font-size: 14px;
  opacity: 0.9;
  background: rgba(255,255,255,0.2);
  padding: 6px 16px;
  border-radius: 20px;
  display: inline-block;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 16px;
  color: #333;
}

.info-section {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 4px 20px rgba(45, 90, 39, 0.08);
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-label {
  font-size: 12px;
  color: #999;
}

.info-value {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.quick-actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.action-card {
  display: flex;
  align-items: center;
  gap: 16px;
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 4px 20px rgba(45, 90, 39, 0.08);
  text-decoration: none;
  transition: all 0.2s ease;
}

.action-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 30px rgba(45, 90, 39, 0.15);
}

.action-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.action-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.action-desc {
  font-size: 12px;
  color: #999;
}

.loading-state,
.error-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 300px;
}

.mt-4 {
  margin-top: 16px;
}

.mt-6 {
  margin-top: 24px;
}
</style>

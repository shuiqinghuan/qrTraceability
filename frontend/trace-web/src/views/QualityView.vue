<script setup lang="ts">
import { computed } from 'vue'
import { useProductStore } from '@/stores/product'
import { storeToRefs } from 'pinia'
import QualityCard from '@/components/QualityCard.vue'

const productStore = useProductStore()
const { currentProduct: product, loading } = storeToRefs(productStore)

const qualityItems = computed(() => {
  if (!product.value) return []
  return [
    {
      icon: 'i-ph-chart-line-up',
      label: '糖度',
      value: product.value.sugar_content?.toString() || '--',
      unit: 'Brix'
    },
    {
      icon: 'i-ph-scales',
      label: '单果重量',
      value: product.value.weight?.toString() || '--',
      unit: '克'
    }
  ]
})

const harvestPeriod = computed(() => {
  if (!product.value) return '--'
  const start = product.value.harvest_start_date || ''
  const end = product.value.harvest_end_date || ''
  if (start && end) {
    return `${start} 至 ${end}`
  }
  return start || '--'
})
</script>

<template>
  <div class="quality-view">
    <div class="container">
      <div v-if="loading" class="loading-state">
        <div class="skeleton" style="height: 200px; margin-bottom: 16px;"></div>
        <div class="skeleton" style="height: 200px; margin-bottom: 16px;"></div>
        <div class="skeleton" style="height: 150px;"></div>
      </div>

      <div v-else-if="product">
        <section class="harvest-section">
          <h2 class="section-title">
            <span class="i-ph-calendar-check text-accent"></span>
            采收时间
          </h2>
          <div class="harvest-card">
            <div class="harvest-icon">
              <span class="i-ph-flower-lotus text-4xl text-primary"></span>
            </div>
            <div class="harvest-info">
              <span class="harvest-label">最佳采收期</span>
              <span class="harvest-value">{{ harvestPeriod }}</span>
            </div>
          </div>
        </section>

        <section class="quality-section mt-6">
          <h2 class="section-title">
            <span class="i-ph-chart-bar text-primary"></span>
            质量指标
          </h2>
          <div class="quality-grid">
            <QualityCard 
              v-for="item in qualityItems" 
              :key="item.label"
              :icon="item.icon"
              :label="item.label"
              :value="item.value"
              :unit="item.unit"
            />
          </div>
        </section>

        <section class="detail-section mt-6">
          <h2 class="section-title">
            <span class="i-ph-list-checks text-accent"></span>
            详细报告
          </h2>
          <div class="detail-card">
            <div class="detail-item">
              <span class="detail-label">口感描述</span>
              <span class="detail-value">{{ product.taste || '--' }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">品质等级</span>
              <span class="detail-badge" v-if="product.quality">{{ product.quality }}</span>
              <span class="detail-value" v-else>--</span>
            </div>
          </div>
        </section>

        <section class="suitable-section mt-6" v-if="product.suitable_for?.length">
          <h2 class="section-title">
            <span class="i-ph-users text-primary"></span>
            适应人群
          </h2>
          <div class="suitable-tags">
            <span 
              v-for="tag in product.suitable_for" 
              :key="tag"
              class="suitable-tag"
            >
              {{ tag }}
            </span>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<style scoped>
.quality-view {
  padding-top: 24px;
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

.harvest-section,
.quality-section,
.detail-section,
.suitable-section {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 4px 20px rgba(45, 90, 39, 0.08);
}

.harvest-card {
  display: flex;
  align-items: center;
  gap: 20px;
  background: linear-gradient(135deg, rgba(45, 90, 39, 0.05) 0%, rgba(232, 168, 56, 0.05) 100%);
  border-radius: 12px;
  padding: 24px;
}

.harvest-icon {
  width: 64px;
  height: 64px;
  background: white;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(45, 90, 39, 0.1);
}

.harvest-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.harvest-label {
  font-size: 12px;
  color: #999;
}

.harvest-value {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.quality-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.detail-card {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.detail-item:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.detail-label {
  font-size: 14px;
  color: #666;
}

.detail-value {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.detail-badge {
  font-size: 12px;
  padding: 4px 12px;
  background: linear-gradient(135deg, #2D5A27 0%, #4A8A44 100%);
  color: white;
  border-radius: 20px;
  font-weight: 500;
}

.suitable-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.suitable-tag {
  padding: 8px 16px;
  background: rgba(45, 90, 39, 0.08);
  color: #2D5A27;
  border-radius: 20px;
  font-size: 14px;
}

.loading-state {
  padding: 24px 0;
}

.mt-6 {
  margin-top: 24px;
}
</style>

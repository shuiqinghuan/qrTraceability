<script setup lang="ts">
defineProps<{
  summary: string
  suitableFor: string[]
}>()
</script>

<template>
  <div class="quality-summary">
    <div class="summary-header">
      <span class="i-ph-quotes text-accent"></span>
      <h3 class="summary-title">品质小结</h3>
    </div>
    <p class="summary-content">{{ summary }}</p>
    
    <div class="suitable-section" v-if="suitableFor?.length">
      <span class="suitable-label">适宜人群</span>
      <div class="suitable-tags">
        <span 
          v-for="tag in suitableFor" 
          :key="tag"
          class="suitable-tag"
        >
          {{ tag }}
        </span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.quality-summary {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 4px 20px rgba(45, 90, 39, 0.08);
  border-left: 4px solid #2D5A27;
}

.summary-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.summary-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.summary-content {
  font-size: 15px;
  line-height: 1.8;
  color: #555;
  margin-bottom: 16px;
}

.suitable-section {
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
}

.suitable-label {
  font-size: 12px;
  color: #999;
  margin-bottom: 8px;
  display: block;
}

.suitable-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.suitable-tag {
  padding: 6px 12px;
  background: rgba(232, 168, 56, 0.15);
  color: #C68A20;
  border-radius: 16px;
  font-size: 12px;
}
</style>

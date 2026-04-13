<template>
  <div class="interaction-item related-module">
    <h3 class="interaction-title">同地点产品</h3>
    <div class="interaction-content">
      <div v-if="relatedProducts.length === 0" class="no-related">
        <p>暂无同地点其他产品</p>
      </div>
      <div v-else class="related-list">
        <div v-for="related in relatedProducts" :key="related.id" class="related-item">
          <div class="related-name">{{ related.name }}</div>
          <div class="related-variety">{{ related.variety }}</div>
          <button 
            class="btn-view-related" 
            @click="goToProduct(related.id)"
          >
            查看详情
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue'

const props = defineProps({
  relatedProducts: {
    type: Array,
    required: true,
    default: () => []
  }
})

const emit = defineEmits(['view-related'])

const goToProduct = (id) => {
  emit('view-related', id)
}
</script>

<style scoped>
.related-module {
  background: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.interaction-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin: 0 0 16px 0;
  padding-bottom: 12px;
  border-bottom: 2px solid #4CAF50;
}

.no-related {
  text-align: center;
  padding: 20px;
  color: #999;
}

.related-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.related-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background: #f9f9f9;
  border-radius: 6px;
  border-left: 4px solid #4CAF50;
}

.related-name {
  font-weight: 500;
  color: #333;
}

.related-variety {
  font-size: 12px;
  color: #666;
}

.btn-view-related {
  background: #4CAF50;
  border: none;
  color: white;
  padding: 6px 12px;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-view-related:hover {
  background: #45a049;
}

@media (max-width: 768px) {
  .related-module {
    margin-bottom: 16px;
  }
  
  .related-item {
    flex-direction: column;
    gap: 8px;
    text-align: center;
  }
}

@media (max-width: 480px) {
  .related-list {
    gap: 8px;
  }
}
</style>
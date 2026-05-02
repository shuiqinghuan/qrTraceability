<script setup lang="ts">
defineProps<{
  icon: string
  label: string
  value: string
  unit: string
}>()
</script>

<template>
  <div class="quality-card">
    <div class="card-icon">
      <span :class="icon"></span>
    </div>
    <div class="card-content">
      <span class="card-value">{{ value }}</span>
      <span class="card-unit">{{ unit }}</span>
    </div>
    <span class="card-label">{{ label }}</span>
  </div>
</template>

<style scoped>
.quality-card {
  background: linear-gradient(135deg, rgba(45, 90, 39, 0.08) 0%, rgba(74, 138, 68, 0.05) 100%);
  border-radius: 12px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  transition: all 0.2s ease;
}

.quality-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(45, 90, 39, 0.12);
}

.card-icon {
  width: 48px;
  height: 48px;
  background: white;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: #2D5A27;
  margin-bottom: 12px;
  box-shadow: 0 4px 12px rgba(45, 90, 39, 0.1);
}

.card-content {
  display: flex;
  align-items: baseline;
  gap: 4px;
  margin-bottom: 4px;
}

.card-value {
  font-size: 28px;
  font-weight: 700;
  color: #333;
  font-family: 'Roboto Mono', monospace;
}

.card-unit {
  font-size: 12px;
  color: #999;
}

.card-label {
  font-size: 12px;
  color: #666;
}
</style>

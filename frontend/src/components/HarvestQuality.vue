<template>
  <div class="harvest-quality card">
    <h2 class="module-title">
      <span class="icon">🌾</span>
      采收质量信息
    </h2>

    <div class="quality-content">
      <div class="harvest-time">
        <div class="time-item">
          <span class="icon">📅</span>
          <div class="time-info">
            <span class="label">采收起始时间</span>
            <span class="value">{{ harvest.startDate }}</span>
          </div>
        </div>
        <div class="time-arrow">→</div>
        <div class="time-item">
          <span class="icon">📅</span>
          <div class="time-info">
            <span class="label">采收终止时间</span>
            <span class="value">{{ harvest.endDate }}</span>
          </div>
        </div>
      </div>

      <div class="quality-metrics">
        <div class="metric-item">
          <div class="metric-header">
            <span class="metric-label">糖度 (Brix)</span>
            <span class="metric-value">{{ harvest.sugarContent }}</span>
          </div>
          <div class="progress-bar">
            <div
              class="progress-fill sugar"
              :style="{ width: `${(harvest.sugarContent / 20) * 100}%` }"
            ></div>
          </div>
        </div>

        <div class="metric-item">
          <div class="metric-header">
            <span class="metric-label">单果重量 (克)</span>
            <span class="metric-value">{{ harvest.weight }}g</span>
          </div>
          <div class="progress-bar">
            <div
              class="progress-fill weight"
              :style="{ width: `${(harvest.weight / 30) * 100}%` }"
            ></div>
          </div>
        </div>
      </div>

      <div class="quality-details">
        <div class="detail-item">
          <span class="detail-icon">👅</span>
          <div class="detail-content">
            <span class="detail-label">口感描述</span>
            <p class="detail-text">{{ harvest.taste }}</p>
          </div>
        </div>

        <div class="detail-item">
          <span class="detail-icon">👥</span>
          <div class="detail-content">
            <span class="detail-label">适应人群</span>
            <p class="detail-text">{{ harvest.suitableCrowd }}</p>
          </div>
        </div>

        <div class="detail-item summary">
          <span class="detail-icon">⭐</span>
          <div class="detail-content">
            <span class="detail-label">品质小结</span>
            <p class="detail-text highlight">{{ harvest.qualitySummary }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
/**
 * 采收质量信息展示组件
 * 展示农产品的采收时间区间、质量指标（糖度、单果重量）及口感、适应人群等详情
 */

// 接收父组件传入的采收质量数据
defineProps({
  harvest: {
    type: Object,
    required: true,
    default: () => ({
      startDate: '',
      endDate: '',
      sugarContent: 0,
      weight: 0,
      taste: '',
      suitableCrowd: '',
      qualitySummary: ''
    })
  }
})
</script>

<style scoped>
/* 组件根容器 */
.harvest-quality {
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

.quality-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* 采收时间区间展示区域 */
.harvest-time {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: linear-gradient(135deg, #e8f5e9, #c8e6c9);
  padding: 16px;
  border-radius: 8px;
}

.time-item {
  display: flex;
  align-items: center;
  gap: 10px;
}

.time-item .icon {
  font-size: 24px;
}

.time-info {
  display: flex;
  flex-direction: column;
}

.time-info .label {
  font-size: 12px;
  color: var(--secondary-text);
}

.time-info .value {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-color);
}

.time-arrow {
  font-size: 20px;
  color: var(--primary-color);
}

/* 质量指标区域（进度条展示糖度和重量） */
.quality-metrics {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.metric-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.metric-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.metric-label {
  font-size: 14px;
  color: var(--text-color);
}

.metric-value {
  font-size: 16px;
  font-weight: bold;
  color: var(--primary-color);
}

/* 进度条基础样式 */
.progress-bar {
  height: 8px;
  background: var(--background-color);
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.5s ease;
}

.progress-fill.sugar {
  background: linear-gradient(90deg, var(--secondary-color), var(--primary-color));
}

.progress-fill.weight {
  background: linear-gradient(90deg, var(--accent-color), #ffb74d);
}

/* 品质详情卡片区域（口感、适应人群、品质小结） */
.quality-details {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.detail-item {
  display: flex;
  gap: 12px;
  padding: 12px;
  background: var(--background-color);
  border-radius: 8px;
}

/* 品质小结特殊高亮样式 */
.detail-item.summary {
  background: linear-gradient(135deg, #fff3e0, #ffe0b2);
  border-left: 3px solid var(--accent-color);
}

.detail-icon {
  font-size: 20px;
}

.detail-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.detail-label {
  font-size: 12px;
  color: var(--secondary-text);
}

.detail-text {
  font-size: 14px;
  color: var(--text-color);
  line-height: 1.5;
}

.detail-text.highlight {
  color: var(--accent-color);
  font-weight: 500;
}
</style>

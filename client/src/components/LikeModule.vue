<template>
  <div class="interaction-item like-module">
    <h3 class="interaction-title">点赞</h3>
    <div class="interaction-content">
      <button 
        class="btn-like" 
        :class="{ 'liked': likeStatus.user_liked, 'disabled': !likeStatus.can_like }"
        @click="handleLike"
        :disabled="!likeStatus.can_like"
      >
        <i class="icon-heart" :class="{ 'filled': likeStatus.user_liked }"></i>
        <span>{{ likeStatus.user_liked ? '已点赞' : '点赞' }}</span>
      </button>
      <div class="like-count">
        <i class="icon-heart"></i>
        <span>{{ likeCount }} 次点赞</span>
      </div>
      <div v-if="likeStatus.cooldown > 0" class="like-cooldown">
        冷却时间: {{ formatCooldown(likeStatus.cooldown) }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, defineProps, defineEmits } from 'vue'

const props = defineProps({
  likeStatus: {
    type: Object,
    required: true,
    default: () => ({
      can_like: true,
      cooldown: 0,
      user_liked: false,
      ip_address: ''
    })
  },
  likeCount: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits(['like'])

const handleLike = () => {
  if (!props.likeStatus.can_like) return
  emit('like')
}

const formatCooldown = (seconds) => {
  const minutes = Math.floor(seconds / 60)
  const remainingSeconds = seconds % 60
  return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`
}
</script>

<style scoped>
.like-module {
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

.btn-like {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #f5f5f5;
  border: 2px solid #ddd;
  color: #666;
  padding: 10px 20px;
  border-radius: 25px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  margin-bottom: 12px;
}

.btn-like:hover:not(.disabled) {
  background: #ffebee;
  border-color: #f44336;
  color: #f44336;
}

.btn-like.liked {
  background: #ffebee;
  border-color: #f44336;
  color: #f44336;
}

.btn-like.disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.like-count {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #666;
  font-size: 14px;
  margin-bottom: 8px;
}

.like-cooldown {
  font-size: 12px;
  color: #ff9800;
  font-weight: 500;
}

.icon-heart {
  display: inline-block;
  width: 16px;
  height: 16px;
  border: 2px solid currentColor;
  border-radius: 50%;
  position: relative;
}

.icon-heart.filled::before {
  content: '';
  position: absolute;
  top: 2px;
  left: 2px;
  right: 2px;
  bottom: 2px;
  background: currentColor;
  border-radius: 50%;
}

@media (max-width: 768px) {
  .like-module {
    margin-bottom: 16px;
  }
}

@media (max-width: 480px) {
  .btn-like {
    width: 100%;
    justify-content: center;
  }
}
</style>
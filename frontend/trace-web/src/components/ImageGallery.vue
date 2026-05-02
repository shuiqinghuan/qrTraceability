<script setup lang="ts">
defineProps<{
  images: string[]
}>()

const emit = defineEmits<{
  preview: [index: number]
}>()
</script>

<template>
  <div class="image-gallery">
    <div 
      v-for="(image, index) in images" 
      :key="index"
      class="gallery-item"
      @click="emit('preview', index)"
    >
      <img :src="image" :alt="`图片 ${index + 1}`" loading="lazy" />
    </div>
  </div>
</template>

<style scoped>
.image-gallery {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

@media (min-width: 768px) {
  .image-gallery {
    grid-template-columns: repeat(3, 1fr);
  }
}

.gallery-item {
  aspect-ratio: 4/3;
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
}

.gallery-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.gallery-item:hover img {
  transform: scale(1.05);
}
</style>

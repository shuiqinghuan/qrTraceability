<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useProductStore } from '@/stores/product'
import HomeView from './HomeView.vue'

const props = defineProps<{
  code: string
}>()

const route = useRoute()
const productStore = useProductStore()

const loadProduct = async () => {
  if (props.code) {
    await productStore.fetchByCode(props.code)
  }
}

onMounted(loadProduct)

watch(() => props.code, loadProduct)
</script>

<template>
  <HomeView />
</template>

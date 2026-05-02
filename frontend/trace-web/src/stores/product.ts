import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { productApi } from '@/api'
import type { Product } from '@/types'

export const useProductStore = defineStore('product', () => {
  const currentProduct = ref<Product | null>(null)
  const productList = ref<Product[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const hasProduct = computed(() => !!currentProduct.value)

  const fetchProduct = async (id: number) => {
    loading.value = true
    error.value = null
    try {
      const { data } = await productApi.getProduct(id)
      currentProduct.value = data
    } catch (e) {
      error.value = '获取产品信息失败'
      console.error(e)
    } finally {
      loading.value = false
    }
  }

  const fetchByCode = async (code: string) => {
    loading.value = true
    error.value = null
    try {
      const { data } = await productApi.getProductByCode(code)
      currentProduct.value = data
    } catch (e) {
      error.value = '产品不存在或加载失败'
      console.error(e)
    } finally {
      loading.value = false
    }
  }

  const fetchProducts = async () => {
    loading.value = true
    try {
      const { data } = await productApi.getProducts()
      productList.value = data
    } catch (e) {
      console.error(e)
    } finally {
      loading.value = false
    }
  }

  const clearProduct = () => {
    currentProduct.value = null
    error.value = null
  }

  return {
    currentProduct,
    productList,
    loading,
    error,
    hasProduct,
    fetchProduct,
    fetchByCode,
    fetchProducts,
    clearProduct
  }
})

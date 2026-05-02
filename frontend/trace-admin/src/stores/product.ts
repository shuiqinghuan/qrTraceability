import { defineStore } from 'pinia'
import { ref } from 'vue'
import { productApi } from '@/api'
import type { Product } from '@/types'

export const useProductStore = defineStore('product', () => {
  const products = ref<Product[]>([])
  const currentProduct = ref<Product | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const fetchProducts = async () => {
    loading.value = true
    error.value = null
    try {
      const { data } = await productApi.getProducts()
      products.value = data.results || data
    } catch (e) {
      error.value = '获取产品列表失败'
      console.error(e)
    } finally {
      loading.value = false
    }
  }

  const fetchProduct = async (id: number) => {
    loading.value = true
    error.value = null
    try {
      const { data } = await productApi.getProduct(id)
      currentProduct.value = data
    } catch (e) {
      error.value = '获取产品详情失败'
      console.error(e)
    } finally {
      loading.value = false
    }
  }

  const createProduct = async (productData: Partial<Product>) => {
    loading.value = true
    error.value = null
    try {
      const { data } = await productApi.createProduct(productData)
      products.value.unshift(data)
      return data
    } catch (e) {
      error.value = '创建产品失败'
      console.error(e)
      return null
    } finally {
      loading.value = false
    }
  }

  const updateProduct = async (id: number, productData: Partial<Product>) => {
    loading.value = true
    error.value = null
    try {
      const { data } = await productApi.updateProduct(id, productData)
      const index = products.value.findIndex(p => p.id === id)
      if (index !== -1) {
        products.value[index] = data
      }
      currentProduct.value = data
      return data
    } catch (e) {
      error.value = '更新产品失败'
      console.error(e)
      return null
    } finally {
      loading.value = false
    }
  }

  const deleteProduct = async (id: number) => {
    loading.value = true
    error.value = null
    try {
      await productApi.deleteProduct(id)
      products.value = products.value.filter(p => p.id !== id)
      return true
    } catch (e) {
      error.value = '删除产品失败'
      console.error(e)
      return false
    } finally {
      loading.value = false
    }
  }

  const clearCurrent = () => {
    currentProduct.value = null
    error.value = null
  }

  return {
    products,
    currentProduct,
    loading,
    error,
    fetchProducts,
    fetchProduct,
    createProduct,
    updateProduct,
    deleteProduct,
    clearCurrent
  }
})

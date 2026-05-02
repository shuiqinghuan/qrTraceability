import axios from 'axios'
import type { Product } from '@/types'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

api.interceptors.response.use(
  response => response,
  error => {
    console.error('API Error:', error)
    return Promise.reject(error)
  }
)

export const productApi = {
  getProducts: () => api.get<Product[]>('/products/'),
  
  getProduct: (id: number) => api.get<Product>(`/products/${id}/`),
  
  getProductByCode: (code: string) => api.get<Product>(`/products/public/code/${code}/`)
}

export default api

import axios from 'axios'
import type { Product, LoginDTO, TokenResponse } from '@/types'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

api.interceptors.request.use(config => {
  const token = localStorage.getItem('access_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  response => response,
  async error => {
    const originalRequest = error.config
    
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true
      
      const refreshToken = localStorage.getItem('refresh_token')
      if (refreshToken) {
        try {
          const { data } = await axios.post(`${api.defaults.baseURL}/auth/refresh/`, {
            refresh: refreshToken
          })
          
          localStorage.setItem('access_token', data.access)
          originalRequest.headers.Authorization = `Bearer ${data.access}`
          return api(originalRequest)
        } catch {
          localStorage.removeItem('access_token')
          localStorage.removeItem('refresh_token')
          window.location.href = '/admin/login'
        }
      }
    }
    
    return Promise.reject(error)
  }
)

export const authApi = {
  login: (data: LoginDTO) => api.post<TokenResponse>('/auth/login/', data),
  logout: (refreshToken: string) => api.post('/auth/logout/', { refresh: refreshToken })
}

export const productApi = {
  getProducts: () => api.get<{ results: Product[] }>('/products/'),
  
  getProduct: (id: number) => api.get<Product>(`/products/${id}/`),
  
  createProduct: (data: Partial<Product>) => api.post<Product>('/products/', data),
  
  updateProduct: (id: number, data: Partial<Product>) => api.put<Product>(`/products/${id}/`, data),
  
  deleteProduct: (id: number) => api.delete(`/products/${id}/`)
}

export const mediaApi = {
  uploadMedia: (productId: number, file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('product_id', productId.toString())
    return api.post('/media/', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },
  
  deleteMedia: (id: number) => api.delete(`/media/${id}/`)
}

export default api

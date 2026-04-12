import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    if (error.response && error.response.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// 认证相关
export const authAPI = {
  login: (data) => api.post('/auth/login', data),
  register: (data) => api.post('/auth/register', data)
}

// 种子信息相关
export const seedAPI = {
  list: () => api.get('/seed'),
  get: (id) => api.get(`/seed/${id}`),
  create: (data) => api.post('/seed', data),
  update: (id, data) => api.put(`/seed/${id}`, data),
  delete: (id) => api.delete(`/seed/${id}`)
}

// 种植管理相关
export const plantingAPI = {
  list: () => api.get('/planting'),
  get: (id) => api.get(`/planting/${id}`),
  create: (data) => api.post('/planting', data),
  update: (id, data) => api.put(`/planting/${id}`, data),
  delete: (id) => api.delete(`/planting/${id}`),
  generateQRCode: (id) => api.get(`/planting/${id}/qrcode`)
}

// 生长媒体相关
export const growthAPI = {
  list: () => api.get('/growth'),
  get: (id) => api.get(`/growth/${id}`),
  create: (data) => api.post('/growth', data),
  update: (id, data) => api.put(`/growth/${id}`, data),
  delete: (id) => api.delete(`/growth/${id}`)
}

// 产品品质相关
export const qualityAPI = {
  list: () => api.get('/quality'),
  get: (id) => api.get(`/quality/${id}`),
  create: (data) => api.post('/quality', data),
  update: (id, data) => api.put(`/quality/${id}`, data),
  delete: (id) => api.delete(`/quality/${id}`)
}

// 产品相关
export const productAPI = {
  list: () => api.get('/products'),
  get: (id) => api.get(`/products/${id}`),
  favorite: (id) => api.post(`/products/${id}/favorite`),
  like: (id) => api.post(`/products/${id}/like`)
}

export default api
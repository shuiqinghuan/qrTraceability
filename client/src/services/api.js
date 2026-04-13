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

// 产品相关（新API）
export const productAPI = {
  // 产品列表
  list: (params) => api.get('/products', { params }),
  
  // 产品详情
  getDetail: (id) => api.get(`/products/${id}`),
  
  // 产品媒体
  getMedia: (id) => api.get(`/products/${id}/media`),
  
  // 产品品质
  getQuality: (id) => api.get(`/products/${id}/quality`),
  
  // 产品统计
  getStats: (id) => api.get(`/products/${id}/stats`),
  
  // 相关产品（同地点）
  getRelatedProducts: (id) => api.get(`/products/${id}/related`),
  
  // 相似产品
  getSimilarProducts: (id) => api.get(`/products/${id}/similar`),
  
  // 点赞统计
  getLikeCount: (id) => api.get(`/products/${id}/likes`),
  
  // 点赞状态
  getLikeStatus: (id) => api.get(`/products/${id}/like-status`),
  
  // 点赞产品（带IP限制）
  likeProduct: (id) => api.post(`/products/${id}/like`),
  
  // 取消点赞（需要认证）
  removeLike: (id) => api.post(`/user/likes/${id}/remove`),
  
  // 收藏功能（保持兼容）
  favorite: (id) => api.post(`/products/${id}/favorite`),
  like: (id) => api.post(`/products/${id}/like`) // 旧版兼容
}

// 地点相关
export const locationAPI = {
  // 获取指定地点产品
  getProductsByLocation: (location, params) => api.get(`/locations/${location}/products`, { params }),
  
  // 获取地点统计
  getLocationStats: (location) => api.get(`/locations/${location}/stats`),
  
  // 搜索地点
  searchLocations: (query, params) => api.get('/locations/search', { params: { q: query, ...params } })
}

// 标签相关
export const tagAPI = {
  // 搜索标签
  searchTags: (query, params) => api.get('/tags/search', { params: { q: query, ...params } }),
  
  // 获取种植标签
  getPlantingTags: (plantingId) => api.get(`/planting/${plantingId}/tags`),
  
  // 添加种植标签（需要认证）
  addPlantingTag: (plantingId, data) => api.post(`/planting/${plantingId}/tags`, data),
  
  // 批量添加标签（需要认证）
  batchAddPlantingTags: (plantingId, data) => api.post(`/planting/${plantingId}/tags/batch`, data),
  
  // 更新标签（需要认证）
  updatePlantingTag: (plantingId, tagId, data) => api.put(`/planting/${plantingId}/tags/${tagId}`, data),
  
  // 删除标签（需要认证）
  deletePlantingTag: (plantingId, tagId) => api.delete(`/planting/${plantingId}/tags/${tagId}`)
}

// 媒体管理相关
export const mediaAPI = {
  // 获取产品媒体列表
  list: (productId) => api.get(`/products/${productId}/media`),
  
  // 上传媒体
  upload: (formData) => api.post('/admin/media/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  }),
  
  // 删除媒体
  delete: (mediaId) => api.delete(`/growth/${mediaId}`)
}

// 管理员相关
export const adminAPI = {
  // 管理员登录
  login: (data) => api.post('/admin/login', data),
  
  // 产品管理
  listProducts: (params) => api.get('/admin/products', { params }),
  createProduct: (data) => api.post('/admin/products', data),
  getProductDetail: (id) => api.get(`/admin/products/${id}`),
  updateProduct: (id, data) => api.put(`/admin/products/${id}`, data),
  deleteProduct: (id) => api.delete(`/admin/products/${id}`),
  
  // 媒体上传
  uploadMedia: (data) => api.post('/admin/media/upload', data),
  
  // 品质管理
  updateQuality: (data) => api.post('/admin/quality', data),
  
  // 系统管理
  resetIPLikeRestriction: (id, ip) => api.post(`/admin/likes/${id}/reset-ip`, null, { params: { ip } })
}

export default api
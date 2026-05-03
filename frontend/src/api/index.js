/**
 * Axios实例配置文件
 * 创建统一的HTTP请求实例，配置基础地址、超时时间和拦截器
 */
import axios from 'axios'

// 创建Axios实例，使用环境变量中的API基础地址，未设置时使用默认本地地址
const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8000/api',
  timeout: 10000, // 请求超时时间：10秒
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器：可在发送请求前统一处理（如添加token等）
api.interceptors.request.use(
  (config) => {
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器：统一解包响应数据，直接返回data字段
api.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    return Promise.reject(error)
  }
)

export default api

/**
 * 产品相关API接口模块
 * 封装所有与产品、媒体资源、采收质量相关的HTTP请求
 */
import api from './index'

/** 根据产品溯源码获取产品详情 */
export function getProductByCode(code) {
  return api.get(`/products/${code}/`)
}

/** 获取产品列表（支持分页） */
export function getProductList(params) {
  return api.get('/products/', { params })
}

/** 创建新产品 */
export function createProduct(data) {
  return api.post('/products/', data)
}

/** 更新已有产品信息 */
export function updateProduct(id, data) {
  return api.put(`/products/${id}/`, data)
}

/** 删除产品 */
export function deleteProduct(id) {
  return api.delete(`/products/${id}/`)
}

/** 获取产品的媒体资源列表 */
export function getMediaList(productId, type) {
  return api.get(`/products/${productId}/media/`, { params: { type } })
}

/** 为产品添加媒体资源 */
export function addMedia(productId, data) {
  return api.post(`/products/${productId}/media/`, data)
}

/** 删除媒体资源 */
export function deleteMedia(id) {
  return api.delete(`/media/${id}/`)
}

/** 获取产品的采收质量信息 */
export function getHarvestQuality(productId) {
  return api.get(`/products/${productId}/harvest/`)
}

/** 保存或更新产品的采收质量信息 */
export function saveHarvestQuality(productId, data) {
  return api.post(`/products/${productId}/harvest/`, data)
}

import api from './index'

export function getProductByCode(code) {
  return api.get(`/products/${code}/`)
}

export function getProductList(params) {
  return api.get('/products/', { params })
}

export function createProduct(data) {
  return api.post('/products/', data)
}

export function updateProduct(id, data) {
  return api.put(`/products/${id}/`, data)
}

export function deleteProduct(id) {
  return api.delete(`/products/${id}/`)
}

export function getMediaList(productId, type) {
  return api.get(`/products/${productId}/media/`, { params: { type } })
}

export function addMedia(productId, data) {
  return api.post(`/products/${productId}/media/`, data)
}

export function deleteMedia(id) {
  return api.delete(`/media/${id}/`)
}

export function getHarvestQuality(productId) {
  return api.get(`/products/${productId}/harvest/`)
}

export function saveHarvestQuality(productId, data) {
  return api.post(`/products/${productId}/harvest/`, data)
}

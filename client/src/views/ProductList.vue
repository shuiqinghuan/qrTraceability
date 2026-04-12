<template>
  <div class="product-list">
    <h2>产品浏览</h2>
    
    <div class="product-grid">
      <div v-for="product in products" :key="product.id" class="product-card">
        <div class="product-image">
          <img :src="product.image_url || 'https://via.placeholder.com/200'" :alt="product.name">
        </div>
        <div class="product-info">
          <h3>{{ product.name }}</h3>
          <p class="product-description">{{ product.description }}</p>
          <div class="product-meta">
            <span class="product-sugar">糖度: {{ product.sugar_content }}</span>
            <span class="product-weight">重量: {{ product.weight }}g</span>
          </div>
          <div class="product-actions">
            <button class="favorite-button" @click="toggleFavorite(product.id)" :class="{ 'favorited': favoritedProducts.includes(product.id) }">
              {{ favoritedProducts.includes(product.id) ? '已收藏' : '收藏' }}
            </button>
            <button class="like-button" @click="toggleLike(product.id)" :class="{ 'liked': likedProducts.includes(product.id) }">
              {{ likedProducts.includes(product.id) ? '已点赞' : '点赞' }}
            </button>
            <button class="share-button" @click="shareProduct(product.id)">转发</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { productAPI } from '../services/api'

const products = ref([])
const favoritedProducts = ref([])
const likedProducts = ref([])

onMounted(() => {
  fetchProducts()
  loadUserInteractions()
})

const fetchProducts = async () => {
  try {
    const data = await productAPI.list()
    products.value = data
  } catch (error) {
    console.error('获取产品列表失败', error)
  }
}

const loadUserInteractions = () => {
  // 从本地存储加载用户交互数据
  const savedFavorites = localStorage.getItem('favorites')
  const savedLikes = localStorage.getItem('likes')
  if (savedFavorites) {
    favoritedProducts.value = JSON.parse(savedFavorites)
  }
  if (savedLikes) {
    likedProducts.value = JSON.parse(savedLikes)
  }
}

const toggleFavorite = async (productId) => {
  try {
    if (favoritedProducts.value.includes(productId)) {
      // 取消收藏
      favoritedProducts.value = favoritedProducts.value.filter(id => id !== productId)
    } else {
      // 添加收藏
      favoritedProducts.value.push(productId)
      await productAPI.favorite(productId)
    }
    // 保存到本地存储
    localStorage.setItem('favorites', JSON.stringify(favoritedProducts.value))
  } catch (error) {
    console.error('操作收藏失败', error)
  }
}

const toggleLike = async (productId) => {
  try {
    if (likedProducts.value.includes(productId)) {
      // 取消点赞
      likedProducts.value = likedProducts.value.filter(id => id !== productId)
    } else {
      // 添加点赞
      likedProducts.value.push(productId)
      await productAPI.like(productId)
    }
    // 保存到本地存储
    localStorage.setItem('likes', JSON.stringify(likedProducts.value))
  } catch (error) {
    console.error('操作点赞失败', error)
  }
}

const shareProduct = (productId) => {
  // 转发功能实现
  console.log('转发产品', productId)
  alert('转发功能已触发')
}
</script>

<style scoped>
.product-list {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

h2 {
  margin-bottom: 20px;
  color: #333;
}

.product-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.product-card {
  border: 1px solid #ddd;
  border-radius: 8px;
  overflow: hidden;
  transition: transform 0.2s, box-shadow 0.2s;
}

.product-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
}

.product-image {
  height: 200px;
  overflow: hidden;
}

.product-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.product-info {
  padding: 15px;
}

.product-info h3 {
  margin-bottom: 10px;
  color: #333;
}

.product-description {
  margin-bottom: 15px;
  color: #666;
  font-size: 14px;
  line-height: 1.4;
}

.product-meta {
  display: flex;
  justify-content: space-between;
  margin-bottom: 15px;
  font-size: 14px;
  color: #888;
}

.product-actions {
  display: flex;
  gap: 10px;
}

.favorite-button, .like-button, .share-button {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background-color: white;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.favorite-button:hover {
  background-color: #f5f5f5;
  border-color: #ff4081;
  color: #ff4081;
}

.favorite-button.favorited {
  background-color: #ff4081;
  color: white;
  border-color: #ff4081;
}

.like-button:hover {
  background-color: #f5f5f5;
  border-color: #ff9800;
  color: #ff9800;
}

.like-button.liked {
  background-color: #ff9800;
  color: white;
  border-color: #ff9800;
}

.share-button:hover {
  background-color: #f5f5f5;
  border-color: #2196F3;
  color: #2196F3;
}
</style>
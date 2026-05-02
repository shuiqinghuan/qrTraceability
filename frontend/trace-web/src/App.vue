<script setup lang="ts">
import { useProductStore } from '@/stores/product'
import { storeToRefs } from 'pinia'

const productStore = useProductStore()
const { hasProduct } = storeToRefs(productStore)
</script>

<template>
  <div class="app-wrapper">
    <header class="app-header" v-if="hasProduct">
      <div class="container">
        <h1 class="logo">🌿 农产品溯源</h1>
      </div>
    </header>
    
    <main class="main-content">
      <router-view v-slot="{ Component }">
        <transition name="page" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>

    <nav class="bottom-nav" v-if="hasProduct">
      <router-link to="/" class="nav-item" active-class="active">
        <span class="i-ph-house text-xl"></span>
        <span>首页</span>
      </router-link>
      <router-link to="/media" class="nav-item" active-class="active">
        <span class="i-ph-images text-xl"></span>
        <span>媒体</span>
      </router-link>
      <router-link to="/quality" class="nav-item" active-class="active">
        <span class="i-ph-chart-bar text-xl"></span>
        <span>质量</span>
      </router-link>
    </nav>
  </div>
</template>

<style scoped>
.app-wrapper {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.app-header {
  background: linear-gradient(135deg, #2D5A27 0%, #4A8A44 100%);
  color: white;
  padding: 16px 0;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 2px 12px rgba(45, 90, 39, 0.15);
}

.logo {
  font-size: 20px;
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 8px;
}

.main-content {
  flex: 1;
  padding-bottom: 80px;
}

.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: white;
  display: flex;
  justify-content: space-around;
  padding: 8px 0 16px;
  box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.08);
  z-index: 100;
}

.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 8px 16px;
  color: #666;
  text-decoration: none;
  font-size: 12px;
  transition: color 0.2s ease;
}

.nav-item.active {
  color: #2D5A27;
}

.nav-item:active {
  transform: scale(0.95);
}
</style>

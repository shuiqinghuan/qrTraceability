<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const activeTab = ref('products')

const navItems = [
  { path: '/admin/products', label: '产品管理', icon: 'goods', key: 'products' },
  { path: '/admin/batches', label: '批次管理', icon: 'list', key: 'batches' }
]

const handleNavClick = (path, key) => {
  activeTab.value = key
  router.push(path)
}
</script>

<template>
  <div class="admin-container">
    <!-- 侧边栏 -->
    <aside class="sidebar">
      <div class="sidebar-header">
        <h1 class="logo">农产品溯源系统</h1>
      </div>
      <nav class="sidebar-nav">
        <ul>
          <li 
            v-for="item in navItems" 
            :key="item.key"
            :class="{ active: activeTab === item.key }"
            @click="handleNavClick(item.path, item.key)"
          >
            <van-icon :name="item.icon" />
            <span>{{ item.label }}</span>
          </li>
        </ul>
      </nav>
    </aside>
    
    <!-- 主内容区 -->
    <main class="main-content">
      <div class="content-header">
        <h2 class="page-title">
          {{ activeTab === 'products' ? '产品管理' : '批次管理' }}
        </h2>
      </div>
      <div class="content-body">
        <router-view />
      </div>
    </main>
  </div>
</template>

<style scoped>
.admin-container {
  display: flex;
  min-height: 100vh;
  background-color: #f5f5f5;
}

/* 侧边栏样式 */
.sidebar {
  width: 240px;
  background-color: #1989fa;
  color: white;
  display: flex;
  flex-direction: column;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.1);
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.logo {
  font-size: 18px;
  font-weight: bold;
  margin: 0;
}

.sidebar-nav {
  flex: 1;
  padding: 20px 0;
}

.sidebar-nav ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.sidebar-nav li {
  display: flex;
  align-items: center;
  padding: 12px 20px;
  cursor: pointer;
  transition: all 0.3s ease;
  border-left: 3px solid transparent;
}

.sidebar-nav li:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

.sidebar-nav li.active {
  background-color: rgba(255, 255, 255, 0.2);
  border-left-color: white;
}

.sidebar-nav li van-icon {
  margin-right: 10px;
  font-size: 18px;
}

/* 主内容区样式 */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.content-header {
  background-color: white;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  margin: 0;
  color: #333;
}

.content-body {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .admin-container {
    flex-direction: column;
  }
  
  .sidebar {
    width: 100%;
    height: auto;
    flex-direction: row;
  }
  
  .sidebar-header {
    border-bottom: none;
    border-right: 1px solid rgba(255, 255, 255, 0.1);
  }
  
  .sidebar-nav {
    flex: 1;
    padding: 0;
  }
  
  .sidebar-nav ul {
    display: flex;
  }
  
  .sidebar-nav li {
    flex: 1;
    justify-content: center;
    border-left: none;
    border-bottom: 3px solid transparent;
  }
  
  .sidebar-nav li.active {
    border-left-color: transparent;
    border-bottom-color: white;
  }
  
  .content-body {
    padding: 10px;
  }
}
</style>

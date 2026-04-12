<template>
  <div class="dashboard">
    <header class="dashboard-header">
      <h1>农产品二维码追溯系统</h1>
      <div class="user-info">
        <span>{{ user?.phone_number }}</span>
        <button class="logout-button" @click="handleLogout">退出登录</button>
      </div>
    </header>
    <div class="dashboard-body">
      <aside class="sidebar">
        <nav>
          <ul>
            <li v-if="user?.role === 'serverseed'">
              <router-link to="/dashboard/seed">种子管理</router-link>
            </li>
            <li v-if="user?.role === 'servergrow'">
              <router-link to="/dashboard/growth">生长管理</router-link>
            </li>
            <li v-if="user?.role === 'servermanager'">
              <router-link to="/dashboard/quality">品质管理</router-link>
            </li>
            <li v-if="user?.role === 'clentcustomer'">
              <router-link to="/dashboard/products">产品浏览</router-link>
            </li>
          </ul>
        </nav>
      </aside>
      <main class="main-content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const user = ref(null)

onMounted(() => {
  const userStr = localStorage.getItem('user')
  if (userStr) {
    user.value = JSON.parse(userStr)
  }
})

const handleLogout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  router.push('/login')
}
</script>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.dashboard-header {
  background-color: #4CAF50;
  color: white;
  padding: 15px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.dashboard-header h1 {
  font-size: 20px;
  margin: 0;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 15px;
}

.logout-button {
  background-color: rgba(255, 255, 255, 0.2);
  border: none;
  color: white;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
}

.logout-button:hover {
  background-color: rgba(255, 255, 255, 0.3);
}

.dashboard-body {
  display: flex;
  flex: 1;
}

.sidebar {
  width: 200px;
  background-color: #f0f0f0;
  padding: 20px 0;
}

.sidebar nav ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.sidebar nav ul li {
  margin-bottom: 10px;
}

.sidebar nav ul li a {
  display: block;
  padding: 10px 20px;
  color: #333;
  text-decoration: none;
  transition: background-color 0.2s;
}

.sidebar nav ul li a:hover {
  background-color: rgba(76, 175, 80, 0.1);
}

.sidebar nav ul li a.router-link-active {
  background-color: #4CAF50;
  color: white;
}

.main-content {
  flex: 1;
  padding: 20px;
  background-color: #f5f5f5;
}
</style>
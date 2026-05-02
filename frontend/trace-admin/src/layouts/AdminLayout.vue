<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const activeMenu = computed(() => router.currentRoute.value.name?.toString())

const handleCommand = async (command: string) => {
  if (command === 'logout') {
    await authStore.logout()
  }
}
</script>

<template>
  <div class="admin-layout">
    <el-container>
      <el-header class="header">
        <div class="header-left">
          <h2 class="logo">农产品溯源系统</h2>
        </div>
        <div class="header-right">
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-icon><User /></el-icon>
              <span>管理员</span>
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      
      <el-container>
        <el-aside width="200px" class="aside">
          <el-menu
            :default-active="activeMenu"
            router
            class="aside-menu"
          >
            <el-menu-item index="/admin">
              <el-icon><Goods /></el-icon>
              <span>产品管理</span>
            </el-menu-item>
          </el-menu>
        </el-aside>
        
        <el-main class="main">
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script lang="ts">
import { User, ArrowDown, Goods } from '@element-plus/icons-vue'
export default {
  components: { User, ArrowDown, Goods }
}
</script>

<style lang="scss" scoped>
.admin-layout {
  height: 100vh;
}

.el-container {
  height: 100%;
}

.header {
  background: white;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.logo {
  font-size: 18px;
  font-weight: 600;
  color: #2D5A27;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  color: #666;
}

.aside {
  background: white;
  border-right: 1px solid #E8E8E8;
}

.aside-menu {
  border-right: none;
  
  :deep(.el-menu-item) {
    &.is-active {
      background: rgba(45, 90, 39, 0.1);
      color: #2D5A27;
    }
    
    &:hover {
      background: rgba(45, 90, 39, 0.05);
    }
  }
}

.main {
  background: #F8F9FA;
  padding: 24px;
  overflow-y: auto;
}
</style>

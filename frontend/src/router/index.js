/**
 * 路由配置文件
 * 定义应用的所有页面路由，使用懒加载方式引入视图组件
 */
import { createRouter, createWebHistory } from 'vue-router'

// 路由表定义
const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/Home.vue') // 首页：产品列表
  },
  {
    path: '/about',
    name: 'About',
    component: () => import('../views/About.vue') // 关于页面
  },
  {
    path: '/product/:code?', // code参数可选，支持通过二维码直接访问
    name: 'ProductTrace',
    component: () => import('../views/ProductTrace.vue') // 产品溯源详情页
  }
]

// 创建路由实例，使用HTML5 History模式
const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router

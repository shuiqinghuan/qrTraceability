import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue')
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/Register.vue')
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('../views/Dashboard.vue'),
    children: [
      {
        path: 'seed',
        name: 'SeedManagement',
        component: () => import('../views/SeedManagement.vue')
      },
      {
        path: 'growth',
        name: 'GrowthManagement',
        component: () => import('../views/GrowthManagement.vue')
      },
      {
        path: 'quality',
        name: 'QualityManagement',
        component: () => import('../views/QualityManagement.vue')
      },
      {
        path: 'products',
        name: 'ProductList',
        component: () => import('../views/ProductList.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.path !== '/login' && to.path !== '/register' && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router
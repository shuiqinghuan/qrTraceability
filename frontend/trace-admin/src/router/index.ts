import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  {
    path: '/admin/login',
    name: 'Login',
    component: () => import('@/views/LoginView.vue'),
    meta: { guest: true }
  },
  {
    path: '/admin',
    name: 'Dashboard',
    component: () => import('@/views/DashboardView.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'ProductList',
        component: () => import('@/views/ProductListView.vue')
      },
      {
        path: 'product/new',
        name: 'ProductCreate',
        component: () => import('@/views/ProductEditView.vue')
      },
      {
        path: 'product/:id/edit',
        name: 'ProductEdit',
        component: () => import('@/views/ProductEditView.vue'),
        props: true
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/admin/login')
  } else if (to.meta.guest && authStore.isAuthenticated) {
    next('/admin')
  } else {
    next()
  }
})

export default router

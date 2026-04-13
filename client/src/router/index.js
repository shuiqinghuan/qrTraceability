import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/dashboard/products'
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
  },
  // 产品详情页面（公共页面，无需登录）
  {
    path: '/products/:id',
    name: 'ProductDetail',
    component: () => import('../components/ProductDetail.vue')
  },
  // 地点产品列表页面
  {
    path: '/locations/:location/products',
    name: 'LocationProducts',
    component: () => import('../components/LocationProducts.vue')
  },
  // 后台管理路由
  {
    path: '/admin/login',
    name: 'AdminLogin',
    component: () => import('../views/AdminLogin.vue')
  },
  {
    path: '/admin',
    name: 'AdminDashboard',
    component: () => import('../views/AdminDashboard.vue'),
    children: [
      {
        path: 'products',
        name: 'AdminProductManagement',
        component: () => import('../views/AdminProductManagement.vue')
      },
      {
        path: 'media',
        name: 'AdminMediaManagement',
        component: () => import('../views/AdminMediaManagement.vue')
      },
      {
        path: 'quality',
        name: 'AdminQualityManagement',
        component: () => import('../views/AdminQualityManagement.vue')
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
  const adminToken = localStorage.getItem('admin_token')
  
  // 公开路由（无需登录）
  const publicRoutes = [
    '/login',
    '/register',
    '/products/:id',
    '/locations/:location/products',
    '/admin/login',
    '/dashboard/products'
  ]
  
  // 检查是否是公开路由
  const isPublicRoute = publicRoutes.some(route => {
    if (route.includes(':')) {
      // 处理动态路由
      const routePattern = route.replace(/:[^/]+/g, '[^/]+')
      const regex = new RegExp(`^${routePattern}$`)
      return regex.test(to.path)
    }
    return to.path === route
  })
  
  // 检查是否是后台管理路由
  const isAdminRoute = to.path.startsWith('/admin/') && to.path !== '/admin/login'
  
  if (isAdminRoute) {
    // 后台管理路由需要admin_token
    if (!adminToken) {
      next('/admin/login')
    } else {
      next()
    }
  } else if (!isPublicRoute && !token) {
    // 普通用户路由需要token
    next('/login')
  } else {
    next()
  }
})

export default router
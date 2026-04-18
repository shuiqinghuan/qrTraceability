import { createRouter, createWebHistory } from 'vue-router'
import ProductDetail from './views/ProductDetail.vue'
import Admin from './views/Admin.vue'
import AdminProducts from './views/AdminProducts.vue'
import AdminBatches from './views/AdminBatches.vue'
import AdminAddProduct from './views/AdminAddProduct.vue'
import AdminAddBatch from './views/AdminAddBatch.vue'

const routes = [
  {
    path: '/',
    redirect: '/admin'
  },
  {
    path: '/product/:id',
    name: 'ProductDetail',
    component: ProductDetail
  },
  {
    path: '/admin',
    name: 'Admin',
    component: Admin,
    children: [
      {
        path: '',
        redirect: 'products'
      },
      {
        path: 'products',
        name: 'AdminProducts',
        component: AdminProducts
      },
      {
        path: 'batches',
        name: 'AdminBatches',
        component: AdminBatches
      },
      {
        path: 'products/add',
        name: 'AdminAddProduct',
        component: AdminAddProduct
      },
      {
        path: 'batches/add',
        name: 'AdminAddBatch',
        component: AdminAddBatch
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router

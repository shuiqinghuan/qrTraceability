import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/HomeView.vue')
  },
  {
    path: '/media',
    name: 'Media',
    component: () => import('@/views/MediaView.vue')
  },
  {
    path: '/quality',
    name: 'Quality',
    component: () => import('@/views/QualityView.vue')
  },
  {
    path: '/p/:code',
    name: 'Trace',
    component: () => import('@/views/TraceView.vue'),
    props: true
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router

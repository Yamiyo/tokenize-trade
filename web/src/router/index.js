import { createRouter, createWebHistory } from 'vue-router'
// import HomeView from '@src/views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/index',
      component: () => import('@src/components/TradeView.vue')
    },
  ]
})

export default router

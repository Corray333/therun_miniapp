import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/farm',
      name: 'farm',
      component: () => import('../views/FarmView.vue')
    },
    {
      path: '/chibi',
      name: 'chibi',
      component: () => import('../views/ChibiView.vue')
    },
    {
      path: '/battles',
      name: 'battles',
      component: () => import('../views/BattlesView.vue')
    },
    {
      path: '/friens',
      name: 'friens',
      component: () => import('../views/FriensView.vue')
    },
    {
      path: '/tasks',
      name: 'tasks',
      component: () => import('../views/TasksView.vue')
    },
    {
      path: '/more',
      name: 'more',
      component: () => import('../views/MoreView.vue')
    }
  ]
})

export default router

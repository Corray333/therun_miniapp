import { createRouter, createWebHistory } from 'vue-router'
declare const Telegram: any

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'farm',
      component: () => import('../views/FarmView.vue'),
    },
    {
      path: '/chibi',
      name: 'chibi',
      component: () => import('../views/ChibiView.vue')
    },
    {
      path: '/chibi/cases',
      name: 'cases',
      component: () => import('../views/chibi/CasesView.vue')
    },
    {
      path: '/battles',
      name: 'battles',
      component: () => import('../views/BattlesView.vue')
    },
    {
      path: '/friens/activated',
      name: 'activated',
      component: () => import('../views/friens/ActivatedView.vue')
    },
    {
      path: '/friens/not-activated',
      name: 'not-activated',
      component: () => import('../views/friens/NotActivatedView.vue')
    },
    {
      path: '/friens',
      name: 'friens',
      component: () => import('../views/FriensView.vue')
    },
    {
      path: '/tasks',
      name: 'tasks',
      component: () => import('../views/TasksView.vue'),
      children: [
        {
          path: 'tasks',
          component: () => import('../views/tasks/TasksView.vue')
        },
        {
          path: 'quests',
          component: () => import('../views/tasks/QuestsView.vue')
        },
        {
          path: 'special',
          component: () => import('../views/tasks/SpecialView.vue')
        }
      ]
    },
    {
      path: '/settings',
      name: 'settings',
      component: () => import('../views/SettingsView.vue')
    },
    {
      path: '/onboarding',
      name: 'onboarding',
      component: () => import('../views/OnboardingView.vue')
    }
  ]
})

router.beforeEach(async (to, from, next) => {
  const tg = Telegram.WebApp
	var BackButton = Telegram.WebApp.BackButton
  if (to.name === 'farm') {
    BackButton.hide()
  } else{
    BackButton.show() 
  }

  next()
})

export default router

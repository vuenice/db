import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes: RouteRecordRaw[] = [
  { name: 'login', path: '/login', component: () => import('../views/LoginView.vue') },
  { name: 'register', path: '/register', component: () => import('../views/RegisterView.vue') },
  {
    name: 'workbench',
    path: '/workbench',
    component: () => import('../views/workbench/index.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/workbench/export',
    component: () => import('../views/workbench/export/index.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/workbench/import',
    component: () => import('../views/workbench/import/index.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/workbench/query',
    component: () => import('../views/workbench/query/index.vue'),
    meta: { requiresAuth: true },
  },
  { path: '/', redirect: '/workbench' },
  { path: '/:pathMatch(.*)*', redirect: '/workbench' },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (to.meta?.requiresAuth && !auth.token) {
    if (auth.hasUsers === false) {
      return { name: 'register', query: { redirect: to.fullPath } }
    }
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if ((to.name === 'login' || to.name === 'register') && auth.token) {
    return { name: 'workbench' }
  }
  if (to.name === 'login' && !auth.token && auth.hasUsers === false) {
    return { name: 'register', query: to.query }
  }
})

export default router

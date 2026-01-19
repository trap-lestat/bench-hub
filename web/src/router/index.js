import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

import LoginView from '../views/LoginView.vue'
import DashboardView from '../views/DashboardView.vue'
import UsersView from '../views/UsersView.vue'
import ScriptsView from '../views/ScriptsView.vue'
import TasksView from '../views/TasksView.vue'
import ReportsView from '../views/ReportsView.vue'
import SettingsView from '../views/SettingsView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/dashboard' },
    { path: '/login', component: LoginView, meta: { requiresAuth: false } },
    { path: '/dashboard', component: DashboardView, meta: { requiresAuth: true } },
    { path: '/users', component: UsersView, meta: { requiresAuth: true } },
    { path: '/scripts', component: ScriptsView, meta: { requiresAuth: true } },
    { path: '/tasks', component: TasksView, meta: { requiresAuth: true } },
    { path: '/reports', component: ReportsView, meta: { requiresAuth: true } },
    { path: '/settings', component: SettingsView, meta: { requiresAuth: true } },
  ],
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (to.meta.requiresAuth === false && auth.isAuthenticated) {
    return '/dashboard'
  }
  if (to.meta.requiresAuth !== false && !auth.isAuthenticated) {
    return '/login'
  }
  return true
})

export default router

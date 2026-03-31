import { createRouter, createWebHistory } from 'vue-router'

import HomePage from '@/pages/home.vue'
import LibraryPage from '@/pages/library.vue'
import ProfilePage from '@/pages/profile.vue'
import SettingsPage from '@/pages/settings.vue'
import LoginPage from '@/pages/login.vue'
import DashboardPage from '@/pages/dashboard.vue'
import CreateJobPage from '@/pages/create-job.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomePage,
    },
    {
      path: '/library',
      name: 'library',
      component: LibraryPage,
    },
    {
      path: '/profile',
      name: 'profile',
      component: ProfilePage,
    },
    {
      path: '/settings',
      name: 'settings',
      component: SettingsPage,
    },
    {
      path: '/login',
      name: 'login',
      component: LoginPage,
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: DashboardPage,
    },
    {
      path: '/create-job',
      name: 'create-job',
      component: CreateJobPage,
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
})

export default router

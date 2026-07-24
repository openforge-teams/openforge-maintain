import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/',
    component: () => import('@/layout/MainLayout.vue'),
    meta: { requiresAuth: true },
    redirect: '/dashboard',
    children: [
      { path: 'dashboard', name: 'Dashboard', component: () => import('@/views/Dashboard.vue'), meta: { title: 'menu.dashboard' } },
      { path: 'host', name: 'Host', component: () => import('@/views/host/Index.vue'), meta: { title: 'menu.host' } },
      { path: 'files', name: 'Files', component: () => import('@/views/files/Index.vue'), meta: { title: 'menu.files' } },
      { path: 'containers', name: 'Containers', component: () => import('@/views/containers/Index.vue'), meta: { title: 'menu.containers' } },
      { path: 'images', name: 'Images', component: () => import('@/views/images/Index.vue'), meta: { title: 'menu.images' } },
      { path: 'volumes', name: 'Volumes', component: () => import('@/views/volumes/Index.vue'), meta: { title: 'menu.volumes' } },
      { path: 'networks', name: 'Networks', component: () => import('@/views/networks/Index.vue'), meta: { title: 'menu.networks' } },
      { path: 'compose', name: 'Compose', component: () => import('@/views/compose/Index.vue'), meta: { title: 'menu.compose' } },
      { path: 'websites', name: 'Websites', component: () => import('@/views/websites/Index.vue'), meta: { title: 'menu.websites' } },
      { path: 'databases/mysql', name: 'MySQL', component: () => import('@/views/databases/MySQL.vue'), meta: { title: 'menu.mysql' } },
      { path: 'databases/postgres', name: 'Postgres', component: () => import('@/views/databases/Postgres.vue'), meta: { title: 'menu.postgres' } },
      { path: 'databases/redis', name: 'Redis', component: () => import('@/views/databases/Redis.vue'), meta: { title: 'menu.redis' } },
      { path: 'appstore', name: 'AppStore', component: () => import('@/views/appstore/Index.vue'), meta: { title: 'menu.appstore' } },
      { path: 'appstore/:id', name: 'AppDetail', component: () => import('@/views/appstore/Detail.vue'), meta: { title: 'menu.appstore' } },
      { path: 'cron', name: 'Cron', component: () => import('@/views/cron/Index.vue'), meta: { title: 'menu.cron' } },
      { path: 'ssl', name: 'SSL', component: () => import('@/views/ssl/Index.vue'), meta: { title: 'menu.ssl' } },
      { path: 'backup', name: 'Backup', component: () => import('@/views/backup/Index.vue'), meta: { title: 'menu.backup' } },
      { path: 'firewall', name: 'Firewall', component: () => import('@/views/firewall/Index.vue'), meta: { title: 'menu.firewall' } },
      { path: 'settings', name: 'Settings', component: () => import('@/views/settings/Index.vue'), meta: { title: 'menu.settings' } },
      { path: 'terminal', name: 'Terminal', component: () => import('@/views/terminal/Index.vue'), meta: { title: 'menu.terminal' } },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth !== false && !token) {
    next('/login')
  } else if (to.path === '/login' && token) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router

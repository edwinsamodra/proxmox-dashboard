import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/LoginView.vue'),
    meta: { public: true },
  },
  {
    path: '/',
    component: () => import('@/components/layout/AppLayout.vue'),
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('@/views/DashboardView.vue'),
      },
      {
        path: 'nodes',
        name: 'Nodes',
        component: () => import('@/views/NodesView.vue'),
      },
      {
        path: 'vms',
        name: 'VMs',
        component: () => import('@/views/VMsView.vue'),
      },
      {
        path: 'containers',
        name: 'Containers',
        component: () => import('@/views/ContainersView.vue'),
      },
      {
        path: 'storage',
        name: 'Storage',
        component: () => import('@/views/StorageView.vue'),
      },
      {
        path: 'network',
        name: 'Network',
        component: () => import('@/views/NetworkView.vue'),
      },
      {
        path: 'firewall',
        name: 'Firewall',
        component: () => import('@/views/FirewallView.vue'),
      },
      {
        path: 'backup',
        name: 'Backup',
        component: () => import('@/views/BackupView.vue'),
      },
      {
        path: 'tofu',
        name: 'OpenTofu',
        component: () => import('@/views/TofuView.vue'),
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router

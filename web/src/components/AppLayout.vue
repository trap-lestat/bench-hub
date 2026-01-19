<template>
  <div class="layout">
    <aside class="sidebar">
      <div class="brand">
        <div class="brand-mark">G</div>
        <div>
          <p class="brand-title">Bench Hub</p>
          <p class="brand-subtitle">Admin Console</p>
        </div>
      </div>
      <nav class="nav">
        <router-link to="/dashboard">概览</router-link>
        <router-link to="/users">用户</router-link>
        <router-link to="/scripts">脚本</router-link>
        <router-link to="/tasks">任务</router-link>
        <router-link to="/reports">报告</router-link>
        <router-link to="/settings">设置</router-link>
      </nav>
      <div class="sidebar-foot">
        <span class="badge">v0.1</span>
      </div>
    </aside>

    <div class="main">
      <header class="topbar">
        <div class="topbar-left">
          <h1 class="page-title">{{ pageTitle }}</h1>
          <p class="page-subtitle">关注核心指标与运行状态</p>
        </div>
        <div class="topbar-right">
          <div class="status-dot"></div>
          <span class="status-text">API 在线</span>
          <button class="ghost" type="button" @click="logout">退出</button>
        </div>
      </header>
      <main class="content">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const titleMap = {
  '/dashboard': '仪表盘',
  '/users': '用户管理',
  '/scripts': '脚本管理',
  '/tasks': '压测任务',
  '/reports': '报告中心',
  '/settings': '系统设置',
}

const pageTitle = computed(() => titleMap[route.path] || '仪表盘')

function logout() {
  auth.clearToken()
  router.push('/login')
}
</script>

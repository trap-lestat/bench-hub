<template>
  <div class="login">
    <div class="login-card">
      <div class="login-header">
        <p class="kicker">Bench Hub</p>
        <h1>欢迎回来</h1>
        <p class="muted">登录后进入管理后台</p>
      </div>

      <form class="login-form" @submit.prevent="submit">
        <label>
          用户名
          <input v-model.trim="username" type="text" placeholder="admin" required />
        </label>
        <label>
          密码
          <input v-model.trim="password" type="password" placeholder="••••••" required />
        </label>
        <button class="primary" type="submit" :disabled="loading">
          {{ loading ? '登录中...' : '登录' }}
        </button>
        <p v-if="error" class="error">{{ error }}</p>
      </form>

      <div class="login-footer">
        <span>默认测试账号：admin / admin123</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '../lib/api'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const auth = useAuthStore()

const username = ref('admin')
const password = ref('admin123')
const loading = ref(false)
const error = ref('')

async function submit() {
  error.value = ''
  loading.value = true
  try {
    const response = await api.post('/api/v1/auth/login', {
      username: username.value,
      password: password.value,
    })
    const token = response?.data?.data?.access_token
    if (!token) {
      throw new Error('登录失败')
    }
    auth.setToken(token)
    router.push('/dashboard')
  } catch (err) {
    error.value = '登录失败，请检查账号或服务状态'
  } finally {
    loading.value = false
  }
}
</script>

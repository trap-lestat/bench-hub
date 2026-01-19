<template>
  <section class="list-page">
    <div class="panel">
      <div class="panel-head">
        <h3>用户列表</h3>
        <button class="primary" type="button">新建用户</button>
      </div>
      <div v-if="error" class="error">{{ error }}</div>
      <div v-else>
        <table class="table">
          <thead>
            <tr>
              <th>ID</th>
              <th>用户名</th>
              <th>创建时间</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="user in users" :key="user.id">
              <td>{{ user.id }}</td>
              <td>{{ user.username }}</td>
              <td>{{ formatDate(user.created_at) }}</td>
            </tr>
            <tr v-if="users.length === 0">
              <td colspan="3" class="empty">暂无数据</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </section>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import api from '../lib/api'

const users = ref([])
const error = ref('')

function formatDate(value) {
  if (!value) return '-'
  const date = new Date(value)
  return date.toLocaleString()
}

async function load() {
  try {
    const response = await api.get('/api/v1/users', {
      params: { page: 1, page_size: 20 },
    })
    users.value = response?.data?.data?.items || []
  } catch (err) {
    error.value = '无法获取用户数据'
  }
}

onMounted(load)
</script>

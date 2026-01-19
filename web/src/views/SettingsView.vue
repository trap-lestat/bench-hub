<template>
  <section class="settings">
    <div class="panel">
      <h3>系统设置</h3>
      <form class="settings-form" @submit.prevent="save">
        <label>
          API 地址
          <input type="text" placeholder="http://localhost:8080" disabled />
        </label>
        <label>
          P95 延迟基线
          <input v-model.trim="p95Baseline" type="text" placeholder="P95 < 300ms" />
        </label>
        <button class="primary" type="submit" :disabled="saving">
          {{ saving ? '保存中...' : '保存配置' }}
        </button>
        <p v-if="message" class="muted">{{ message }}</p>
      </form>
    </div>
  </section>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import api from '../lib/api'

const p95Baseline = ref('P95 < 300ms')
const saving = ref(false)
const message = ref('')

async function load() {
  try {
    const response = await api.get('/api/v1/settings/p95-baseline')
    p95Baseline.value = response?.data?.data?.value || p95Baseline.value
  } catch (err) {
    message.value = '无法获取配置'
  }
}

async function save() {
  saving.value = true
  message.value = ''
  try {
    await api.put('/api/v1/settings/p95-baseline', { value: p95Baseline.value })
    message.value = '已保存'
  } catch (err) {
    message.value = '保存失败'
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

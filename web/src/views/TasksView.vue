<template>
  <section class="list-page">
    <div class="panel">
      <div class="panel-head">
        <h3>压测任务</h3>
        <button class="primary" type="button" @click="toggleForm">
          {{ showForm ? '收起' : '新建任务' }}
        </button>
      </div>

      <div v-if="showForm" class="form-grid">
        <label>
          任务名称
          <input v-model.trim="form.name" type="text" placeholder="login-smoke" />
        </label>
        <label>
          脚本
          <select v-model="form.script_id">
            <option value="" disabled>请选择脚本</option>
            <option v-for="script in scriptOptions" :key="script.id" :value="script.id">
              {{ script.name }}
            </option>
          </select>
        </label>
        <label>
          用户数
          <input v-model.number="form.users_count" type="number" min="1" />
        </label>
        <label>
          每秒增用户数
          <input v-model.number="form.spawn_rate" type="number" min="1" />
        </label>
        <label>
          时长（秒）
          <input v-model.number="form.duration_seconds" type="number" min="10" />
        </label>
        <div class="form-actions">
          <button class="primary" type="button" @click="createTask">保存</button>
          <button class="ghost" type="button" @click="resetForm">清空</button>
        </div>
      </div>

      <div v-if="error" class="error">{{ error }}</div>
      <div v-else>
        <table class="table">
          <thead>
            <tr>
              <th>名称</th>
              <th>脚本</th>
              <th>用户数</th>
              <th>状态</th>
              <th>创建时间</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="task in tasks" :key="task.id">
              <td>{{ task.name }}</td>
              <td>{{ scriptName(task.script_id) }}</td>
              <td>{{ task.users_count }}</td>
              <td><span class="status" :class="task.status">{{ task.status }}</span></td>
              <td>{{ formatDate(task.created_at) }}</td>
              <td>
                <button class="ghost" type="button" @click="runTask(task.id)">运行</button>
                <button class="ghost" type="button" @click="stopTask(task.id)">停止</button>
              </td>
            </tr>
            <tr v-if="tasks.length === 0">
              <td colspan="6" class="empty">暂无任务</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </section>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import api from '../lib/api'

const tasks = ref([])
const error = ref('')
const showForm = ref(false)
const scriptOptions = ref([])

const form = reactive({
  name: '',
  script_id: '',
  users_count: 50,
  spawn_rate: 5,
  duration_seconds: 300,
})

function toggleForm() {
  showForm.value = !showForm.value
}

function resetForm() {
  form.name = ''
  form.script_id = ''
  form.users_count = 50
  form.spawn_rate = 5
  form.duration_seconds = 300
}

function formatDate(value) {
  if (!value) return '-'
  return new Date(value).toLocaleString()
}

function scriptName(id) {
  const match = scriptOptions.value.find((script) => script.id === id)
  return match ? match.name : id
}

async function load() {
  try {
    const response = await api.get('/api/v1/tasks', {
      params: { page: 1, page_size: 50 },
    })
    tasks.value = response?.data?.data?.items || []
  } catch (err) {
    error.value = '无法获取任务数据'
  }
}

async function loadScripts() {
  try {
    const response = await api.get('/api/v1/scripts', {
      params: { page: 1, page_size: 100 },
    })
    scriptOptions.value = response?.data?.data?.items || []
  } catch (err) {
    error.value = '无法获取脚本列表'
  }
}

async function createTask() {
  try {
    await api.post('/api/v1/tasks', {
      name: form.name,
      script_id: form.script_id,
      users_count: form.users_count,
      spawn_rate: form.spawn_rate,
      duration_seconds: form.duration_seconds,
    })
    resetForm()
    await load()
  } catch (err) {
    error.value = '创建任务失败'
  }
}

async function stopTask(id) {
  try {
    await api.post(`/api/v1/tasks/${id}/stop`)
    await load()
  } catch (err) {
    error.value = '停止任务失败'
  }
}

async function runTask(id) {
  try {
    await api.post(`/api/v1/tasks/${id}/run`)
    await load()
  } catch (err) {
    error.value = '运行任务失败'
  }
}

onMounted(load)
onMounted(loadScripts)
</script>

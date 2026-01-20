<template>
  <section class="list-page">
    <div class="panel">
      <div class="panel-head">
        <h3>压测任务</h3>
        <button class="primary" type="button" @click="toggleForm">
          {{ showForm ? '收起' : (form.id ? '编辑任务' : '新建任务') }}
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
        <label>
          目标地址（可选）
          <input v-model.trim="form.target_host" type="text" placeholder="http://api:8080" />
        </label>
        <label>
          JMeter 吞吐（每分钟）
          <input v-model.number="form.jmeter_tpm" type="number" min="1" placeholder="60" />
        </label>
        <div class="form-actions">
          <button class="primary" type="button" @click="createTask">保存</button>
          <button class="ghost" type="button" @click="resetForm">清空</button>
        </div>
      </div>

      <div class="filter-bar">
        <label>
          脚本类型筛选
          <select v-model="filterType">
            <option value="">全部</option>
            <option value="locust">Locust</option>
            <option value="jmeter">JMeter</option>
          </select>
        </label>
      </div>

      <div v-if="error" class="error">{{ error }}</div>
      <div v-else>
        <table class="table">
          <thead>
            <tr>
              <th>名称</th>
              <th>脚本</th>
              <th>类型</th>
              <th>用户数</th>
              <th>时长</th>
              <th>目标地址</th>
              <th>JMeter 吞吐</th>
              <th>状态</th>
              <th>创建时间</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="task in filteredTasks" :key="task.id">
              <td>{{ task.name }}</td>
              <td>{{ scriptName(task.script_id) }}</td>
              <td><span class="tag" :class="scriptType(task.script_id)">{{ scriptType(task.script_id) }}</span></td>
              <td>{{ task.users_count }}</td>
              <td>{{ task.duration_seconds }}s</td>
              <td>{{ task.target_host || "-" }}</td>
              <td>{{ task.jmeter_tpm ?? "-" }}</td>
              <td><span class="status" :class="task.status">{{ task.status }}</span></td>
              <td>{{ formatDate(task.created_at) }}</td>
              <td>
                <button class="ghost" type="button" @click="openRunModal(task)">运行</button>
                <button class="ghost" type="button" @click="stopTask(task.id)">停止</button>
              </td>
            </tr>
            <tr v-if="filteredTasks.length === 0">
              <td colspan="10" class="empty">暂无任务</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div v-if="showRunModal" class="modal-mask" @click.self="closeRunModal">
      <div class="modal">
        <div class="modal-head">
          <div>
            <h4>运行任务</h4>
            <p class="muted">任务：{{ runningTaskName || '未命名' }}</p>
          </div>
          <button class="ghost" type="button" @click="closeRunModal">关闭</button>
        </div>
        <div class="modal-body">
          <label>
            目标地址
            <input v-model.trim="runTargetHost" type="text" :placeholder="defaultTargetHint" />
            <span class="hint">留空则使用默认目标：{{ defaultTargetHint }}</span>
          </label>
        </div>
        <div class="modal-actions">
          <button class="ghost" type="button" @click="closeRunModal">取消</button>
          <button class="primary" type="button" @click="confirmRun">开始运行</button>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import api from '../lib/api'

const tasks = ref([])
const error = ref('')
const showForm = ref(false)
const scriptOptions = ref([])
const filterType = ref("")

const showRunModal = ref(false)
const runTargetHost = ref('')
const defaultTargetHost = import.meta.env.VITE_DEFAULT_TARGET_HOST || 'http://api:8080'
const defaultTargetHint = `默认 ${defaultTargetHost}`
const runningTaskId = ref('')
const runningTaskName = ref('')


const filteredTasks = computed(() => {
  if (!filterType.value) return tasks.value
  return tasks.value.filter((task) => scriptType(task.script_id) === filterType.value)
})

const form = reactive({
  name: '',
  script_id: '',
  users_count: 50,
  spawn_rate: 5,
  duration_seconds: 300,
  target_host: "",
  jmeter_tpm: null,
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
  form.target_host = ""
  form.jmeter_tpm = null
}

function formatDate(value) {
  if (!value) return '-'
  return new Date(value).toLocaleString()
}

function scriptName(id) {
  const match = scriptOptions.value.find((script) => script.id === id)
  return match ? match.name : id
}


function scriptType(id) {
  const match = scriptOptions.value.find((script) => script.id === id)
  return match ? match.type || 'locust' : '-'
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


function editTask(task) {
  showForm.value = true
  form.id = task.id
  form.name = task.name
  form.script_id = task.script_id
  form.users_count = task.users_count
  form.spawn_rate = task.spawn_rate
  form.duration_seconds = task.duration_seconds
  form.target_host = task.target_host || ''
  form.jmeter_tpm = task.jmeter_tpm ?? null
}

async function createTask() {
  try {
    const payload = {
      name: form.name,
      script_id: form.script_id,
      users_count: form.users_count,
      spawn_rate: form.spawn_rate,
      duration_seconds: form.duration_seconds,
      target_host: form.target_host || '',
      jmeter_tpm: form.jmeter_tpm || undefined,
    }
    if (form.id) {
      await api.put(`/api/v1/tasks/${form.id}`, payload)
    } else {
      await api.post('/api/v1/tasks', payload)
    }
    resetForm()
    await load()
  } catch (err) {
    error.value = form.id ? '更新任务失败' : '创建任务失败'
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

function openRunModal(task) {
  runningTaskId.value = task.id
  runningTaskName.value = task.name || ''
  runTargetHost.value = ''
  showRunModal.value = true
}

function closeRunModal() {
  showRunModal.value = false
  runningTaskId.value = ''
  runningTaskName.value = ''
  runTargetHost.value = ''
}

async function confirmRun() {
  if (!runningTaskId.value) return
  const trimmedHost = runTargetHost.value.trim()
  const payload = trimmedHost ? { target_host: trimmedHost } : undefined
  try {
    await api.post(`/api/v1/tasks/${runningTaskId.value}/run`, payload)
    closeRunModal()
    await load()
  } catch (err) {
    error.value = '运行任务失败'
  }
}

onMounted(load)
onMounted(loadScripts)
</script>

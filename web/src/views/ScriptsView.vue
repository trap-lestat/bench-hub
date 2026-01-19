<template>
  <section class="list-page">
    <div class="panel">
      <div class="panel-head">
        <h3>压测脚本</h3>
        <button class="primary" type="button" @click="toggleForm">
          {{ showForm ? '收起' : '新增脚本' }}
        </button>
      </div>

      <div v-if="showForm" class="form-grid">
        <label>
          脚本名称
          <input v-model.trim="form.name" type="text" placeholder="login_flow" />
        </label>
        <label>
          描述
          <input v-model.trim="form.description" type="text" placeholder="登录 + 列表" />
        </label>
        <label class="full">
          脚本内容
          <textarea v-model="form.content" rows="6" placeholder="from locust import ..."></textarea>
        </label>
        <div class="form-actions">
          <button class="primary" type="button" @click="submitScript">
            {{ form.id ? '更新' : '保存' }}
          </button>
          <button class="ghost" type="button" @click="resetForm">清空</button>
        </div>
      </div>

      <div class="import-box">
        <label>
          导入脚本（.py）
          <input type="file" @change="handleFile" />
        </label>
        <button class="ghost" type="button" @click="importScript" :disabled="!importFile">导入</button>
      </div>

      <div v-if="error" class="error">{{ error }}</div>
      <div v-else>
        <table class="table">
          <thead>
            <tr>
              <th>名称</th>
              <th>描述</th>
              <th>更新时间</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="script in scripts" :key="script.id">
              <td>{{ script.name }}</td>
              <td>{{ script.description || '-' }}</td>
              <td>{{ formatDate(script.updated_at) }}</td>
              <td>
                <button class="ghost" type="button" @click="editScript(script)">编辑</button>
                <button class="ghost" type="button" @click="removeScript(script.id)">删除</button>
              </td>
            </tr>
            <tr v-if="scripts.length === 0">
              <td colspan="4" class="empty">暂无脚本</td>
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

const scripts = ref([])
const error = ref('')
const showForm = ref(false)
const importFile = ref(null)

const form = reactive({
  id: '',
  name: '',
  description: '',
  content: '',
})

function toggleForm() {
  showForm.value = !showForm.value
}

function resetForm() {
  form.id = ''
  form.name = ''
  form.description = ''
  form.content = ''
}

function formatDate(value) {
  if (!value) return '-'
  return new Date(value).toLocaleString()
}

async function load() {
  try {
    const response = await api.get('/api/v1/scripts', {
      params: { page: 1, page_size: 50 },
    })
    scripts.value = response?.data?.data?.items || []
  } catch (err) {
    error.value = '无法获取脚本数据'
  }
}

async function submitScript() {
  try {
    if (form.id) {
      await api.put(`/api/v1/scripts/${form.id}`, {
        name: form.name,
        description: form.description,
        content: form.content,
      })
    } else {
      await api.post('/api/v1/scripts', {
        name: form.name,
        description: form.description,
        content: form.content,
      })
    }
    resetForm()
    await load()
  } catch (err) {
    error.value = form.id ? '更新脚本失败' : '创建脚本失败'
  }
}

function editScript(script) {
  showForm.value = true
  form.id = script.id
  form.name = script.name
  form.description = script.description || ''
  form.content = script.content || ''
}

function handleFile(event) {
  const file = event.target.files[0]
  importFile.value = file || null
}

async function importScript() {
  if (!importFile.value) return
  const data = new FormData()
  data.append('name', importFile.value.name.replace(/\.py$/, ''))
  data.append('file', importFile.value)

  try {
    await api.post('/api/v1/scripts/import', data)
    importFile.value = null
    await load()
  } catch (err) {
    error.value = '导入脚本失败'
  }
}

async function removeScript(id) {
  try {
    await api.delete(`/api/v1/scripts/${id}`)
    await load()
  } catch (err) {
    error.value = '删除脚本失败'
  }
}

onMounted(load)
</script>

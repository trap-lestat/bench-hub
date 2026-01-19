<template>
  <section class="list-page">
    <div class="panel">
      <div class="panel-head">
        <h3>压测报告</h3>
      </div>

      <div v-if="error" class="error">{{ error }}</div>
      <div v-else>
        <table class="table">
          <thead>
            <tr>
              <th>名称</th>
              <th>类型</th>
              <th>任务</th>
              <th>生成时间</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="report in reports" :key="report.id">
              <td>{{ report.name }}</td>
              <td>{{ report.type }}</td>
              <td>{{ report.task_name || report.task_id || '-' }}</td>
              <td>{{ formatDate(report.created_at) }}</td>
              <td>
                <button
                  v-if="report.type === 'html'"
                  class="ghost"
                  type="button"
                  @click="previewReport(report)"
                >预览</button>
                <button class="ghost" type="button" @click="downloadReport(report)">下载</button>
              </td>
            </tr>
            <tr v-if="reports.length === 0">
              <td colspan="5" class="empty">暂无报告</td>
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

const reports = ref([])
const error = ref('')

function formatDate(value) {
  if (!value) return '-'
  return new Date(value).toLocaleString()
}

async function load() {
  try {
    const response = await api.get('/api/v1/reports', {
      params: { page: 1, page_size: 50 },
    })
    reports.value = response?.data?.data?.items || []
  } catch (err) {
    error.value = '无法获取报告数据'
  }
}

async function previewReport(report) {
  try {
    const response = await api.get(`/api/v1/reports/${report.id}/preview`, {
      responseType: 'blob',
    })
    const url = URL.createObjectURL(response.data)
    window.open(url, '_blank', 'noopener,noreferrer')
  } catch (err) {
    error.value = '无法预览报告'
  }
}

async function downloadReport(report) {
  try {
    const response = await api.get(`/api/v1/reports/${report.id}/download`, {
      responseType: 'blob',
    })
    const url = URL.createObjectURL(response.data)
    const link = document.createElement('a')
    link.href = url
    link.download = report.name || 'report'
    document.body.appendChild(link)
    link.click()
    link.remove()
    URL.revokeObjectURL(url)
  } catch (err) {
    error.value = '无法下载报告'
  }
}

onMounted(load)
</script>

<template>
  <section class="dashboard">
    <div class="stats">
      <div class="card" style="animation-delay: 0ms;">
        <p class="card-label">当前用户</p>
        <h2>{{ summary.users_count }}</h2>
        <p class="card-meta">已注册用户</p>
      </div>
      <div class="card" style="animation-delay: 120ms;">
        <p class="card-label">脚本总量</p>
        <h2>{{ summary.scripts_count }}</h2>
        <p class="card-meta">已维护脚本</p>
      </div>
      <div class="card" style="animation-delay: 240ms;">
        <p class="card-label">报告数量</p>
        <h2>{{ summary.reports_count }}</h2>
        <p class="card-meta">已生成报告</p>
      </div>
      <div class="card" style="animation-delay: 360ms;">
        <p class="card-label">P95 延迟</p>
        <h2>{{ summary.p95_baseline }}</h2>
        <p class="card-meta">压测基线</p>
      </div>
    </div>

    <div class="grid">
      <div class="panel" style="animation-delay: 320ms;">
        <h3>运行摘要</h3>
        <ul class="summary">
          <li>用户数：{{ summary.users_count }}</li>
          <li>脚本数：{{ summary.scripts_count }}</li>
          <li>报告数：{{ summary.reports_count }}</li>
          <li>当前基线：{{ summary.p95_baseline }}</li>
          <li>常规目标：P95 &lt; 300ms，错误率 &lt; 1%</li>
          <li>稳定性：30 分钟压测无明显波动</li>
        </ul>
      </div>
      <div class="panel" style="animation-delay: 420ms;">
        <h3>下一步建议</h3>
        <ol class="steps">
          <li v-for="item in suggestions" :key="item">{{ item }}</li>
        </ol>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed, onMounted, reactive } from 'vue'
import api from '../lib/api'

const summary = reactive({
  users_count: 0,
  scripts_count: 0,
  reports_count: 0,
  p95_baseline: 'P95 < 300ms',
})

const suggestions = computed(() => {
  const items = []
  if (summary.scripts_count === 0) {
    items.push('创建一条压测脚本')
  }
  if (summary.reports_count === 0) {
    items.push('运行一次压测任务并生成报告')
  }
  if (items.length === 0) {
    items.push('查看最新报告并确认 P95 是否达标')
  }
  items.push('根据当前基线优化压测目标')
  return items
})

async function load() {
  try {
    const response = await api.get('/api/v1/dashboard/summary')
    const data = response?.data?.data
    if (data) {
      summary.users_count = data.users_count
      summary.scripts_count = data.scripts_count
      summary.reports_count = data.reports_count
      summary.p95_baseline = data.p95_baseline
    }
  } catch (err) {
    // keep defaults if backend unavailable
  }
}

onMounted(load)
</script>

<template>
  <div class="app-shell">
    <AppLayout v-if="showLayout">
      <router-view />
    </AppLayout>
    <router-view v-else />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import AppLayout from './components/AppLayout.vue'

const route = useRoute()
const showLayout = computed(() => route.meta.requiresAuth !== false)
</script>

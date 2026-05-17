<script setup lang="ts">
import { NFlex, NText } from "naive-ui"
import { onMounted } from "vue"
import { useAuthStore } from "@/stores/auth"
import { useApiCall } from "@/composables/useApiCall"

const authStore = useAuthStore()
const apiCall = useApiCall()

onMounted(async () => {
  if (authStore.user != null) return
  await apiCall(() => authStore.ensureLoaded())
})
</script>

<template>
  <n-flex align="center" :size="8">
    <n-text strong>tech-generator</n-text>
    <n-text v-if="authStore.user != null" depth="3">· {{ authStore.user.name }}</n-text>
  </n-flex>
</template>

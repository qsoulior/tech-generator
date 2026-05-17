<script setup lang="ts">
import { NLayout, NLayoutContent, NFlex, NResult, NButton } from "naive-ui"
import { useRouter } from "vue-router"

const router = useRouter()

function goHome() {
  router.push({ name: "auth" })
}
</script>

<template>
  <n-layout style="height: 100vh">
    <n-layout-content content-class="layout-content" embedded>
      <n-flex vertical align="center" justify="center" style="height: 100vh">
        <n-result status="404" title="Страница не найдена" description="Такого адреса здесь нет">
          <template #footer>
            <n-button @click="goHome">На главную</n-button>
          </template>
        </n-result>
      </n-flex>
    </n-layout-content>
  </n-layout>
</template>

<style scoped>
:deep(.layout-content) {
  padding: 1.5rem;
}
</style>

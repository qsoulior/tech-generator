<script setup lang="ts">
import { RouterView } from "vue-router"
import { NConfigProvider, NGlobalStyle, NMessageProvider, NLayout, NLayoutContent, NBackTop } from "naive-ui"
</script>

<template>
  <n-config-provider>
    <n-global-style />
    <n-message-provider>
      <n-layout style="height: 100vh">
        <n-back-top />
        <n-layout-content content-class="layout-content" embedded>
          <RouterView />
        </n-layout-content>
      </n-layout>
    </n-message-provider>
  </n-config-provider>
</template>

<style scoped>
:deep(.layout-content) {
  padding: 1rem;
  min-height: 100vh;
}

@media screen and (min-width: 768px) {
  :deep(.layout-content) {
    padding: 2rem;
    min-height: 100vh;
  }
}
</style>

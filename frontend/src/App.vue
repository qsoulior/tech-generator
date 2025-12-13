<script setup lang="ts">
import { RouterView } from "vue-router"
import { NConfigProvider, NGlobalStyle, NMessageProvider } from "naive-ui"
</script>

<template>
  <n-config-provider>
    <n-global-style />
    <n-message-provider>
      <RouterView />
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { NLayout, NLayoutContent, NFlex, NCard, NTabs, NTabPane } from "naive-ui"
import SignInForm from "@/components/SignInForm.vue"
import SignUpForm from "@/components/SignUpForm.vue"
</script>

<template>
  <n-layout class="layout">
    <n-layout-content content-class="layout-content" embedded>
      <n-flex vertical align="center" justify="center" class="content">
        <n-card class="card">
          <n-tabs type="line" size="large">
            <n-tab-pane name="Вход">
              <SignInForm />
            </n-tab-pane>
            <n-tab-pane name="Регистрация">
              <SignUpForm />
            </n-tab-pane>
          </n-tabs>
        </n-card>
      </n-flex>
    </n-layout-content>
  </n-layout>
</template>

<style scoped>
.layout,
.content {
  min-height: 100dvh;
}

.card {
  width: 40rem;
}

:deep(.layout-content) {
  padding: 1rem;
  min-height: 100dvh;
}

@media screen and (min-width: 768px) {
  :deep(.layout-content) {
    padding: 2rem;
  }
}
</style>

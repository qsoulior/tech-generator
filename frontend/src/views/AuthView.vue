<script setup lang="ts">
import { NLayout, NLayoutContent, NFlex, NCard, NTabs, NTabPane } from "naive-ui"
import SignInForm from "@/components/SignInForm.vue"
import SignUpForm from "@/components/SignUpForm.vue"
</script>

<template>
  <n-layout style="height: 100vh">
    <n-layout-content content-class="layout-content" embedded>
      <n-flex vertical align="center" justify="center" style="height: 100vh">
        <n-card style="width: 40rem">
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

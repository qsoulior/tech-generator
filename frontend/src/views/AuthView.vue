<script setup lang="ts">
import { NFlex, NCard, NTabs, NTabPane } from "naive-ui"
import SignInForm from "@/components/SignInForm.vue"
import SignUpForm from "@/components/SignUpForm.vue"
</script>

<template>
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
</template>

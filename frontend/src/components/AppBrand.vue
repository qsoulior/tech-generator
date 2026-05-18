<script setup lang="ts">
import { NFlex, NText, NDropdown, useMessage } from "naive-ui"
import { onMounted } from "vue"
import { useRouter } from "vue-router"
import { useAuthStore } from "@/stores/auth"
import { useApiCall } from "@/composables/useApiCall"

const authStore = useAuthStore()
const apiCall = useApiCall()
const router = useRouter()
const message = useMessage()

onMounted(async () => {
  if (authStore.user != null) return
  await apiCall(() => authStore.ensureLoaded())
})

const userMenuOptions = [{ key: "signOut", label: "Выйти" }]

async function onUserMenuSelect(key: string) {
  if (key !== "signOut") return
  const r = await apiCall(() => authStore.signOut())
  if (!r.ok) return
  message.success("Вы вышли из аккаунта")
  router.push({ name: "auth" })
}
</script>

<template>
  <n-flex align="center" :size="8">
    <n-text strong>tech-generator</n-text>
    <n-dropdown
      v-if="authStore.user != null"
      :options="userMenuOptions"
      trigger="click"
      placement="bottom-start"
      @select="onUserMenuSelect"
    >
      <n-text depth="3" class="user-trigger">· {{ authStore.user.name }}</n-text>
    </n-dropdown>
  </n-flex>
</template>

<style scoped>
.user-trigger {
  cursor: pointer;
}
</style>

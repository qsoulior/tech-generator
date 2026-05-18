<script setup lang="ts">
import { NDropdown, NFlex, NIcon, NText, useMessage } from "naive-ui"
import { onMounted } from "vue"
import { useRouter } from "vue-router"
import { useAuthStore } from "@/stores/auth"
import { useApiCall } from "@/composables/useApiCall"
import IconUserOutlined from "@/components/icons/IconUserOutlined.vue"

const authStore = useAuthStore()
const apiCall = useApiCall()
const router = useRouter()
const message = useMessage()

onMounted(async () => {
  if (authStore.user != null) return
  await apiCall(() => authStore.ensureLoaded())
})

const options = [{ key: "signOut", label: "Выйти" }]

async function onSelect(key: string) {
  if (key !== "signOut") return
  const r = await apiCall(() => authStore.signOut())
  if (!r.ok) return
  message.success("Вы вышли из аккаунта")
  router.push({ name: "auth" })
}
</script>

<template>
  <n-dropdown
    v-if="authStore.user != null"
    :options="options"
    trigger="click"
    placement="bottom-end"
    @select="onSelect"
  >
    <n-flex align="center" :size="6" :wrap="false" class="trigger">
      <n-icon :size="18">
        <IconUserOutlined />
      </n-icon>
      <n-text>{{ authStore.user.name }}</n-text>
    </n-flex>
  </n-dropdown>
</template>

<style scoped>
.trigger {
  cursor: pointer;
}
</style>

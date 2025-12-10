<script setup lang="ts">
import { NFlex, NCard, NText, NButton, NIcon, NPopconfirm, useMessage } from "naive-ui"
import IconDeleteOutlined from "@/components/icons/IconDeleteOutlined.vue"
import { useRouter } from "vue-router"

const router = useRouter()
const message = useMessage()

const props = defineProps<{
  id: number
  name: string
  authorName: string
  createdAt: Date
}>()

const emit = defineEmits<{
  delete: []
}>()

async function templateDelete() {
  const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/template/delete/${props.id}`, {
    method: "DELETE",
    credentials: "include",
  })

  if (!response.ok) {
    if (response.status === 401 || response.status === 403) {
      router.push({ name: "auth" })
      return
    }

    const result = await response.json()
    message.error(result.message)
    return
  }

  emit("delete")
}

async function onPositiveClick() {
  await templateDelete()
}
</script>

<template>
  <router-link :to="{}" style="width: 100%; text-decoration: none">
    <n-card>
      <n-flex align="center" justify="space-between">
        <n-flex vertical size="small">
          <n-text strong>{{ props.name }}</n-text>
          <n-text>Автор: {{ props.authorName }}</n-text>
          <n-text>Создан: {{ props.createdAt.toLocaleString() }}</n-text>
        </n-flex>
        <n-popconfirm positive-text="Да" negative-text="Нет" @positive-click="onPositiveClick">
          <template #trigger>
            <n-button secondary @click.prevent>
              <template #icon>
                <n-icon>
                  <IconDeleteOutlined />
                </n-icon>
              </template>
            </n-button>
          </template>
          <template #default>Вы точно хотите удалить шаблон?</template>
        </n-popconfirm>
      </n-flex>
    </n-card>
  </router-link>
</template>

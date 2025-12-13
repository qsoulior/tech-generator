<script setup lang="ts">
import { NFlex, NCard, NText, NTag } from "naive-ui"

const props = defineProps<{
  projectId: number
  templateId: number
  taskId: number
  status: string
  versionNumber: number
  creatorName: string
  createdAt: Date
  updatedAt: Date
}>()

const statusToString = new Map([
  ["created", "В очереди"],
  ["in_progress", "В процессе"],
  ["succeed", "Успешно"],
  ["failed", "Ошибка"],
])

type TagType = "default" | "primary" | "success" | "info" | "warning" | "error"

const statusToType = new Map<string, TagType>([
  ["created", "primary"],
  ["in_progress", "primary"],
  ["succeed", "success"],
  ["failed", "error"],
])
</script>

<template>
  <router-link
    :to="{ name: 'task', params: { projectID: projectId, templateID: templateId, taskID: taskId } }"
    style="width: 100%; text-decoration: none"
  >
    <n-card>
      <n-flex vertical size="small">
        <n-flex align="center" justify="space-between">
          <n-text strong>v{{ props.versionNumber }}</n-text>
          <n-tag :type="statusToType.get(props.status)">{{ statusToString.get(props.status) }}</n-tag>
        </n-flex>
        <n-text>Автор: {{ props.creatorName }}</n-text>
        <n-text>
          Создан: {{ props.createdAt.toLocaleString() }} · Обновлен: {{ props.createdAt.toLocaleString() }}
        </n-text>
      </n-flex>
    </n-card>
  </router-link>
</template>

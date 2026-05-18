<script setup lang="ts">
import { NFlex, NCard, NText, NTag, NIcon, NTooltip } from "naive-ui"
import { computed, type Component } from "vue"
import IconClockCircleOutlined from "@/components/icons/IconClockCircleOutlined.vue"
import IconSyncOutlined from "@/components/icons/IconSyncOutlined.vue"
import IconCheckCircleOutlined from "@/components/icons/IconCheckCircleOutlined.vue"
import IconCloseCircleOutlined from "@/components/icons/IconCloseCircleOutlined.vue"
import { formatRelativeTime } from "@/utils/relativeTime"

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

const statusToIcon = new Map<string, Component>([
  ["created", IconClockCircleOutlined],
  ["in_progress", IconSyncOutlined],
  ["succeed", IconCheckCircleOutlined],
  ["failed", IconCloseCircleOutlined],
])

const createdRelative = computed(() => formatRelativeTime(props.createdAt))
const updatedRelative = computed(() => formatRelativeTime(props.updatedAt))
</script>

<template>
  <router-link
    :to="{ name: 'task', params: { projectID: projectId, templateID: templateId, taskID: taskId } }"
    style="width: 100%; text-decoration: none"
  >
    <n-card>
      <n-flex vertical size="small">
        <n-flex align="center" justify="space-between">
          <n-flex align="baseline" :size="8" :wrap="false">
            <n-text strong>Результат #{{ props.taskId }}</n-text>
            <n-text depth="3">· v{{ props.versionNumber }}</n-text>
          </n-flex>
          <n-tag :type="statusToType.get(props.status)">
            <template #icon>
              <n-icon>
                <component :is="statusToIcon.get(props.status)" />
              </n-icon>
            </template>
            {{ statusToString.get(props.status) }}
          </n-tag>
        </n-flex>
        <n-text>Автор: {{ props.creatorName }}</n-text>
        <n-text>
          Создан:
          <n-tooltip trigger="hover">
            <template #trigger>
              <span class="time">{{ createdRelative }}</span>
            </template>
            {{ props.createdAt.toLocaleString() }}
          </n-tooltip>
          · Обновлён:
          <n-tooltip trigger="hover">
            <template #trigger>
              <span class="time">{{ updatedRelative }}</span>
            </template>
            {{ props.updatedAt.toLocaleString() }}
          </n-tooltip>
        </n-text>
      </n-flex>
    </n-card>
  </router-link>
</template>

<style scoped>
.time {
  border-bottom: 1px dashed currentColor;
  cursor: help;
}
</style>

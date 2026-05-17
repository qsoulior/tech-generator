<script setup lang="ts">
import { NFlex, NCard, NText, NButton, NIcon, NPopconfirm } from "naive-ui"
import IconDeleteOutlined from "@/components/icons/IconDeleteOutlined.vue"
import { templateDelete } from "@/api/template"
import { useApiCall } from "@/composables/useApiCall"

const apiCall = useApiCall()

const props = defineProps<{
  projectId: number
  templateId: number
  name: string
  authorName: string
  createdAt: Date
}>()

const emit = defineEmits<{
  delete: []
}>()

async function onPositiveClick() {
  const r = await apiCall(() => templateDelete(props.templateId))
  if (!r.ok) return
  emit("delete")
}
</script>

<template>
  <router-link
    :to="{ name: 'template', params: { projectID: projectId, templateID: templateId } }"
    style="width: 100%; text-decoration: none"
  >
    <n-card>
      <n-flex align="center" justify="space-between">
        <n-flex vertical size="small">
          <n-text strong>{{ props.name }}</n-text>
          <n-text>Автор: {{ props.authorName }}</n-text>
          <n-text>Создан: {{ props.createdAt.toLocaleString() }}</n-text>
        </n-flex>
        <n-popconfirm positive-text="Да" negative-text="Нет" @positive-click="onPositiveClick">
          <template #trigger>
            <n-button secondary aria-label="Удалить шаблон" title="Удалить шаблон" @click.prevent>
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

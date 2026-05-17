<script setup lang="ts">
import { NFlex, NCard, NText, NButton, NIcon, NPopconfirm } from "naive-ui"
import { ref } from "vue"
import IconDeleteOutlined from "@/components/icons/IconDeleteOutlined.vue"
import IconEditOutlined from "@/components/icons/IconEditOutlined.vue"
import TemplateUpdateModal from "@/components/TemplateUpdateModal.vue"
import { templateDelete } from "@/api/template"
import { useApiCall } from "@/composables/useApiCall"
import { useTemplateStore } from "@/stores/template"

const apiCall = useApiCall()
const templateStore = useTemplateStore()

const props = defineProps<{
  projectId: number
  templateId: number
  name: string
  authorName: string
  createdAt: Date
}>()

const emit = defineEmits<{
  delete: []
  update: [id: number, name: string]
}>()

const showUpdateModal = ref(false)

async function onPositiveClick() {
  const r = await apiCall(() => templateDelete(props.templateId))
  if (!r.ok) return
  emit("delete")
}

function onEditClick() {
  showUpdateModal.value = true
}

function onUpdateSubmit(name: string) {
  templateStore.invalidate(props.templateId)
  emit("update", props.templateId, name)
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
        <n-flex>
          <n-button
            secondary
            aria-label="Редактировать шаблон"
            title="Редактировать шаблон"
            @click.prevent="onEditClick"
          >
            <template #icon>
              <n-icon>
                <IconEditOutlined />
              </n-icon>
            </template>
          </n-button>
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
      </n-flex>
    </n-card>
    <TemplateUpdateModal
      v-model:show-modal="showUpdateModal"
      :template-id="templateId"
      :initial-name="name"
      @submit="onUpdateSubmit"
    />
  </router-link>
</template>

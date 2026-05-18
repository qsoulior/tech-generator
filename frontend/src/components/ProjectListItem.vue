<script setup lang="ts">
import { NFlex, NCard, NText, NButton, NIcon, NPopconfirm } from "naive-ui"
import { ref } from "vue"
import IconDeleteOutlined from "@/components/icons/IconDeleteOutlined.vue"
import IconEditOutlined from "@/components/icons/IconEditOutlined.vue"
import ProjectUpdateModal from "@/components/ProjectUpdateModal.vue"
import { projectDelete } from "@/api/project"
import { useApiCall } from "@/composables/useApiCall"
import { useProjectStore } from "@/stores/project"

const apiCall = useApiCall()
const projectStore = useProjectStore()

const props = defineProps<{
  id: number
  name: string
  authorName: string
}>()

const emit = defineEmits<{
  delete: [id: number]
  update: [id: number, name: string]
}>()

const showUpdateModal = ref(false)

async function onPositiveClick() {
  const r = await apiCall(() => projectDelete(props.id))
  if (!r.ok) return
  emit("delete", props.id)
}

function onEditClick() {
  showUpdateModal.value = true
}

function onUpdateSubmit(name: string) {
  projectStore.put(props.id, { name, authorName: props.authorName })
  emit("update", props.id, name)
}
</script>

<template>
  <router-link :to="{ name: 'project', params: { projectID: id } }" style="width: 100%; text-decoration: none">
    <n-card>
      <n-flex align="center" justify="space-between">
        <n-flex vertical size="small">
          <n-text strong>{{ props.name }}</n-text>
          <n-text>Автор: {{ props.authorName }}</n-text>
        </n-flex>
        <n-flex>
          <n-button
            secondary
            aria-label="Редактировать проект"
            title="Редактировать проект"
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
              <n-button secondary aria-label="Удалить проект" title="Удалить проект" @click.prevent>
                <template #icon>
                  <n-icon>
                    <IconDeleteOutlined />
                  </n-icon>
                </template>
              </n-button>
            </template>
            <template #default>Вы точно хотите удалить проект?</template>
          </n-popconfirm>
        </n-flex>
      </n-flex>
    </n-card>
    <ProjectUpdateModal
      v-model:show-modal="showUpdateModal"
      :project-id="id"
      :initial-name="name"
      @submit="onUpdateSubmit"
    />
  </router-link>
</template>

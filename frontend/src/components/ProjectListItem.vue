<script setup lang="ts">
import { NFlex, NCard, NText, NButton, NIcon, NPopconfirm } from "naive-ui"
import IconDeleteOutlined from "@/components/icons/IconDeleteOutlined.vue"
import IconEditOutlined from "@/components/icons/IconEditOutlined.vue"
import { projectDelete } from "@/api/project"
import { useApiCall } from "@/composables/useApiCall"

const apiCall = useApiCall()

const props = defineProps<{
  id: number
  name: string
  authorName: string
}>()

const emit = defineEmits<{
  delete: [id: number]
}>()

async function onPositiveClick() {
  const r = await apiCall(() => projectDelete(props.id))
  if (!r.ok) return
  emit("delete", props.id)
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
          <n-button secondary @click.prevent>
            <template #icon>
              <n-icon>
                <IconEditOutlined />
              </n-icon>
            </template>
          </n-button>
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
            <template #default>Вы точно хотите удалить проект?</template>
          </n-popconfirm>
        </n-flex>
      </n-flex>
    </n-card>
  </router-link>
</template>

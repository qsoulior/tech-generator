<script setup lang="ts">
import TaskListItem from "@/components/TaskListItem.vue"
import AppHeader from "@/components/AppHeader.vue"
import { NLayout, NLayoutContent, NFlex, NPagination, NText, NEmpty, NSpin } from "naive-ui"
import { onMounted, ref } from "vue"
import { taskList as fetchTasks, type TaskStatus } from "@/api/task"
import { useApiCall } from "@/composables/useApiCall"
import { usePagination } from "@/composables/usePagination"

const props = defineProps<{
  templateID: number
  projectID: number
}>()

const apiCall = useApiCall()

const { page, pageSize, totalPages, pageSizes } = usePagination("результатов")
const totalTasks = ref(0)
const loading = ref(true)

interface Task {
  id: number
  status: TaskStatus
  versionNumber: number
  creatorName: string
  createdAt: Date
  updatedAt: Date
}

const tasks = ref<Task[]>([])

async function taskList() {
  loading.value = true
  const r = await apiCall(() =>
    fetchTasks({
      templateID: props.templateID,
      page: page.value,
      size: pageSize.value,
    }),
  )
  loading.value = false
  if (!r.ok) return

  totalTasks.value = r.value.totalTasks
  totalPages.value = r.value.totalPages

  tasks.value = r.value.tasks.map((task) => ({
    id: task.id,
    status: task.status,
    versionNumber: task.versionNumber,
    creatorName: task.creatorName,
    createdAt: new Date(task.createdAt),
    updatedAt: new Date(task.updatedAt ?? task.createdAt),
  }))
}

onMounted(async () => {
  await taskList()
})

async function onUpdatePage() {
  await taskList()
}

async function onUpdatePageSize() {
  await taskList()
}
</script>

<template>
  <n-layout>
    <AppHeader />
    <n-layout content-style="height: calc(100vh - 59px)">
      <n-layout-content content-class="layout-content" embedded style="height: 100%">
        <n-flex v-if="loading" justify="center" align="center" style="height: 100%">
          <n-spin size="large" />
        </n-flex>
        <n-flex v-else-if="totalTasks === 0" justify="center" align="center" style="height: 100%">
          <n-empty description="Результатов генерации пока нет" />
        </n-flex>
        <n-flex v-else vertical align="center" style="max-width: 50rem; margin: auto">
          <n-text depth="3" style="width: 100%">Всего: {{ totalTasks }}</n-text>
          <TaskListItem
            v-for="task in tasks"
            :key="task.id"
            :project-id="projectID"
            :template-id="templateID"
            :task-id="task.id"
            :status="task.status"
            :version-number="task.versionNumber"
            :creator-name="task.creatorName"
            :created-at="task.createdAt"
            :updated-at="task.updatedAt"
          />
          <n-pagination
            v-model:page="page"
            v-model:page-size="pageSize"
            :page-count="totalPages"
            show-size-picker
            :page-sizes="pageSizes"
            @update:page="onUpdatePage"
            @update:page-size="onUpdatePageSize"
          />
        </n-flex>
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<style scoped>
:deep(.layout-content) {
  padding: 1.5rem;
}
</style>

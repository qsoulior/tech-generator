<script setup lang="ts">
import TaskListItem from "@/components/TaskListItem.vue"
import {
  NLayout,
  NLayoutHeader,
  NLayoutContent,
  NFlex,
  NText,
  NMenu,
  NPagination,
  type MenuOption,
  useMessage,
} from "naive-ui"
import { h, onMounted, ref } from "vue"
import { RouterLink, useRouter } from "vue-router"

const props = defineProps<{
  templateID: number
  projectID: number
}>()

const router = useRouter()
const message = useMessage()

const totalTasks = ref(0)
const totalPages = ref(0)
const page = ref(1)
const pageSize = ref(50)

const pageSizes = [
  { label: "10 результатов", value: 10 },
  { label: "50 результатов", value: 50 },
  { label: "100 результатов", value: 100 },
  { label: "500 результатов", value: 500 },
]

interface Task {
  id: number
  status: string
  versionNumber: number
  creatorName: string
  createdAt: Date
  updatedAt: Date
}

const templateName = ref("")
const tasks = ref<Task[]>([])

// запрос списка задач
interface TaskListResultTask {
  id: number
  status: string
  versionNumber: number
  creatorName: string
  createdAt: string
  updatedAt: string
}

interface TaskListResult {
  tasks: TaskListResultTask[]
  totalTasks: number
  totalPages: number
}

async function taskList() {
  const params = new URLSearchParams({
    page: page.value.toString(),
    size: pageSize.value.toString(),
    templateID: props.templateID.toString(),
  })

  const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/task/list?${params}`, {
    method: "GET",
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

  const result: TaskListResult = await response.json()

  totalTasks.value = result.totalTasks
  totalPages.value = result.totalPages

  tasks.value = result.tasks.map((task) => {
    return {
      id: task.id,
      status: task.status,
      versionNumber: task.versionNumber,
      creatorName: task.creatorName,
      createdAt: new Date(task.createdAt),
      updatedAt: new Date(task.updatedAt),
    }
  })
}

// запрос шаблона
interface TemplateGetResultVersion {
  id: number
  number: number
}

interface TemplateGetResult {
  name: string
  version?: TemplateGetResultVersion
}

async function templateGet() {
  const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/template/get/${props.templateID}`, {
    method: "GET",
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

  const result: TemplateGetResult = await response.json()
  templateName.value = result.name
}

onMounted(async () => {
  await templateGet()
  await taskList()
})

async function onUpdatePage() {
  await taskList()
}

async function onUpdatePageSize() {
  await taskList()
}

const menuOptionsCenter: MenuOption[] = [
  {
    label: () =>
      h(
        RouterLink,
        {
          to: {
            name: "template",
            params: {
              projectID: props.projectID,
              templateID: props.templateID,
            },
          },
        },
        { default: () => templateName.value },
      ),
    key: "template",
  },
]

const menuOptionsRight: MenuOption[] = [
  {
    label: () =>
      h(
        RouterLink,
        {
          to: {
            name: "projectList",
          },
        },
        { default: () => "Проекты" },
      ),
    key: "projectList",
  },
  {
    label: () =>
      h(
        RouterLink,
        {
          to: {
            name: "project",
            params: { projectID: props.projectID },
          },
        },
        { default: () => "Шаблоны" },
      ),
    key: "project",
  },
]
</script>

<template>
  <n-layout>
    <n-layout-header bordered style="padding: 0.5rem 1rem">
      <n-flex align="center" justify="space-between">
        <n-text strong>tech-generator</n-text>
        <n-flex>
          <n-menu mode="horizontal" :options="menuOptionsCenter" />
        </n-flex>
        <n-flex>
          <n-menu mode="horizontal" :options="menuOptionsRight" />
        </n-flex>
      </n-flex>
    </n-layout-header>
    <n-layout content-style="height: calc(100vh - 59px)">
      <n-layout-content content-class="layout-content" embedded style="height: 100%">
        <n-flex vertical align="center" style="max-width: 50rem; margin: auto">
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

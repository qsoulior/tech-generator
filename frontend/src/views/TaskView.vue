<script setup lang="ts">
import {
  NLayout,
  NLayoutHeader,
  NLayoutContent,
  NFlex,
  NIcon,
  NText,
  NButton,
  NCard,
  NTag,
  NSpin,
  NEmpty,
  NPagination,
  NAlert,
} from "naive-ui"
import { onMounted, ref, computed, watch, type Component } from "vue"
import { MdEditor, config, type ToolbarNames } from "md-editor-v3"
import RU from "@vavt/cm-extension/dist/locale/ru"
import "md-editor-v3/lib/style.css"
import AppHeader from "@/components/AppHeader.vue"
import TaskErrorListItem from "@/components/TaskErrorListItem.vue"
import TaskErrorListSearch from "@/components/TaskErrorListSearch.vue"
import IconDownloadOutlined from "@/components/icons/IconDownloadOutlined.vue"
import IconClockCircleOutlined from "@/components/icons/IconClockCircleOutlined.vue"
import IconSyncOutlined from "@/components/icons/IconSyncOutlined.vue"
import IconCheckCircleOutlined from "@/components/icons/IconCheckCircleOutlined.vue"
import IconCloseCircleOutlined from "@/components/icons/IconCloseCircleOutlined.vue"
import { taskGet, type TaskGetError, type TaskGetVariableError, type TaskStatus } from "@/api/task"
import { useApiCall } from "@/composables/useApiCall"
import { usePagination } from "@/composables/usePagination"
import { useTemplateStore } from "@/stores/template"
import { fromBase64 } from "@/utils/base64"

const apiCall = useApiCall()
const templateStore = useTemplateStore()

const props = defineProps<{
  templateID: number
  projectID: number
  taskID: number
}>()

const templateName = ref("")
const data = ref<string | null>(null)
const error = ref<TaskGetError | null>(null)
const status = ref<TaskStatus | null>(null)
const loading = ref(true)

type TagType = "default" | "primary" | "success" | "info" | "warning" | "error"

const statusToString = new Map<string, string>([
  ["created", "В очереди"],
  ["in_progress", "В процессе"],
  ["succeed", "Успешно"],
  ["failed", "Ошибка"],
])

const statusToType = new Map<string, TagType>([
  ["created", "default"],
  ["in_progress", "default"],
  ["succeed", "success"],
  ["failed", "error"],
])

const statusToIcon = new Map<string, Component>([
  ["created", IconClockCircleOutlined],
  ["in_progress", IconSyncOutlined],
  ["succeed", IconCheckCircleOutlined],
  ["failed", IconCloseCircleOutlined],
])

const isPending = computed(() => status.value === "created" || status.value === "in_progress")

const pendingMessage = computed(() => {
  if (status.value === "created") return "Задача ожидает обработки в очереди"
  if (status.value === "in_progress") return "Задача обрабатывается"
  return ""
})

const hasSystemMessage = computed(() => error.value?.message != null && error.value.message !== "")
const variableErrors = computed<TaskGetVariableError[]>(() => error.value?.variableErrors ?? [])
const hasErrorDetails = computed(() => hasSystemMessage.value || variableErrors.value.length > 0)

const errorSearch = ref("")
const { page, pageSize, totalPages, pageSizes } = usePagination("ошибок", 10)

const filteredVariableErrors = computed<TaskGetVariableError[]>(() => {
  const q = errorSearch.value.trim().toLowerCase()
  if (q === "") return variableErrors.value
  return variableErrors.value.filter((ve) => {
    if (ve.title?.toLowerCase().includes(q)) return true
    if (ve.name?.toLowerCase().includes(q)) return true
    if (ve.message?.toLowerCase().includes(q)) return true
    return (
      ve.constraintErrors?.some((ce) => ce.name?.toLowerCase().includes(q) || ce.message?.toLowerCase().includes(q)) ??
      false
    )
  })
})

const totalVariableErrors = computed(() => filteredVariableErrors.value.length)

const pagedVariableErrors = computed(() => {
  const start = (page.value - 1) * pageSize.value
  return filteredVariableErrors.value.slice(start, start + pageSize.value)
})

watch(
  [totalVariableErrors, pageSize],
  () => {
    totalPages.value = Math.max(1, Math.ceil(totalVariableErrors.value / pageSize.value))
    if (page.value > totalPages.value) page.value = 1
  },
  { immediate: true },
)

watch(errorSearch, () => {
  page.value = 1
})

async function loadTask() {
  loading.value = true
  const r = await apiCall(() => taskGet(props.taskID))
  loading.value = false
  if (!r.ok) return

  status.value = r.value.task.status

  if (r.value.result != null) {
    data.value = fromBase64(r.value.result)
  } else {
    error.value = r.value.task.error ?? null
  }
}

async function loadTemplate() {
  const r = await apiCall(() => templateStore.ensureLoaded(props.templateID))
  if (!r.ok) return
  templateName.value = r.value.name
}

function download() {
  if (data.value == null) return

  const file = new Blob([data.value], { type: "text/markdown" })

  const url = URL.createObjectURL(file)

  const a = document.createElement("a")
  a.href = url
  a.download = `${templateName.value}.md`

  document.body.appendChild(a)
  a.click()

  document.body.removeChild(a)
  window.URL.revokeObjectURL(url)
}

onMounted(async () => {
  await loadTemplate()
  await loadTask()
})

config({
  editorConfig: {
    languageUserDefined: {
      ru: RU,
    },
  },
})

const toolbars: ToolbarNames[] = ["preview", "previewOnly"]
</script>

<template>
  <div class="page">
    <AppHeader />
    <n-layout-header bordered class="toolbar">
      <n-flex align="center" justify="space-between" :wrap="false">
        <n-tag v-if="status != null" :type="statusToType.get(status)" size="small">
          <template #icon>
            <n-icon>
              <component :is="statusToIcon.get(status)" />
            </n-icon>
          </template>
          {{ statusToString.get(status) }}
        </n-tag>
        <n-button v-if="data != null" size="small" secondary @click="download()">
          <template #icon>
            <n-icon>
              <IconDownloadOutlined />
            </n-icon>
          </template>
          Скачать
        </n-button>
      </n-flex>
    </n-layout-header>
    <n-layout class="page-body">
      <n-layout-content content-class="layout-content" embedded style="height: 100%">
        <n-flex v-if="loading" justify="center" align="center" style="height: 100%">
          <n-spin size="large" />
        </n-flex>
        <MdEditor
          v-else-if="data != null"
          v-model="data"
          language="ru"
          :toolbars="toolbars"
          style="height: 100%"
          :read-only="true"
        />
        <div v-else-if="error != null" class="error-scroll">
          <n-flex vertical align="center" class="error-content">
            <n-alert v-if="hasSystemMessage" type="error" style="width: 100%">
              <div class="system-error-message">{{ error.message }}</div>
            </n-alert>
            <template v-if="variableErrors.length > 0">
              <TaskErrorListSearch v-model:value="errorSearch" />
              <n-text depth="3" style="width: 100%">Всего: {{ totalVariableErrors }}</n-text>
              <TaskErrorListItem
                v-for="variableError in pagedVariableErrors"
                :key="variableError.id"
                :variable-error="variableError"
              />
              <n-flex v-if="totalVariableErrors === 0" justify="center" style="width: 100%; padding: 2rem 0">
                <n-empty description="Ошибок по запросу не найдено" />
              </n-flex>
              <n-pagination
                v-model:page="page"
                v-model:page-size="pageSize"
                :page-count="totalPages"
                show-size-picker
                :page-sizes="pageSizes"
              />
            </template>
            <n-flex v-if="!hasErrorDetails" justify="center" style="width: 100%; padding: 2rem 0">
              <n-empty description="Подробности ошибки недоступны" />
            </n-flex>
          </n-flex>
        </div>
        <n-flex v-else-if="isPending" justify="center" align="center" style="height: 100%">
          <n-card style="max-width: 30rem; width: 100%">
            <n-flex vertical align="center" :size="16">
              <n-icon :size="48" :depth="3">
                <component :is="statusToIcon.get(status!)" />
              </n-icon>
              <n-text>{{ pendingMessage }}</n-text>
            </n-flex>
          </n-card>
        </n-flex>
        <n-flex v-else justify="center" align="center" style="height: 100%">
          <n-empty description="Нет данных для отображения" />
        </n-flex>
      </n-layout-content>
    </n-layout>
  </div>
</template>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  height: 100vh;
}

.page-body {
  flex: 1;
  min-height: 0;
}

.toolbar {
  padding: 0.5rem 1rem;
  flex-shrink: 0;
}

:deep(.layout-content) {
  padding: 1.5rem;
}

.system-error-message {
  white-space: pre-wrap;
  word-break: break-word;
}

.error-scroll {
  height: 100%;
  overflow-y: auto;
}

.error-content {
  max-width: 50rem;
  margin: auto;
}
</style>

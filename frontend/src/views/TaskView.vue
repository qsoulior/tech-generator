<script setup lang="ts">
import {
  NLayout,
  NLayoutHeader,
  NLayoutContent,
  NFlex,
  NText,
  NMenu,
  NButton,
  NCard,
  NTable,
  useMessage,
  type MenuOption,
} from "naive-ui"
import { h, onMounted, ref } from "vue"
import { RouterLink, useRouter } from "vue-router"
import { MdEditor, config, type ToolbarNames } from "md-editor-v3"
import RU from "@vavt/cm-extension/dist/locale/ru"
import "md-editor-v3/lib/preview.css"

const router = useRouter()
const message = useMessage()

const props = defineProps<{
  templateID: number
  projectID: number
  taskID: number
}>()

const templateName = ref("")
const data = ref<string | null>(null)
const error = ref<TaskGetResultError | null>(null)

// запрос задачи
interface TaskGetResultConstraintError {
  id: string
  name: string
  message: string
}

interface TaskGetResultVariableError {
  id: number
  name: string
  message?: string
  constraintErrors: TaskGetResultConstraintError[]
}

interface TaskGetResultError {
  message?: string
  variableErrors: TaskGetResultVariableError[]
}

interface TaskGetResultTask {
  id: number
  versionID: number
  status: string
  payload: Record<string, string>
  error: TaskGetResultError
  creatorName: string
  createdAt: string
  updatedAt: string
}

interface TaskGetResult {
  task: TaskGetResultTask
  result: string | null
}

async function taskGet() {
  const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/task/get/${props.taskID}`, {
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

  const result: TaskGetResult = await response.json()

  if (result.result != null) {
    data.value = fromBase64(result.result)
  } else {
    error.value = result.task.error
  }
}

function fromBase64(data: string) {
  const bin = atob(data)
  const base64 = Uint8Array.from(bin, (m) => m.codePointAt(0) ?? 0)
  return new TextDecoder().decode(base64)
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

// триггеры
onMounted(async () => {
  await templateGet()
  await taskGet()
})

config({
  editorConfig: {
    languageUserDefined: {
      ru: RU,
    },
  },
})

const toolbars: ToolbarNames[] = ["preview", "previewOnly"]

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
  {
    label: () =>
      h(
        RouterLink,
        {
          to: {
            name: "taskList",
            params: { projectID: props.projectID, templateID: props.templateID },
          },
        },
        { default: () => "Результаты" },
      ),
    key: "taskList",
  },
]
</script>

<template>
  <n-layout>
    <n-layout-header bordered style="padding: 0.5rem 1rem">
      <n-flex align="center" justify="space-between">
        <n-text strong>tech-generator</n-text>
        <n-flex align="center">
          <n-flex>
            <n-menu mode="horizontal" :options="menuOptionsCenter" />
          </n-flex>
          <n-button v-if="data != null" secondary @click="download()">Скачать</n-button>
        </n-flex>
        <n-flex>
          <n-menu mode="horizontal" :options="menuOptionsRight" />
        </n-flex>
      </n-flex>
    </n-layout-header>
    <n-layout content-style="height: calc(100vh - 59px)">
      <n-layout-content content-class="layout-content" embedded style="height: 100%">
        <MdEditor
          v-if="data != null"
          v-model="data"
          language="ru"
          :toolbars="toolbars"
          style="height: 100%"
          :read-only="true"
        />
        <n-flex v-if="error != null" justify="center">
          <n-card style="width: 50rem">
            <template #header>Ошибки</template>
            <template #default>
              <n-flex vertical size="large">
                <n-text v-if="error.message != undefined">Общая ошибка: {{ error.message }}</n-text>
                <n-flex v-for="variableError in error.variableErrors" :key="variableError.id" vertical>
                  <n-text strong>{{ variableError.name }}</n-text>
                  <n-text v-if="variableError.message != undefined">{{ variableError.message }}</n-text>
                  <n-table :single-line="false" size="small">
                    <thead>
                      <tr>
                        <th>Ограничение</th>
                        <th>Ошибка</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="constraintError in variableError.constraintErrors" :key="constraintError.id">
                        <td>
                          {{ constraintError.name }}
                        </td>
                        <td>
                          {{ constraintError.message }}
                        </td>
                      </tr>
                    </tbody>
                  </n-table>
                </n-flex>
              </n-flex>
            </template>
          </n-card>
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

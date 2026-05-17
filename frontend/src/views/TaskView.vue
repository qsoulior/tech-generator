<script setup lang="ts">
import { NLayout, NLayoutHeader, NLayoutContent, NFlex, NText, NButton, NCard, NTable } from "naive-ui"
import { computed, onMounted, ref } from "vue"
import { MdEditor, config, type ToolbarNames } from "md-editor-v3"
import RU from "@vavt/cm-extension/dist/locale/ru"
import "md-editor-v3/lib/preview.css"
import HeaderMenu, { type HeaderMenuItem } from "@/components/HeaderMenu.vue"
import { taskGet, type TaskGetError } from "@/api/task"
import { templateGet } from "@/api/template"
import { useApiCall } from "@/composables/useApiCall"
import { fromBase64 } from "@/utils/base64"

const apiCall = useApiCall()

const props = defineProps<{
  templateID: number
  projectID: number
  taskID: number
}>()

const templateName = ref("")
const data = ref<string | null>(null)
const error = ref<TaskGetError | null>(null)

async function loadTask() {
  const r = await apiCall(() => taskGet(props.taskID))
  if (!r.ok) return

  if (r.value.result != null) {
    data.value = fromBase64(r.value.result)
  } else {
    error.value = r.value.task.error ?? null
  }
}

async function loadTemplate() {
  const r = await apiCall(() => templateGet(props.templateID))
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

const menuItemsCenter = computed<HeaderMenuItem[]>(() => [
  {
    key: "template",
    label: templateName.value,
    to: { name: "template", params: { projectID: props.projectID, templateID: props.templateID } },
  },
])

const menuItemsRight: HeaderMenuItem[] = [
  { key: "projectList", label: "Проекты", to: { name: "projectList" } },
  { key: "project", label: "Шаблоны", to: { name: "project", params: { projectID: props.projectID } } },
  {
    key: "taskList",
    label: "Результаты",
    to: { name: "taskList", params: { projectID: props.projectID, templateID: props.templateID } },
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
            <HeaderMenu :items="menuItemsCenter" />
          </n-flex>
          <n-button v-if="data != null" secondary @click="download()">Скачать</n-button>
        </n-flex>
        <n-flex>
          <HeaderMenu :items="menuItemsRight" />
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

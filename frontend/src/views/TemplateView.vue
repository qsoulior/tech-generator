<script setup lang="ts">
import {
  NLayout,
  NLayoutHeader,
  NLayoutSider,
  NLayoutContent,
  NScrollbar,
  NFlex,
  NButton,
  NIcon,
  NText,
  NDivider,
  NTooltip,
  useMessage,
} from "naive-ui"
import { computed, onMounted, ref } from "vue"
import { MdEditor, config, type ToolbarNames } from "md-editor-v3"
import RU from "@vavt/cm-extension/dist/locale/ru"
import "md-editor-v3/lib/style.css"
import VariableListSearch from "@/components/VariableListSearch.vue"
import VariableCreateModal from "@/components/VariableCreateModal.vue"
import VariableUpdateModal from "@/components/VariableUpdateModal.vue"
import IconDeleteOutlined from "@/components/icons/IconDeleteOutlined.vue"
import IconEditOutlined from "@/components/icons/IconEditOutlined.vue"
import IconSaveOutlined from "@/components/icons/IconSaveOutlined.vue"
import IconPlayCircleOutlined from "@/components/icons/IconPlayCircleOutlined.vue"
import IconUnorderedListOutlined from "@/components/icons/IconUnorderedListOutlined.vue"
import IconDownloadOutlined from "@/components/icons/IconDownloadOutlined.vue"
import TaskCreateModal from "@/components/TaskCreateModal.vue"
import TemplateUpdateModal from "@/components/TemplateUpdateModal.vue"
import AppHeader from "@/components/AppHeader.vue"
import { versionCreate, type VersionCreateVariable } from "@/api/version"
import type { TemplateImportPayload, TemplateImportVariable } from "@/api/template"
import { useApiCall } from "@/composables/useApiCall"
import { useTemplateStore } from "@/stores/template"
import { fromBase64, toBase64 } from "@/utils/base64"
import { formatRelativeTime } from "@/utils/relativeTime"
import { useRouter } from "vue-router"

const message = useMessage()
const apiCall = useApiCall()
const templateStore = useTemplateStore()
const router = useRouter()

function goToResults() {
  router.push({ name: "taskList", params: { projectID: props.projectID, templateID: props.templateID } })
}

const props = defineProps<{
  templateID: number
  projectID: number
}>()

const showCreateModal = ref(false)
const showUpdateModal = ref(false)
const showTaskCreateModal = ref(false)
const showTemplateUpdateModal = ref(false)

const inputTypeToString = new Map([
  ["input", "Входная"],
  ["computed", "Вычисляемая"],
])

const typeToString = new Map([
  ["string", "Строка"],
  ["integer", "Целое число"],
  ["float", "Вещественное число"],
])

const name = ref("")
const versionNumber = ref(0)
const versionID = ref<number>()
const versionCreatedAt = ref<Date>()
const data = ref("")

const versionCreatedRelative = computed(() =>
  versionCreatedAt.value != undefined ? formatRelativeTime(versionCreatedAt.value) : "",
)

interface Constraint {
  name: string
  expression: string
  isActive: boolean
}

interface Variable {
  name: string
  title: string
  type: string
  expression: string
  inputType: string
  constraints: Constraint[]
}

const variables = ref<Variable[]>([])
const variableUpdating = ref<Variable>()
const variableUpdatingIndex = ref<number>()
const variableSearch = ref("")

const filteredVariables = computed(() => {
  const query = variableSearch.value.trim().toLowerCase()
  const indexed = variables.value.map((variable, index) => ({ variable, index }))
  if (query === "") return indexed
  return indexed.filter(
    ({ variable }) => variable.title.toLowerCase().includes(query) || variable.name.toLowerCase().includes(query),
  )
})

const createOccupiedSlugs = computed(() => variables.value.map((v) => v.name))

const updateOccupiedSlugs = computed(() =>
  variables.value.filter((_, index) => index !== variableUpdatingIndex.value).map((v) => v.name),
)

function escapeRegExp(value: string): string {
  return value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")
}

function renameSlugInTemplate(oldSlug: string, newSlug: string, excludeIndex: number) {
  const wordRe = new RegExp(`\\b${escapeRegExp(oldSlug)}\\b`, "g")
  data.value = data.value.replace(wordRe, newSlug)

  variables.value.forEach((variable, index) => {
    if (index === excludeIndex) return
    variable.expression = variable.expression.replace(wordRe, newSlug)
  })
}

function rewriteOwnConstraints(variable: Variable, oldSlug: string, newSlug: string): Variable {
  const wordRe = new RegExp(`\\b${escapeRegExp(oldSlug)}\\b`, "g")
  return {
    ...variable,
    constraints: variable.constraints.map((c) => ({
      ...c,
      expression: c.expression.replace(wordRe, newSlug),
    })),
  }
}

const savedSnapshot = ref({ data: "", variables: "[]" })

const isDirty = computed(
  () => savedSnapshot.value.data !== data.value || savedSnapshot.value.variables !== JSON.stringify(variables.value),
)

function saveSnapshot() {
  savedSnapshot.value = {
    data: data.value,
    variables: JSON.stringify(variables.value),
  }
}

function variableCreate(variable: Variable) {
  variables.value.push(variable)
}

function variableUpdate(variable: Variable, index: number) {
  variableUpdating.value = {
    name: variable.name,
    title: variable.title,
    type: variable.type,
    expression: variable.expression,
    inputType: variable.inputType,
    constraints: variable.constraints.map((constraint) => ({
      name: constraint.name,
      expression: constraint.expression,
      isActive: constraint.isActive,
    })),
  }
  variableUpdatingIndex.value = index
  showUpdateModal.value = true
}

function handleVariableUpdate(variable: Variable) {
  const index = variableUpdatingIndex.value
  if (index == undefined || index >= variables.value.length) return

  const previous = variables.value[index]
  let next = variable
  if (previous != undefined && previous.name !== "" && previous.name !== variable.name) {
    renameSlugInTemplate(previous.name, variable.name, index)
    next = rewriteOwnConstraints(variable, previous.name, variable.name)
  }

  variables.value[index] = next
}

function variableDelete(index: number) {
  variables.value.splice(index, 1)
}

async function loadTemplate() {
  const r = await apiCall(() => templateStore.ensureLoaded(props.templateID))
  if (!r.ok) return

  name.value = r.value.name

  if (r.value.version != undefined) {
    versionID.value = r.value.version.id
    versionNumber.value = r.value.version.number
    versionCreatedAt.value = new Date(r.value.version.createdAt)
    data.value = fromBase64(r.value.version.data)
    variables.value = r.value.version.variables.map((variable) => ({
      name: variable.name,
      title: variable.title,
      type: variable.type,
      expression: variable.expression ?? "",
      inputType: variable.isInput ? "input" : "computed",
      constraints: variable.constraints.map((constraint) => ({
        name: constraint.name,
        expression: constraint.expression,
        isActive: constraint.isActive,
      })),
    }))
  }

  saveSnapshot()
}

function onTemplateRename(newName: string) {
  name.value = newName
  templateStore.invalidate(props.templateID)
  templateStore.setMeta(props.templateID, { name: newName })
}

function handleExport() {
  const payload: TemplateImportPayload = {
    name: name.value,
    version: {
      data: toBase64(data.value),
      variables: variables.value.map<TemplateImportVariable>((variable) => ({
        name: variable.name,
        title: variable.title,
        type: variable.type as TemplateImportVariable["type"],
        expression: variable.expression || undefined,
        isInput: variable.inputType == "input",
        constraints: variable.constraints.map((constraint) => ({
          name: constraint.name,
          expression: constraint.expression,
          isActive: constraint.isActive,
        })),
      })),
    },
  }

  const blob = new Blob([JSON.stringify(payload, null, 2)], { type: "application/json" })
  const url = URL.createObjectURL(blob)
  const link = document.createElement("a")
  link.href = url
  link.download = `${name.value || "template"}.json`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

async function saveVersion() {
  const r = await apiCall(() =>
    versionCreate({
      templateID: props.templateID,
      data: toBase64(data.value),
      variables: variables.value.map<VersionCreateVariable>((variable) => ({
        name: variable.name,
        title: variable.title,
        type: variable.type as VersionCreateVariable["type"],
        expression: variable.expression,
        isInput: variable.inputType == "input",
        constraints: variable.constraints,
      })),
    }),
  )
  if (!r.ok) return

  versionID.value = r.value.id
  versionNumber.value++
  versionCreatedAt.value = new Date()
  templateStore.invalidate(props.templateID)
  saveSnapshot()
  message.success("Шаблон сохранен")
}

onMounted(async () => {
  await loadTemplate()
})

config({
  editorConfig: {
    languageUserDefined: {
      ru: RU,
    },
  },
})

const toolbars: ToolbarNames[] = [
  "revoke",
  "next",
  "-",
  "bold",
  "underline",
  "italic",
  "strikeThrough",
  "-",
  "title",
  "sub",
  "sup",
  "quote",
  "unorderedList",
  "orderedList",
  "task",
  "-",
  "codeRow",
  "code",
  "link",
  "image",
  "table",
  "mermaid",
  "katex",
  "=",
  "prettier",
  "preview",
  "previewOnly",
]
</script>

<template>
  <div class="page">
    <AppHeader />
    <n-layout-header bordered class="toolbar">
      <n-flex align="center" justify="space-between" :wrap="false">
        <n-flex align="baseline" :size="8" :wrap="false">
          <n-text>
            Версия {{ versionNumber }}
            <template v-if="versionCreatedAt != undefined">
              ·
              <n-tooltip trigger="hover">
                <template #trigger>
                  <span class="time">{{ versionCreatedRelative }}</span>
                </template>
                {{ versionCreatedAt.toLocaleString() }}
              </n-tooltip>
            </template>
          </n-text>
        </n-flex>
        <n-flex align="center" :size="8" :wrap="false">
          <n-button size="small" secondary @click="showTemplateUpdateModal = true">
            <template #icon>
              <n-icon>
                <IconEditOutlined />
              </n-icon>
            </template>
            Редактировать
          </n-button>
          <n-tooltip :disabled="isDirty">
            <template #trigger>
              <span>
                <n-button size="small" secondary :disabled="!isDirty" @click="saveVersion">
                  <template #icon>
                    <n-icon>
                      <IconSaveOutlined />
                    </n-icon>
                  </template>
                  Сохранить
                </n-button>
              </span>
            </template>
            Нет изменений для сохранения
          </n-tooltip>
          <n-tooltip :disabled="!isDirty">
            <template #trigger>
              <span>
                <n-button
                  size="small"
                  secondary
                  :disabled="versionID == undefined || isDirty"
                  @click="showTaskCreateModal = true"
                >
                  <template #icon>
                    <n-icon>
                      <IconPlayCircleOutlined />
                    </n-icon>
                  </template>
                  Выполнить
                </n-button>
              </span>
            </template>
            Сохраните шаблон, чтобы выполнить
          </n-tooltip>
          <n-button size="small" secondary @click="handleExport">
            <template #icon>
              <n-icon>
                <IconDownloadOutlined />
              </n-icon>
            </template>
            Экспорт
          </n-button>
          <n-button size="small" secondary @click="goToResults">
            <template #icon>
              <n-icon>
                <IconUnorderedListOutlined />
              </n-icon>
            </template>
            Результаты
          </n-button>
        </n-flex>
      </n-flex>
    </n-layout-header>
    <TemplateUpdateModal
      v-model:show-modal="showTemplateUpdateModal"
      :template-id="templateID"
      :initial-name="name"
      @submit="onTemplateRename"
    />
    <TaskCreateModal
      v-if="versionID != undefined"
      v-model:show-modal="showTaskCreateModal"
      :version-id="versionID"
      :variables="variables.filter((v) => v.inputType == 'input')"
    />
    <n-layout has-sider class="page-body">
      <n-layout-sider collapse-mode="width" width="25%" :collapsed-width="0" show-trigger="bar" bordered>
        <n-flex vertical class="sider">
          <n-flex vertical class="sider-section">
            <n-text>Переменные</n-text>
            <VariableListSearch v-model:value="variableSearch" />
            <n-button secondary class="full-width" @click="showCreateModal = true">Добавить переменную</n-button>
            <VariableCreateModal
              v-model:show-modal="showCreateModal"
              :occupied-slugs="createOccupiedSlugs"
              @submit="variableCreate"
            />
            <VariableUpdateModal
              v-model:show-modal="showUpdateModal"
              v-model:variable="variableUpdating"
              :occupied-slugs="updateOccupiedSlugs"
              @submit="handleVariableUpdate"
            />
          </n-flex>
          <n-divider class="flush-divider" />
          <n-scrollbar class="variables" content-style="padding: 1rem">
            <n-flex vertical size="large">
              <n-flex
                v-for="{ variable, index } in filteredVariables"
                :key="index"
                align="center"
                justify="space-between"
              >
                <n-flex vertical class="variable-row" :size="0" @click="variableUpdate(variable, index)">
                  <n-text>
                    {{ variable.title }}
                  </n-text>
                  <n-text depth="3" code style="font-size: 0.75rem">
                    {{ variable.name }}
                  </n-text>
                  <n-text depth="3">
                    Тип: {{ inputTypeToString.get(variable.inputType) }} · Значение:
                    {{ typeToString.get(variable.type) }} · Ограничений:
                    {{ variable.constraints.length }}
                  </n-text>
                </n-flex>
                <n-button secondary @click="variableDelete(index)">
                  <template #icon>
                    <n-icon>
                      <IconDeleteOutlined />
                    </n-icon>
                  </template>
                </n-button>
              </n-flex>
            </n-flex>
          </n-scrollbar>
        </n-flex>
      </n-layout-sider>
      <n-layout-content content-class="layout-content" embedded>
        <MdEditor v-model="data" language="ru" :toolbars="toolbars" class="editor" />
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

.sider {
  height: 100%;
}

.sider-section {
  padding: 1rem;
}

.full-width {
  width: 100%;
}

.flush-divider {
  margin: 0;
}

.variables {
  max-height: 100%;
}

.variable-row {
  flex-grow: 1;
  cursor: pointer;
}

.editor {
  height: 100%;
}

.time {
  border-bottom: 1px dashed currentColor;
  cursor: help;
}

:deep(.layout-content) {
  padding-left: 1.5rem;
}
</style>

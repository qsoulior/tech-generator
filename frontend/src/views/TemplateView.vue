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
  useMessage,
} from "naive-ui"
import { onMounted, ref } from "vue"
import { MdEditor, config, type ToolbarNames } from "md-editor-v3"
import RU from "@vavt/cm-extension/dist/locale/ru"
import "md-editor-v3/lib/style.css"
import VariableListSearch from "@/components/VariableListSearch.vue"
import VariableCreateModal from "@/components/VariableCreateModal.vue"
import VariableUpdateModal from "@/components/VariableUpdateModal.vue"
import IconDeleteOutlined from "@/components/icons/IconDeleteOutlined.vue"
import TaskCreateModal from "@/components/TaskCreateModal.vue"
import HeaderMenu, { type HeaderMenuItem } from "@/components/HeaderMenu.vue"
import AppBrand from "@/components/AppBrand.vue"
import { versionCreate, type VersionCreateVariable } from "@/api/version"
import { useApiCall } from "@/composables/useApiCall"
import { useTemplateStore } from "@/stores/template"
import { fromBase64, toBase64 } from "@/utils/base64"

const message = useMessage()
const apiCall = useApiCall()
const templateStore = useTemplateStore()

const props = defineProps<{
  templateID: number
  projectID: number
}>()

const showCreateModal = ref(false)
const showUpdateModal = ref(false)
const showTaskCreateModal = ref(false)

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
const data = ref("")

interface Constraint {
  name: string
  expression: string
  isActive: boolean
}

interface Variable {
  name: string
  type: string
  expression: string
  inputType: string
  constraints: Constraint[]
}

const variables = ref<Variable[]>([])
const variableUpdating = ref<Variable>()
const variableUpdatingIndex = ref<number>()

function variableCreate(variable: Variable) {
  variables.value.push(variable)
}

function variableUpdate(variable: Variable, index: number) {
  variableUpdating.value = {
    name: variable.name,
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
  if (variableUpdatingIndex.value != undefined && variableUpdatingIndex.value < variables.value.length) {
    variables.value[variableUpdatingIndex.value] = variable
  }
}

function variableDelete(index: number) {
  variables.value.splice(index, 1)
}

async function loadTemplate() {
  const r = await apiCall(() => templateStore.ensureLoaded(props.templateID))
  if (!r.ok) return

  name.value = r.value.name

  if (r.value.version == undefined) {
    return
  }

  versionID.value = r.value.version.id
  versionNumber.value = r.value.version.number
  data.value = fromBase64(r.value.version.data)
  variables.value = r.value.version.variables.map((variable) => ({
    name: variable.name,
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

async function saveVersion() {
  const r = await apiCall(() =>
    versionCreate({
      templateID: props.templateID,
      data: toBase64(data.value),
      variables: variables.value.map<VersionCreateVariable>((variable) => ({
        name: variable.name,
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
  templateStore.invalidate(props.templateID)
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

const menuItems: HeaderMenuItem[] = [
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
        <AppBrand />
        <n-flex align="center" :wrap="false">
          <n-text>{{ name }}</n-text>
          <n-text>v{{ versionNumber }}</n-text>
          <n-button secondary @click="saveVersion">Сохранить</n-button>
          <n-button secondary @click="showTaskCreateModal = true">Выполнить</n-button>
          <TaskCreateModal
            v-model:show-modal="showTaskCreateModal"
            :version-id="versionID ?? 0"
            :variables="variables.filter((v) => v.inputType == 'input')"
          />
        </n-flex>
        <n-flex>
          <HeaderMenu :items="menuItems" />
        </n-flex>
      </n-flex>
    </n-layout-header>
    <n-layout has-sider content-style="height: calc(100vh - 59px)">
      <n-layout-sider collapse-mode="width" width="25%" :collapsed-width="0" show-trigger="bar" bordered>
        <n-flex vertical style="height: 100%">
          <n-flex vertical style="padding: 1rem">
            <n-text>Переменные</n-text>
            <VariableListSearch />
            <n-button secondary style="width: 100%" @click="showCreateModal = true">Добавить переменную</n-button>
            <VariableCreateModal
              :template-id="templateID"
              v-model:show-modal="showCreateModal"
              @submit="variableCreate"
            />
            <VariableUpdateModal
              :template-id="templateID"
              v-model:show-modal="showUpdateModal"
              v-model:variable="variableUpdating"
              @submit="handleVariableUpdate"
            />
          </n-flex>
          <n-divider style="margin: 0" />
          <n-scrollbar style="max-height: 100%" content-style="padding: 1rem">
            <n-flex vertical size="large">
              <n-flex v-for="(variable, index) in variables" :key="index" align="center" justify="space-between">
                <n-flex
                  vertical
                  @click="variableUpdate(variable, index)"
                  style="flex-grow: 1; cursor: pointer"
                  :size="0"
                >
                  <n-text>
                    {{ variable.name }}
                  </n-text>
                  <n-text depth="3">
                    Тип: {{ inputTypeToString.get(variable.inputType) }} · Значение:
                    {{ typeToString.get(variable.type) }} · Ограничений:
                    {{ variable.constraints.length }}
                  </n-text>
                  <n-text depth="3"> </n-text>
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
        <MdEditor v-model="data" language="ru" :toolbars="toolbars" style="height: 100%" />
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<style scoped>
:deep(.layout-content) {
  padding-left: 1.5rem;
}
</style>

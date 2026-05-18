<script setup lang="ts">
import {
  NModal,
  NForm,
  NFormItem,
  NInput,
  NButton,
  NTabs,
  NTabPane,
  NUpload,
  NUploadDragger,
  NText,
  NIcon,
  useMessage,
} from "naive-ui"
import type { FormRules, FormInst, UploadFileInfo } from "naive-ui"
import { ref } from "vue"
import {
  templateCreate,
  templateImport,
  type TemplateImportPayload,
  type TemplateImportVariable,
  type TemplateImportConstraint,
} from "@/api/template"
import { useApiCall } from "@/composables/useApiCall"
import IconDownloadOutlined from "@/components/icons/IconDownloadOutlined.vue"

const apiCall = useApiCall()
const message = useMessage()

const showModal = defineModel("showModal", { default: false })

const props = defineProps<{
  projectId: number
}>()

const emit = defineEmits<{
  submit: []
}>()

const formRef = ref<FormInst | null>(null)
const loading = ref(false)
const importing = ref(false)
const fileList = ref<UploadFileInfo[]>([])

interface Model {
  name: string
}

const modelRef = ref<Model>({
  name: "",
})

const rules: FormRules = {
  name: {
    required: true,
    message: "Название шаблона не может быть пустым",
  },
}

function handleValidateClick(e: MouseEvent) {
  e.preventDefault()
  formRef.value?.validate(async (errors) => {
    if (!errors) {
      await submit(modelRef.value)
    }
  })
}

async function submit(model: Model) {
  loading.value = true
  try {
    const r = await apiCall(() => templateCreate({ name: model.name, projectID: props.projectId }))
    if (!r.ok) return

    emit("submit")
    showModal.value = false
  } finally {
    loading.value = false
  }
}

function normalizePayload(raw: unknown): TemplateImportPayload | null {
  if (typeof raw !== "object" || raw === null) return null
  const obj = raw as Record<string, unknown>
  if (typeof obj.name !== "string" || obj.name.length === 0) return null

  const payload: TemplateImportPayload = { name: obj.name }

  if (obj.version != null && typeof obj.version === "object") {
    const version = obj.version as Record<string, unknown>
    if (typeof version.data !== "string" || !Array.isArray(version.variables)) return null

    const variables: TemplateImportVariable[] = version.variables.map((rawVariable) => {
      const v = rawVariable as Record<string, unknown>
      const rawConstraints = Array.isArray(v.constraints) ? v.constraints : []
      const constraints: TemplateImportConstraint[] = rawConstraints.map((rawConstraint) => {
        const c = rawConstraint as Record<string, unknown>
        return {
          name: String(c.name ?? ""),
          expression: String(c.expression ?? ""),
          isActive: Boolean(c.isActive),
        }
      })
      return {
        name: String(v.name ?? ""),
        type: v.type as TemplateImportVariable["type"],
        expression: typeof v.expression === "string" ? v.expression : undefined,
        isInput: Boolean(v.isInput),
        constraints,
      }
    })

    payload.version = {
      data: version.data,
      variables,
    }
  }

  return payload
}

async function handleImportFile({ file }: { file: UploadFileInfo }) {
  const raw = file.file
  if (!raw) return

  importing.value = true
  try {
    const text = await raw.text()
    let parsed: unknown
    try {
      parsed = JSON.parse(text)
    } catch {
      message.error("Файл не является валидным JSON")
      return
    }

    const payload = normalizePayload(parsed)
    if (payload == null) {
      message.error("Структура файла не соответствует ожидаемой")
      return
    }

    const r = await apiCall(() => templateImport({ projectID: props.projectId, template: payload }))
    if (!r.ok) return

    message.success("Шаблон импортирован")
    emit("submit")
    showModal.value = false
  } finally {
    importing.value = false
    fileList.value = []
  }
}

function beforeUpload({ file }: { file: UploadFileInfo }): boolean {
  if (!file.file) return false
  const isJson = file.file.type === "application/json" || /\.json$/i.test(file.file.name)
  if (!isJson) {
    message.error("Допустимы только JSON-файлы")
    return false
  }
  return true
}
</script>

<template>
  <n-modal v-model:show="showModal" preset="card" style="width: 50rem">
    <template #header>Добавление шаблона</template>
    <template #default>
      <n-tabs type="line" animated default-value="create">
        <n-tab-pane name="create" tab="Создать пустой">
          <n-form ref="formRef" :model="modelRef" :rules="rules">
            <n-form-item path="name" label="Название шаблона">
              <n-input v-model:value="modelRef.name" placeholder="Введите название шаблона" />
            </n-form-item>
            <n-form-item>
              <n-button
                style="width: 100%"
                secondary
                type="primary"
                :loading="loading"
                :disabled="loading"
                @click="handleValidateClick"
              >
                Добавить
              </n-button>
            </n-form-item>
          </n-form>
        </n-tab-pane>
        <n-tab-pane name="import" tab="Импортировать">
          <n-upload
            v-model:file-list="fileList"
            accept=".json,application/json"
            :default-upload="false"
            :show-file-list="false"
            :disabled="importing"
            :max="1"
            @before-upload="beforeUpload"
            @change="handleImportFile"
          >
            <n-upload-dragger>
              <div style="margin-bottom: 0.75rem">
                <n-icon size="48" :depth="3">
                  <IconDownloadOutlined />
                </n-icon>
              </div>
              <n-text style="font-size: 16px">Нажмите или перетащите JSON-файл шаблона</n-text>
              <n-text depth="3" style="display: block; margin-top: 0.5rem">
                Будет создан новый шаблон со всеми переменными и ограничениями
              </n-text>
            </n-upload-dragger>
          </n-upload>
        </n-tab-pane>
      </n-tabs>
    </template>
  </n-modal>
</template>

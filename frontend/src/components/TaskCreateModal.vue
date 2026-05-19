<script setup lang="ts">
import { NModal, NForm, NFormItem, NInput, NInputNumber, NButton, NFlex, NText, useMessage } from "naive-ui"
import type { FormInst, FormItemRule } from "naive-ui"
import { ref, watch } from "vue"
import { taskCreate } from "@/api/task"
import { useApiCall } from "@/composables/useApiCall"

const message = useMessage()
const apiCall = useApiCall()

const showModal = defineModel("showModal", { default: false })

const props = defineProps<{
  versionId: number
  variables: Variable[]
}>()

interface Variable {
  name: string
  title: string
  type: string
}

const formRef = ref<FormInst | null>(null)
const loading = ref(false)

type VariableValue = string | number | null

interface ModelVariable {
  name: string
  title: string
  type: string
  value: VariableValue
}

interface Model {
  variables: ModelVariable[]
}

const modelRef = ref<Model>({
  variables: [],
})

const valueRule: FormItemRule = {
  required: true,
  trigger: "submit",
  validator: (_rule, value: VariableValue) => {
    if (value === null || value === undefined || value === "") {
      return new Error("Значение не может быть пустым")
    }
    return true
  },
}

function initialValue(type: string): VariableValue {
  return type === "integer" || type === "float" ? null : ""
}

function asNumber(value: VariableValue): number | null {
  return typeof value === "number" ? value : null
}

function asString(value: VariableValue): string {
  return typeof value === "string" ? value : ""
}

function handleValidateClick(e: MouseEvent) {
  e.preventDefault()
  formRef.value?.validate(async (errors) => {
    if (!errors) {
      await submit()
    }
  })
}

async function submit() {
  loading.value = true
  try {
    const r = await apiCall(() =>
      taskCreate({
        versionID: props.versionId,
        payload: Object.fromEntries(
          modelRef.value.variables.map((v) => [v.name, v.value === null ? "" : String(v.value)]),
        ),
      }),
    )
    if (!r.ok) return

    message.success("Запущена генерация документа")
    showModal.value = false
  } finally {
    loading.value = false
  }
}

watch(showModal, (value) => {
  if (value) {
    modelRef.value.variables = props.variables.map((variable) => ({
      name: variable.name,
      title: variable.title,
      type: variable.type,
      value: initialValue(variable.type),
    }))
  }
})
</script>

<template>
  <n-modal v-model:show="showModal" preset="card" style="width: 50rem">
    <template #header>Выполнение генерации</template>
    <template #default>
      <n-form ref="formRef" :model="modelRef" :show-require-mark="false">
        <n-form-item
          v-for="(variable, index) in modelRef.variables"
          :key="index"
          :path="`variables[${index}].value`"
          :rule="valueRule"
        >
          <template #label>
            <n-flex align="baseline" :size="6">
              <n-text>{{ variable.title }}</n-text>
              <n-text depth="3" code style="font-size: 0.75rem">{{ variable.name }}</n-text>
            </n-flex>
          </template>
          <n-input-number
            v-if="variable.type === 'integer'"
            :value="asNumber(variable.value)"
            :precision="0"
            placeholder="Введите целое число"
            style="width: 100%"
            @update:value="(v) => (variable.value = v)"
          />
          <n-input-number
            v-else-if="variable.type === 'float'"
            :value="asNumber(variable.value)"
            placeholder="Введите число"
            style="width: 100%"
            @update:value="(v) => (variable.value = v)"
          />
          <n-input
            v-else
            :value="asString(variable.value)"
            placeholder="Введите строку"
            @update:value="(v) => (variable.value = v)"
          />
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
            Запустить генерацию
          </n-button>
        </n-form-item>
      </n-form>
    </template>
  </n-modal>
</template>

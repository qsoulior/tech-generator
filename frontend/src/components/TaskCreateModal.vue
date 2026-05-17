<script setup lang="ts">
import { NModal, NForm, NFormItem, NInput, NButton, useMessage } from "naive-ui"
import type { FormRules, FormInst } from "naive-ui"
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
  type: string
}

const formRef = ref<FormInst | null>(null)
const loading = ref(false)

interface ModelVariable {
  name: string
  type: string
  value: string
}

interface Model {
  variables: ModelVariable[]
}

const modelRef = ref<Model>({
  variables: [],
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
        payload: Object.fromEntries(modelRef.value.variables.map((v) => [v.name, v.value])),
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
      type: variable.type,
      value: "",
    }))
  }
})
</script>

<template>
  <n-modal v-model:show="showModal" preset="card" style="width: 50rem">
    <template #header>Выполнение генерации</template>
    <template #default>
      <n-form ref="formRef" :model="modelRef" :rules="rules">
        <n-form-item
          v-for="(variable, index) in modelRef.variables"
          :key="index"
          :path="`variables[${index}].value`"
          :label="variable.name"
        >
          <n-input v-model:value="variable.value" placeholder="Введите значение" />
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

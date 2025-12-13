<script setup lang="ts">
import { NModal, NForm, NFormItem, NInput, NButton, useMessage } from "naive-ui"
import type { FormRules, FormInst } from "naive-ui"
import { ref, watch } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()

const message = useMessage()

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
      await taskCreate()
    }
  })
}

async function taskCreate() {
  if (props.versionId == undefined) return

  const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/task/create`, {
    method: "POST",
    body: JSON.stringify({
      versionID: props.versionId,
      payload: Object.fromEntries(new Map(modelRef.value.variables.map((v) => [v.name, v.value]))),
    }),
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
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

  message.success("Запущена генерация документа")
  showModal.value = false
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
          <n-button style="width: 100%" secondary type="primary" @click="handleValidateClick">
            Запустить генерацию
          </n-button>
        </n-form-item>
      </n-form>
    </template>
  </n-modal>
</template>

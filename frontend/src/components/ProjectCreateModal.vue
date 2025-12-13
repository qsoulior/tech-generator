<script setup lang="ts">
import { NModal, NForm, NFormItem, NInput, NButton, useMessage } from "naive-ui"
import type { FormRules, FormInst } from "naive-ui"
import { ref } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()

const message = useMessage()

const showModal = defineModel("showModal", { default: false })

const emit = defineEmits<{
  submit: []
}>()

const formRef = ref<FormInst | null>(null)

interface Model {
  name: string
}

const modelRef = ref<Model>({
  name: "",
})

const rules: FormRules = {
  name: {
    required: true,
    message: "Название проекта не может быть пустым",
  },
}

function handleValidateClick(e: MouseEvent) {
  e.preventDefault()
  formRef.value?.validate(async (errors) => {
    if (!errors) {
      await projectCreate(modelRef.value)
    }
  })
}

async function projectCreate(model: Model) {
  const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/project/create`, {
    method: "POST",
    body: JSON.stringify({
      name: model.name,
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

  emit("submit")
  showModal.value = false
}
</script>

<template>
  <n-modal v-model:show="showModal" preset="card" style="width: 50rem">
    <template #header>Добавление проекта</template>
    <template #default>
      <n-form ref="formRef" :model="modelRef" :rules="rules">
        <n-form-item path="name" label="Название проекта">
          <n-input v-model:value="modelRef.name" placeholder="Введите название проекта" />
        </n-form-item>
        <n-form-item>
          <n-button style="width: 100%" secondary type="primary" @click="handleValidateClick">Добавить</n-button>
        </n-form-item>
      </n-form>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { NModal, NForm, NFormItem, NInput, NButton } from "naive-ui"
import type { FormRules, FormInst } from "naive-ui"
import { ref } from "vue"
import { projectCreate } from "@/api/project"
import { useApiCall } from "@/composables/useApiCall"

const apiCall = useApiCall()

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
      await submit(modelRef.value)
    }
  })
}

async function submit(model: Model) {
  const r = await apiCall(() => projectCreate({ name: model.name }))
  if (!r.ok) return

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

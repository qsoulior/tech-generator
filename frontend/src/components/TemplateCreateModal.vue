<script setup lang="ts">
import { NModal, NForm, NFormItem, NInput, NButton } from "naive-ui"
import type { FormRules, FormInst } from "naive-ui"
import { ref } from "vue"
import { templateCreate } from "@/api/template"
import { useApiCall } from "@/composables/useApiCall"

const apiCall = useApiCall()

const showModal = defineModel("showModal", { default: false })

const props = defineProps<{
  projectId: number
}>()

const emit = defineEmits<{
  submit: []
}>()

const formRef = ref<FormInst | null>(null)
const loading = ref(false)

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
</script>

<template>
  <n-modal v-model:show="showModal" preset="card" style="width: 50rem">
    <template #header>Добавление шаблона</template>
    <template #default>
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
    </template>
  </n-modal>
</template>

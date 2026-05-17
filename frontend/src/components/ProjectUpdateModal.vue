<script setup lang="ts">
import { NModal, NForm, NFormItem, NInput, NButton } from "naive-ui"
import type { FormRules, FormInst } from "naive-ui"
import { ref, watch } from "vue"
import { projectUpdate } from "@/api/project"
import { useApiCall } from "@/composables/useApiCall"

const apiCall = useApiCall()

const showModal = defineModel("showModal", { default: false })

const props = defineProps<{
  projectId: number
  initialName: string
}>()

const emit = defineEmits<{
  submit: [name: string]
}>()

const formRef = ref<FormInst | null>(null)
const loading = ref(false)

interface Model {
  name: string
}

const modelRef = ref<Model>({
  name: props.initialName,
})

watch(
  () => [showModal.value, props.initialName],
  ([show]) => {
    if (show) {
      modelRef.value.name = props.initialName
    }
  },
)

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
  loading.value = true
  try {
    const r = await apiCall(() => projectUpdate(props.projectId, { name: model.name }))
    if (!r.ok) return

    emit("submit", model.name)
    showModal.value = false
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <n-modal v-model:show="showModal" preset="card" style="width: 50rem">
    <template #header>Редактирование проекта</template>
    <template #default>
      <n-form ref="formRef" :model="modelRef" :rules="rules">
        <n-form-item path="name" label="Название проекта">
          <n-input v-model:value="modelRef.name" placeholder="Введите название проекта" />
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
            Сохранить
          </n-button>
        </n-form-item>
      </n-form>
    </template>
  </n-modal>
</template>

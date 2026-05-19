<script setup lang="ts">
import { NModal } from "naive-ui"
import VariableForm from "@/components/VariableForm.vue"
import { ref, watch } from "vue"

const showModal = defineModel("showModal", { default: false })

const props = defineProps<{
  occupiedSlugs: string[]
}>()

const emit = defineEmits<{
  submit: [variable: Variable]
}>()

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

function emptyModel(): Variable {
  return {
    name: "",
    title: "",
    type: "string",
    expression: "",
    inputType: "input",
    constraints: [],
  }
}

const model = ref<Variable>(emptyModel())

watch(showModal, (value) => {
  if (value) model.value = emptyModel()
})

function handleSubmit() {
  emit("submit", model.value)
  showModal.value = false
}
</script>

<template>
  <n-modal v-model:show="showModal" preset="card" style="width: 50rem">
    <template #header>Добавление переменной</template>
    <template #default>
      <VariableForm
        submit-text="Добавить переменную"
        :occupied-slugs="props.occupiedSlugs"
        v-model:model="model"
        @submit="handleSubmit"
      />
    </template>
  </n-modal>
</template>

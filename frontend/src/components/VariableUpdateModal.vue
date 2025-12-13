<script setup lang="ts">
import { NModal } from "naive-ui"
import VariableForm from "@/components/VariableForm.vue"

const showModal = defineModel("showModal", { default: false })
const variable = defineModel<Variable>("variable", {
  default: {
    name: "",
    type: "string",
    expression: "",
    inputType: "input",
    constraints: [],
  },
})

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
  type: string
  expression: string
  inputType: string
  constraints: Constraint[]
}

function handleSubmit() {
  emit("submit", variable.value)
  showModal.value = false
}
</script>

<template>
  <n-modal v-model:show="showModal" preset="card" style="width: 50rem">
    <template #header>Изменение переменной</template>
    <template #default>
      <VariableForm submit-text="Сохранить изменения" v-model:model="variable" @submit="handleSubmit" />
    </template>
  </n-modal>
</template>

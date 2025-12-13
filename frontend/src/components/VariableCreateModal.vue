<script setup lang="ts">
import { NModal } from "naive-ui"
import VariableForm from "@/components/VariableForm.vue"
import { ref } from "vue"

const showModal = defineModel("showModal", { default: false })

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

const model = ref<Variable>({
  name: "",
  type: "string",
  expression: "",
  inputType: "input",
  constraints: [],
})

function handleSubmit() {
  emit("submit", model.value)
  showModal.value = false
  model.value = {
    name: "",
    type: "string",
    expression: "",
    inputType: "input",
    constraints: [],
  }
}
</script>

<template>
  <n-modal v-model:show="showModal" preset="card" style="width: 50rem">
    <template #header>Добавление переменной</template>
    <template #default>
      <VariableForm submit-text="Добавить переменную" v-model:model="model" @submit="handleSubmit" />
    </template>
  </n-modal>
</template>

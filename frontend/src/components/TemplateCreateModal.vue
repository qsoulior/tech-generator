<script setup lang="ts">
import { NModal, NTabs, NTabPane } from "naive-ui"
import TemplateCreateEmptyForm from "@/components/TemplateCreateEmptyForm.vue"
import TemplateCreateImportForm from "@/components/TemplateCreateImportForm.vue"
import TemplateCreateLibraryForm from "@/components/TemplateCreateLibraryForm.vue"

const showModal = defineModel("showModal", { default: false })

defineProps<{
  projectId: number
}>()

const emit = defineEmits<{
  submit: []
}>()

function onSubmit() {
  emit("submit")
  showModal.value = false
}
</script>

<template>
  <n-modal v-model:show="showModal" preset="card" style="width: 50rem">
    <template #header>Добавление шаблона</template>
    <template #default>
      <n-tabs type="line" default-value="create">
        <n-tab-pane name="create" tab="Создать пустой">
          <TemplateCreateEmptyForm :project-id="projectId" @submit="onSubmit" />
        </n-tab-pane>
        <n-tab-pane name="library" tab="Из библиотеки">
          <TemplateCreateLibraryForm :project-id="projectId" @submit="onSubmit" />
        </n-tab-pane>
        <n-tab-pane name="import" tab="Импортировать">
          <TemplateCreateImportForm :project-id="projectId" @submit="onSubmit" />
        </n-tab-pane>
      </n-tabs>
    </template>
  </n-modal>
</template>

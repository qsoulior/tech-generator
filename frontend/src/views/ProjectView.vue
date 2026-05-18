<script setup lang="ts">
import { NLayout, NLayoutContent, NFlex, NPagination, NButton, NText } from "naive-ui"
import { onMounted, ref } from "vue"
import TemplateListItem from "@/components/TemplateListItem.vue"
import TemplateListSearch from "@/components/TemplateListSearch.vue"
import AppHeader from "@/components/AppHeader.vue"
import TemplateCreateModal from "@/components/TemplateCreateModal.vue"
import { templateList as fetchTemplates } from "@/api/template"
import { useApiCall } from "@/composables/useApiCall"
import { usePagination } from "@/composables/usePagination"

const props = defineProps<{
  projectID: number
}>()

const apiCall = useApiCall()

const { page, pageSize, totalPages, pageSizes } = usePagination("шаблонов")
const totalTemplates = ref(0)

const templateName = ref<string>("")

const showModal = ref(false)

interface Template {
  id: number
  name: string
  authorName: string
  createdAt: Date
  updatedAt: Date
}

const templates = ref<Template[]>([])

async function templateList() {
  const r = await apiCall(() =>
    fetchTemplates({
      projectID: props.projectID,
      page: page.value,
      size: pageSize.value,
      templateName: templateName.value || undefined,
    }),
  )
  if (!r.ok) return

  totalTemplates.value = r.value.totalTemplates
  totalPages.value = r.value.totalPages

  templates.value = r.value.templates.map((template) => ({
    id: template.id,
    name: template.name,
    authorName: template.authorName,
    createdAt: new Date(template.createdAt),
    updatedAt: new Date(template.updatedAt ?? template.createdAt),
  }))
}

onMounted(async () => {
  await templateList()
})

async function onUpdatePage() {
  await templateList()
}

async function onUpdatePageSize() {
  await templateList()
}

async function onSubmitModal() {
  await templateList()
}

async function onSubmitSearch() {
  await templateList()
}

async function onDeleteTemplate() {
  await templateList()
}

function onUpdateTemplate(id: number, name: string) {
  const template = templates.value.find((t) => t.id === id)
  if (template) template.name = name
}
</script>

<template>
  <n-layout>
    <AppHeader />
    <n-layout content-style="height: calc(100vh - 59px)">
      <n-layout-content content-class="layout-content" embedded style="height: 100%">
        <n-flex vertical align="center" style="max-width: 50rem; margin: auto">
          <TemplateListSearch v-model:value="templateName" @submit="onSubmitSearch" />
          <n-button secondary style="width: 100%" @click="showModal = true">Добавить шаблон</n-button>
          <TemplateCreateModal :project-id="projectID" v-model:show-modal="showModal" @submit="onSubmitModal" />
          <n-text depth="3" style="width: 100%">Всего: {{ totalTemplates }}</n-text>
          <TemplateListItem
            v-for="template in templates"
            :project-id="projectID"
            :template-id="template.id"
            :key="template.id"
            :name="template.name"
            :author-name="template.authorName"
            :created-at="template.createdAt"
            :updated-at="template.updatedAt"
            @delete="onDeleteTemplate"
            @update="onUpdateTemplate"
          >
          </TemplateListItem>
          <n-pagination
            v-model:page="page"
            v-model:page-size="pageSize"
            :page-count="totalPages"
            show-size-picker
            :page-sizes="pageSizes"
            @update:page="onUpdatePage"
            @update:page-size="onUpdatePageSize"
          />
        </n-flex>
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<style scoped>
:deep(.layout-content) {
  padding: 1.5rem;
}
</style>

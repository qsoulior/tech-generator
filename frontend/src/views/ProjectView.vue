<script setup lang="ts">
import {
  NLayout,
  NText,
  NLayoutContent,
  NLayoutHeader,
  NMenu,
  NFlex,
  NPagination,
  NButton,
  type MenuOption,
} from "naive-ui"
import { h, onMounted, ref } from "vue"
import TemplateListItem from "@/components/TemplateListItem.vue"
import TemplateListSearch from "@/components/TemplateListSearch.vue"
import { RouterLink } from "vue-router"
import TemplateCreateModal from "@/components/TemplateCreateModal.vue"
import { templateList as fetchTemplates } from "@/api/template"
import { useApiCall } from "@/composables/useApiCall"

const props = defineProps<{
  projectID: number
}>()

const apiCall = useApiCall()

const totalTemplates = ref(0)
const totalPages = ref(0)
const page = ref(1)
const pageSize = ref(50)

const templateName = ref<string>("")

const showModal = ref(false)

const pageSizes = [
  { label: "10 шаблонов", value: 10 },
  { label: "50 шаблонов", value: 50 },
  { label: "100 шаблонов", value: 100 },
  { label: "500 шаблонов", value: 500 },
]

interface Template {
  id: number
  name: string
  authorName: string
  createdAt: Date
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

const menuOptions: MenuOption[] = [
  {
    label: () =>
      h(
        RouterLink,
        {
          to: {
            name: "projectList",
          },
        },
        { default: () => "Проекты" },
      ),
    key: "projectList",
  },
]
</script>

<template>
  <n-layout>
    <n-layout-header bordered style="padding: 0.5rem 1rem">
      <n-flex align="center" justify="space-between">
        <n-text strong>tech-generator</n-text>
        <n-flex>
          <n-menu mode="horizontal" :options="menuOptions" />
        </n-flex>
      </n-flex>
    </n-layout-header>
    <n-layout content-style="height: calc(100vh - 59px)">
      <n-layout-content content-class="layout-content" embedded style="height: 100%">
        <n-flex vertical align="center" style="max-width: 50rem; margin: auto">
          <TemplateListSearch v-model:value="templateName" @submit="onSubmitSearch" />
          <n-button secondary style="width: 100%" @click="showModal = true">Добавить шаблон</n-button>
          <TemplateCreateModal :project-id="projectID" v-model:show-modal="showModal" @submit="onSubmitModal" />
          <TemplateListItem
            v-for="template in templates"
            :project-id="projectID"
            :template-id="template.id"
            :key="template.id"
            :name="template.name"
            :author-name="template.authorName"
            :created-at="template.createdAt"
            @delete="onDeleteTemplate"
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

<script setup lang="ts">
import { NFlex, NPagination, NButton, useMessage } from "naive-ui"
import { onMounted, ref } from "vue"
import TemplateListItem from "@/components/TemplateListItem.vue"
import TemplateListSearch from "@/components/TemplateListSearch.vue"
import { useRouter } from "vue-router"
import TemplateCreateModal from "@/components/TemplateCreateModal.vue"

const props = defineProps<{
  id: number
}>()

const router = useRouter()
const message = useMessage()

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

// запрос списка шаблонов
interface TemplateListResultTemplate {
  id: number
  name: string
  authorName: string
  createdAt: string
}

interface TemplateListResult {
  templates: TemplateListResultTemplate[]
  totalTemplates: number
  totalPages: number
}

async function templateList() {
  const params = new URLSearchParams({
    page: page.value.toString(),
    size: pageSize.value.toString(),
  })

  if (templateName.value !== "") {
    params.append("templateName", templateName.value)
  }

  const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/template/list/${props.id}?${params}`, {
    method: "GET",
    credentials: "include",
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

  const result: TemplateListResult = await response.json()

  totalTemplates.value = result.totalTemplates
  totalPages.value = result.totalPages

  templates.value = result.templates.map((template) => {
    return {
      id: template.id,
      name: template.name,
      authorName: template.authorName,
      createdAt: new Date(template.createdAt),
    }
  })
}

// триггеры для обновления списка шаблонов
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
</script>

<template>
  <n-flex vertical align="center" style="max-width: 50rem; margin: auto">
    <TemplateListSearch v-model:value="templateName" @submit="onSubmitSearch" />
    <n-button secondary style="width: 100%" @click="showModal = true">Добавить проект</n-button>
    <TemplateCreateModal :project-id="id" v-model:show-modal="showModal" @submit="onSubmitModal" />
    <TemplateListItem
      v-for="template in templates"
      :id="template.id"
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
</template>

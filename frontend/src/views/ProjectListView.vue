<script setup lang="ts">
import { NLayout, NLayoutContent, NFlex, NPagination, NButton, NText, NLayoutHeader, useMessage } from "naive-ui"
import { onMounted, ref } from "vue"
import ProjectListItem from "@/components/ProjectListItem.vue"
import ProjectListSearch from "@/components/ProjectListSearch.vue"
import { useRouter } from "vue-router"
import ProjectCreateModal from "@/components/ProjectCreateModal.vue"

const router = useRouter()
const message = useMessage()

const totalProjects = ref(0)
const totalPages = ref(0)
const page = ref(1)
const pageSize = ref(50)

const projectName = ref<string>("")

const showModal = ref(false)

const pageSizes = [
  { label: "10 проектов", value: 10 },
  { label: "50 проектов", value: 50 },
  { label: "100 проектов", value: 100 },
  { label: "500 проектов", value: 500 },
]

interface Project {
  id: number
  name: string
  authorName: string
}

const projects = ref<Project[]>([])

// запрос списка проектов
interface ProjectListResultProject {
  id: number
  name: string
  authorName: string
}

interface ProjectListResult {
  projects: ProjectListResultProject[]
  totalProjects: number
  totalPages: number
}

async function projectList() {
  const params = new URLSearchParams({
    page: page.value.toString(),
    size: pageSize.value.toString(),
  })

  if (projectName.value !== "") {
    params.append("projectName", projectName.value)
  }

  const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/project/list?${params}`, {
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

  const result: ProjectListResult = await response.json()

  totalProjects.value = result.totalProjects
  totalPages.value = result.totalPages

  projects.value = result.projects.map((project) => {
    return {
      id: project.id,
      name: project.name,
      authorName: project.authorName,
    }
  })
}

// триггеры для обновления списка проектов
onMounted(async () => {
  await projectList()
})

async function onUpdatePage() {
  await projectList()
}

async function onUpdatePageSize() {
  await projectList()
}

async function onSubmitModal() {
  await projectList()
}

async function onSubmitSearch() {
  await projectList()
}

async function onDeleteProject() {
  await projectList()
}
</script>

<template>
  <n-layout>
    <n-layout-header bordered style="padding: 0.5rem 1rem">
      <n-flex align="center" justify="start" style="padding: 10px 0">
        <n-text strong>tech-generator</n-text>
      </n-flex>
    </n-layout-header>
    <n-layout content-style="height: calc(100vh - 59px)">
      <n-layout-content content-class="layout-content" embedded style="height: 100%">
        <n-flex vertical align="center" style="max-width: 50rem; margin: auto">
          <ProjectListSearch v-model:value="projectName" @submit="onSubmitSearch" />
          <n-button secondary style="width: 100%" @click="showModal = true">Добавить проект</n-button>
          <ProjectCreateModal v-model:show-modal="showModal" @submit="onSubmitModal" />
          <ProjectListItem
            v-for="project in projects"
            :id="project.id"
            :key="project.id"
            :name="project.name"
            :author-name="project.authorName"
            @delete="onDeleteProject"
          >
          </ProjectListItem>
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

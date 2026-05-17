<script setup lang="ts">
import { NLayout, NLayoutContent, NFlex, NPagination, NButton, NLayoutHeader, NText } from "naive-ui"
import { onMounted, ref } from "vue"
import ProjectListItem from "@/components/ProjectListItem.vue"
import ProjectListSearch from "@/components/ProjectListSearch.vue"
import ProjectCreateModal from "@/components/ProjectCreateModal.vue"
import AppBrand from "@/components/AppBrand.vue"
import { projectList as fetchProjects, type ProjectListItem as Project } from "@/api/project"
import { useApiCall } from "@/composables/useApiCall"
import { usePagination } from "@/composables/usePagination"
import { useProjectStore } from "@/stores/project"

const apiCall = useApiCall()
const projectStore = useProjectStore()

const { page, pageSize, totalPages, pageSizes } = usePagination("проектов")
const totalProjects = ref(0)

const projectName = ref<string>("")

const showModal = ref(false)

const projects = ref<Project[]>([])

async function projectList() {
  const r = await apiCall(() =>
    fetchProjects({
      page: page.value,
      size: pageSize.value,
      projectName: projectName.value || undefined,
    }),
  )
  if (!r.ok) return

  totalProjects.value = r.value.totalProjects
  totalPages.value = r.value.totalPages
  projects.value = r.value.projects
  for (const project of r.value.projects) {
    projectStore.put(project.id, { name: project.name, authorName: project.authorName })
  }
}

onMounted(async () => {
  await projectList()
})

async function onDeleteProject(id: number) {
  projectStore.invalidate(id)
  await projectList()
}
</script>

<template>
  <n-layout>
    <n-layout-header bordered class="header">
      <n-flex align="center" justify="start" class="header-inner">
        <AppBrand />
      </n-flex>
    </n-layout-header>
    <n-layout content-style="height: calc(100vh - 59px)">
      <n-layout-content content-class="layout-content" embedded class="content">
        <n-flex vertical align="center" class="content-inner">
          <ProjectListSearch v-model:value="projectName" @submit="projectList" />
          <n-button secondary class="full-width" @click="showModal = true">Добавить проект</n-button>
          <ProjectCreateModal v-model:show-modal="showModal" @submit="projectList" />
          <n-text depth="3" class="full-width">Всего: {{ totalProjects }}</n-text>
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
            @update:page="projectList"
            @update:page-size="projectList"
          />
        </n-flex>
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<style scoped>
.header {
  padding: 0.5rem 1rem;
}

.header-inner {
  padding: 10px 0;
}

.content {
  height: 100%;
}

.content-inner {
  max-width: 50rem;
  margin: auto;
}

.full-width {
  width: 100%;
}

:deep(.layout-content) {
  padding: 1.5rem;
}
</style>

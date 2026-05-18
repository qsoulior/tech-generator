<script setup lang="ts">
import { NLayoutHeader, NFlex, NBreadcrumb, NBreadcrumbItem, NIcon } from "naive-ui"
import { computed, watch, type Component } from "vue"
import { useRoute, useRouter, type RouteLocationRaw } from "vue-router"
import AppBrand from "@/components/AppBrand.vue"
import AppUserMenu from "@/components/AppUserMenu.vue"
import IconAppstoreOutlined from "@/components/icons/IconAppstoreOutlined.vue"
import IconFolderOutlined from "@/components/icons/IconFolderOutlined.vue"
import IconFileTextOutlined from "@/components/icons/IconFileTextOutlined.vue"
import IconUnorderedListOutlined from "@/components/icons/IconUnorderedListOutlined.vue"
import IconFileDoneOutlined from "@/components/icons/IconFileDoneOutlined.vue"
import { useProjectStore } from "@/stores/project"
import { useTemplateStore } from "@/stores/template"
import { useApiCall } from "@/composables/useApiCall"

interface Crumb {
  key: string
  label: string
  icon: Component
  to?: RouteLocationRaw
}

const route = useRoute()
const router = useRouter()
const projectStore = useProjectStore()
const templateStore = useTemplateStore()
const apiCall = useApiCall()

const projectID = computed(() => {
  const v = route.params.projectID
  return typeof v === "string" ? Number(v) : undefined
})

const templateID = computed(() => {
  const v = route.params.templateID
  return typeof v === "string" ? Number(v) : undefined
})

const taskID = computed(() => {
  const v = route.params.taskID
  return typeof v === "string" ? Number(v) : undefined
})

const projectName = computed(() =>
  projectID.value != undefined ? (projectStore.get(projectID.value)?.name ?? "") : "",
)
const templateName = computed(() =>
  templateID.value != undefined ? (templateStore.getMeta(templateID.value)?.name ?? "") : "",
)

watch(
  projectID,
  async (id) => {
    if (id == undefined) return
    await apiCall(() => projectStore.ensureLoaded(id))
  },
  { immediate: true },
)

watch(
  templateID,
  async (id) => {
    if (id == undefined) return
    await apiCall(() => templateStore.ensureMetaLoaded(id))
  },
  { immediate: true },
)

const crumbs = computed<Crumb[]>(() => {
  const items: Crumb[] = []
  const isLeaf = (name: string) => route.name === name

  items.push({
    key: "projects",
    label: "Проекты",
    icon: IconAppstoreOutlined,
    to: isLeaf("projectList") ? undefined : { name: "projectList" },
  })

  if (projectID.value != undefined) {
    items.push({
      key: `project-${projectID.value}`,
      label: projectName.value || "…",
      icon: IconFolderOutlined,
      to: isLeaf("project") ? undefined : { name: "project", params: { projectID: projectID.value } },
    })
  }

  if (templateID.value != undefined) {
    items.push({
      key: `template-${templateID.value}`,
      label: templateName.value || "…",
      icon: IconFileTextOutlined,
      to: isLeaf("template")
        ? undefined
        : { name: "template", params: { projectID: projectID.value, templateID: templateID.value } },
    })

    if (route.name === "taskList" || route.name === "task") {
      items.push({
        key: "tasks",
        label: "Результаты",
        icon: IconUnorderedListOutlined,
        to: isLeaf("taskList")
          ? undefined
          : { name: "taskList", params: { projectID: projectID.value, templateID: templateID.value } },
      })
    }

    if (taskID.value != undefined) {
      items.push({
        key: `task-${taskID.value}`,
        label: `Результат #${taskID.value}`,
        icon: IconFileDoneOutlined,
      })
    }
  }

  return items
})

function resolveHref(to: RouteLocationRaw | undefined): string | undefined {
  if (to == undefined) return undefined
  return router.resolve(to).href
}

function onCrumbClick(event: MouseEvent, to: RouteLocationRaw | undefined) {
  if (to == undefined) return
  if (event.metaKey || event.ctrlKey || event.shiftKey || event.altKey || event.button !== 0) return
  event.preventDefault()
  router.push(to)
}
</script>

<template>
  <n-layout-header bordered class="header">
    <n-flex align="center" justify="space-between" :wrap="false">
      <n-flex align="center" :size="16" :wrap="false">
        <AppBrand />
        <n-breadcrumb>
          <n-breadcrumb-item
            v-for="crumb in crumbs"
            :key="crumb.key"
            :href="resolveHref(crumb.to)"
            @click="onCrumbClick($event, crumb.to)"
          >
            <span class="crumb-content">
              <n-icon :size="16">
                <component :is="crumb.icon" />
              </n-icon>
              {{ crumb.label }}
            </span>
          </n-breadcrumb-item>
        </n-breadcrumb>
      </n-flex>
      <AppUserMenu />
    </n-flex>
  </n-layout-header>
</template>

<style scoped>
.header {
  padding: 0.5rem 1rem;
}

.crumb-content {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}
</style>

import { defineStore } from "pinia"
import { ref } from "vue"
import type { ProjectListItem } from "@/api/project"

export const useProjectStore = defineStore("project", () => {
  const cache = ref(new Map<number, ProjectListItem>())

  function put(project: ProjectListItem) {
    cache.value.set(project.id, project)
  }

  function get(projectID: number): ProjectListItem | undefined {
    return cache.value.get(projectID)
  }

  function invalidate(projectID: number) {
    cache.value.delete(projectID)
  }

  return { put, get, invalidate }
})

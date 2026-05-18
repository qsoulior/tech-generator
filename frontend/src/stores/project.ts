import { defineStore } from "pinia"
import { ref } from "vue"
import { projectGet, type ProjectGetResult } from "@/api/project"

export const useProjectStore = defineStore("project", () => {
  const cache = ref(new Map<number, ProjectGetResult>())

  function put(projectID: number, project: ProjectGetResult) {
    cache.value.set(projectID, project)
  }

  function get(projectID: number): ProjectGetResult | undefined {
    return cache.value.get(projectID)
  }

  async function ensureLoaded(projectID: number): Promise<ProjectGetResult> {
    const existing = cache.value.get(projectID)
    if (existing != null) return existing

    const fetched = await projectGet(projectID)
    cache.value.set(projectID, fetched)
    return fetched
  }

  function invalidate(projectID: number) {
    cache.value.delete(projectID)
  }

  return { put, get, ensureLoaded, invalidate }
})

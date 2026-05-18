import { defineStore } from "pinia"
import { ref } from "vue"
import { templateGet, type TemplateGetResult } from "@/api/template"

export const useTemplateStore = defineStore("template", () => {
  const cache = ref(new Map<number, TemplateGetResult>())

  function get(templateID: number): TemplateGetResult | undefined {
    return cache.value.get(templateID)
  }

  async function ensureLoaded(templateID: number): Promise<TemplateGetResult> {
    const existing = cache.value.get(templateID)
    if (existing != null) return existing

    const fetched = await templateGet(templateID)
    cache.value.set(templateID, fetched)
    return fetched
  }

  function invalidate(templateID: number) {
    cache.value.delete(templateID)
  }

  return { get, ensureLoaded, invalidate }
})

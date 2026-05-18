import { defineStore } from "pinia"
import { ref } from "vue"
import {
  templateGet,
  templateGetMeta,
  type TemplateGetMetaResult,
  type TemplateGetResult,
} from "@/api/template"

export const useTemplateStore = defineStore("template", () => {
  const cache = ref(new Map<number, TemplateGetResult>())
  const metaCache = ref(new Map<number, TemplateGetMetaResult>())

  function get(templateID: number): TemplateGetResult | undefined {
    return cache.value.get(templateID)
  }

  function getMeta(templateID: number): TemplateGetMetaResult | undefined {
    const full = cache.value.get(templateID)
    if (full != null) return { name: full.name }
    return metaCache.value.get(templateID)
  }

  async function ensureLoaded(templateID: number): Promise<TemplateGetResult> {
    const existing = cache.value.get(templateID)
    if (existing != null) return existing

    const fetched = await templateGet(templateID)
    cache.value.set(templateID, fetched)
    return fetched
  }

  async function ensureMetaLoaded(templateID: number): Promise<TemplateGetMetaResult> {
    const full = cache.value.get(templateID)
    if (full != null) return { name: full.name }

    const meta = metaCache.value.get(templateID)
    if (meta != null) return meta

    const fetched = await templateGetMeta(templateID)
    metaCache.value.set(templateID, fetched)
    return fetched
  }

  function invalidate(templateID: number) {
    cache.value.delete(templateID)
  }

  function setMeta(templateID: number, meta: TemplateGetMetaResult) {
    metaCache.value.set(templateID, meta)
  }

  return { get, getMeta, ensureLoaded, ensureMetaLoaded, invalidate, setMeta }
})

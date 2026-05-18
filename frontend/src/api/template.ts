import { apiDelete, apiGet, apiPost } from "./client"
import type { components } from "./schema.gen"

export type TemplateListResult = components["schemas"]["TemplateListResponse"]
export type TemplateListItem = TemplateListResult["templates"][number]
export type TemplateGetResult = components["schemas"]["TemplateGetByIDResponse"]
export type TemplateGetMetaResult = components["schemas"]["TemplateGetMetaByIDResponse"]
export type TemplateGetVersion = components["schemas"]["TemplateGetByIDVersion"]
export type TemplateGetVariable = TemplateGetVersion["variables"][number]
export type TemplateGetConstraint = TemplateGetVariable["constraints"][number]
export type TemplateCreateInput = components["schemas"]["TemplateCreateRequest"]
export type TemplateUpdateInput = components["schemas"]["TemplateUpdateRequest"]
export type TemplateImportInput = components["schemas"]["TemplateImportRequest"]
export type TemplateImportPayload = components["schemas"]["TemplateImportPayload"]
export type TemplateImportVersion = components["schemas"]["TemplateImportVersion"]
export type TemplateImportVariable = TemplateImportVersion["variables"][number]
export type TemplateImportConstraint = TemplateImportVariable["constraints"][number]
export type TemplateImportResult = components["schemas"]["TemplateImportResponse"]
export type TemplateDefaultListResult = components["schemas"]["TemplateDefaultListResponse"]
export type TemplateDefaultListItem = TemplateDefaultListResult["templates"][number]
export type TemplateCreateFromDefaultInput = components["schemas"]["TemplateCreateFromDefaultRequest"]
export type TemplateCreateFromDefaultResult = components["schemas"]["TemplateCreateFromDefaultResponse"]

export interface TemplateListParams {
  projectID: number
  page: number
  size: number
  templateName?: string
}

export function templateList(params: TemplateListParams): Promise<TemplateListResult> {
  const search = new URLSearchParams({
    page: params.page.toString(),
    size: params.size.toString(),
  })
  if (params.templateName) {
    search.append("templateName", params.templateName)
  }
  return apiGet<TemplateListResult>(`/template/list/${params.projectID}?${search}`)
}

export function templateGet(templateID: number): Promise<TemplateGetResult> {
  return apiGet<TemplateGetResult>(`/template/get/${templateID}`)
}

export function templateGetMeta(templateID: number): Promise<TemplateGetMetaResult> {
  return apiGet<TemplateGetMetaResult>(`/template/get_meta/${templateID}`)
}

export function templateCreate(input: TemplateCreateInput): Promise<void> {
  return apiPost<void>(`/template/create`, input)
}

export function templateUpdate(templateID: number, input: TemplateUpdateInput): Promise<void> {
  return apiPost<void>(`/template/update/${templateID}`, input)
}

export function templateDelete(templateID: number): Promise<void> {
  return apiDelete<void>(`/template/delete/${templateID}`)
}

export function templateImport(input: TemplateImportInput): Promise<TemplateImportResult> {
  return apiPost<TemplateImportResult>(`/template/import`, input)
}

export interface TemplateDefaultListParams {
  page: number
  size: number
  templateName?: string
}

export function templateDefaultList(params: TemplateDefaultListParams): Promise<TemplateDefaultListResult> {
  const search = new URLSearchParams({
    page: params.page.toString(),
    size: params.size.toString(),
  })
  if (params.templateName) {
    search.append("templateName", params.templateName)
  }
  return apiGet<TemplateDefaultListResult>(`/template/default/list?${search}`)
}

export function templateCreateFromDefault(
  input: TemplateCreateFromDefaultInput,
): Promise<TemplateCreateFromDefaultResult> {
  return apiPost<TemplateCreateFromDefaultResult>(`/template/create_from_default`, input)
}

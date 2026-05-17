import { apiDelete, apiGet, apiPost } from "./client"

export interface TemplateListItem {
  id: number
  name: string
  authorName: string
  createdAt: string
}

export interface TemplateListResult {
  templates: TemplateListItem[]
  totalTemplates: number
  totalPages: number
}

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

export interface TemplateGetConstraint {
  name: string
  expression: string
  isActive: boolean
}

export interface TemplateGetVariable {
  name: string
  type: string
  isInput: boolean
  expression: string
  constraints: TemplateGetConstraint[]
}

export interface TemplateGetVersion {
  id: number
  number: number
  data: string
  createdAt: string
  variables: TemplateGetVariable[]
}

export interface TemplateGetResult {
  name: string
  version?: TemplateGetVersion
}

export function templateGet(templateID: number): Promise<TemplateGetResult> {
  return apiGet<TemplateGetResult>(`/template/get/${templateID}`)
}

export interface TemplateCreateInput {
  name: string
  projectID: number
}

export function templateCreate(input: TemplateCreateInput): Promise<void> {
  return apiPost<void>(`/template/create`, input)
}

export function templateDelete(templateID: number): Promise<void> {
  return apiDelete<void>(`/template/delete/${templateID}`)
}

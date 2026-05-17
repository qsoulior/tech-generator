import { apiDelete, apiGet, apiPost } from "./client"

export interface ProjectListItem {
  id: number
  name: string
  authorName: string
}

export interface ProjectListResult {
  projects: ProjectListItem[]
  totalProjects: number
  totalPages: number
}

export interface ProjectListParams {
  page: number
  size: number
  projectName?: string
}

export function projectList(params: ProjectListParams): Promise<ProjectListResult> {
  const search = new URLSearchParams({
    page: params.page.toString(),
    size: params.size.toString(),
  })
  if (params.projectName) {
    search.append("projectName", params.projectName)
  }
  return apiGet<ProjectListResult>(`/project/list?${search}`)
}

export interface ProjectCreateInput {
  name: string
}

export function projectCreate(input: ProjectCreateInput): Promise<void> {
  return apiPost<void>(`/project/create`, input)
}

export function projectDelete(id: number): Promise<void> {
  return apiDelete<void>(`/project/delete/${id}`)
}

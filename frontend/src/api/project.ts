import { apiDelete, apiGet, apiPost } from "./client"
import type { components } from "./schema.gen"

export type ProjectListResult = components["schemas"]["ProjectListResponse"]
export type ProjectListItem = ProjectListResult["projects"][number]
export type ProjectGetResult = components["schemas"]["ProjectGetByIDResponse"]
export type ProjectCreateInput = components["schemas"]["ProjectCreateRequest"]
export type ProjectUpdateInput = components["schemas"]["ProjectUpdateRequest"]
export type ProjectUsersResult = components["schemas"]["ProjectUsersResponse"]
export type ProjectUsersItem = ProjectUsersResult["users"][number]
export type ProjectUserRole = ProjectUsersItem["role"]
export type ProjectUpdateUsersInput = components["schemas"]["ProjectUpdateUsersRequest"]

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

export function projectGet(projectID: number): Promise<ProjectGetResult> {
  return apiGet<ProjectGetResult>(`/project/get/${projectID}`)
}

export function projectCreate(input: ProjectCreateInput): Promise<void> {
  return apiPost<void>(`/project/create`, input)
}

export function projectUpdate(id: number, input: ProjectUpdateInput): Promise<void> {
  return apiPost<void>(`/project/update/${id}`, input)
}

export function projectDelete(id: number): Promise<void> {
  return apiDelete<void>(`/project/delete/${id}`)
}

export function projectUsers(projectID: number): Promise<ProjectUsersResult> {
  return apiGet<ProjectUsersResult>(`/project/users/${projectID}`)
}

export function projectUpdateUsers(projectID: number, input: ProjectUpdateUsersInput): Promise<void> {
  return apiPost<void>(`/project/update_users/${projectID}`, input)
}

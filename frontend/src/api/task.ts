import { apiGet, apiPost } from "./client"
import type { components } from "./schema.gen"

export type TaskListResult = components["schemas"]["TaskListResponse"]
export type TaskListItem = TaskListResult["tasks"][number]
export type TaskStatus = components["schemas"]["TaskStatus"]
export type TaskGetResult = components["schemas"]["TaskGetByIDResponse"]
export type TaskGetTask = TaskGetResult["task"]
export type TaskGetError = NonNullable<TaskGetTask["error"]>
export type TaskGetVariableError = NonNullable<TaskGetError["variableErrors"]>[number]
export type TaskGetConstraintError = NonNullable<TaskGetVariableError["constraintErrors"]>[number]
export type TaskCreateInput = components["schemas"]["TaskCreateRequest"]

export interface TaskListParams {
  templateID: number
  page: number
  size: number
}

export function taskList(params: TaskListParams): Promise<TaskListResult> {
  const search = new URLSearchParams({
    page: params.page.toString(),
    size: params.size.toString(),
    templateID: params.templateID.toString(),
  })
  return apiGet<TaskListResult>(`/task/list?${search}`)
}

export function taskGet(taskID: number): Promise<TaskGetResult> {
  return apiGet<TaskGetResult>(`/task/get/${taskID}`)
}

export function taskCreate(input: TaskCreateInput): Promise<void> {
  return apiPost<void>(`/task/create`, input)
}

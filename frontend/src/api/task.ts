import { apiGet, apiPost } from "./client"

export interface TaskListItem {
  id: number
  status: string
  versionNumber: number
  creatorName: string
  createdAt: string
  updatedAt: string
}

export interface TaskListResult {
  tasks: TaskListItem[]
  totalTasks: number
  totalPages: number
}

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

export interface TaskGetConstraintError {
  id: string
  name: string
  message: string
}

export interface TaskGetVariableError {
  id: number
  name: string
  message?: string
  constraintErrors: TaskGetConstraintError[]
}

export interface TaskGetError {
  message?: string
  variableErrors: TaskGetVariableError[]
}

export interface TaskGetTask {
  id: number
  versionID: number
  status: string
  payload: Record<string, string>
  error: TaskGetError
  creatorName: string
  createdAt: string
  updatedAt: string
}

export interface TaskGetResult {
  task: TaskGetTask
  result: string | null
}

export function taskGet(taskID: number): Promise<TaskGetResult> {
  return apiGet<TaskGetResult>(`/task/get/${taskID}`)
}

export interface TaskCreateInput {
  versionID: number
  payload: Record<string, string>
}

export function taskCreate(input: TaskCreateInput): Promise<void> {
  return apiPost<void>(`/task/create`, input)
}

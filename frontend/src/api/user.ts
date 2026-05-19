import { apiDelete, apiGet, apiPost } from "./client"
import type { components } from "./schema.gen"

export type UserCreateInput = components["schemas"]["UserCreateRequest"]
export type UserTokenCreateInput = components["schemas"]["UserTokenCreateRequest"]
export type UserGetResult = components["schemas"]["UserGetByIDResponse"]
export type UserListResult = components["schemas"]["UserListResponse"]
export type UserListItem = UserListResult["users"][number]

export interface UserListParams {
  page: number
  size: number
  userName?: string
}

export function userCreate(input: UserCreateInput): Promise<void> {
  return apiPost<void>(`/user/create`, input)
}

export function userTokenCreate(input: UserTokenCreateInput): Promise<void> {
  return apiPost<void>(`/user/token/create`, input)
}

export function userTokenDelete(): Promise<void> {
  return apiDelete(`/user/token/delete`)
}

export function userGet(): Promise<UserGetResult> {
  return apiGet<UserGetResult>(`/user/get`)
}

export function userList(params: UserListParams): Promise<UserListResult> {
  const search = new URLSearchParams({
    page: params.page.toString(),
    size: params.size.toString(),
  })
  if (params.userName) {
    search.append("userName", params.userName)
  }
  return apiGet<UserListResult>(`/user/list?${search}`)
}

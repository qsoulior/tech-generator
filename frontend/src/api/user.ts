import { apiPost } from "./client"
import type { components } from "./schema.gen"

export type UserCreateInput = components["schemas"]["UserCreateRequest"]
export type UserTokenCreateInput = components["schemas"]["UserTokenCreateRequest"]

export function userCreate(input: UserCreateInput): Promise<void> {
  return apiPost<void>(`/user/create`, input)
}

export function userTokenCreate(input: UserTokenCreateInput): Promise<void> {
  return apiPost<void>(`/user/token/create`, input)
}

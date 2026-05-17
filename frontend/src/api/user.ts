import { apiPost } from "./client"

export interface UserCreateInput {
  name: string
  email: string
  password: string
}

export function userCreate(input: UserCreateInput): Promise<void> {
  return apiPost<void>(`/user/create`, input)
}

export interface UserTokenCreateInput {
  name: string
  password: string
}

export function userTokenCreate(input: UserTokenCreateInput): Promise<void> {
  return apiPost<void>(`/user/token/create`, input)
}

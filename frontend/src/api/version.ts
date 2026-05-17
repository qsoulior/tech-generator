import { apiPost } from "./client"

export interface VersionCreateConstraint {
  name: string
  expression: string
  isActive: boolean
}

export interface VersionCreateVariable {
  name: string
  type: string
  expression: string
  isInput: boolean
  constraints: VersionCreateConstraint[]
}

export interface VersionCreateInput {
  templateID: number
  data: string
  variables: VersionCreateVariable[]
}

export interface VersionCreateResult {
  id: number
}

export function versionCreate(input: VersionCreateInput): Promise<VersionCreateResult> {
  return apiPost<VersionCreateResult>(`/version/create`, input)
}

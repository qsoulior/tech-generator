import { apiPost } from "./client"
import type { components } from "./schema.gen"

export type VersionCreateInput = components["schemas"]["VersionCreateRequest"]
export type VersionCreateVariable = VersionCreateInput["variables"][number]
export type VersionCreateConstraint = VersionCreateVariable["constraints"][number]
export type VersionCreateResult = components["schemas"]["VersionCreateResponse"]

export function versionCreate(input: VersionCreateInput): Promise<VersionCreateResult> {
  return apiPost<VersionCreateResult>(`/version/create`, input)
}

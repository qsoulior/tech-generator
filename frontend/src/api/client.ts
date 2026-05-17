export class ApiError extends Error {
  status: number

  constructor(status: number, message: string) {
    super(message)
    this.name = "ApiError"
    this.status = status
  }
}

export class UnauthorizedApiError extends ApiError {
  constructor(status: number, message: string) {
    super(status, message)
    this.name = "UnauthorizedApiError"
  }
}

async function readErrorMessage(response: Response): Promise<string> {
  try {
    const body = await response.json()
    if (body && typeof body.message === "string") return body.message
  } catch {
    // тело не JSON
  }
  return `HTTP ${response.status}`
}

async function request<T>(path: string, init: RequestInit): Promise<T> {
  const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}${path}`, {
    credentials: "include",
    ...init,
  })

  if (!response.ok) {
    const message = await readErrorMessage(response)
    if (response.status === 401 || response.status === 403) {
      throw new UnauthorizedApiError(response.status, message)
    }
    throw new ApiError(response.status, message)
  }

  if (response.status === 204) return undefined as T
  const text = await response.text()
  return (text ? JSON.parse(text) : undefined) as T
}

export function apiGet<T>(path: string): Promise<T> {
  return request<T>(path, { method: "GET" })
}

export function apiPost<T>(path: string, body?: unknown): Promise<T> {
  const hasBody = body !== undefined
  return request<T>(path, {
    method: "POST",
    body: hasBody ? JSON.stringify(body) : undefined,
    headers: hasBody ? { "Content-Type": "application/json" } : undefined,
  })
}

export function apiDelete<T = void>(path: string): Promise<T> {
  return request<T>(path, { method: "DELETE" })
}

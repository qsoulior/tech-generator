import { useMessage } from "naive-ui"
import { useRouter } from "vue-router"
import { ApiError, UnauthorizedApiError } from "@/api/client"

export type ApiCallResult<T> = { ok: true; value: T } | { ok: false }

/**
 * Оборачивает вызов API и берёт на себя обработку ошибок:
 *  - UnauthorizedApiError → редирект на /auth (если ещё не там);
 *  - любая другая ApiError → message.error со server-side сообщением;
 *  - всё прочее → message.error("Неизвестная ошибка").
 *
 * View не пишет try/catch — получает discriminated result и решает, что делать
 * дальше по успеху.
 */
export function useApiCall() {
  const message = useMessage()
  const router = useRouter()

  return async function call<T>(fn: () => Promise<T>): Promise<ApiCallResult<T>> {
    try {
      return { ok: true, value: await fn() }
    } catch (e) {
      if (e instanceof UnauthorizedApiError) {
        if (router.currentRoute.value.name !== "auth") {
          router.push({ name: "auth" })
          return { ok: false }
        }
        message.error(e.message)
        return { ok: false }
      }
      if (e instanceof ApiError) {
        message.error(e.message)
        return { ok: false }
      }
      message.error("Неизвестная ошибка")
      return { ok: false }
    }
  }
}

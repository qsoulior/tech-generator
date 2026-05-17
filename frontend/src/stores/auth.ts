import { defineStore } from "pinia"
import { ref } from "vue"
import { userGet, type UserGetResult } from "@/api/user"

export const useAuthStore = defineStore("auth", () => {
  const user = ref<UserGetResult | null>(null)

  async function ensureLoaded(): Promise<UserGetResult> {
    if (user.value != null) return user.value
    const fetched = await userGet()
    user.value = fetched
    return fetched
  }

  function clear() {
    user.value = null
  }

  return { user, ensureLoaded, clear }
})

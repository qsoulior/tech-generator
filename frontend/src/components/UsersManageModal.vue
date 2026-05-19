<script setup lang="ts">
import {
  NModal,
  NInput,
  NInputGroup,
  NSelect,
  NButton,
  NIcon,
  NFlex,
  NText,
  NEmpty,
  NSpin,
  NCard,
  type SelectOption,
} from "naive-ui"
import { computed, ref, watch } from "vue"
import IconDeleteOutlined from "@/components/icons/IconDeleteOutlined.vue"
import IconSearchOutlined from "@/components/icons/IconSearchOutlined.vue"
import { projectUsers, projectUpdateUsers, type ProjectUserRole } from "@/api/project"
import { templateUsers, templateUpdateUsers, type TemplateUserRole } from "@/api/template"
import { userList, type UserListItem } from "@/api/user"
import { useApiCall } from "@/composables/useApiCall"

type Kind = "project" | "template"

const props = defineProps<{
  kind: Kind
  entityId: number
}>()

const emit = defineEmits<{
  submit: []
}>()

const showModal = defineModel<boolean>("show", { default: false })

const apiCall = useApiCall()

const PAGE_SIZE = 20

interface ManagedUser {
  id: number
  name: string
  email: string
  role: ProjectUserRole | TemplateUserRole
}

const projectRoleOptions: SelectOption[] = [
  { label: "Чтение", value: "read" },
  { label: "Запись", value: "write" },
  { label: "Управление", value: "maintain" },
]

const templateRoleOptions: SelectOption[] = [
  { label: "Чтение", value: "read" },
  { label: "Запись", value: "write" },
]

const roleOptions = computed<SelectOption[]>(() =>
  props.kind === "project" ? projectRoleOptions : templateRoleOptions,
)

const defaultRole = "read" as const

const loadingCurrent = ref(false)
const saving = ref(false)
const currentUsers = ref<ManagedUser[]>([])

const searchInput = ref("")
const searchQuery = ref("")
const searchResults = ref<UserListItem[]>([])
const searchPage = ref(1)
const searchHasMore = ref(false)
const searchLoading = ref(false)
const searchInitialLoaded = ref(false)
let searchToken = 0

const currentIDs = computed(() => new Set(currentUsers.value.map((u) => u.id)))

const availableResults = computed(() => searchResults.value.filter((u) => !currentIDs.value.has(u.id)))

async function loadCurrentUsers() {
  loadingCurrent.value = true
  try {
    if (props.kind === "project") {
      const r = await apiCall(() => projectUsers(props.entityId))
      if (!r.ok) return
      currentUsers.value = r.value.users.map((u) => ({
        id: u.id,
        name: u.name,
        email: u.email,
        role: u.role,
      }))
    } else {
      const r = await apiCall(() => templateUsers(props.entityId))
      if (!r.ok) return
      currentUsers.value = r.value.users.map((u) => ({
        id: u.id,
        name: u.name,
        email: u.email,
        role: u.role,
      }))
    }
  } finally {
    loadingCurrent.value = false
  }
}

async function fetchUsers(page: number, query: string, token: number) {
  searchLoading.value = true
  try {
    const r = await apiCall(() =>
      userList({
        page,
        size: PAGE_SIZE,
        userName: query || undefined,
      }),
    )
    if (token !== searchToken) return
    if (!r.ok) {
      searchHasMore.value = false
      return
    }
    if (page === 1) {
      searchResults.value = r.value.users
    } else {
      searchResults.value = [...searchResults.value, ...r.value.users]
    }
    searchPage.value = page
    searchHasMore.value = page < r.value.totalPages
    searchInitialLoaded.value = true
  } finally {
    if (token === searchToken) {
      searchLoading.value = false
    }
  }
}

function resetAndSearch() {
  searchToken++
  searchQuery.value = searchInput.value.trim()
  searchResults.value = []
  searchPage.value = 1
  searchHasMore.value = false
  searchInitialLoaded.value = false
  void fetchUsers(1, searchQuery.value, searchToken)
}

function loadMore() {
  if (searchLoading.value || !searchHasMore.value) return
  void fetchUsers(searchPage.value + 1, searchQuery.value, searchToken)
}

function addUser(user: UserListItem) {
  if (currentIDs.value.has(user.id)) return
  currentUsers.value.push({
    id: user.id,
    name: user.name,
    email: user.email,
    role: defaultRole,
  })
}

function removeUser(id: number) {
  currentUsers.value = currentUsers.value.filter((u) => u.id !== id)
}

async function save() {
  saving.value = true
  try {
    if (props.kind === "project") {
      const users = currentUsers.value.map((u) => ({
        id: u.id,
        role: u.role as ProjectUserRole,
      }))
      const r = await apiCall(() => projectUpdateUsers(props.entityId, { users }))
      if (!r.ok) return
    } else {
      const users = currentUsers.value.map((u) => ({
        id: u.id,
        role: u.role as TemplateUserRole,
      }))
      const r = await apiCall(() => templateUpdateUsers(props.entityId, { users }))
      if (!r.ok) return
    }

    emit("submit")
    showModal.value = false
  } finally {
    saving.value = false
  }
}

const scrollSentinel = ref<HTMLElement | null>(null)
let observer: IntersectionObserver | undefined

function attachObserver(el: HTMLElement | null) {
  observer?.disconnect()
  observer = undefined
  if (el == null) return
  observer = new IntersectionObserver(
    (entries) => {
      if (entries.some((e) => e.isIntersecting)) loadMore()
    },
    { threshold: 0.5 },
  )
  observer.observe(el)
}

watch(scrollSentinel, attachObserver)

watch(showModal, async (value) => {
  if (!value) {
    observer?.disconnect()
    observer = undefined
    return
  }
  searchInput.value = ""
  await loadCurrentUsers()
  resetAndSearch()
})

const title = computed(() => (props.kind === "project" ? "Доступ к проекту" : "Доступ к шаблону"))
</script>

<template>
  <n-modal v-model:show="showModal" preset="card" style="width: 50rem" :auto-focus="false">
    <template #header>{{ title }}</template>
    <template #default>
      <n-flex vertical :size="16">
        <n-flex vertical :size="8">
          <n-text strong>Пользователи с доступом</n-text>
          <n-spin :show="loadingCurrent">
            <n-empty v-if="!loadingCurrent && currentUsers.length === 0" description="Пока никого нет" />
            <n-flex v-else vertical :size="8">
              <n-card v-for="user in currentUsers" :key="user.id" size="small">
                <n-flex justify="space-between" align="center" :wrap="false" :size="8">
                  <n-flex vertical align="start" :size="2" style="min-width: 0; flex: 1">
                    <n-text strong>{{ user.name }}</n-text>
                    <n-text depth="3" style="font-size: 0.85rem">{{ user.email }}</n-text>
                  </n-flex>
                  <n-select
                    v-model:value="user.role"
                    :options="roleOptions"
                    size="small"
                    style="width: 10rem"
                  />
                  <n-button
                    secondary
                    size="small"
                    aria-label="Убрать пользователя"
                    title="Убрать пользователя"
                    @click="removeUser(user.id)"
                  >
                    <template #icon>
                      <n-icon>
                        <IconDeleteOutlined />
                      </n-icon>
                    </template>
                  </n-button>
                </n-flex>
              </n-card>
            </n-flex>
          </n-spin>
        </n-flex>

        <n-flex vertical :size="8">
          <n-text strong>Добавить пользователя</n-text>
          <n-input-group>
            <n-input
              v-model:value="searchInput"
              placeholder="Поиск по имени"
              clearable
              @keyup.enter="resetAndSearch"
            />
            <n-button secondary aria-label="Найти" @click="resetAndSearch">
              <template #icon>
                <n-icon>
                  <IconSearchOutlined />
                </n-icon>
              </template>
            </n-button>
          </n-input-group>
          <div class="scroll-area">
            <n-empty
              v-if="searchInitialLoaded && availableResults.length === 0 && !searchLoading"
              description="Ничего не найдено"
              style="margin: 1rem 0"
            />
            <n-flex v-else vertical :size="8">
              <n-card v-for="user in availableResults" :key="user.id" size="small">
                <n-flex justify="space-between" align="center" :wrap="false" :size="8">
                  <n-flex vertical align="start" :size="2" style="min-width: 0; flex: 1">
                    <n-text strong>{{ user.name }}</n-text>
                    <n-text depth="3" style="font-size: 0.85rem">{{ user.email }}</n-text>
                  </n-flex>
                  <n-button size="small" secondary type="primary" @click="addUser(user)">Добавить</n-button>
                </n-flex>
              </n-card>
              <div ref="scrollSentinel" style="height: 1px" />
              <n-flex v-if="searchLoading" justify="center" style="padding: 0.5rem">
                <n-spin size="small" />
              </n-flex>
            </n-flex>
          </div>
        </n-flex>

        <n-button secondary type="primary" :loading="saving" :disabled="saving || loadingCurrent" @click="save">
          Сохранить
        </n-button>
      </n-flex>
    </template>
  </n-modal>
</template>

<style scoped>
.scroll-area {
  max-height: 20rem;
  overflow-y: auto;
}
</style>

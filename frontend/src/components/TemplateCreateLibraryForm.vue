<script setup lang="ts">
import {
  NInput,
  NInputGroup,
  NButton,
  NIcon,
  NCard,
  NFlex,
  NText,
  NPagination,
  NSkeleton,
  NEmpty,
  NForm,
  NFormItem,
  NTooltip,
  useMessage,
} from "naive-ui"
import { onMounted, ref } from "vue"
import { templateDefaultList, templateCreateFromDefault, type TemplateDefaultListItem } from "@/api/template"
import { useApiCall } from "@/composables/useApiCall"
import { usePagination } from "@/composables/usePagination"
import { formatRelativeTime } from "@/utils/relativeTime"
import IconSearchOutlined from "@/components/icons/IconSearchOutlined.vue"

const apiCall = useApiCall()
const message = useMessage()

const props = defineProps<{
  projectId: number
}>()

const emit = defineEmits<{
  submit: []
}>()

const templates = ref<TemplateDefaultListItem[]>([])
const total = ref(0)
const searchInput = ref("")
const search = ref("")
const loading = ref(true)

const selected = ref<TemplateDefaultListItem | null>(null)
const copyName = ref("")
const copyLoading = ref(false)

const { page, pageSize, totalPages, pageSizes } = usePagination("шаблонов", 10)

async function load() {
  loading.value = true
  try {
    const r = await apiCall(() =>
      templateDefaultList({
        page: page.value,
        size: pageSize.value,
        templateName: search.value || undefined,
      }),
    )
    if (!r.ok) return

    templates.value = r.value.templates
    total.value = r.value.totalTemplates
    totalPages.value = r.value.totalPages
  } finally {
    loading.value = false
  }
}

onMounted(load)

function onSearchSubmit() {
  search.value = searchInput.value
  page.value = 1
  void load()
}

async function onPageUpdate() {
  await load()
}

async function onPageSizeUpdate() {
  page.value = 1
  await load()
}

function onSelect(template: TemplateDefaultListItem) {
  selected.value = template
  copyName.value = template.name
}

function onBack() {
  selected.value = null
  copyName.value = ""
}

async function onCreate() {
  if (!selected.value) return
  if (copyName.value.trim() === "") {
    message.error("Название шаблона не может быть пустым")
    return
  }

  copyLoading.value = true
  try {
    const r = await apiCall(() =>
      templateCreateFromDefault({
        sourceTemplateID: selected.value!.id,
        projectID: props.projectId,
        name: copyName.value.trim(),
      }),
    )
    if (!r.ok) return

    message.success("Шаблон создан на основе библиотечного")
    emit("submit")
  } finally {
    copyLoading.value = false
  }
}
</script>

<template>
  <template v-if="selected == null">
    <n-flex vertical>
      <n-input-group>
        <n-input v-model:value="searchInput" placeholder="Название шаблона" @keyup.enter="onSearchSubmit" />
        <n-button @click="onSearchSubmit">
          <template #icon>
            <n-icon>
              <IconSearchOutlined />
            </n-icon>
          </template>
        </n-button>
      </n-input-group>
      <n-text depth="3">Всего: {{ total }}</n-text>
      <n-flex v-if="loading" vertical size="small">
        <n-card v-for="i in 3" :key="i">
          <n-flex vertical size="small">
            <n-skeleton text :repeat="1" style="width: 50%" />
            <n-skeleton text :repeat="1" style="width: 30%" />
          </n-flex>
        </n-card>
      </n-flex>
      <n-flex v-else vertical size="small">
        <n-empty v-if="templates.length === 0" description="Шаблоны не найдены" />
        <n-card v-for="template in templates" :key="template.id" style="cursor: pointer" @click="onSelect(template)">
          <n-flex vertical size="small">
            <n-text strong>{{ template.name }}</n-text>
            <n-text depth="3">
              Добавлен:
              <n-tooltip trigger="hover">
                <template #trigger>
                  <span class="time">{{ formatRelativeTime(new Date(template.createdAt)) }}</span>
                </template>
                {{ new Date(template.createdAt).toLocaleString() }}
              </n-tooltip>
              · Обновлён:
              <n-tooltip trigger="hover">
                <template #trigger>
                  <span class="time">{{
                    formatRelativeTime(new Date(template.updatedAt ?? template.createdAt))
                  }}</span>
                </template>
                {{ new Date(template.updatedAt ?? template.createdAt).toLocaleString() }}
              </n-tooltip>
            </n-text>
          </n-flex>
        </n-card>
      </n-flex>
      <n-flex justify="center" style="min-height: 34px">
        <n-skeleton v-if="loading" width="280px" height="34px" :sharp="false" />
        <n-pagination
          v-else
          v-model:page="page"
          v-model:page-size="pageSize"
          :page-count="totalPages"
          show-size-picker
          :page-sizes="pageSizes"
          @update:page="onPageUpdate"
          @update:page-size="onPageSizeUpdate"
        />
      </n-flex>
    </n-flex>
  </template>
  <template v-else>
    <n-flex vertical>
      <n-text
        >На базе шаблона <n-text strong>«{{ selected.name }}»</n-text></n-text
      >
      <n-form>
        <n-form-item label="Название нового шаблона">
          <n-input v-model:value="copyName" placeholder="Введите название шаблона" />
        </n-form-item>
      </n-form>
      <n-flex justify="space-between">
        <n-button :disabled="copyLoading" @click="onBack">Назад</n-button>
        <n-button type="primary" secondary :loading="copyLoading" :disabled="copyLoading" @click="onCreate">
          Добавить
        </n-button>
      </n-flex>
    </n-flex>
  </template>
</template>

<style scoped>
.time {
  border-bottom: 1px dashed currentColor;
  cursor: help;
}
</style>

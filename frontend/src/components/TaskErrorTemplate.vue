<script setup lang="ts">
import { NAlert, NFlex, NText } from "naive-ui"
import { computed } from "vue"
import type { TaskGetTemplateError } from "@/api/task"

const props = defineProps<{
  message: string
  templateError: TaskGetTemplateError
}>()

const location = computed(() => {
  const { line, column } = props.templateError
  return column != null && column > 0 ? `Строка ${line}, столбец ${column}` : `Строка ${line}`
})

const lineLabel = computed(() => `${props.templateError.line} │ `)
const caretIndent = computed(() => {
  const column = props.templateError.column
  if (column == null || column <= 0) return null
  return " ".repeat(lineLabel.value.length + column - 1) + "^"
})

const hasSnippet = computed(() => {
  const s = props.templateError.snippet
  return s != null && s !== ""
})
const hasDetail = computed(() => {
  const d = props.templateError.detail
  return d != null && d !== ""
})
</script>

<template>
  <n-alert type="error" style="width: 100%">
    <n-flex vertical :size="8">
      <n-text strong>{{ message }}</n-text>
      <n-text depth="3">{{ location }}</n-text>
      <pre v-if="hasSnippet" class="snippet">{{ lineLabel + templateError.snippet }}<template v-if="caretIndent">
{{ caretIndent }}</template></pre>
      <n-text v-if="hasDetail" depth="3" class="detail">{{ templateError.detail }}</n-text>
    </n-flex>
  </n-alert>
</template>

<style scoped>
.snippet {
  margin: 0;
  padding: 0.5rem 0.75rem;
  background: rgba(0, 0, 0, 0.04);
  border-radius: 4px;
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, "Liberation Mono", monospace;
  font-size: 0.85rem;
  white-space: pre;
  overflow-x: auto;
}

.detail {
  white-space: pre-wrap;
  word-break: break-word;
  font-size: 0.85rem;
}
</style>

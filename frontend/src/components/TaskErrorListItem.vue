<script setup lang="ts">
import { NCard, NCollapse, NCollapseItem, NFlex, NText, NTag, NTable, NAlert } from "naive-ui"
import { computed } from "vue"
import type { TaskGetVariableError } from "@/api/task"

const props = defineProps<{
  variableError: TaskGetVariableError
}>()

const constraintErrorsCount = computed(() => props.variableError.constraintErrors?.length ?? 0)
const hasMessage = computed(() => props.variableError.message != null && props.variableError.message !== "")
</script>

<template>
  <n-card class="item" size="small">
    <n-collapse :trigger-areas="['main', 'arrow', 'extra']">
      <n-collapse-item :name="String(variableError.id)">
        <template #header>
          <n-flex align="baseline" :size="6">
            <n-text strong>{{ variableError.title }}</n-text>
            <n-text depth="3" code style="font-size: 0.75rem">{{ variableError.name }}</n-text>
          </n-flex>
        </template>
        <template v-if="constraintErrorsCount > 0" #header-extra>
          <n-tag size="small" type="error">Ограничений: {{ constraintErrorsCount }}</n-tag>
        </template>
        <n-flex vertical size="small">
          <n-alert v-if="hasMessage" type="error">
            <div class="message">{{ variableError.message }}</div>
          </n-alert>
          <n-table v-if="constraintErrorsCount > 0" :single-line="false" size="small">
            <thead>
              <tr>
                <th>Ограничение</th>
                <th>Ошибка</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="constraintError in variableError.constraintErrors" :key="constraintError.id">
                <td>{{ constraintError.name }}</td>
                <td>{{ constraintError.message }}</td>
              </tr>
            </tbody>
          </n-table>
          <n-text v-if="!hasMessage && constraintErrorsCount === 0" depth="3">Подробности недоступны</n-text>
        </n-flex>
      </n-collapse-item>
    </n-collapse>
  </n-card>
</template>

<style scoped>
.item {
  width: 100%;
}

.message {
  white-space: pre-wrap;
  word-break: break-word;
}
</style>

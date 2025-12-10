<script setup lang="ts">
import { NInput, NInputGroup, NButton, NIcon } from "naive-ui"
import IconSearchOutlined from "@/components/icons/IconSearchOutlined.vue"

const search = defineModel<string>("value", { default: "" })
const emit = defineEmits<{
  submit: []
}>()

function submit() {
  emit("submit")
}
</script>

<template>
  <n-input-group>
    <n-input v-model:value="search" placeholder="Название шаблона" @keyup.enter="submit" />
    <n-button @click="submit">
      <template #icon>
        <n-icon>
          <IconSearchOutlined />
        </n-icon>
      </template>
    </n-button>
  </n-input-group>
</template>

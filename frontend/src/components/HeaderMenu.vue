<script setup lang="ts">
import { NMenu, type MenuOption } from "naive-ui"
import { h, computed } from "vue"
import { RouterLink, type RouteLocationRaw } from "vue-router"

export interface HeaderMenuItem {
  key: string
  label: string
  to: RouteLocationRaw
}

const props = defineProps<{
  items: HeaderMenuItem[]
}>()

const options = computed<MenuOption[]>(() =>
  props.items.map((item) => ({
    key: item.key,
    label: () => h(RouterLink, { to: item.to }, { default: () => item.label }),
  })),
)
</script>

<template>
  <n-menu mode="horizontal" :options="options" />
</template>

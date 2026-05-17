import { ref } from "vue"

export function usePagination(entityLabel: string, defaultPageSize = 50) {
  const page = ref(1)
  const pageSize = ref(defaultPageSize)
  const totalPages = ref(0)

  const pageSizes = [10, 50, 100, 500].map((value) => ({
    label: `${value} ${entityLabel}`,
    value,
  }))

  return { page, pageSize, totalPages, pageSizes }
}

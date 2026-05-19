<script setup lang="ts">
import {
  NForm,
  NFormItem,
  NInput,
  NInputGroup,
  NButton,
  NSelect,
  NTabs,
  NTabPane,
  NFlex,
  NCheckbox,
  NCard,
  NIcon,
  NEmpty,
  NText,
  NTooltip,
  useMessage,
} from "naive-ui"
import type { FormRules, FormInst, FormItemInst, SelectOption, FormItemRule } from "naive-ui"
import { computed, ref, watch } from "vue"
import IconDeleteOutlined from "@/components/icons/IconDeleteOutlined.vue"
import IconAddOutlined from "@/components/icons/IconAddOutlined.vue"
import IconCopyOutlined from "@/components/icons/IconCopyOutlined.vue"
import IconSyncOutlined from "@/components/icons/IconSyncOutlined.vue"
import { SLUG_PATTERN, SLUG_MAX_LEN, slugify } from "@/utils/slug"

const props = defineProps<{
  submitText: string
  occupiedSlugs: string[]
}>()

const emit = defineEmits<{
  submit: []
}>()

const message = useMessage()
const formRef = ref<FormInst | null>(null)
const slugFormItemRef = ref<FormItemInst | null>(null)

const typeOptions: SelectOption[] = [
  { label: "Строка", value: "string" },
  { label: "Целое число", value: "integer" },
  { label: "Вещественное число", value: "float" },
]

const isInputOptions: SelectOption[] = [
  { label: "Входная", value: "input" },
  { label: "Вычисляемая", value: "computed" },
]

interface ModelConstraint {
  name: string
  expression: string
  isActive: boolean
}

interface Model {
  name: string
  title: string
  type: string
  expression: string
  inputType: string
  constraints: ModelConstraint[]
}

const modelRef = defineModel<Model>("model", { required: true })

const placeholderSnippet = computed(() => `{{ .${modelRef.value.name} }}`)

const autoSlug = computed(() => slugify(modelRef.value.title, props.occupiedSlugs))
const isSlugAuto = computed(() => modelRef.value.name === autoSlug.value)

watch(
  () => modelRef.value.title,
  (_, oldTitle) => {
    const previousAuto = slugify(oldTitle, props.occupiedSlugs)
    if (modelRef.value.name === "" || modelRef.value.name === previousAuto) {
      modelRef.value.name = autoSlug.value
    }
  },
)

function regenerateSlug() {
  modelRef.value.name = autoSlug.value
  slugFormItemRef.value?.restoreValidation()
}

async function copySlug() {
  try {
    await navigator.clipboard.writeText(modelRef.value.name)
    message.success("Идентификатор скопирован")
  } catch {
    message.error("Не удалось скопировать")
  }
}

const rules: FormRules = {
  title: {
    required: true,
    message: "Название не может быть пустым",
  },
  name: {
    required: true,
    validator: (_rule: FormItemRule, value: string) => {
      if (value === "") return new Error("Идентификатор не может быть пустым")
      if (value.length > SLUG_MAX_LEN) return new Error(`Идентификатор длиннее ${SLUG_MAX_LEN} символов`)
      if (!SLUG_PATTERN.test(value)) return new Error("Только латиница, цифры и _ (не с цифры)")
      if (props.occupiedSlugs.includes(value)) return new Error("Идентификатор уже используется")
      return true
    },
    trigger: ["input", "blur"],
  },
  expression: {
    validator: (_rule: FormItemRule, value: string) => modelRef.value.inputType == "input" || value != "",
    message: "Выражение не может быть пустым",
  },
}

const rulesConstraintName = {
  required: true,
  message: "Название не может быть пустым",
}

const rulesConstraintExpression = {
  required: true,
  message: "Выражение не может быть пустым",
}

function handleValidateClick(e: MouseEvent) {
  e.preventDefault()
  formRef.value?.validate(async (errors) => {
    if (!errors) {
      emit("submit")
    }
  })
}

function handleAddClick(e: MouseEvent, index: number) {
  e.preventDefault()
  modelRef.value.constraints.splice(index + 1, 0, {
    name: "",
    expression: "",
    isActive: true,
  })
}

function handleDeleteClick(e: MouseEvent, index: number) {
  e.preventDefault()
  modelRef.value.constraints.splice(index, 1)
}
</script>

<template>
  <n-form ref="formRef" :model="modelRef" :rules="rules">
    <n-tabs type="line" size="large">
      <n-tab-pane name="Переменная" display-directive="show">
        <n-form-item path="title" label="Название">
          <n-input v-model:value="modelRef.title" placeholder="Введите название переменной" />
        </n-form-item>
        <n-form-item ref="slugFormItemRef" path="name" label="Идентификатор">
          <n-flex vertical :size="6" style="width: 100%">
            <n-input-group>
              <n-input v-model:value="modelRef.name" placeholder="slug" />
              <n-tooltip>
                <template #trigger>
                  <n-button @click="copySlug">
                    <template #icon>
                      <n-icon><IconCopyOutlined /></n-icon>
                    </template>
                  </n-button>
                </template>
                Скопировать
              </n-tooltip>
              <n-tooltip>
                <template #trigger>
                  <n-button :disabled="isSlugAuto" @click="regenerateSlug">
                    <template #icon>
                      <n-icon><IconSyncOutlined /></n-icon>
                    </template>
                  </n-button>
                </template>
                Сгенерировать из названия
              </n-tooltip>
            </n-input-group>
            <n-text depth="3" style="font-size: 0.75rem">
              В шаблоне используйте плейсхолдер <code>{{ placeholderSnippet }}</code>
            </n-text>
          </n-flex>
        </n-form-item>
        <n-form-item path="type" label="Тип значения">
          <n-select v-model:value="modelRef.type" :options="typeOptions" />
        </n-form-item>
        <n-form-item path="inputType" label="Тип переменной">
          <n-select v-model:value="modelRef.inputType" :options="isInputOptions" />
        </n-form-item>
        <n-form-item
          v-if="modelRef.inputType == 'computed'"
          :required="modelRef.inputType == 'computed'"
          path="expression"
          label="Выражение"
        >
          <n-input v-model:value="modelRef.expression" placeholder="Введите выражение" />
        </n-form-item>
        <n-form-item>
          <n-button style="width: 100%" secondary type="primary" @click="handleValidateClick">
            {{ submitText }}
          </n-button>
        </n-form-item>
      </n-tab-pane>
      <n-tab-pane name="Ограничения" display-directive="show">
        <n-empty v-if="modelRef.constraints.length == 0" description="Нет ограничений" style="margin: 2rem 0">
          <template #extra>
            <n-button size="small" @click="handleAddClick($event, 0)">Создать новое ограничение</n-button>
          </template>
        </n-empty>
        <n-flex v-else vertical>
          <n-card size="small" v-for="(constraint, index) in modelRef.constraints" :key="index">
            <n-form-item label="Название" :path="`constraints[${index}].name`" :rule="rulesConstraintName">
              <n-input v-model:value="constraint.name" placeholder="Введите название" />
            </n-form-item>
            <n-form-item label="Выражение" :path="`constraints[${index}].expression`" :rule="rulesConstraintExpression">
              <n-input v-model:value="constraint.expression" placeholder="Введите Выражение" />
            </n-form-item>
            <n-form-item path="isActive" :show-label="false" :show-feedback="false">
              <n-flex align="center" justify="space-between" style="width: 100%">
                <n-checkbox v-model:checked="constraint.isActive">Активно</n-checkbox>
                <n-flex>
                  <n-button secondary @click="handleAddClick($event, index)">
                    <template #icon>
                      <n-icon>
                        <IconAddOutlined />
                      </n-icon>
                    </template>
                  </n-button>
                  <n-button secondary @click="handleDeleteClick($event, index)">
                    <template #icon>
                      <n-icon>
                        <IconDeleteOutlined />
                      </n-icon>
                    </template>
                  </n-button>
                </n-flex>
              </n-flex>
            </n-form-item>
          </n-card>
        </n-flex>
        <n-form-item>
          <n-button style="width: 100%" secondary type="primary" @click="handleValidateClick">
            {{ submitText }}
          </n-button>
        </n-form-item>
      </n-tab-pane>
    </n-tabs>
  </n-form>
</template>

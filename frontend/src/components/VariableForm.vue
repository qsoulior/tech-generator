<script setup lang="ts">
import {
  NForm,
  NFormItem,
  NInput,
  NButton,
  NSelect,
  NTabs,
  NTabPane,
  NFlex,
  NCheckbox,
  NCard,
  NIcon,
  NEmpty,
} from "naive-ui"
import type { FormRules, FormInst, SelectOption, FormItemRule } from "naive-ui"
import { ref } from "vue"
import IconDeleteOutlined from "@/components/icons/IconDeleteOutlined.vue"
import IconAddOutlined from "@/components/icons/IconAddOutlined.vue"

defineProps<{
  submitText: string
}>()

const emit = defineEmits<{
  submit: []
}>()

const formRef = ref<FormInst | null>(null)

const typeOptions: SelectOption[] = [
  {
    label: "Строка",
    value: "string",
  },
  {
    label: "Целое число",
    value: "integer",
  },
  {
    label: "Вещественное число",
    value: "float",
  },
]

const isInputOptions: SelectOption[] = [
  {
    label: "Входная",
    value: "input",
  },
  {
    label: "Вычисляемая",
    value: "computed",
  },
]

interface ModelConstraint {
  name: string
  expression: string
  isActive: boolean
}

interface Model {
  name: string
  type: string
  expression: string
  inputType: string
  constraints: ModelConstraint[]
}

const modelRef = defineModel<Model>("model", { required: true })

const rules: FormRules = {
  name: {
    required: true,
    message: "Название не может быть пустым",
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
        <n-form-item path="name" label="Название">
          <n-input v-model:value="modelRef.name" placeholder="Введите название переменной" />
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

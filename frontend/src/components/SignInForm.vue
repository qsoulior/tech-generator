<script setup lang="ts">
import { useMessage, NForm, NFormItem, NInput, NCheckbox, NButton } from "naive-ui"
import type { FormInst, FormRules } from "naive-ui"
import { ref } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()
const message = useMessage()

const formRef = ref<FormInst | null>(null)

interface Model {
  name: string
  password: string
  remember: boolean
}

const modelRef = ref<Model>({
  name: "",
  password: "",
  remember: false,
})

const rules: FormRules = {
  name: {
    required: true,
    message: "Имя пользователя не может быть пустым",
    trigger: "blur",
  },
  password: {
    required: true,
    message: "Пароль не может быть пустым",
    trigger: "blur",
  },
}

function handleValidateClick(e: MouseEvent) {
  e.preventDefault()
  formRef.value?.validate(async (errors) => {
    if (!errors) {
      await userTokenCreate(modelRef.value)
    }
  })
}

async function userTokenCreate(model: Model) {
  const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/user/token/create`, {
    method: "POST",
    body: JSON.stringify({
      name: model.name,
      password: model.password,
    }),
    headers: {
      "Content-Type": "application/json",
    },
  })

  if (response.ok) {
    router.push({ name: "projectList" })
    message.success("Вы успешно вошли")
    return
  }

  const result = await response.json()
  message.error(result.message)
}
</script>

<template>
  <n-form ref="formRef" :model="modelRef" :rules="rules">
    <n-form-item path="name" label="Имя пользователя">
      <n-input v-model:value="modelRef.name" placeholder="Введите имя пользователя" />
    </n-form-item>
    <n-form-item path="password" label="Пароль">
      <n-input
        v-model:value="modelRef.password"
        type="password"
        placeholder="Введите пароль"
        show-password-on="click"
      />
    </n-form-item>
    <n-form-item path="remember" :show-label="false" :show-feedback="false">
      <n-checkbox v-model:checked="modelRef.remember"> Запомнить меня </n-checkbox>
    </n-form-item>
    <n-form-item>
      <n-button type="primary" @click="handleValidateClick">Войти</n-button>
    </n-form-item>
  </n-form>
</template>

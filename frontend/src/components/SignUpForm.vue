<script setup lang="ts">
import { useMessage, NForm, NFormItem, NInput, NButton } from "naive-ui"
import type { FormInst, FormRules, FormItemRule } from "naive-ui"
import { ref } from "vue"
import { useRouter } from "vue-router"

const router = useRouter()
const message = useMessage()

const formRef = ref<FormInst | null>(null)

interface Model {
  name: string
  email: string
  password: string
  passwordConfirm: string
}

const modelRef = ref<Model>({
  name: "",
  email: "",
  password: "",
  passwordConfirm: "",
})

const rules: FormRules = {
  name: {
    required: true,
    message: "Имя пользователя не может быть пустым",
    trigger: "blur",
  },
  email: {
    required: true,
    message: "Email не может быть пустым",
    trigger: "blur",
  },
  password: {
    required: true,
    message: "Пароль не может быть пустым",
    trigger: "blur",
  },
  passwordConfirm: [
    {
      required: true,
      message: "Пароль не может быть пустым",
      trigger: "blur",
    },
    {
      validator: (_rule: FormItemRule, value: string) => value === modelRef.value.password,
      message: "Пароли должны совпадать",
      trigger: "blur",
    },
  ],
}

function handleValidateClick(e: MouseEvent) {
  e.preventDefault()
  formRef.value?.validate(async (errors) => {
    if (!errors) {
      await userCreate(modelRef.value)
    }
  })
}

async function userCreate(model: Model) {
  const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/user/create`, {
    method: "POST",
    body: JSON.stringify({
      name: model.name,
      email: model.email,
      password: model.password,
    }),
    headers: {
      "Content-Type": "application/json",
    },
  })

  if (response.ok) {
    return userTokenCreate(model.name, model.password)
  }

  const result = await response.json()
  message.error(result.message)
}

async function userTokenCreate(name: string, password: string) {
  const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/user/token/create`, {
    method: "POST",
    body: JSON.stringify({
      name: name,
      password: password,
    }),
    headers: {
      "Content-Type": "application/json",
    },
  })

  if (response.ok) {
    router.push({ name: "projectList" })
    message.success("Вы успешно зарегистрировались")
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
    <n-form-item path="email" label="Email">
      <n-input v-model:value="modelRef.email" placeholder="Введите электронный адрес" />
    </n-form-item>
    <n-form-item path="password" label="Пароль">
      <n-input
        v-model:value="modelRef.password"
        type="password"
        placeholder="Введите пароль"
        show-password-on="click"
      />
    </n-form-item>
    <n-form-item path="passwordConfirm" label="Повторите пароль">
      <n-input
        v-model:value="modelRef.passwordConfirm"
        type="password"
        placeholder="Введите пароль ещё раз"
        show-password-on="click"
      />
    </n-form-item>
    <n-form-item>
      <n-button type="primary" @click="handleValidateClick">Зарегистрироваться</n-button>
    </n-form-item>
  </n-form>
</template>

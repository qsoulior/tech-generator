<script setup lang="ts">
import { useMessage, NLayout, NFlex, NCard, NForm, NFormItem, NInput, NCheckbox, NButton } from "naive-ui"
import type { FormInst, FormRules } from "naive-ui"
import { ref } from "vue"

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
      await userTokenCreate(modelRef.value.name, modelRef.value.password)
    }
  })
}

async function userTokenCreate(name: string, password: string) {
  try {
    const response = await fetch("http://127.0.0.1:3000/user/token/create", {
      method: "POST",
      body: JSON.stringify({
        name: name,
        password: password,
      }),
      headers: {
        "Content-Type": "application/json",
      },
    })

    message.info(response.status.toString())
  } catch (error) {
    console.log(error)
  }
}
</script>

<template>
  <n-layout embedded>
    <n-flex style="min-height: 100vh" vertical align="center" justify="center">
      <n-card style="width: 30rem">
        <template #header>Вход</template>
        <template #default>
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
      </n-card>
    </n-flex>
  </n-layout>
</template>

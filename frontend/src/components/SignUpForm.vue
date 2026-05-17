<script setup lang="ts">
import { useMessage, NForm, NFormItem, NInput, NButton } from "naive-ui"
import type { FormInst, FormRules, FormItemRule } from "naive-ui"
import { ref } from "vue"
import { useRouter } from "vue-router"
import { userCreate, userTokenCreate } from "@/api/user"
import { useApiCall } from "@/composables/useApiCall"
import { useAuthStore } from "@/stores/auth"

const router = useRouter()
const message = useMessage()
const apiCall = useApiCall()
const authStore = useAuthStore()

const formRef = ref<FormInst | null>(null)
const loading = ref(false)

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
      await signUp(modelRef.value)
    }
  })
}

async function signUp(model: Model) {
  loading.value = true
  try {
    const created = await apiCall(() =>
      userCreate({ name: model.name, email: model.email, password: model.password }),
    )
    if (!created.ok) return

    const signedIn = await apiCall(() => userTokenCreate({ name: model.name, password: model.password }))
    if (!signedIn.ok) return

    authStore.clear()
    router.push({ name: "projectList" })
    message.success("Вы успешно зарегистрировались")
  } finally {
    loading.value = false
  }
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
      <n-button secondary type="primary" :loading="loading" :disabled="loading" @click="handleValidateClick">
        Зарегистрироваться
      </n-button>
    </n-form-item>
  </n-form>
</template>

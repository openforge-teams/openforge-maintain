<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <CloudServerOutlined style="font-size: 40px; color: #1890ff" />
        <h1>{{ t('login.title') }}</h1>
      </div>
      <a-form :model="formState" :rules="rules" ref="formRef" @finish="handleLogin" layout="vertical">
        <a-form-item :label="t('login.username')" name="username">
          <a-input v-model:value="formState.username" size="large" :placeholder="t('login.usernameRequired')">
            <template #prefix><UserOutlined /></template>
          </a-input>
        </a-form-item>
        <a-form-item :label="t('login.password')" name="password">
          <a-input-password v-model:value="formState.password" size="large" :placeholder="t('login.passwordRequired')">
            <template #prefix><LockOutlined /></template>
          </a-input-password>
        </a-form-item>
        <a-form-item :label="t('login.totp')" name="totp_code">
          <a-input v-model:value="formState.totp_code" size="large" :placeholder="t('login.totpPlaceholder')">
            <template #prefix><SafetyOutlined /></template>
          </a-input>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" size="large" block :loading="loading">
            {{ t('login.loginBtn') }}
          </a-button>
        </a-form-item>
      </a-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { message } from 'ant-design-vue'
import { useUserStore } from '@/store/user'
import { CloudServerOutlined, UserOutlined, LockOutlined, SafetyOutlined } from '@ant-design/icons-vue'
import type { FormInstance, Rule } from 'ant-design-vue/es/form'

const router = useRouter()
const { t } = useI18n()
const userStore = useUserStore()
const formRef = ref<FormInstance>()
const loading = ref(false)

const formState = reactive({ username: '', password: '', totp_code: '' })

const rules: Record<string, Rule[]> = {
  username: [{ required: true, message: () => t('login.usernameRequired') }],
  password: [{ required: true, message: () => t('login.passwordRequired') }],
}

async function handleLogin() {
  loading.value = true
  try {
    await userStore.login(formState.username, formState.password, formState.totp_code || undefined)
    message.success(t('common.success'))
    router.push('/dashboard')
  } catch { /* handled by interceptor */ } finally { loading.value = false }
}
</script>

<style scoped>
.login-container {
  height: 100vh; display: flex; align-items: center; justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}
.login-card {
  width: 400px; padding: 40px; background: #fff; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,.15);
}
.login-header { text-align: center; margin-bottom: 32px; }
.login-header h1 { margin-top: 12px; font-size: 24px; color: rgba(0,0,0,.85); }
</style>

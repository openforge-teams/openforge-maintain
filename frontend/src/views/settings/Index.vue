<template>
  <div>
    <PageHeader :title="t('settings.title')" />
    <a-tabs v-model:activeKey="activeTab">
      <a-tab-pane key="profile" :tab="t('settings.profile')">
        <a-card style="max-width: 600px">
          <a-form :model="profileForm" layout="vertical" @finish="handleSaveProfile">
            <a-form-item :label="t('settings.username')"><a-input v-model:value="profileForm.username" /></a-form-item>
            <a-form-item :label="t('settings.email')"><a-input v-model:value="profileForm.email" /></a-form-item>
            <a-form-item><a-button type="primary" html-type="submit" :loading="saving">{{ t('common.save') }}</a-button></a-form-item>
          </a-form>
        </a-card>
      </a-tab-pane>
      <a-tab-pane key="password" :tab="t('settings.changePassword')">
        <a-card style="max-width: 600px">
          <a-form :model="pwdForm" layout="vertical" @finish="handleChangePassword">
            <a-form-item :label="t('settings.oldPassword')"><a-input-password v-model:value="pwdForm.old_password" /></a-form-item>
            <a-form-item :label="t('settings.newPassword')"><a-input-password v-model:value="pwdForm.new_password" /></a-form-item>
            <a-form-item :label="t('settings.confirmPassword')"><a-input-password v-model:value="pwdForm.confirm_password" /></a-form-item>
            <a-form-item><a-button type="primary" html-type="submit" :loading="changing">{{ t('common.save') }}</a-button></a-form-item>
          </a-form>
        </a-card>
      </a-tab-pane>
      <a-tab-pane key="language" :tab="t('settings.languageSetting')">
        <a-card style="max-width: 400px">
          <a-radio-group v-model:value="appStore.language" @change="onLangChange">
            <a-radio value="zh-CN">中文</a-radio>
            <a-radio value="en-US">English</a-radio>
          </a-radio-group>
        </a-card>
      </a-tab-pane>
    </a-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { message } from 'ant-design-vue'
import { useUserStore } from '@/store/user'
import { updateProfile } from '@/api/user'
import { changePassword } from '@/api/auth'
import { useAppStore } from '@/store/app'
import PageHeader from '@/components/PageHeader.vue'

const { t, locale } = useI18n()
const userStore = useUserStore()
const appStore = useAppStore()
const activeTab = ref('profile')
const saving = ref(false)
const changing = ref(false)

const profileForm = reactive({ username: '', email: '' })
const pwdForm = reactive({ old_password: '', new_password: '', confirm_password: '' })

onMounted(() => {
  if (userStore.userInfo) { profileForm.username = userStore.userInfo.username; profileForm.email = userStore.userInfo.email || '' }
})

async function handleSaveProfile() {
  saving.value = true
  try { await updateProfile(profileForm); message.success(t('common.success')); userStore.getUserInfo() } catch {} finally { saving.value = false }
}

async function handleChangePassword() {
  if (pwdForm.new_password !== pwdForm.confirm_password) { message.error('Passwords do not match'); return }
  changing.value = true
  try { await changePassword(pwdForm); message.success(t('settings.passwordChanged')); pwdForm.old_password = ''; pwdForm.new_password = ''; pwdForm.confirm_password = '' } catch {} finally { changing.value = false }
}

function onLangChange() { locale.value = appStore.language }
</script>
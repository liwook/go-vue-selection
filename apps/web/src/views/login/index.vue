<template>
  <div class="login">
    <el-form ref="loginForms" :model="loginForm" :rules="rules" class="login-form">
      <h1 class="title">登录</h1>
      <el-form-item prop="username">
        <el-input v-model="loginForm.username" :prefix-icon="User" type="text" placeholder="用户名" />
      </el-form-item>
      <el-form-item prop="password">
        <el-input
          v-model="loginForm.password"
          :prefix-icon="Lock"
          type="password"
          placeholder="密码"
          show-password
          @keyup.enter="login"
        />
      </el-form-item>
      <el-form-item>
        <el-button
          :loading="loading"
          class="login-btn"
          type="primary"
          size="default"
          @click="login"
        >
          登录
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { User, Lock } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { useUserStore } from '@/store/modules/user'
import { ElNotification } from 'element-plus'

const useStore = useUserStore()
const $router = useRouter()
const $route = useRoute()

const loading = ref(false)
const loginForms = ref<FormInstance>()

const loginForm = ref({
  username: 'admin',
  password: '123456',
})

const rules = ref<FormRules<typeof loginForm.value>>({
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度 3~20 个字符', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度 6~20 个字符', trigger: 'blur' },
  ],
})

async function login() {
  // 表单校验
  const formEl = loginForms.value
  if (!formEl) return
  const valid = await formEl.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    // 调用 store 的 userLogin（内部调 client.POST，失败会 throw）
    await useStore.userLogin({
      username: loginForm.value.username,
      password: loginForm.value.password,
    })
    // 登录成功：跳转到 redirect 或首页
    const redirect = ($route.query.redirect as string) || '/home'
    $router.push(redirect)
    ElNotification({ type: 'success', message: '登录成功' })
  } catch (error) {
    // 登录失败（业务码非200 / HTTP 错误）：client 已弹过提示，这里可补一条
    ElNotification({
      type: 'error',
      message: error instanceof Error ? error.message : '登录失败',
    })
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="scss">
.login {
  width: 100%;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1f4037 0%, #99f2c8 100%);

  .login-form {
    width: 360px;
    padding: 30px;
    border-radius: 12px;
    background: rgba(255, 255, 255, 0.9);

    .title {
      text-align: center;
      margin-bottom: 20px;
      color: #333;
    }

    .login-btn {
      width: 100%;
    }
  }
}
</style>

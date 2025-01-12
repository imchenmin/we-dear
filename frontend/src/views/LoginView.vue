<template>
  <div class="login-view">
    <el-card class="login-card">
      <template #header>
        <div class="login-header">
          <h2>医生登录</h2>
        </div>
      </template>
      
      <el-form
        ref="loginForm"
        :model="form"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名" />
        </el-form-item>
        
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="form.password"
            type="password"
            show-password
            placeholder="请输入密码"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="handleLogin" :loading="loading">
            登录
          </el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="text" @click="goToPatientLogin">
            我是患者，去患者登录 >
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { authApi } from '@/api/auth'

const router = useRouter()
const userStore = useUserStore()

const form = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能小于6位', trigger: 'blur' }
  ]
}

const loading = ref(false)
const loginForm = ref()

const handleLogin = async () => {
  if (!loginForm.value) return
  
  await loginForm.value.validate(async (valid: boolean) => {
    if (valid) {
      loading.value = true
      try {
        const response = await authApi.login({
          username: form.username,
          password: form.password,
          role: 'doctor'
        })
        userStore.setToken(response.token)
        userStore.setUser(response.user)
        ElMessage.success('登录成功')
        router.push('/doctor')
      } catch (error) {
        // 错误已在request拦截器中处理
      } finally {
        loading.value = false
      }
    }
  })
}

const goToPatientLogin = () => {
  router.push('/patient/login')
}
</script>

<style scoped>
.login-view {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f7fa;
}

.login-card {
  width: 400px;
}

.login-header {
  text-align: center;
}

.login-header h2 {
  margin: 0;
  color: #303133;
}
</style> 
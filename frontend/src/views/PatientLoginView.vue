<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <div class="card-header">
          <span>{{ isLogin ? '患者登录' : '患者注册' }}</span>
          <el-button
            type="text"
            @click="toggleMode"
          >
            {{ isLogin ? '没有账号？立即注册' : '已有账号？立即登录' }}
          </el-button>
        </div>
      </template>

      <!-- 登录表单 -->
      <el-form
        v-if="isLogin"
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        label-width="80px"
      >
        <el-form-item label="用户名" prop="username">
          <el-input v-model="loginForm.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            show-password
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleLogin" :loading="loading">
            登录
          </el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="text" @click="goToDoctorLogin">
            我是医生，去医生登录 >
          </el-button>
        </el-form-item>
      </el-form>

      <!-- 注册表单 -->
      <el-form
        v-else
        ref="registerFormRef"
        :model="registerForm"
        :rules="registerRules"
        label-width="80px"
      >
        <el-form-item label="姓名" prop="name">
          <el-input v-model="registerForm.name" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="用户名" prop="username">
          <el-input v-model="registerForm.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="registerForm.password"
            type="password"
            placeholder="请输入密码"
            show-password
          />
        </el-form-item>
        <el-form-item label="选择医生" prop="doctorId">
          <el-select v-model="registerForm.doctorId" placeholder="请选择医生">
            <el-option
              v-for="doctor in doctors"
              :key="doctor.id"
              :label="doctor.name"
              :value="doctor.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleRegister" :loading="loading">
            注册
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { authApi } from '@/api/auth'
import { doctorApi } from '@/api/doctor'
import type { Doctor } from '@/types'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()
const isLogin = ref(true)
const loading = ref(false)
const doctors = ref<Doctor[]>([])

const loginForm = ref({
  username: '',
  password: ''
})

const registerForm = ref({
  name: '',
  username: '',
  password: '',
  doctorId: ''
})

const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能小于6位', trigger: 'blur' }
  ]
}

const registerRules = {
  name: [
    { required: true, message: '请输入姓名', trigger: 'blur' }
  ],
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能小于6位', trigger: 'blur' }
  ],
  doctorId: [
    { required: true, message: '请选择医生', trigger: 'change' }
  ]
}

const toggleMode = () => {
  isLogin.value = !isLogin.value
}

const handleLogin = async () => {
  loading.value = true
  try {
    const response = await authApi.login({
      username: loginForm.value.username,
      password: loginForm.value.password,
      role: 'patient'
    })
    
    // 使用 userStore 保存用户信息
    userStore.setToken(response.token)
    localStorage.setItem('userRole', 'patient')
    userStore.setUser(response.user)
    
    // 患者登录后跳转到患者聊天页面
    router.push('/patient/chat')
    ElMessage.success('登录成功')
  } catch (error: any) {
    ElMessage.error(error.message || '登录失败')
  } finally {
    loading.value = false
  }
}

const handleRegister = async () => {
  loading.value = true
  try {
    await authApi.register({
      name: registerForm.value.name,
      username: registerForm.value.username,
      password: registerForm.value.password,
      doctorId: registerForm.value.doctorId
    })
    ElMessage.success('注册成功，请登录')
    isLogin.value = true
  } catch (error: any) {
    ElMessage.error(error.message || '注册失败')
  } finally {
    loading.value = false
  }
}

const loadDoctors = async () => {
  try {
    doctors.value = await doctorApi.getAllDoctors()
  } catch (error: any) {
    ElMessage.error('加载医生列表失败')
  }
}

const goToDoctorLogin = () => {
  router.push('/doctor/login')
}

onMounted(() => {
  loadDoctors()
})
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f5f7fa;
}

.login-card {
  width: 400px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.el-form-item {
  margin-bottom: 20px;
}

.el-button {
  width: 100%;
}
</style> 
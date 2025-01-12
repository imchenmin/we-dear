import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { request } from '@/utils/request'
import { useRouter } from 'vue-router'

interface User {
  id: string
  username: string
  name: string
  role: 'doctor' | 'patient' | 'admin'
  doctorId?: string
  avatar?: string
}

export const useUserStore = defineStore('user', () => {
  const token = ref('')
  const user = ref<User | null>(null)
  const router = useRouter()

  // 计算属性
  const isAdmin = computed(() => user.value?.role === 'admin')
  const isDoctor = computed(() => user.value?.role === 'doctor')
  const isPatient = computed(() => user.value?.role === 'patient')

  function setToken(newToken: string) {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  function setUser(newUser: User) {
    user.value = newUser
    localStorage.setItem('user', JSON.stringify(newUser))
  }

  async function logout() {
    try {
      // 调用登出接口
      await request.post('/logout')
    } catch (error) {
      console.error('Logout failed:', error)
    } finally {
      // 清除本地存储和状态
      token.value = ''
      user.value = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      
      // 根据角色重定向到相应的登录页
      const currentRole = user.value?.role
      router.push(currentRole === 'doctor' ? '/doctor/login' : '/patient/login')
    }
  }

  // 修改密码
  async function changePassword(oldPassword: string, newPassword: string, userId?: string) {
    try {
      const data: any = {
        newPassword
      }
      
      // 如果是管理员修改他人密码
      if (userId && user.value?.role === 'admin') {
        data.userId = userId
      } else {
        // 普通用户修改自己的密码需要旧密码
        data.oldPassword = oldPassword
      }
      
      await request.post('/change-password', data)
    } catch (error) {
      throw error
    }
  }

  // 初始化状态
  function init() {
    const storedToken = localStorage.getItem('token')
    const storedUser = localStorage.getItem('user')
    
    if (storedToken) {
      token.value = storedToken
    }
    
    if (storedUser) {
      try {
        user.value = JSON.parse(storedUser)
      } catch (error) {
        console.error('Failed to parse stored user:', error)
        // 如果解析失败，清除存储
        localStorage.removeItem('user')
        user.value = null
      }
    }
  }

  // 获取用户角色
  function getRole() {
    return user.value?.role
  }

  // 检查是否已登录
  function isLoggedIn() {
    return !!token.value && !!user.value
  }

  // 初始化
  init()

  return {
    token,
    user,
    isAdmin,
    isDoctor,
    isPatient,
    setToken,
    setUser,
    logout,
    changePassword,
    getRole,
    isLoggedIn
  }
}) 
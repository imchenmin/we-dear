import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { request } from '@/utils/request'

// 从 cookie 中获取 token
function getTokenFromCookie(): string | null {
  const cookies = document.cookie.split(';')
  for (const cookie of cookies) {
    const [name, value] = cookie.trim().split('=')
    if (name === 'token') {
      return value
    }
  }
  return null
}

export const useUserStore = defineStore('user', () => {
  const token = ref('')
  const user = ref<any>(null)

  // 计算属性：是否是管理员
  const isAdmin = computed(() => user.value?.role === 'admin')

  function setToken(newToken: string) {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  function setUser(newUser: any) {
    user.value = newUser
    localStorage.setItem('user', JSON.stringify(newUser))
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    // 清除 cookie
    document.cookie = 'token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT'
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
  const cookieToken = getTokenFromCookie()
  const storedToken = localStorage.getItem('token')
  const storedUser = localStorage.getItem('user')
  
  // 优先使用 cookie 中的 token
  if (cookieToken) {
    token.value = cookieToken
  } else if (storedToken) {
    token.value = storedToken
  }
  
  if (storedUser) {
    user.value = JSON.parse(storedUser)
  }

  return {
    token,
    user,
    isAdmin,
    setToken,
    setUser,
    logout,
    changePassword
  }
}) 
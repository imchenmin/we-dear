import { defineStore } from 'pinia'
import { ref } from 'vue'

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
    setToken,
    setUser,
    logout
  }
}) 
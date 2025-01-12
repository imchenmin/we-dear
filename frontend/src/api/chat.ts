import { request } from '@/utils/request'
import type { Message } from '@/types'
import { useUserStore } from '@/stores/user'

export const chatApi = {
  // 获取聊天历史
  async getChatHistory(patientId: string): Promise<Message[]> {
    const userStore = useUserStore()
    console.log('Fetching chat history for patient:', patientId)
    return request.get(`/chat/${patientId}`, {
      params: {
        role: userStore.user?.role,
        userId: userStore.user?.id
      }
    })
  },

  // 发送消息
  async sendMessage(patientId: string, content: string, role: 'doctor' | 'patient'): Promise<Message> {
    const userStore = useUserStore()
    console.log(`Sending ${role} message:`, { patientId, content })
    return request.post(`/chat/${patientId}/${role}`, {
      content,
      sender: userStore.user?.id,
      type: 'text',
      role
    })
  }
} 
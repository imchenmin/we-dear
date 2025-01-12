import type { Patient, Message, AISuggestion, MessageFeedback, FeedbackStats } from '@/types'
import { request } from '@/utils/request'

export const patientApi = {
  // 获取患者列表
  async getPatients(): Promise<Patient[]> {
    console.log('Fetching patients list')
    return request.get('/patients')
  },

  // 获取患者详情
  async getPatientById(id: string): Promise<Patient> {
    console.log('Fetching patient details:', id)
    return request.get(`/patients/${id}`)
  },

  // 获取聊天历史
  async getChatHistory(patientId: string): Promise<Message[]> {
    console.log('Fetching chat history for patient:', patientId)
    return request.get(`/chat/${patientId}`, {
      params: {
        role: 'doctor',
        userId: localStorage.getItem('userId')
      }
    })
  },

  // 发送医生消息
  async sendDoctorMessage(patientId: string, content: string, sender: string): Promise<Message> {
    console.log('Sending doctor message:', { patientId, content })
    return request.post(`/chat/${patientId}/doctor`, {
      content,
      sender,
      type: 'text',
      role: 'doctor'
    })
  },

  // 发送患者消息
  async sendPatientMessage(patientId: string, content: string, sender: string): Promise<Message> {
    console.log('Sending patient message:', { patientId, content })
    return request.post(`/chat/${patientId}/patient`, {
      content,
      sender,
      type: 'text',
      role: 'patient'
    })
  },

  // 获取 AI 建议
  async getAISuggestions(patientId: string, messageId: string): Promise<AISuggestion[]> {
    console.log('Fetching AI suggestions:', { patientId, messageId })
    return request.get(`/chat/${patientId}/suggestions`, {
      params: { messageId }
    })
  },

  // 创建消息评价
  async createMessageFeedback(suggestionId: string, data: Partial<MessageFeedback>): Promise<MessageFeedback> {
    const response = await request.post(`/ai-suggestions/${suggestionId}/feedback`, {
      ...data,
      suggestionId
    })
    return response.data
  },

  // 更新消息评价
  async updateMessageFeedback(feedbackId: string, data: Partial<MessageFeedback>): Promise<MessageFeedback> {
    const response = await request.put(`/ai-suggestions/feedback/${feedbackId}`, data)
    return response.data
  },

  // 获取消息评价列表
  async getMessageFeedbacks(suggestionId: string): Promise<MessageFeedback[]> {
    const response = await request.get(`/ai-suggestions/feedback`, {
      params: { suggestionId }
    })
    return response.data
  },

  // 获取消息评价统计
  async getMessageFeedbackStats(suggestionId: string): Promise<FeedbackStats> {
    const response = await request.get(`/ai-suggestions/feedback/stats`, {
      params: { suggestionId }
    })
    return response.data
  }
} 
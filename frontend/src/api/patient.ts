import type { Patient, Message, AISuggestion } from '@/types'
import axios from 'axios'

const API_BASE = '/api'

// 创建axios实例
const apiClient = axios.create({
  baseURL: API_BASE,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 响应拦截器
apiClient.interceptors.response.use(
  response => response.data,
  error => {
    console.error('API Error:', error.response?.data || error.message)
    return Promise.reject(error)
  }
)

export const patientApi = {
  // 获取患者列表
  async getPatients(): Promise<Patient[]> {
    console.log('Fetching patients list')
    return apiClient.get('/patients')
  },

  // 获取患者详情
  async getPatientById(id: string): Promise<Patient> {
    console.log('Fetching patient details:', id)
    return apiClient.get(`/patients/${id}`)
  },

  // 获取聊天历史
  async getChatHistory(patientId: string): Promise<Message[]> {
    console.log('Fetching chat history for patient:', patientId)
    return apiClient.get(`/chat/${patientId}`)
  },

  // 发送医生消息
  async sendDoctorMessage(patientId: string, content: string, sender: string): Promise<Message> {
    console.log('Sending doctor message:', { patientId, content })
    return apiClient.post(`/chat/${patientId}/doctor`, {
      content,
      sender,
      type: 'text',
      role: 'doctor'
    })
  },

  // 发送患者消息
  async sendPatientMessage(patientId: string, content: string, sender: string): Promise<Message> {
    console.log('Sending patient message:', { patientId, content })
    return apiClient.post(`/chat/${patientId}/patient`, {
      content,
      sender,
      type: 'text',
      role: 'patient'
    })
  },

  // 获取 AI 建议
  async getAISuggestions(patientId: string, messageId: string): Promise<AISuggestion[]> {
    console.log('Fetching AI suggestions:', { patientId, messageId })
    return apiClient.get(`/chat/${patientId}/suggestions`, {
      params: { messageId }
    })
  }
} 
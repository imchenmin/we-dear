import type { Patient, Message, AISuggestion } from '@/types'
import axios from 'axios'

const API_BASE = '/api'

export const patientApi = {
  // 获取患者列表
  async getPatients(): Promise<Patient[]> {
    const response = await axios.get(`${API_BASE}/patients`)
    return response.data
  },

  // 获取患者详情
  async getPatientById(id: string): Promise<Patient> {
    const response = await axios.get(`${API_BASE}/patients/${id}`)
    return response.data
  },

  // 获取聊天历史
  async getChatHistory(patientId: string): Promise<Message[]> {
    const response = await axios.get(`${API_BASE}/chat/${patientId}`)
    return response.data
  },

  // 发送医生消息
  async sendDoctorMessage(patientId: string, content: string, sender: string): Promise<Message> {
    const response = await axios.post(`${API_BASE}/chat/${patientId}/doctor`, {
      content,
      sender,
      type: 'text',
      role: 'doctor'
    })
    return response.data
  },

  // 发送患者消息
  async sendPatientMessage(patientId: string, content: string, sender: string): Promise<Message> {
    const response = await axios.post(`${API_BASE}/chat/${patientId}/patient`, {
      content,
      sender,
      type: 'text',
      role: 'patient'
    })
    return response.data
  },

  // 获取 AI 建议
  async getAISuggestions(patientId: string, messageId: string): Promise<AISuggestion[]> {
    const response = await axios.get(`${API_BASE}/chat/${patientId}/suggestions`, {
      params: { messageId }
    })
    return response.data
  }
} 
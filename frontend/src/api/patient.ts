import type { Patient, Message, AISuggestion } from '@/types'
import axios from 'axios'
import { request } from '@/utils/request'

const API_BASE = '/api'


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
    return request.get(`/chat/${patientId}`)
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
  }
} 
import type { Doctor } from '@/types'
import { request } from '@/utils/request'

export const doctorApi = {
  // 获取所有医生列表
  async getAllDoctors(): Promise<Doctor[]> {
    return request.get('/doctors')
  },

  // 获取医生详情
  async getDoctorById(id: string): Promise<Doctor> {
    return request.get(`/doctors/${id}`)
  },

  // 创建医生
  async createDoctor(data: Partial<Doctor>): Promise<Doctor> {
    return request.post('/doctors', data)
  },

  // 更新医生信息
  async updateDoctor(id: string, data: Partial<Doctor>): Promise<Doctor> {
    return request.put(`/doctors/${id}`, data)
  },

  // 删除医生
  async deleteDoctor(id: string): Promise<void> {
    return request.delete(`/doctors/${id}`)
  }
} 
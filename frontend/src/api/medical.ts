import { request } from '@/utils/request'
import type { FollowUpRecord, MedicalRecord } from '@/types'

export const medicalApi = {
  // 随访记录相关接口
  getFollowUpRecords(patientId: string) {
    return request.get<FollowUpRecord[]>(`/patients/${patientId}/followup`)
  },

  createFollowUpRecord(record: Partial<FollowUpRecord>) {
    return request.post<FollowUpRecord>('/followup', record)
  },

  updateFollowUpRecord(record: FollowUpRecord) {
    return request.put<FollowUpRecord>(`/followup/${record.id}`, record)
  },

  deleteFollowUpRecord(id: string) {
    return request.delete(`/followup/${id}`)
  },

  // 医疗记录相关接口
  getMedicalRecords(patientId: string) {
    return request.get<MedicalRecord[]>(`/patients/${patientId}/medical`)
  },

  createMedicalRecord(record: Partial<MedicalRecord>) {
    return request.post<MedicalRecord>('/medical', record)
  },

  updateMedicalRecord(record: MedicalRecord) {
    return request.put<MedicalRecord>(`/medical/${record.id}`, record)
  },

  deleteMedicalRecord(id: string) {
    return request.delete(`/medical/${id}`)
  }
} 
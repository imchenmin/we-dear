interface BaseModel {
  id: string
  createdAt: string
  updatedAt: string
}

export interface Department extends BaseModel {
  name: string
  description: string
  code: string
}

export interface Doctor extends BaseModel {
  name: string
  title: string
  departmentId: string
  department: Department
  license: string
  specialty: string
  avatar: string
  status: string
}

export interface Patient extends BaseModel {
  name: string
  gender: string
  age: number
  birthday: string
  phone: string
  emergencyPhone: string
  address: string
  idCard: string
  bloodType: string
  allergies: string[]
  chronicDiseases: string[]
  avatar: string
}

export interface Message extends BaseModel {
  patientId: string
  doctorId: string
  recordId: string
  content: string
  type: string
  role: string
  status: string
  replyTo: string
}

export interface AISuggestion extends BaseModel {
  messageId: string
  patientId: string
  content: string
  confidence: number
  category: string
  priority: number
  status: string
  reviewedBy: string
  reviewedAt: string
  reviewNotes: string
}

export interface MedicalRecord extends BaseModel {
  patientId: string
  doctorId: string
  diagnosisDate: string
  symptoms: string[]
  diagnosis: string
  treatment: string
  prescription: string
  notes: string
  status: string
  type: string
  department: string
  cost: number
  attachments: string[]
  createdAt: string
  updatedAt: string
}

export interface Attachment extends BaseModel {
  messageId: string
  recordId: string
  type: string
  url: string
  name: string
  size: number
  contentType: string
  uploadedBy: string
}

export interface FollowUpRecord {
  id: string
  patientId: string
  doctorId: string
  title: string
  content: string
  followUpDate: string
  nextFollowUp: string
  status: 'completed' | 'pending'
  type: string
  attachments: string[]
  createdAt: string
  updatedAt: string
} 
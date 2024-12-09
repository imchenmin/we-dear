export interface Message {
  id: string
  content: string
  timestamp: number
  role: 'doctor' | 'patient'
  sender: string
  avatar?: string
}

export interface Patient {
  id: string
  name: string
  gender: string
  age: number
  phone: string
  diagnosis: string
  doctor: string
  avatar: string
  messages: Message[]
} 
import { request } from '@/utils/request'

interface LoginRequest {
  username: string
  password: string
  role: 'doctor' | 'patient'
}

interface RegisterRequest {
  name: string
  username: string
  password: string
  doctorId: string
}

interface LoginResponse {
  token: string
  user: {
    id: string
    username: string
    name: string
    role: 'doctor' | 'patient'
    doctorId?: string
    avatar?: string
  }
}

export const authApi = {
  async login(data: LoginRequest): Promise<LoginResponse> {
    return request.post('/login', data)
  },

  async register(data: RegisterRequest): Promise<void> {
    return request.post('/register', data)
  },

  async changePassword(data: { oldPassword: string; newPassword: string }): Promise<void> {
    return request.post('/change-password', data)
  }
} 
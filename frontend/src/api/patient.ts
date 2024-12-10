import type { Patient, Message } from '@/types'

class PatientApi {
  private baseUrl = 'http://localhost:8080/api'

  async getAllPatients(): Promise<Patient[]> {
    try {
      console.log('Fetching all patients from:', `${this.baseUrl}/patients`)
      const response = await fetch(`${this.baseUrl}/patients`)
      
      if (!response.ok) {
        const errorText = await response.text()
        console.error('Failed to fetch patients:', errorText)
        throw new Error(`Failed to fetch patients: ${response.status} ${response.statusText}`)
      }

      const data = await response.json()
      console.log('Received patients data:', data)
      return data
    } catch (error) {
      console.error('Error in getAllPatients:', error)
      throw error
    }
  }

  async getPatientById(id: string): Promise<Patient> {
    try {
      console.log('Fetching patient details:', id)
      const response = await fetch(`${this.baseUrl}/patients/${id}`)
      
      if (!response.ok) {
        const errorText = await response.text()
        console.error(`Failed to fetch patient ${id}:`, errorText)
        throw new Error(`Failed to fetch patient: ${response.status} ${response.statusText}`)
      }

      const data = await response.json()
      console.log('Received patient data:', data)
      return data
    } catch (error) {
      console.error('Error in getPatientById:', error)
      throw error
    }
  }

  async getChatHistory(patientId: string): Promise<Message[]> {
    try {
      console.log('Fetching chat history for patient:', patientId)
      const response = await fetch(`${this.baseUrl}/chat/${patientId}`)
      
      if (!response.ok) {
        const errorText = await response.text()
        console.error('Failed to fetch chat history:', errorText)
        throw new Error(`Failed to fetch chat history: ${response.status} ${response.statusText}`)
      }

      const data = await response.json()
      console.log('Received chat history:', data)
      return data
    } catch (error) {
      console.error('Error in getChatHistory:', error)
      throw error
    }
  }

  async sendDoctorMessage(patientId: string, content: string, sender: string, avatar?: string): Promise<Message> {
    try {
      console.log('Sending doctor message to patient:', patientId)
      const response = await fetch(`${this.baseUrl}/chat/${patientId}/doctor`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          content,
          sender,
          avatar,
          timestamp: Date.now(),
        }),
      })
      
      if (!response.ok) {
        const errorText = await response.text()
        console.error('Failed to send doctor message:', errorText)
        throw new Error(`Failed to send message: ${response.status} ${response.statusText}`)
      }

      const data = await response.json()
      console.log('Message sent successfully:', data)
      return data
    } catch (error) {
      console.error('Error in sendDoctorMessage:', error)
      throw error
    }
  }

  async sendPatientMessage(patientId: string, content: string, sender: string, avatar?: string): Promise<Message> {
    try {
      console.log('Sending patient message:', patientId)
      const response = await fetch(`${this.baseUrl}/chat/${patientId}/patient`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          content,
          sender,
          avatar,
          timestamp: Date.now(),
        }),
      })
      
      if (!response.ok) {
        const errorText = await response.text()
        console.error('Failed to send patient message:', errorText)
        throw new Error(`Failed to send message: ${response.status} ${response.statusText}`)
      }

      const data = await response.json()
      console.log('Message sent successfully:', data)
      return data
    } catch (error) {
      console.error('Error in sendPatientMessage:', error)
      throw error
    }
  }
}

export const patientApi = new PatientApi() 
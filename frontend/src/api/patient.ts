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

  async sendMessage(patientId: string, message: Message): Promise<void> {
    try {
      console.log('Sending message to patient:', patientId, message)
      const response = await fetch(`${this.baseUrl}/patients/${patientId}/messages`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(message),
      })
      
      if (!response.ok) {
        const errorText = await response.text()
        console.error('Failed to send message:', errorText)
        throw new Error(`Failed to send message: ${response.status} ${response.statusText}`)
      }

      console.log('Message sent successfully')
    } catch (error) {
      console.error('Error in sendMessage:', error)
      throw error
    }
  }
}

export const patientApi = new PatientApi() 
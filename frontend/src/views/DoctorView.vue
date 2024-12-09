<template>
  <div class="doctor-view">
    <div class="patient-section">
      <PatientList
        v-model="activePatientId"
        :patients="patients"
      />
    </div>
    
    <div class="chat-section">
      <div class="chat-header">
        <h2>{{ selectedPatient?.name || '无选中患者' }}</h2>
      </div>
      
      <div class="chat-messages" ref="chatContainer">
        <template v-if="selectedPatient">
          <ChatMessage
            v-for="message in selectedPatient.messages"
            :key="message.id"
            :message="message"
          />
          <div v-if="selectedPatient.messages.length === 0" class="no-messages">
            暂无消息记录
          </div>
        </template>
        <div v-else class="no-patient-selected">
          请选择一位患者开始对话
        </div>
      </div>
      
      <ChatInput @send="handleSendMessage" :disabled="!selectedPatient" />
    </div>
    
    <div class="profile-section">
      <PatientProfile :patient="selectedPatient" />
    </div>

    <el-loading 
      v-model:visible="loading"
      :lock="true"
      text="加载中..."
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage, ElLoading } from 'element-plus'
import PatientList from '@/components/PatientList.vue'
import PatientProfile from '@/components/PatientProfile.vue'
import ChatMessage from '@/components/ChatMessage.vue'
import ChatInput from '@/components/ChatInput.vue'
import { patientApi } from '@/api/patient'
import type { Patient } from '@/types'

const patients = ref<Patient[]>([])
const activePatientId = ref('')
const chatContainer = ref<HTMLElement | null>(null)
const loading = ref(false)

const selectedPatient = computed(() => {
  return patients.value.find(p => p.id === activePatientId.value)
})

// 获取所有患者信息
const fetchPatients = async () => {
  loading.value = true
  try {
    console.log('Fetching all patients...')
    const response = await patientApi.getAllPatients()
    console.log('Received patients:', response)
    patients.value = response || []
  } catch (error) {
    console.error('Error fetching patients:', error)
    ElMessage.error('获取患者列表失败，请刷新页面重试')
    patients.value = []
  } finally {
    loading.value = false
  }
}

// 获取患者详细信息
const fetchPatientDetail = async (id: string) => {
  if (!id) return
  
  loading.value = true
  try {
    console.log('Fetching patient details:', id)
    const response = await patientApi.getPatientById(id)
    console.log('Received patient details:', response)
    
    const patientIndex = patients.value.findIndex(p => p.id === id)
    if (patientIndex !== -1) {
      patients.value[patientIndex] = {
        ...response,
        messages: response.messages || []
      }
    }
  } catch (error) {
    console.error('Error fetching patient details:', error)
    ElMessage.error('获取患者详情失败，请重试')
  } finally {
    loading.value = false
  }
}

// 发送消息
const handleSendMessage = async (content: string) => {
  if (!activePatientId.value || !content.trim()) {
    ElMessage.warning('请先选择患者并输入消息')
    return
  }

  loading.value = true
  try {
    const message = {
      id: Date.now().toString(),
      content: content.trim(),
      timestamp: Date.now(),
      role: 'doctor' as const,
      sender: '医生',
    }

    await patientApi.sendMessage(activePatientId.value, message)
    await fetchPatientDetail(activePatientId.value)
    scrollToBottom()
  } catch (error) {
    console.error('Error sending message:', error)
    ElMessage.error('发送消息失败，请重试')
  } finally {
    loading.value = false
  }
}

// 滚动到底部
const scrollToBottom = () => {
  if (chatContainer.value) {
    setTimeout(() => {
      if (chatContainer.value) {
        chatContainer.value.scrollTop = chatContainer.value.scrollHeight
      }
    }, 100)
  }
}

// 监听患者切换
watch(activePatientId, (newId) => {
  if (newId) {
    console.log('Patient selected:', newId)
    fetchPatientDetail(newId)
  }
})

// 监听消息变化，自动滚动
watch(() => selectedPatient?.value?.messages, () => {
  scrollToBottom()
}, { deep: true })

onMounted(async () => {
  console.log('DoctorView mounted')
  await fetchPatients()
})
</script>

<style scoped>
.doctor-view {
  display: grid;
  grid-template-columns: 300px 1fr 300px;
  height: 100vh;
  background: #f5f7fa;
  position: relative;
}

.patient-section {
  background: #fff;
  border-right: 1px solid #dcdfe6;
}

.chat-section {
  display: flex;
  flex-direction: column;
  background: #fff;
  margin: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.chat-header {
  padding: 16px;
  border-bottom: 1px solid #dcdfe6;
}

.chat-header h2 {
  margin: 0;
  color: #303133;
  font-size: 18px;
}

.chat-messages {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

.no-patient-selected,
.no-messages {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #909399;
  font-size: 16px;
}

.profile-section {
  padding: 20px;
  background: #f5f7fa;
}
</style> 
<template>
  <div class="patient-chat">
    <el-card class="chat-card">
      <template #header>
        <div class="card-header">
          <el-select
            v-model="selectedPatientId"
            placeholder="选择患者"
            class="patient-select"
          >
            <el-option
              v-for="patient in patients"
              :key="patient.id"
              :label="patient.name"
              :value="patient.id"
            />
          </el-select>
        </div>
      </template>

      <div class="chat-messages" ref="chatContainer">
        <template v-if="currentPatient">
          <ChatMessage
            v-for="message in currentPatient.messages"
            :key="message.id"
            :message="message"
          />
          <div v-if="currentPatient.messages?.length === 0" class="no-messages">
            暂无消息记录
          </div>
        </template>
        <div v-else class="no-patient-selected">
          请选择一位患者开始对话
        </div>
      </div>

      <div class="chat-input">
        <ChatInput 
          @send="handleSendMessage" 
          :disabled="!currentPatient"
        />
        <div class="upload-actions">
          <el-upload
            action="/api/upload"
            :show-file-list="false"
            :on-success="handleImageUpload"
            :disabled="!currentPatient"
          >
            <el-button :icon="Picture" :disabled="!currentPatient">图片</el-button>
          </el-upload>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { Picture } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import ChatMessage from '@/components/ChatMessage.vue'
import ChatInput from '@/components/ChatInput.vue'
import { patientApi } from '@/api/patient'
import type { Patient } from '@/types'

const patients = ref<Patient[]>([])
const selectedPatientId = ref('')
const chatContainer = ref<HTMLElement | null>(null)

// 获取当前选中的患者信息
const currentPatient = computed(() => {
  return patients.value.find(p => p.id === selectedPatientId.value)
})

// 加载患者列表
const loadPatients = async () => {
  try {
    const response = await patientApi.getAllPatients()
    patients.value = response
  } catch (error) {
    console.error('Failed to load patients:', error)
    ElMessage.error('加载患者列表失败')
  }
}

// 加载患者消息
const loadPatientDetail = async () => {
  if (!selectedPatientId.value) return

  try {
    const response = await patientApi.getPatientById(selectedPatientId.value)
    const patientIndex = patients.value.findIndex(p => p.id === selectedPatientId.value)
    if (patientIndex !== -1) {
      patients.value[patientIndex] = response
    }
  } catch (error) {
    console.error('Failed to load patient details:', error)
    ElMessage.error('加载患者信息失败')
  }
}

// 发送文本消息
const handleSendMessage = async (content: string) => {
  if (!currentPatient.value || !content.trim()) return

  try {
    await patientApi.sendPatientMessage(
      selectedPatientId.value,
      content,
      currentPatient.value.name,
      currentPatient.value.avatar
    )
    await loadPatientDetail()
    scrollToBottom()
  } catch (error) {
    console.error('Failed to send message:', error)
    ElMessage.error('发送消息失败')
  }
}

// 处理图片上传
const handleImageUpload = async (response: { url: string }) => {
  if (!currentPatient.value) return

  try {
    await patientApi.sendPatientMessage(
      selectedPatientId.value,
      `[图片] ${response.url}`,
      currentPatient.value.name,
      currentPatient.value.avatar
    )
    await loadPatientDetail()
    scrollToBottom()
  } catch (error) {
    console.error('Failed to send image message:', error)
    ElMessage.error('发送图片失败')
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

// 监听患者选择变化
watch(selectedPatientId, () => {
  if (selectedPatientId.value) {
    loadPatientDetail()
  }
})

// 组件挂载时加载数据
onMounted(() => {
  loadPatients()
})
</script>

<style scoped>
.patient-chat {
  max-width: 800px;
  margin: 20px auto;
  height: calc(100vh - 40px);
}

.chat-card {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.patient-select {
  width: 200px;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background: #f5f7fa;
  border-radius: 4px;
  margin-bottom: 20px;
}

.no-messages,
.no-patient-selected {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #909399;
  font-size: 16px;
}

.chat-input {
  border-top: 1px solid #dcdfe6;
  padding-top: 20px;
}

.upload-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 12px;
  gap: 8px;
}
</style> 
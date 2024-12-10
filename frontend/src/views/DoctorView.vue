<template>
  <div class="doctor-view">
    <!-- 左侧患者列表 -->
    <div class="patient-list-container">
      <PatientList
        :patients="patients"
        v-model="activePatientId"
      />
    </div>

    <!-- 右侧聊天区域 -->
    <div class="chat-container" v-if="selectedPatient">
      <!-- 患者信息卡片 -->
      <div class="patient-info-card">
        <PatientProfile :patient="selectedPatient" />
      </div>

      <!-- 聊天区域 -->
      <div class="chat-area">
        <div class="chat-messages" ref="messagesContainer">
          <template v-if="messages.length > 0">
            <ChatMessage
              v-for="message in messages"
              :key="message.id"
              :message="message"
              :showAISuggestion="true"
            />
          </template>
          <div v-else class="no-messages">
            暂无消息记录
          </div>
        </div>

        <!-- 输入区域 -->
        <div class="chat-input-area">
          <ChatInput
            @send="handleSendMessage"
            :disabled="!selectedPatient"
          />
        </div>
      </div>
    </div>

    <!-- 未选择患者时的提示 -->
    <div v-else class="no-patient-selected">
      <el-empty description="请选择一位患者开始对话" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import PatientList from '@/components/PatientList.vue'
import PatientProfile from '@/components/PatientProfile.vue'
import ChatMessage from '@/components/ChatMessage.vue'
import ChatInput from '@/components/ChatInput.vue'
import { patientApi } from '@/api/patient'
import type { Patient, Message } from '@/types'

const patients = ref<Patient[]>([])
const activePatientId = ref('')
const messages = ref<Message[]>([])
const messagesContainer = ref<HTMLElement | null>(null)

// 当前选中的患者
const selectedPatient = computed(() => {
  return patients.value.find(p => p.id === activePatientId.value)
})

// 获取所有患者信息
const loadPatients = async () => {
  try {
    patients.value = await patientApi.getPatients()
  } catch (error) {
    console.error('Failed to load patients:', error)
    ElMessage.error('加载患者列表失败')
  }
}

// 获取聊天记录
const loadMessages = async () => {
  if (!activePatientId.value) return
  
  try {
    console.log('Loading chat history for patient:', activePatientId.value)
    messages.value = await patientApi.getChatHistory(activePatientId.value)
    console.log('Loaded messages:', messages.value)
    scrollToBottom()
  } catch (error) {
    console.error('Failed to load messages:', error)
    ElMessage.error('加载聊天记录失败')
  }
}

// 发送消息
const handleSendMessage = async (content: string) => {
  if (!selectedPatient.value) return
  
  try {
    console.log('Sending message for patient:', activePatientId.value)
    const message = await patientApi.sendDoctorMessage(
      activePatientId.value,
      content,
      'doctor'
    )
    console.log('Message sent:', message)
    messages.value.push(message)
    scrollToBottom()
  } catch (error) {
    console.error('Failed to send message:', error)
    ElMessage.error('发送消息失败')
  }
}

// 滚动到底部
const scrollToBottom = () => {
  setTimeout(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  }, 100)
}

// 监听患者切换
watch(activePatientId, () => {
  messages.value = []
  loadMessages()
})

// 组件挂载时加载数据
onMounted(() => {
  loadPatients()
})
</script>

<style scoped>
.doctor-view {
  display: flex;
  height: 100vh;
  background-color: #f5f7fa;
}

.patient-list-container {
  width: 300px;
  flex-shrink: 0;
  background-color: #fff;
  border-right: 1px solid #e4e7ed;
}

.chat-container {
  flex: 1;
  display: flex;
  gap: 20px;
  padding: 20px;
  overflow: hidden;
}

.patient-info-card {
  width: 300px;
  flex-shrink: 0;
}

.chat-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.chat-messages {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

.chat-input-area {
  flex-shrink: 0;
  border-top: 1px solid #e4e7ed;
}

.no-messages {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  color: #909399;
  font-size: 14px;
}

.no-patient-selected {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #fff;
}

/* 自定义滚动条样式 */
.chat-messages::-webkit-scrollbar {
  width: 6px;
}

.chat-messages::-webkit-scrollbar-thumb {
  background-color: #dcdfe6;
  border-radius: 3px;
}

.chat-messages::-webkit-scrollbar-track {
  background-color: #f5f7fa;
}
</style> 
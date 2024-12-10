<template>
  <div class="patient-chat">
    <!-- 左侧患者列表 -->
    <div class="patient-list-container">
      <PatientList
        :patients="patients"
        v-model="selectedPatientId"
      />
    </div>

    <!-- 右侧聊天区域 -->
    <div class="chat-container" v-if="currentPatient">
      <!-- 聊天区域 -->
      <div class="chat-area">
        <div class="chat-header">
          <el-avatar :size="40" :src="currentPatient.avatar" />
          <div class="patient-info">
            <div class="patient-name">{{ currentPatient.name }}</div>
            <div class="patient-details">
              {{ currentPatient.age }}岁 | {{ currentPatient.gender === 'male' ? '男' : '女' }}
              <template v-if="currentPatient.chronicDiseases.length > 0">
                | {{ currentPatient.chronicDiseases.join('、') }}
              </template>
            </div>
          </div>
        </div>

        <div class="chat-messages" ref="messagesContainer">
          <template v-if="messages.length > 0">
            <ChatMessage
              v-for="message in messages"
              :key="message.id"
              :message="message"
            />
          </template>
          <div v-else class="no-messages">
            暂无消息记录
          </div>
        </div>

        <div class="chat-input-area">
          <ChatInput
            @send="handleSendMessage"
            :disabled="!currentPatient"
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
import ChatMessage from '@/components/ChatMessage.vue'
import ChatInput from '@/components/ChatInput.vue'
import PatientList from '@/components/PatientList.vue'
import { patientApi } from '@/api/patient'
import type { Patient, Message } from '@/types'

const patients = ref<Patient[]>([])
const selectedPatientId = ref('')
const messages = ref<Message[]>([])
const messagesContainer = ref<HTMLElement | null>(null)

// 获取当前选中的患者信息
const currentPatient = computed(() => {
  return patients.value.find(p => p.id === selectedPatientId.value)
})

// 加载患者列表
const loadPatients = async () => {
  try {
    patients.value = await patientApi.getPatients()
  } catch (error) {
    console.error('Failed to load patients:', error)
    ElMessage.error('加载患者列表失败')
  }
}

// 加载聊天记录
const loadMessages = async () => {
  if (!selectedPatientId.value) return

  try {
    messages.value = await patientApi.getChatHistory(selectedPatientId.value)
    scrollToBottom()
  } catch (error) {
    console.error('Failed to load messages:', error)
    ElMessage.error('加载聊天记录失败')
  }
}

// 发送消息
const handleSendMessage = async (content: string) => {
  if (!currentPatient.value) return

  try {
    const message = await patientApi.sendPatientMessage(
      selectedPatientId.value,
      content,
      currentPatient.value.name
    )
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
watch(selectedPatientId, () => {
  messages.value = []
  loadMessages()
})

// 组件挂载时加载数据
onMounted(() => {
  loadPatients()
})
</script>

<style scoped>
.patient-chat {
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
  padding: 20px;
  overflow: hidden;
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

.chat-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  border-bottom: 1px solid #e4e7ed;
  background-color: #fff;
}

.patient-info {
  flex: 1;
}

.patient-name {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
}

.patient-details {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
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
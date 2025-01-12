<template>
  <div class="patient-chat">
    <!-- 聊天区域 -->
    <div class="chat-container">
      <div class="chat-area">
        <div class="chat-header">
          <template v-if="selectedDoctor">
            <el-avatar :size="40" :src="selectedDoctor.avatar" />
            <div class="doctor-info">
              <div class="doctor-name">{{ selectedDoctor.name }}</div>
              <div class="doctor-details">
                {{ selectedDoctor.title }} | {{ selectedDoctor.department?.name }}
              </div>
            </div>
          </template>
          <template v-else>
            <div class="doctor-info">
              <div class="doctor-name">正在加载医生信息...</div>
            </div>
          </template>
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
            :disabled="!selectedDoctor"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import ChatMessage from '@/components/ChatMessage.vue'
import ChatInput from '@/components/ChatInput.vue'
import { doctorApi } from '@/api/doctor'
import type { Doctor, Message } from '@/types'
import { useRouter } from 'vue-router'

const router = useRouter()
const userStore = useUserStore()
const selectedDoctor = ref<Doctor | null>(null)
const messages = ref<Message[]>([])
const messagesContainer = ref<HTMLElement | null>(null)

// 检查用户登录状态
const checkUserLogin = () => {
  if (!userStore.token || !userStore.user) {
    ElMessage.error('请先登录')
    router.push('/patient/login')
    return false
  }
  return true
}

// 加载医生信息
const loadDoctor = async () => {
  if (!checkUserLogin()) return

  try {
    // 直接从用户信息中获取医生ID
    const doctorId = userStore.user.doctorId
    if (!doctorId) {
      ElMessage.error('未分配医生')
      return
    }

    // 使用医生信息接口
    const doctorResponse = await fetch(`/api/doctors/${doctorId}`, {
      headers: {
        'Authorization': `Bearer ${userStore.token}`
      }
    })
    if (!doctorResponse.ok) {
      throw new Error('Failed to fetch doctor info')
    }
    selectedDoctor.value = await doctorResponse.json()
  } catch (error) {
    console.error('Failed to load doctor:', error)
    ElMessage.error('加载医生信息失败')
  }
}

// 加载聊天记录
const loadMessages = async () => {
  if (!userStore.user) {
    ElMessage.error('用户未登录')
    return
  }

  try {
    const response = await fetch(`/api/chat/${userStore.user.id}`, {
      headers: {
        'Authorization': `Bearer ${userStore.token}`
      },
      params: {
        role: 'patient',
        userId: userStore.user.id
      }
    })
    if (!response.ok) throw new Error('Failed to fetch messages')
    messages.value = await response.json()
    scrollToBottom()
  } catch (error) {
    console.error('Failed to load messages:', error)
    ElMessage.error('加载聊天记录失败')
  }
}

// 发送消息
const handleSendMessage = async (content: string) => {
  if (!userStore.user || !selectedDoctor.value) {
    ElMessage.error('用户未登录或未分配医生')
    return
  }

  try {
    const response = await fetch(`/api/chat/${userStore.user.id}/patient`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${userStore.token}`
      },
      body: JSON.stringify({
        content,
        sender: userStore.user.id,
        doctorId: selectedDoctor.value.id
      })
    })

    if (!response.ok) throw new Error('Failed to send message')
    const message = await response.json()
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

// 组件挂载时加载数据
onMounted(() => {
  if (checkUserLogin()) {
    loadDoctor().then(() => {
      loadMessages()
    })
  }
})

// 监听用户状态变化
watch(() => userStore.user, (newUser) => {
  if (newUser) {
    loadDoctor().then(() => {
      loadMessages()
    })
  }
}, { immediate: true })
</script>

<style scoped>
.patient-chat {
  display: flex;
  height: 100vh;
  background-color: #f5f7fa;
  padding: 20px;
}

.chat-container {
  flex: 1;
  display: flex;
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

.doctor-info {
  flex: 1;
}

.doctor-name {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
}

.doctor-details {
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
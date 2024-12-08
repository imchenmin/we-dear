<template>
  <div class="patient-chat">
    <el-card class="chat-card">
      <template #header>
        <div class="card-header">
          <h2>患者聊天测试界面</h2>
        </div>
      </template>

      <!-- 选择患者 -->
      <el-select v-model="selectedPatientId" placeholder="选择患者" class="patient-select">
        <el-option
          v-for="patient in patients"
          :key="patient.id"
          :label="patient.name"
          :value="patient.id"
        />
      </el-select>

      <!-- 聊天记录 -->
      <div class="chat-messages" ref="chatContainer">
        <div
          v-for="(message, index) in currentPatient?.messages"
          :key="index"
          :class="['message', message.type]"
        >
          <div class="message-content">
            <el-avatar
              :size="32"
              :src="message.type === 'patient' ? currentPatient?.avatar : doctorAvatar"
            />
            <div class="message-bubble">
              <template v-if="message.contentType === 'text'">
                {{ message.content }}
              </template>
              <template v-else-if="message.contentType === 'image'">
                <el-image
                  :src="message.content"
                  :preview-src-list="[message.content]"
                  fit="cover"
                  class="message-image"
                />
              </template>
              <template v-else-if="message.contentType === 'audio'">
                <audio controls :src="message.content"></audio>
              </template>
            </div>
          </div>
          <div v-if="message.aiSuggestion" class="ai-suggestion">
            <el-alert
              title="AI建议回复"
              type="info"
              :closable="false"
            >
              {{ message.aiSuggestion }}
            </el-alert>
          </div>
        </div>
      </div>

      <!-- 输入区域 -->
      <div class="chat-input">
        <el-input
          v-model="newMessage"
          type="textarea"
          :rows="3"
          placeholder="输入问题..."
          :disabled="!selectedPatientId"
        />
        <div class="input-actions">
          <el-upload
            action="http://localhost:8080/api/patient/upload"
            :show-file-list="false"
            :on-success="handleImageUpload"
            :disabled="!selectedPatientId"
          >
            <el-button type="primary" :icon="Picture" :disabled="!selectedPatientId">图片</el-button>
          </el-upload>
          <el-button type="primary" :icon="Microphone" :disabled="!selectedPatientId">语音</el-button>
          <el-button 
            type="primary" 
            @click="sendMessage" 
            :disabled="!selectedPatientId || !newMessage.trim()"
          >
            发送
          </el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { Picture, Microphone } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'

const patients = ref([])
const selectedPatientId = ref('')
const newMessage = ref('')
const doctorAvatar = 'https://placeholder.co/100'
const chatContainer = ref(null)

// 获取当前选中的患者信息
const currentPatient = computed(() => {
  return patients.value.find(p => p.id === selectedPatientId.value)
})

// 加载患者列表
const loadPatients = async () => {
  try {
    const response = await axios.get('http://localhost:8080/api/patients')
    patients.value = response.data
  } catch (error) {
    console.error('Failed to load patients:', error)
    ElMessage.error('加载患者列表失败')
  }
}

// 加载患者消息
const loadPatientMessages = async () => {
  if (!selectedPatientId.value) return

  try {
    const response = await axios.get(`http://localhost:8080/api/patient/${selectedPatientId.value}/messages`)
    const patientIndex = patients.value.findIndex(p => p.id === selectedPatientId.value)
    if (patientIndex !== -1) {
      patients.value[patientIndex] = response.data.patient
    }
  } catch (error) {
    console.error('Failed to load messages:', error)
    ElMessage.error('加载消息失败')
  }
}

// 发送消息
const sendMessage = async () => {
  if (!selectedPatientId.value || !newMessage.value.trim()) return

  const message = {
    type: 'patient',
    contentType: 'text',
    content: newMessage.value,
    timestamp: new Date().toISOString()
  }

  try {
    const response = await axios.post(
      `http://localhost:8080/api/patient/${selectedPatientId.value}/question`,
      message
    )
    await loadPatientMessages()
    newMessage.value = ''
    scrollToBottom()
  } catch (error) {
    console.error('Failed to send message:', error)
    ElMessage.error('发送消息失败')
  }
}

// 处理图片上传
const handleImageUpload = async (response: any) => {
  if (!selectedPatientId.value) return

  const message = {
    type: 'patient',
    contentType: 'image',
    content: `http://localhost:8080${response.url}`,
    timestamp: new Date().toISOString()
  }

  try {
    await axios.post(
      `http://localhost:8080/api/patient/${selectedPatientId.value}/question`,
      message
    )
    await loadPatientMessages()
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
      chatContainer.value.scrollTop = chatContainer.value.scrollHeight
    }, 100)
  }
}

// 监听患者选择变化
watch(selectedPatientId, () => {
  loadPatientMessages()
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
  padding: 20px;
}

.chat-card {
  min-height: 600px;
  display: flex;
  flex-direction: column;
}

.card-header {
  text-align: center;
}

.patient-select {
  width: 100%;
  margin-bottom: 20px;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  min-height: 300px;
  border: 1px solid #eee;
  border-radius: 4px;
  margin-bottom: 20px;
}

.message {
  margin-bottom: 20px;
}

.message-content {
  display: flex;
  gap: 10px;
}

.message.patient {
  flex-direction: row;
}

.message.doctor {
  flex-direction: row-reverse;
}

.message-bubble {
  background: #f0f2f5;
  padding: 10px;
  border-radius: 8px;
  max-width: 70%;
}

.message-image {
  max-width: 200px;
  border-radius: 4px;
}

.ai-suggestion {
  margin-top: 10px;
  margin-left: 42px;
}

.chat-input {
  margin-top: auto;
}

.input-actions {
  display: flex;
  gap: 10px;
  margin-top: 10px;
  justify-content: flex-end;
}
</style> 
<template>
  <div class="doctor-view">
    <el-container>
      <!-- 侧边栏 -->
      <el-aside width="300px">
        <div class="patient-list">
          <el-input
            v-model="searchQuery"
            placeholder="搜索患者"
            prefix-icon="Search"
          />
          <el-scrollbar height="calc(100vh - 60px)">
            <el-menu
              :default-active="activePatientId"
              @select="handlePatientSelect"
            >
              <el-menu-item
                v-for="patient in filteredPatients"
                :key="patient.id"
                :index="patient.id"
              >
                <el-avatar :size="32" :src="patient.avatar" />
                <span class="patient-name">{{ patient.name }}</span>
              </el-menu-item>
            </el-menu>
          </el-scrollbar>
        </div>
      </el-aside>

      <!-- 主要内容区域 -->
      <el-container>
        <el-main>
          <div v-if="selectedPatient" class="main-content">
            <!-- 患者档案 -->
            <el-card class="patient-profile">
              <template #header>
                <div class="card-header">
                  <span>患者档案</span>
                  <el-button type="primary">编辑档案</el-button>
                </div>
              </template>
              <el-descriptions :column="2" border>
                <el-descriptions-item label="姓名">{{ selectedPatient.name }}</el-descriptions-item>
                <el-descriptions-item label="性别">{{ selectedPatient.gender }}</el-descriptions-item>
                <el-descriptions-item label="年龄">{{ selectedPatient.age }}</el-descriptions-item>
                <el-descriptions-item label="联系电话">{{ selectedPatient.phone }}</el-descriptions-item>
                <el-descriptions-item label="诊断">{{ selectedPatient.diagnosis }}</el-descriptions-item>
                <el-descriptions-item label="主治医生">{{ selectedPatient.doctor }}</el-descriptions-item>
              </el-descriptions>
            </el-card>

            <!-- 聊天区域 -->
            <el-card class="chat-area">
              <template #header>
                <div class="card-header">
                  <span>聊天记录</span>
                </div>
              </template>
              <div class="chat-messages" ref="chatContainer">
                <div
                  v-for="(message, index) in selectedPatient.messages"
                  :key="index"
                  :class="['message', message.type]"
                >
                  <div class="message-content">
                    <el-avatar
                      :size="32"
                      :src="message.type === 'patient' ? selectedPatient.avatar : doctorAvatar"
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
              <div class="chat-input">
                <el-input
                  v-model="newMessage"
                  type="textarea"
                  :rows="3"
                  placeholder="输入回复内容..."
                />
                <div class="input-actions">
                  <el-upload
                    action="/api/upload"
                    :show-file-list="false"
                    :on-success="handleImageUpload"
                  >
                    <el-button type="primary" :icon="Picture">图片</el-button>
                  </el-upload>
                  <el-button type="primary" :icon="Microphone">语音</el-button>
                  <el-button type="primary" @click="sendMessage">发送</el-button>
                </div>
              </div>
            </el-card>
          </div>
          <div v-else class="no-patient-selected">
            <el-empty description="请选择一个患者" />
          </div>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onUnmounted } from 'vue'
import { Picture, Microphone } from '@element-plus/icons-vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'

// 模拟数据
const patients = ref([
  {
    id: '1',
    name: '张三',
    gender: '男',
    age: 45,
    phone: '13800138000',
    diagnosis: '2型糖尿病',
    doctor: '李医生',
    avatar: 'https://placeholder.co/100',
    messages: [
      {
        type: 'patient',
        contentType: 'text',
        content: '医生，我最近血糖有点高，该怎么调整？',
        timestamp: '2024-01-20 10:00:00',
        aiSuggestion: '建议回复：建议您注意以下几点：1. 严格控制饮食，减少碳水摄入；2. 坚持运动，每天步行30分钟；3. 按时服药；4. 定期监测血糖。'
      }
    ]
  }
])

const doctorAvatar = 'https://placeholder.co/100'
const searchQuery = ref('')
const activePatientId = ref('')
const newMessage = ref('')
const refreshInterval = ref<number | null>(null)

const filteredPatients = computed(() => {
  return patients.value.filter(patient =>
    patient.name.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

const selectedPatient = computed(() => {
  return patients.value.find(patient => patient.id === activePatientId.value)
})

// 获取患者信息
const loadPatientData = async (id: string) => {
  try {
    const response = await axios.get(`http://localhost:8080/api/patients/${id}`)
    const patientIndex = patients.value.findIndex(p => p.id === id)
    if (patientIndex !== -1) {
      patients.value[patientIndex] = response.data
    }
  } catch (error) {
    console.error('Failed to load patient data:', error)
    ElMessage.error('加载患者数据失败')
  }
}

const handlePatientSelect = async (id: string) => {
  activePatientId.value = id
  await loadPatientData(id)
  
  // 清除之前的定时器
  if (refreshInterval.value) {
    clearInterval(refreshInterval.value)
  }
  
  // 设置新的定时器，每5秒刷新一次
  refreshInterval.value = setInterval(() => {
    if (activePatientId.value) {
      loadPatientData(activePatientId.value)
    }
  }, 5000) as unknown as number
}

// 组件卸载时清除定时器
onUnmounted(() => {
  if (refreshInterval.value) {
    clearInterval(refreshInterval.value)
  }
})

const handleImageUpload = (response: any) => {
  // 处理图片上传成功后的逻辑
  console.log('Image uploaded:', response)
}

const sendMessage = async () => {
  if (!newMessage.value.trim()) return

  if (selectedPatient.value) {
    const messageData = {
      type: 'doctor',
      contentType: 'text',
      content: newMessage.value,
      timestamp: new Date().toISOString()
    }

    try {
      // 发送消息到后端
      await axios.post(`http://localhost:8080/api/patients/${selectedPatient.value.id}/messages`, messageData)
      
      // 成功后清空输入
      newMessage.value = ''
      
      // 立即刷新数据以显示新消息
      await loadPatientData(selectedPatient.value.id)
    } catch (error) {
      console.error('Failed to send message:', error)
      ElMessage.error('发送消息失败')
    }
  }
}
</script>

<style scoped>
.doctor-view {
  height: 100vh;
}

.patient-list {
  padding: 10px;
}

.patient-name {
  margin-left: 10px;
}

.main-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
  height: calc(100vh - 40px);
}

.patient-profile {
  flex: 0 0 auto;
}

.chat-area {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
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
  padding: 20px;
  border-top: 1px solid #eee;
}

.input-actions {
  display: flex;
  gap: 10px;
  margin-top: 10px;
  justify-content: flex-end;
}

.no-patient-selected {
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}
</style> 
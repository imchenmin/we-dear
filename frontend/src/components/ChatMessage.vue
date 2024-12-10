<template>
  <div class="message" :class="{ 'message-doctor': message.role === 'doctor', 'message-patient': message.role === 'patient' }">
    <el-avatar
      :size="40"
      :src="message.role === 'doctor' ? '/doctor-avatar.png' : '/patient-avatar.png'"
      class="message-avatar"
    />
    <div class="message-content">
      <div class="message-header">
        <span class="message-role">{{ message.role === 'doctor' ? '医生' : '患者' }}</span>
        <span class="message-time">{{ formatTime(message.createdAt) }}</span>
      </div>
      <div class="message-text">{{ message.content }}</div>
      
      <!-- AI建议，仅在医生视图且是患者消息时显示 -->
      <div v-if="showAISuggestion && message.role === 'patient'" class="ai-suggestion">
        <el-alert
          v-if="aiSuggestion"
          type="info"
          :closable="false"
          show-icon
        >
          <template #title>
            <div class="ai-suggestion-title">
              <el-icon><ChatLineRound /></el-icon>
              <span>AI建议回复</span>
              <el-tag size="small" :type="getSuggestionPriorityType(aiSuggestion.priority)">
                {{ getSuggestionPriorityText(aiSuggestion.priority) }}
              </el-tag>
            </div>
          </template>
          <div class="ai-suggestion-content">{{ aiSuggestion.content }}</div>
        </el-alert>
        <el-alert
          v-else-if="isLoadingAI"
          type="info"
          :closable="false"
          show-icon
        >
          <template #title>
            <div class="ai-suggestion-title">
              <el-icon><Loading /></el-icon>
              <span>正在生成 AI 建议...</span>
            </div>
          </template>
        </el-alert>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ChatLineRound, Loading } from '@element-plus/icons-vue'
import type { Message, AISuggestion } from '@/types'
import axios from 'axios'

const props = defineProps<{
  message: Message
  showAISuggestion?: boolean
}>()

const aiSuggestion = ref<AISuggestion | null>(null)
const isLoadingAI = ref(false)

const formatTime = (timestamp: string) => {
  return new Date(timestamp).toLocaleString()
}

const getSuggestionPriorityType = (priority: number) => {
  switch (priority) {
    case 5: return 'danger'    // 危急
    case 4: return 'warning'   // 紧急
    case 3: return ''          // 普通
    case 2: return 'info'      // 低优先级
    case 1: return 'success'   // 最低优先级
    default: return 'info'
  }
}

const getSuggestionPriorityText = (priority: number) => {
  switch (priority) {
    case 5: return '危急'
    case 4: return '紧急'
    case 3: return '普通'
    case 2: return '低优先级'
    case 1: return '最低优先级'
    default: return '未知'
  }
}

// 加载 AI 建议
const loadAISuggestion = async () => {
  if (!props.showAISuggestion || props.message.role !== 'patient') {
    return
  }

  isLoadingAI.value = true
  try {
    const response = await axios.get(`/api/chat/${props.message.patientId}/suggestions?messageId=${props.message.id}`)
    const suggestions = response.data
    if (suggestions && suggestions.length > 0) {
      aiSuggestion.value = suggestions[0]
    }
  } catch (error) {
    console.error('Failed to load AI suggestion:', error)
  } finally {
    isLoadingAI.value = false
  }
}

// 组件挂载时加载 AI 建议
onMounted(() => {
  loadAISuggestion()
})
</script>

<style scoped>
.message {
  display: flex;
  margin-bottom: 20px;
  gap: 12px;
}

.message-doctor {
  flex-direction: row;
}

.message-patient {
  flex-direction: row-reverse;
}

.message-content {
  max-width: 70%;
}

.message-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 4px;
  font-size: 12px;
  color: #666;
}

.message-text {
  padding: 12px 16px;
  background: #f4f4f5;
  border-radius: 8px;
  word-break: break-word;
  line-height: 1.4;
}

.message-doctor .message-text {
  background: #ecf5ff;
}

.ai-suggestion {
  margin-top: 12px;
}

.ai-suggestion-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.ai-suggestion-content {
  margin-top: 8px;
  white-space: pre-line;
  color: #606266;
  font-size: 14px;
}

:deep(.el-alert) {
  border-radius: 8px;
}

:deep(.el-alert.is-light) {
  background-color: #f0f9ff;
  border: 1px solid #e0f2fe;
}
</style> 
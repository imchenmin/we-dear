<template>
  <div class="message" :class="{ 'message-right': isDoctor }">
    <el-avatar
      :size="40"
      :src="isDoctor ? '/doctor-avatar.png' : message.avatar"
      class="message-avatar"
    />
    <div class="message-content">
      <div class="message-info">
        <span class="message-sender">{{ isDoctor ? '医生' : message.sender }}</span>
        <span class="message-time">{{ formatTime(message.timestamp) }}</span>
      </div>
      <div class="message-bubble">
        {{ message.content }}
      </div>
      <!-- AI建议，仅在医生视图且是患者消息时显示 -->
      <div v-if="showAISuggestion && !isDoctor" class="ai-suggestion">
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
import { ref, computed, onMounted } from 'vue'
import { ChatLineRound, Loading } from '@element-plus/icons-vue'
import { patientApi } from '@/api/patient'
import type { Message } from '@/types'

const props = defineProps<{
  message: Message
  showAISuggestion?: boolean // 是否显示AI建议（仅在医生视图中为true）
  patientId?: string
}>()

const isDoctor = computed(() => props.message.role === 'doctor')
const aiSuggestion = ref<any>(null)
const isLoadingAI = ref(false)

const formatTime = (timestamp: number) => {
  const date = new Date(timestamp)
  return date.toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 加载 AI 建议
const loadAISuggestion = async () => {
  if (!props.showAISuggestion || !props.patientId || isDoctor.value) {
    return
  }

  isLoadingAI.value = true
  try {
    const suggestions = await patientApi.getAISuggestions(props.patientId, props.message.id)
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

.message-right {
  flex-direction: row-reverse;
}

.message-content {
  max-width: 70%;
}

.message-info {
  margin-bottom: 4px;
  font-size: 12px;
}

.message-sender {
  font-weight: 500;
  color: #606266;
}

.message-time {
  margin-left: 8px;
  color: #909399;
}

.message-bubble {
  padding: 12px 16px;
  background: #f4f4f5;
  border-radius: 8px;
  word-break: break-word;
  line-height: 1.4;
}

.message-right .message-bubble {
  background: #ecf5ff;
}

.message-right .message-info {
  text-align: right;
}

.ai-suggestion {
  margin-top: 12px;
}

.ai-suggestion-title {
  display: flex;
  align-items: center;
  gap: 4px;
  font-weight: 500;
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
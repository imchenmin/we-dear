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
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Message } from '@/types'

const props = defineProps<{
  message: Message
}>()

const isDoctor = computed(() => props.message.role === 'doctor')

const formatTime = (timestamp: number) => {
  const date = new Date(timestamp)
  return date.toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit'
  })
}
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
</style> 
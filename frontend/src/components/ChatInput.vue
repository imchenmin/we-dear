<template>
  <div class="chat-input">
    <el-input
      v-model="messageText"
      type="textarea"
      :rows="3"
      placeholder="输入消息..."
      resize="none"
      :disabled="disabled"
      @keydown.enter.prevent="handleSend"
    />
    <div class="input-actions">
      <el-button 
        type="primary" 
        :disabled="disabled || !messageText.trim()" 
        @click="handleSend"
      >
        发送
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

defineProps<{
  disabled?: boolean
}>()

const emit = defineEmits<{
  (e: 'send', message: string): void
}>()

const messageText = ref('')

const handleSend = () => {
  const text = messageText.value.trim()
  if (text) {
    emit('send', text)
    messageText.value = ''
  }
}
</script>

<style scoped>
.chat-input {
  padding: 16px;
  border-top: 1px solid #dcdfe6;
  background: #fff;
}

.input-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 12px;
}

:deep(.el-textarea__inner) {
  resize: none;
  border-radius: 8px;
}

:deep(.el-textarea.is-disabled .el-textarea__inner) {
  background-color: #f5f7fa;
  border-color: #e4e7ed;
  color: #909399;
  cursor: not-allowed;
}
</style> 
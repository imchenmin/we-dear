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
      
      <!-- AI建议区域 -->
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
          <div class="ai-suggestion-meta">
            <span>置信度: {{ (aiSuggestion.confidence * 100).toFixed(0) }}%</span>
            <span>类别: {{ getSuggestionCategory(aiSuggestion.category) }}</span>
            <span>状态: {{ getSuggestionStatus(aiSuggestion.status) }}</span>
          </div>
          
          <!-- AI建议评价区域 -->
          <div class="suggestion-feedback">
            <div class="feedback-stats" v-if="feedbackStats">
              <span class="likes">
                <el-icon :color="userFeedback?.rating === 1 ? '#409EFF' : '#909399'">
                  <ArrowUpBold />
                </el-icon>
                {{ feedbackStats.likes }}
              </span>
              <span class="dislikes">
                <el-icon :color="userFeedback?.rating === -1 ? '#F56C6C' : '#909399'">
                  <ArrowDownBold />
                </el-icon>
                {{ feedbackStats.dislikes }}
              </span>
            </div>
            <div class="feedback-actions">
              <el-button
                type="text"
                size="small"
                :class="{ active: userFeedback?.rating === 1 }"
                @click="handleFeedback(1)"
              >
                <el-icon><ArrowUpBold /></el-icon>
                有帮助
              </el-button>
              <el-button
                type="text"
                size="small"
                :class="{ active: userFeedback?.rating === -1 }"
                @click="handleFeedback(-1)"
              >
                <el-icon><ArrowDownBold /></el-icon>
                没帮助
              </el-button>
              <el-button
                type="text"
                size="small"
                @click="showCommentDialog = true"
              >
                <el-icon><ChatDotRound /></el-icon>
                评论
              </el-button>
            </div>
          </div>
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

    <!-- 评论对话框 -->
    <el-dialog
      v-model="showCommentDialog"
      title="评价AI建议"
      width="500px"
    >
      <el-form
        ref="feedbackFormRef"
        :model="feedbackForm"
        label-width="80px"
      >
        <el-form-item label="评价">
          <el-radio-group v-model="feedbackForm.rating">
            <el-radio :label="1">有帮助</el-radio>
            <el-radio :label="-1">没帮助</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="标签">
          <el-select
            v-model="feedbackForm.tags"
            multiple
            placeholder="选择标签"
          >
            <el-option label="事实错误" value="fact_error" />
            <el-option label="逻辑错误" value="logic_error" />
            <el-option label="全面专业" value="professional" />
            <el-option label="易懂" value="easy" />
          </el-select>
        </el-form-item>
        <el-form-item label="评论">
          <el-input
            v-model="feedbackForm.comment"
            type="textarea"
            :rows="4"
            placeholder="请输入您对AI建议的评价"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showCommentDialog = false">取消</el-button>
          <el-button type="primary" @click="submitFeedback">提交</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {
  ChatLineRound,
  Loading,
  ArrowDownBold,
  ArrowUpBold,
  ChatDotRound
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { Message, AISuggestion, MessageFeedback } from '@/types'
import { AISuggestionFeedbackTag, AISuggestionFeedbackStatus } from '@/types'
import { patientApi } from '@/api/patient'

const props = defineProps<{
  message: Message
  showAISuggestion?: boolean
}>()

const aiSuggestion = ref<AISuggestion | null>(null)
const isLoadingAI = ref(false)
const showCommentDialog = ref(false)
const userFeedback = ref<MessageFeedback | null>(null)
const feedbackStats = ref<{ likes: number; dislikes: number } | null>(null)

const feedbackForm = ref({
  rating: 1,
  tags: [] as string[],
  comment: ''
})

// 加载评价统计信息
const loadFeedbackStats = async () => {
  if (!aiSuggestion.value) return

  try {
    feedbackStats.value = await patientApi.getMessageFeedbackStats(aiSuggestion.value.id)
  } catch (error) {
    console.error('Failed to load feedback stats:', error)
  }
}

// 加载用户的评价
const loadUserFeedback = async () => {
  if (!aiSuggestion.value) return

  try {
    const feedbacks = await patientApi.getMessageFeedbacks(aiSuggestion.value.id)
    if (feedbacks.length > 0) {
      userFeedback.value = feedbacks[0]
    }
  } catch (error) {
    console.error('Failed to load user feedback:', error)
  }
}

// 处理快速评价
const handleFeedback = async (rating: number) => {
  if (!aiSuggestion.value) return

  try {
    if (userFeedback.value) {
      // 如果已经评价过，更新评价
      await patientApi.updateMessageFeedback(userFeedback.value.id, {
        rating,
        suggestionId: aiSuggestion.value.id
      })
    } else {
      // 创建新评价
      await patientApi.createMessageFeedback(aiSuggestion.value.id, {
        rating,
        suggestionId: aiSuggestion.value.id,
        patientId: props.message.patientId,
        status: AISuggestionFeedbackStatus.Pending
      })
    }
    await loadFeedbackStats()
    await loadUserFeedback()
    ElMessage.success('评价成功')
  } catch (error) {
    console.error('Failed to submit feedback:', error)
    ElMessage.error('评价失败')
  }
}

// 提交详细评价
const submitFeedback = async () => {
  if (!aiSuggestion.value) return

  try {
    const feedbackData = {
      ...feedbackForm.value,
      suggestionId: aiSuggestion.value.id,
      patientId: props.message.patientId,
      status: AISuggestionFeedbackStatus.Pending
    }

    if (userFeedback.value) {
      // 更新评价
      await patientApi.updateMessageFeedback(userFeedback.value.id, feedbackData)
    } else {
      // 创建新评价
      await patientApi.createMessageFeedback(aiSuggestion.value.id, feedbackData)
    }
    showCommentDialog.value = false
    await loadFeedbackStats()
    await loadUserFeedback()
    ElMessage.success('评价提交成功')
  } catch (error) {
    console.error('Failed to submit feedback:', error)
    ElMessage.error('评价提交失败')
  }
}

// 格式化时间
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

const getSuggestionCategory = (category: string) => {
  const categories: Record<string, string> = {
    medication: '用药建议',
    visit: '就医建议',
    lifestyle: '生活建议',
    urgent: '紧急建议'
  }
  return categories[category] || category
}

const getSuggestionStatus = (status: string) => {
  const statuses: Record<string, string> = {
    pending: '待审核',
    approved: '已采纳',
    rejected: '已拒绝'
  }
  return statuses[status] || status
}

// 加载 AI 建议
const loadAISuggestion = async () => {
  if (!props.showAISuggestion || props.message.role !== 'patient') {
    return
  }

  isLoadingAI.value = true
  try {
    const suggestions = await patientApi.getAISuggestions(props.message.patientId, props.message.id)
    if (suggestions && suggestions.length > 0) {
      aiSuggestion.value = suggestions[0]
    }
  } catch (error) {
    console.error('Failed to load AI suggestion:', error)
  } finally {
    isLoadingAI.value = false
  }
}

// 组件挂载时加载数据
onMounted(() => {
  if (props.showAISuggestion && props.message.role === 'patient') {
    loadAISuggestion().then(() => {
      if (aiSuggestion.value) {
        loadFeedbackStats()
        loadUserFeedback()
      }
    })
  }
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

.ai-suggestion-meta {
  margin-top: 8px;
  display: flex;
  gap: 16px;
  color: #909399;
  font-size: 12px;
}

:deep(.el-alert) {
  border-radius: 8px;
}

:deep(.el-alert.is-light) {
  background-color: #f0f9ff;
  border: 1px solid #e0f2fe;
}

.message-feedback {
  margin-top: 8px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: #909399;
}

.feedback-stats {
  display: flex;
  gap: 16px;
}

.feedback-stats .likes,
.feedback-stats .dislikes {
  display: flex;
  align-items: center;
  gap: 4px;
}

.feedback-actions {
  display: flex;
  gap: 16px;
}

.feedback-actions .el-button {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
}

.feedback-actions .el-button.active {
  color: #409EFF;
}

.feedback-actions .el-button:hover {
  background-color: #f5f7fa;
}

:deep(.el-dialog__body) {
  padding-top: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.suggestion-feedback {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #e0f2fe;
}
</style> 
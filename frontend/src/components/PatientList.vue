<template>
  <div class="patient-list">
    <div class="search-bar">
      <el-input
        v-model="searchText"
        placeholder="搜索患者"
        prefix-icon="Search"
      />
    </div>
    
    <div class="list-container">
      <el-scrollbar>
        <div
          v-for="patient in filteredPatients"
          :key="patient.patientId"
          class="patient-item"
          :class="{ active: modelValue === patient.patientId }"
          @click="$emit('update:modelValue', patient.patientId)"
        >
          <el-avatar :size="40" :src="patient.patientAvatar" />
          <div class="patient-info">
            <div class="patient-name">{{ patient.patientName }}</div>
            <div class="last-message">{{ patient.lastMessage }}</div>
          </div>
          <div class="meta-info">
            <div class="time">{{ formatTime(patient.lastMessageAt) }}</div>
            <el-badge
              v-if="patient.unreadCount > 0"
              :value="patient.unreadCount"
              class="unread-badge"
            />
          </div>
        </div>
      </el-scrollbar>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Search } from '@element-plus/icons-vue'
import { request } from '@/utils/request'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const props = defineProps<{
  modelValue: string
}>()

defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const patients = ref<any[]>([])
const searchText = ref('')

const filteredPatients = computed(() => {
  if (!searchText.value) return patients.value
  
  const search = searchText.value.toLowerCase()
  return patients.value.filter(patient => 
    patient.patientName.toLowerCase().includes(search)
  )
})

const formatTime = (time: string) => {
  return dayjs(time).fromNow()
}

const loadPatients = async () => {
  try {
    const data = await request('chat/list')
    patients.value = data
  } catch (error) {
    console.error('Failed to load patients:', error)
  }
}

// 定期刷新聊天列表
let refreshInterval: number
onMounted(() => {
  loadPatients()
  refreshInterval = setInterval(loadPatients, 30000) // 每30秒刷新一次
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>

<style scoped>
.patient-list {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.search-bar {
  padding: 16px;
}

.list-container {
  flex: 1;
  overflow: hidden;
}

.patient-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.patient-item:hover {
  background-color: #f5f7fa;
}

.patient-item.active {
  background-color: #ecf5ff;
}

.patient-info {
  flex: 1;
  margin-left: 12px;
  overflow: hidden;
}

.patient-name {
  font-weight: 500;
  margin-bottom: 4px;
}

.last-message {
  font-size: 13px;
  color: #909399;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.meta-info {
  text-align: right;
  margin-left: 12px;
}

.time {
  font-size: 12px;
  color: #909399;
  margin-bottom: 4px;
}

.unread-badge {
  margin-left: auto;
}
</style> 
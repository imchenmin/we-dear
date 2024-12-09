<template>
  <div class="patient-list">
    <el-input
      v-model="searchQuery"
      placeholder="搜索患者..."
      prefix-icon="Search"
      clearable
      class="search-input"
    />
    
    <el-scrollbar height="calc(100vh - 120px)">
      <div
        v-for="patient in filteredPatients"
        :key="patient.id"
        class="patient-item"
        :class="{ active: patient.id === modelValue }"
        @click="$emit('update:modelValue', patient.id)"
      >
        <el-avatar :size="40" :src="patient.avatar" />
        <div class="patient-info">
          <div class="patient-name">{{ patient.name }}</div>
          <div class="patient-diagnosis">{{ patient.diagnosis }}</div>
        </div>
      </div>
    </el-scrollbar>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Patient } from '@/types'

const props = defineProps<{
  patients: Patient[]
  modelValue: string
}>()

defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const searchQuery = ref('')

const filteredPatients = computed(() => {
  if (!searchQuery.value) return props.patients
  const query = searchQuery.value.toLowerCase()
  return props.patients.filter(patient => 
    patient.name.toLowerCase().includes(query) ||
    patient.diagnosis.toLowerCase().includes(query)
  )
})
</script>

<style scoped>
.patient-list {
  height: 100%;
  border-right: 1px solid #dcdfe6;
  padding: 16px;
}

.search-input {
  margin-bottom: 16px;
}

.patient-item {
  display: flex;
  align-items: center;
  padding: 12px;
  cursor: pointer;
  border-radius: 8px;
  margin-bottom: 8px;
  transition: all 0.3s ease;
}

.patient-item:hover {
  background-color: #f5f7fa;
}

.patient-item.active {
  background-color: #ecf5ff;
}

.patient-info {
  margin-left: 12px;
  flex: 1;
}

.patient-name {
  font-weight: 500;
  color: #303133;
}

.patient-diagnosis {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}
</style> 
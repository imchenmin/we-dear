<template>
  <div class="patient-profile" v-if="patient">
    <div class="profile-header">
      <el-avatar :size="64" :src="patient.avatar" />
      <h2>{{ patient.name }}</h2>
    </div>
    
    <div class="profile-info">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="ID">{{ patient.id }}</el-descriptions-item>
        <el-descriptions-item label="性别">{{ patient.gender }}</el-descriptions-item>
        <el-descriptions-item label="年龄">{{ patient.age }}岁</el-descriptions-item>
        <el-descriptions-item label="电话">{{ patient.phone }}</el-descriptions-item>
        <el-descriptions-item label="诊断">{{ patient.diagnosis }}</el-descriptions-item>
        <el-descriptions-item label="主治医生">{{ patient.doctor }}</el-descriptions-item>
      </el-descriptions>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Patient } from '@/types'

defineProps<{
  patient?: Patient
}>()
</script>

<style scoped>
.patient-profile {
  padding: 20px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
}

.profile-header h2 {
  margin: 0;
  color: #303133;
}

.profile-info {
  margin-top: 16px;
}
</style> 
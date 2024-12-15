<template>
  <div class="main-layout">
    <el-container>
      <!-- 顶部导航栏 -->
      <el-header>
        <el-menu
          mode="horizontal"
          :router="true"
          class="main-menu"
        >
          <el-menu-item index="/">
            <el-icon><Monitor /></el-icon>
            医生工作台
          </el-menu-item>
          <el-menu-item index="/patient-chat">
            <el-icon><ChatDotRound /></el-icon>
            患者聊天（调试端）
          </el-menu-item>
          <el-menu-item index="/departments">
            <el-icon><Monitor/></el-icon>
            科室管理
          </el-menu-item>
          
          <div class="menu-right">
            <el-dropdown @command="handleCommand">
              <el-avatar :size="32" :src="userStore.user?.avatar">
                {{ userStore.user?.name?.[0] }}
              </el-avatar>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <el-button type="primary" @click="showAddDoctor">
              <el-icon><Plus /></el-icon>添加医生
            </el-button>
            <el-button type="success" @click="showAddPatient">
              <el-icon><Plus /></el-icon>添加患者
            </el-button>
          </div>
        </el-menu>
      </el-header>

      <!-- 主要内容区域 -->
      <el-main>
        <router-view></router-view>
      </el-main>
    </el-container>

    <!-- 添加医生弹窗 -->
    <el-dialog
      v-model="doctorDialogVisible"
      title="添加新医生"
      width="50%"
      :before-close="handleClose"
    >
      <DoctorForm @success="handleDoctorAdded" />
    </el-dialog>

    <!-- 添加患者弹窗 -->
    <el-dialog
      v-model="patientDialogVisible"
      title="添加新患者"
      width="50%"
      :before-close="handleClose"
    >
      <PatientForm @success="handlePatientAdded" />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Monitor, ChatDotRound, Plus } from '@element-plus/icons-vue'
import { ElMessageBox } from 'element-plus'
import DoctorForm from '@/components/DoctorForm.vue'
import PatientForm from '@/components/PatientForm.vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const doctorDialogVisible = ref(false)
const patientDialogVisible = ref(false)
const router = useRouter()
const userStore = useUserStore()

const showAddDoctor = () => {
  doctorDialogVisible.value = true
}

const showAddPatient = () => {
  patientDialogVisible.value = true
}

const handleClose = (done: () => void) => {
  ElMessageBox.confirm('确认关闭？')
    .then(() => {
      done()
    })
    .catch(() => {})
}

const handleDoctorAdded = () => {
  doctorDialogVisible.value = false
  // 可以在这里刷新医生列表
}

const handlePatientAdded = () => {
  patientDialogVisible.value = false
  // 可以在这里刷新患者列表
}

const handleCommand = (command: string) => {
  if (command === 'logout') {
    userStore.logout()
    router.push('/login')
    ElMessage.success('已退出登录')
  }
}
</script>

<style scoped>
.main-layout {
  min-height: 100vh;
}

.el-header {
  padding: 0;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.12);
}

.main-menu {
  display: flex;
  align-items: center;
  height: 60px;
}

.menu-right {
  margin-left: auto;
  display: flex;
  gap: 12px;
  padding-right: 20px;
}

.el-main {
  padding: 0;
}

:deep(.el-dialog__body) {
  padding-top: 20px;
}
</style> 
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
          <!-- 左侧导航项 -->
          <div class="menu-left">
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
          </div>

          <!-- 右侧用户区域 -->
          <div class="menu-right">
            <!-- 系统管理下拉菜单 -->
            <el-dropdown v-if="userStore.isAdmin" trigger="click">
              <el-button type="primary">
                <el-icon><Setting /></el-icon>
                系统管理
                <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item index="/departments">
                    <el-icon><Monitor /></el-icon>科室管理
                  </el-dropdown-item>
                  <el-dropdown-item @click="showAddDoctor">
                    <el-icon><Plus /></el-icon>添加医生
                  </el-dropdown-item>
                  <el-dropdown-item @click="showAddPatient">
                    <el-icon><Plus /></el-icon>添加患者
                  </el-dropdown-item>
                  <el-dropdown-item @click="showDoctorPasswordChange">
                    <el-icon><Lock /></el-icon>修改医生密码
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>

            <!-- 用户信息和操作区 -->
            <el-dropdown trigger="click">
              <div class="user-info">
                <el-avatar :size="32" :src="userStore.user?.avatar">
                  {{ userStore.user?.name?.[0] }}
                </el-avatar>
                <span class="user-name">{{ userStore.user?.name }}</span>
                <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </div>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile" @click="handleCommand('profile')">
                    <el-icon><User /></el-icon>个人资料
                  </el-dropdown-item>
                  <el-dropdown-item command="changePassword" @click="handleCommand('changePassword')">
                    <el-icon><Lock /></el-icon>修改密码
                  </el-dropdown-item>
                  <el-dropdown-item divided command="logout" @click="handleCommand('logout')">
                    <el-icon><SwitchButton /></el-icon>退出登录
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
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

    <!-- 修改密码弹窗 -->
    <el-dialog
      v-model="passwordDialogVisible"
      title="修改密码"
      width="30%"
    >
      <el-form
        ref="passwordFormRef"
        :model="passwordForm"
        :rules="passwordRules"
        label-width="100px"
      >
        <el-form-item label="原密码" prop="oldPassword">
          <el-input
            v-model="passwordForm.oldPassword"
            type="password"
            show-password
            placeholder="请输入原密码"
          />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input
            v-model="passwordForm.newPassword"
            type="password"
            show-password
            placeholder="请输入新密码"
          />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="passwordForm.confirmPassword"
            type="password"
            show-password
            placeholder="请再次输入新密码"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="passwordDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitPasswordChange">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 修改医生密码弹窗 -->
    <el-dialog
      v-model="doctorPasswordDialogVisible"
      title="修改医生密码"
      width="30%"
    >
      <el-form
        ref="doctorPasswordFormRef"
        :model="doctorPasswordForm"
        :rules="doctorPasswordRules"
        label-width="100px"
      >
        <el-form-item label="选择医生" prop="userId">
          <el-select v-model="doctorPasswordForm.userId" placeholder="请选择医生">
            <el-option
              v-for="doctor in doctors"
              :key="doctor.id"
              :label="doctor.name"
              :value="doctor.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input
            v-model="doctorPasswordForm.newPassword"
            type="password"
            show-password
            placeholder="请输入新密码"
          />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="doctorPasswordForm.confirmPassword"
            type="password"
            show-password
            placeholder="请再次输入新密码"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="doctorPasswordDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitDoctorPasswordChange">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  ChatDotRound,
  Plus,
  User,
  Lock,
  SwitchButton,
  ArrowDown,
  Setting,
  Monitor
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { useUserStore } from '@/stores/user'
import DoctorForm from '@/components/DoctorForm.vue'
import PatientForm from '@/components/PatientForm.vue'
import { request } from '@/utils/request'
const router = useRouter()
const userStore = useUserStore()

const doctorDialogVisible = ref(false)
const patientDialogVisible = ref(false)
const passwordDialogVisible = ref(false)
const doctorPasswordDialogVisible = ref(false)

// 密码表单相关
const passwordFormRef = ref<FormInstance>()
const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const validatePass = (rule: any, value: string, callback: any) => {
  if (value === '') {
    callback(new Error('请再次输入密码'))
  } else if (value !== passwordForm.value.newPassword) {
    callback(new Error('两次输入密码不一致!'))
  } else {
    callback()
  }
}

const passwordRules: FormRules = {
  oldPassword: [
    { required: true, message: '请输入原密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能小于6位', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能小于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, validator: validatePass, trigger: 'blur' }
  ]
}

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

const handleCommand = async (command: string) => {
  switch (command) {
    case 'profile':
      // TODO: 跳转到个人资料页面或打开编辑弹窗
      ElMessage.info('功能开发中')
      break
    case 'changePassword':
      passwordDialogVisible.value = true
      break
    case 'logout':
      await userStore.logout()
      router.push('/login')
      ElMessage.success('已退出登录')
      break
  }
}

const submitPasswordChange = async () => {
  if (!passwordFormRef.value) return

  await passwordFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        await userStore.changePassword(
          passwordForm.value.oldPassword,
          passwordForm.value.newPassword
        )
        ElMessage.success('密码修改成功')
        passwordDialogVisible.value = false
        passwordForm.value = {
          oldPassword: '',
          newPassword: '',
          confirmPassword: ''
        }
      } catch (error) {
        ElMessage.error('密码修改失败')
      }
    }
  })
}

// 添加相关的响应式数据和方法
const doctorPasswordFormRef = ref<FormInstance>()
const doctors = ref<Doctor[]>([])
const doctorPasswordForm = ref({
  userId: '',
  newPassword: '',
  confirmPassword: ''
})

// 加载医生列表
const loadDoctors = async () => {
  try {
    doctors.value = await request.get('/doctors')
  } catch (error) {
    ElMessage.error('加载医生列表失败')
  }
}

// 显示修改医生密码弹窗
const showDoctorPasswordChange = () => {
  doctorPasswordDialogVisible.value = true
  loadDoctors()
}

// 提交修改医生密码
const submitDoctorPasswordChange = async () => {
  if (!doctorPasswordFormRef.value) return

  await doctorPasswordFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        await userStore.changePassword(
          '', // 管理员不需要提供旧密码
          doctorPasswordForm.value.newPassword,
          doctorPasswordForm.value.userId
        )
        ElMessage.success('密码修改成功')
        doctorPasswordDialogVisible.value = false
        doctorPasswordForm.value = {
          userId: '',
          newPassword: '',
          confirmPassword: ''
        }
      } catch (error) {
        ElMessage.error('密码修改失败')
      }
    }
  })
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
  justify-content: space-between;
  height: 60px;
}

.menu-left {
  display: flex;
}

.menu-right {
  display: flex;
  align-items: center;
  gap: 16px;
  padding-right: 20px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.user-info:hover {
  background-color: #f5f7fa;
}

.user-name {
  font-size: 14px;
  color: #606266;
}

.el-main {
  padding: 0;
}

:deep(.el-dropdown-menu__item) {
  display: flex;
  align-items: center;
  gap: 8px;
}

:deep(.el-dialog__body) {
  padding-top: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style> 
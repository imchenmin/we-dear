<template>
  <div class="department-view">
    <div class="department-list">
      <el-card>
        <template #header>
          <div class="card-header">
            <span>科室列表</span>
            <el-button type="primary" @click="showAddDepartment">
              <el-icon><Plus /></el-icon>新增科室
            </el-button>
          </div>
        </template>
        
        <el-table :data="departments" style="width: 100%">
          <el-table-column prop="name" label="科室名称" />
          <el-table-column prop="code" label="科室编码" />
          <el-table-column prop="description" label="描述" />
          <el-table-column label="医生数量" width="100">
            <template #default="{ row }">
              {{ getDoctorCount(row.id) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200">
            <template #default="{ row }">
              <el-button-group>
                <el-button type="primary" @click="showDoctors(row)">
                  查看医生
                </el-button>
                <el-button type="warning" @click="editDepartment(row)">
                  编辑
                </el-button>
                <el-button type="danger" @click="deleteDepartment(row)">
                  删除
                </el-button>
              </el-button-group>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>

    <!-- 科室医生列表抽屉 -->
    <el-drawer
      v-model="doctorsDrawerVisible"
      :title="`${currentDepartment?.name || ''} - 医生列表`"
      size="60%"
    >
      <el-table :data="departmentDoctors" style="width: 100%">
        <el-table-column prop="name" label="姓名" />
        <el-table-column prop="title" label="职称" />
        <el-table-column prop="license" label="执业证号" />
        <el-table-column prop="specialty" label="专长" show-overflow-tooltip />
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
              {{ row.status === 'active' ? '在职' : '离职' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button-group>
              <el-button type="primary" @click="editDoctor(row)">
                编辑
              </el-button>
              <el-button type="danger" @click="deleteDoctor(row)">
                删除
              </el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>
    </el-drawer>

    <!-- 科室表单对话框 -->
    <el-dialog
      v-model="departmentDialogVisible"
      :title="editingDepartment ? '编辑科室' : '新增科室'"
      width="40%"
    >
      <el-form
        ref="departmentForm"
        :model="departmentForm"
        :rules="departmentRules"
        label-width="100px"
      >
        <el-form-item label="科室名称" prop="name">
          <el-input v-model="departmentForm.name" />
        </el-form-item>
        <el-form-item label="科室编码" prop="code">
          <el-input v-model="departmentForm.code" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="departmentForm.description"
            type="textarea"
            :rows="3"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="departmentDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveDepartment">
          确定
        </el-button>
      </template>
    </el-dialog>

    <!-- 医生编辑对话框 -->
    <el-dialog
      v-model="doctorDialogVisible"
      :title="editingDoctor ? '编辑医生信息' : '新增医生'"
      width="50%"
    >
      <DoctorForm
        :initial-data="editingDoctor"
        @success="handleDoctorSaved"
      />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import DoctorForm from '@/components/DoctorForm.vue'
import type { Department, Doctor } from '@/types'

const departments = ref<Department[]>([])
const doctors = ref<Doctor[]>([])
const currentDepartment = ref<Department | null>(null)
const departmentDoctors = computed(() => 
  doctors.value.filter(d => d.departmentId === currentDepartment.value?.id)
)

const doctorsDrawerVisible = ref(false)
const departmentDialogVisible = ref(false)
const doctorDialogVisible = ref(false)

const departmentForm = ref({
  name: '',
  code: '',
  description: ''
})

const editingDepartment = ref<Department | null>(null)
const editingDoctor = ref<Doctor | null>(null)

const departmentRules = {
  name: [{ required: true, message: '请输入科室名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入科室编码', trigger: 'blur' }]
}

// 获取科室列表
const loadDepartments = async () => {
  try {
    const response = await fetch('/api/departments')
    departments.value = await response.json()
  } catch (error) {
    ElMessage.error('加载科室列表失败')
  }
}

// 获取医生列表
const loadDoctors = async () => {
  try {
    const response = await fetch('/api/doctors')
    doctors.value = await response.json()
  } catch (error) {
    ElMessage.error('加载医生列表失败')
  }
}

// 获取科室医生数量
const getDoctorCount = (departmentId: string) => {
  return doctors.value.filter(d => d.departmentId === departmentId).length
}

// 显示科室医生列表
const showDoctors = (department: Department) => {
  currentDepartment.value = department
  doctorsDrawerVisible.value = true
}

// 新增科室
const showAddDepartment = () => {
  editingDepartment.value = null
  departmentForm.value = {
    name: '',
    code: '',
    description: ''
  }
  departmentDialogVisible.value = true
}

// 编辑科室
const editDepartment = (department: Department) => {
  editingDepartment.value = department
  departmentForm.value = { ...department }
  departmentDialogVisible.value = true
}

// 保存科室
const saveDepartment = async () => {
  try {
    const url = editingDepartment.value
      ? `/api/departments/${editingDepartment.value.id}`
      : '/api/departments'
    const method = editingDepartment.value ? 'PUT' : 'POST'
    
    const response = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(departmentForm.value)
    })

    if (response.ok) {
      ElMessage.success(`${editingDepartment.value ? '更新' : '创建'}科室成功`)
      departmentDialogVisible.value = false
      loadDepartments()
    } else {
      const error = await response.json()
      ElMessage.error(error.error)
    }
  } catch (error) {
    ElMessage.error('保存失败')
  }
}

// 删除科室
const deleteDepartment = async (department: Department) => {
  try {
    await ElMessageBox.confirm(
      '删除科室将同时删除该科室下的所有医生，是否继续？',
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const response = await fetch(`/api/departments/${department.id}`, {
      method: 'DELETE'
    })

    if (response.ok) {
      ElMessage.success('删除科室成功')
      loadDepartments()
      loadDoctors()
    } else {
      const error = await response.json()
      ElMessage.error(error.error)
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 编辑医生
const editDoctor = (doctor: Doctor) => {
  editingDoctor.value = doctor
  doctorDialogVisible.value = true
}

// 删除医生
const deleteDoctor = async (doctor: Doctor) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除该医生吗？',
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const response = await fetch(`/api/doctors/${doctor.id}`, {
      method: 'DELETE'
    })

    if (response.ok) {
      ElMessage.success('删除医生成功')
      loadDoctors()
    } else {
      const error = await response.json()
      ElMessage.error(error.error)
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 处理医生保存成功
const handleDoctorSaved = () => {
  doctorDialogVisible.value = false
  loadDoctors()
}

// 初始加载数据
loadDepartments()
loadDoctors()
</script>

<style scoped>
.department-view {
  padding: 20px;
}

.department-list {
  max-width: 1200px;
  margin: 0 auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

:deep(.el-drawer__body) {
  padding: 20px;
}
</style> 
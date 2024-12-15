<template>
  <el-form
    ref="formRef"
    :model="form"
    :rules="rules"
    label-width="100px"
  >
    <el-form-item label="头像">
      <el-upload
        class="avatar-uploader"
        action="/api/upload"
        :show-file-list="false"
        :on-success="handleAvatarSuccess"
        :before-upload="beforeAvatarUpload"
      >
        <img v-if="form.avatar" :src="form.avatar" class="avatar">
        <el-icon v-else><Plus /></el-icon>
      </el-upload>
    </el-form-item>

    <el-form-item label="姓名" prop="name">
      <el-input v-model="form.name" />
    </el-form-item>

    <el-form-item label="职称" prop="title">
      <el-input v-model="form.title" />
    </el-form-item>

    <el-form-item label="所属科室" prop="departmentId">
      <el-select v-model="form.departmentId" placeholder="请选择科室">
        <el-option
          v-for="dept in departments"
          :key="dept.id"
          :label="dept.name"
          :value="dept.id"
        />
      </el-select>
    </el-form-item>

    <el-form-item label="执业证号" prop="license">
      <el-input v-model="form.license" />
    </el-form-item>

    <el-form-item label="专长" prop="specialty">
      <el-input
        v-model="form.specialty"
        type="textarea"
        :rows="3"
      />
    </el-form-item>

    <el-form-item label="状态" prop="status">
      <el-select v-model="form.status">
        <el-option label="在职" value="active" />
        <el-option label="离职" value="inactive" />
      </el-select>
    </el-form-item>

    <el-form-item>
      <el-button type="primary" @click="submitForm">保存</el-button>
      <el-button @click="resetForm">重置</el-button>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import type { Doctor, Department } from '@/types'
import { request } from '@/utils/request'

const props = defineProps<{
  initialData?: Doctor
}>()

const emit = defineEmits<{
  (e: 'success'): void
}>()

const formRef = ref<FormInstance>()
const departments = ref<Department[]>([])

const form = ref({
  name: '',
  title: '',
  departmentId: '',
  license: '',
  specialty: '',
  avatar: '',
  status: 'active'
})

const rules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  title: [{ required: true, message: '请输入职称', trigger: 'blur' }],
  departmentId: [{ required: true, message: '请选择科室', trigger: 'change' }],
  license: [{ required: true, message: '请输入执业证号', trigger: 'blur' }]
}

// 初始化表单数据
onMounted(async () => {
  try {
    // 加载科室列表
    departments.value = await request.get('/departments')
    
    // 如果是编辑模式，填充表单数据
    if (props.initialData) {
      form.value = { ...props.initialData }
    }
  } catch (error) {
    ElMessage.error('加载数据失败')
  }
})

// 头像上传相关
const handleAvatarSuccess = (res: any) => {
  form.value.avatar = res.url
}

const beforeAvatarUpload = (file: File) => {
  const isJPG = file.type === 'image/jpeg' || file.type === 'image/png'
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isJPG) {
    ElMessage.error('头像只能是 JPG/PNG 格式!')
  }
  if (!isLt2M) {
    ElMessage.error('头像大小不能超过 2MB!')
  }
  return isJPG && isLt2M
}

// 提交表单
const submitForm = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        const url = props.initialData
          ? `/doctors/${props.initialData.id}`
          : '/doctors'
        const method = props.initialData ? 'put' : 'post'
        
        await request[method](url, form.value)
        ElMessage.success(`${props.initialData ? '更新' : '创建'}成功`)
        emit('success')
      } catch (error) {
        ElMessage.error(`${props.initialData ? '更新' : '创建'}失败`)
      }
    }
  })
}

// 重置表单
const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
}
</script>

<style scoped>
.avatar-uploader {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  width: 178px;
  height: 178px;
}

.avatar {
  width: 178px;
  height: 178px;
  display: block;
}

:deep(.el-upload) {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}
</style> 
<template>
  <div class="doctor-form">
    <h2>添加新医生</h2>
    <el-form :model="form" :rules="rules" ref="doctorForm" label-width="120px">
      <el-form-item label="姓名" prop="name">
        <el-input v-model="form.name"></el-input>
      </el-form-item>
      
      <el-form-item label="职称" prop="title">
        <el-input v-model="form.title"></el-input>
      </el-form-item>

      <el-form-item label="所属科室" prop="departmentId">
        <el-select v-model="form.departmentId" placeholder="请选择科室">
          <el-option
            v-for="dept in departments"
            :key="dept.id"
            :label="dept.name"
            :value="dept.id">
          </el-option>
        </el-select>
      </el-form-item>

      <el-form-item label="执业证号" prop="license">
        <el-input v-model="form.license"></el-input>
      </el-form-item>

      <el-form-item label="专长" prop="specialty">
        <el-input type="textarea" v-model="form.specialty"></el-input>
      </el-form-item>

      <el-form-item label="头像" prop="avatar">
        <el-upload
          class="avatar-uploader"
          action="/api/upload"
          :show-file-list="false"
          :on-success="handleAvatarSuccess">
          <img v-if="form.avatar" :src="form.avatar" class="avatar">
          <el-icon v-else><Plus /></el-icon>
        </el-upload>
      </el-form-item>

      <el-form-item>
        <el-button type="primary" @click="submitForm">添加医生</el-button>
        <el-button @click="resetForm">重置</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { Department } from '../types'
import { request } from '@/utils/request'


const departments = ref<Department[]>([])

const form = reactive({
  name: '',
  title: '',
  departmentId: '',
  license: '',
  specialty: '',
  avatar: '',
  status: 'active'
})

const rules = {
  name: [{ required: true, message: '请输入医生姓名', trigger: 'blur' }],
  title: [{ required: true, message: '请输入职称', trigger: 'blur' }],
  departmentId: [{ required: true, message: '请选择科室', trigger: 'change' }],
  license: [{ required: true, message: '请输入执业证号', trigger: 'blur' }]
}

const doctorForm = ref()

onMounted(async () => {
  try {
    departments.value = await request.get('/departments')
  } catch (error) {
    ElMessage.error('获取科室列表失败')
  }
})

const handleAvatarSuccess = (res: any) => {
  form.avatar = res.url
}

const submitForm = async () => {
  if (!doctorForm.value) return
  
  await doctorForm.value.validate(async (valid: boolean) => {
    if (valid) {
      try {
        // 使用requets鉴权
        const response = await request.post('/doctors', form)
        
        if (response.ok) {
          ElMessage.success('医生添加成功')
          resetForm()
        } else {
          ElMessage.error('医生添加失败')
        }
      } catch (error) {
        ElMessage.error('提交失败')
      }
    }
  })
}

const resetForm = () => {
  if (doctorForm.value) {
    doctorForm.value.resetFields()
  }
}
</script>

<style scoped>
.doctor-form {
  max-width: 600px;
  margin: 20px auto;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.1);
}

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
</style> 
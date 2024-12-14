<template>
  <div class="patient-form">
    <h2>添加新患者</h2>
    <el-form :model="form" :rules="rules" ref="patientForm" label-width="120px">
      <el-form-item label="姓名" prop="name">
        <el-input v-model="form.name"></el-input>
      </el-form-item>

      <el-form-item label="性别" prop="gender">
        <el-radio-group v-model="form.gender">
          <el-radio label="male">男</el-radio>
          <el-radio label="female">女</el-radio>
          <el-radio label="other">其他</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-form-item label="年龄" prop="age">
        <el-input-number v-model="form.age" :min="0" :max="150"></el-input-number>
      </el-form-item>

      <el-form-item label="出生日期" prop="birthday">
        <el-date-picker v-model="form.birthday" type="date" placeholder="选择日期"></el-date-picker>
      </el-form-item>

      <el-form-item label="联系电话" prop="phone">
        <el-input v-model="form.phone"></el-input>
      </el-form-item>

      <el-form-item label="紧急联系电话" prop="emergencyPhone">
        <el-input v-model="form.emergencyPhone"></el-input>
      </el-form-item>

      <el-form-item label="身份证号" prop="idCard">
        <el-input v-model="form.idCard"></el-input>
      </el-form-item>

      <el-form-item label="血型" prop="bloodType">
        <el-select v-model="form.bloodType" placeholder="请选择血型">
          <el-option label="A" value="A"></el-option>
          <el-option label="B" value="B"></el-option>
          <el-option label="AB" value="AB"></el-option>
          <el-option label="O" value="O"></el-option>
        </el-select>
      </el-form-item>

      <el-form-item label="过敏史">
        <el-select
          v-model="form.allergies"
          multiple
          filterable
          allow-create
          placeholder="请输入或选择过敏史">
        </el-select>
      </el-form-item>

      <el-form-item label="慢性病史">
        <el-select
          v-model="form.chronicDiseases"
          multiple
          filterable
          allow-create
          placeholder="请输入或选择慢性病史">
        </el-select>
      </el-form-item>

      <el-form-item label="主治医生" prop="doctorId">
        <el-select v-model="form.doctorId" placeholder="请选择主治医生">
          <el-option
            v-for="doctor in doctors"
            :key="doctor.id"
            :label="doctor.name"
            :value="doctor.id">
          </el-option>
        </el-select>
      </el-form-item>

      <el-form-item>
        <el-button type="primary" @click="submitForm">添加患者</el-button>
        <el-button @click="resetForm">重置</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { Doctor } from '../types'

const doctors = ref<Doctor[]>([])

const form = reactive({
  name: '',
  gender: '',
  age: 0,
  birthday: '',
  phone: '',
  emergencyPhone: '',
  idCard: '',
  bloodType: '',
  allergies: [],
  chronicDiseases: [],
  doctorId: '',
  avatar: ''
})

const rules = {
  name: [{ required: true, message: '请输入患者姓名', trigger: 'blur' }],
  gender: [{ required: true, message: '请选择性别', trigger: 'change' }],
  age: [{ required: true, message: '请输入年龄', trigger: 'blur' }],
  phone: [{ required: true, message: '请输入联系电话', trigger: 'blur' }],
  idCard: [{ required: true, message: '请输入身份证号', trigger: 'blur' }],
  doctorId: [{ required: true, message: '请选择主治医生', trigger: 'change' }]
}

const patientForm = ref()

onMounted(async () => {
  try {
    const response = await fetch('/api/doctors')
    doctors.value = await response.json()
  } catch (error) {
    ElMessage.error('获取医生列表失败')
  }
})

const submitForm = async () => {
  if (!patientForm.value) return
  
  await patientForm.value.validate(async (valid: boolean) => {
    if (valid) {
      try {
        const formData = {
          ...form,
          allergies: form.allergies || [],
          chronicDiseases: form.chronicDiseases || []
        }

        const response = await fetch('/api/patients', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(formData)
        })
        
        if (response.ok) {
          ElMessage.success('患者添加成功')
          resetForm()
        } else {
          const error = await response.json()
          ElMessage.error(`患者添加失败: ${error.error}`)
        }
      } catch (error) {
        ElMessage.error('提交失败')
      }
    }
  })
}

const resetForm = () => {
  if (patientForm.value) {
    patientForm.value.resetFields()
  }
}
</script>

<style scoped>
.patient-form {
  max-width: 600px;
  margin: 20px auto;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.1);
}
</style> 
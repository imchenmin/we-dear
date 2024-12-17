<template>
  <div class="patient-profile" v-if="patient">
    <div class="profile-header">
      <el-avatar :size="80" :src="patient.avatar" />
      <div class="patient-basic">
        <h2>{{ patient.name }}</h2>
        <span class="patient-tags">
          <el-tag size="small" type="warning" v-if="patient.chronicDiseases?.length">慢性病</el-tag>
          <el-tag size="small" type="danger" v-if="patient.allergies?.length">过敏史</el-tag>
        </span>
      </div>
    </div>

    <!-- 使用固定位置的tabs作为导航 -->
    <div class="profile-nav">
      <el-tabs @tab-click="handleTabClick">
        <el-tab-pane label="基本信息" name="basic" />
        <el-tab-pane label="随访记录" name="followup" />
        <el-tab-pane label="医疗记录" name="medical" />
      </el-tabs>
    </div>
    
    <!-- 内容区域改为垂直布局 -->
    <div class="profile-content">
      <!-- 基本信息部分 -->
      <div id="basic" class="section">
        <h3 class="section-title">基本信息</h3>
        <div class="profile-info">
          <el-descriptions :column="1" border>
            <el-descriptions-item label="性别">{{ formatGender(patient.gender) }}</el-descriptions-item>
            <el-descriptions-item label="年龄">{{ patient.age }}岁</el-descriptions-item>
            <el-descriptions-item label="出生年月（未校验仅展示）">{{ formatDate(patient.birthday) }}</el-descriptions-item>
            <el-descriptions-item label="血型">{{ patient.bloodType || '未知' }}</el-descriptions-item>
            <el-descriptions-item label="联系电话" :span="2">{{ patient.phone }}</el-descriptions-item>
            <el-descriptions-item label="紧急联系电话" :span="2">{{ patient.emergencyPhone }}</el-descriptions-item>
            <el-descriptions-item label="居住地址" :span="2">{{ patient.address }}</el-descriptions-item>
            <el-descriptions-item label="过敏史" :span="2">
              <el-tag 
                v-for="allergy in patient.allergies" 
                :key="allergy"
                class="mx-1"
                type="danger"
                effect="light"
              >
                {{ allergy }}
              </el-tag>
              <span v-if="!patient.allergies?.length">无</span>
            </el-descriptions-item>
            <el-descriptions-item label="慢性病史" :span="2">
              <el-tag
                v-for="disease in patient.chronicDiseases"
                :key="disease"
                class="mx-1"
                type="warning"
                effect="light"
              >
                {{ disease }}
              </el-tag>
              <span v-if="!patient.chronicDiseases?.length">无</span>
            </el-descriptions-item>
          </el-descriptions>
        </div>
      </div>

      <!-- 随访记录部分 -->
      <div id="followup" class="section">
        <div class="section-header">
          <h3 class="section-title">随访记录</h3>
          <el-button type="primary" @click="showAddFollowUp">
            <el-icon><Plus /></el-icon>新增随访
          </el-button>
        </div>
        <div class="follow-up-records">
          <el-timeline>
            <el-timeline-item
              v-for="record in followUpRecords"
              :key="record.id"
              :timestamp="formatDate(record.followUpDate)"
              :type="record.status === 'completed' ? 'success' : 'primary'"
            >
              <h4>{{ record.title }}</h4>
              <p>{{ record.content }}</p>
              <p class="record-meta">
                下次随访: {{ formatDate(record.nextFollowUp) }}
                <el-tag size="small" :type="record.status === 'completed' ? 'success' : 'warning'">
                  {{ record.status === 'completed' ? '已完成' : '待随访' }}
                </el-tag>
              </p>
            </el-timeline-item>
          </el-timeline>
        </div>
      </div>

      <!-- 医疗记录部分 -->
      <div id="medical" class="section">
        <div class="section-header">
          <h3 class="section-title">医疗记录</h3>
          <el-button type="primary" @click="showAddMedical">
            <el-icon><Plus /></el-icon>新增病历
          </el-button>
        </div>
        <div class="medical-records">
          <el-table :data="medicalRecords" style="width: 100%">
            <el-table-column prop="diagnosisDate" label="就诊日期" width="180">
              <template #default="scope">
                {{ formatDate(scope.row.diagnosisDate) }}
              </template>
            </el-table-column>
            <el-table-column prop="diagnosis" label="诊断结果" />
            <el-table-column prop="treatment" label="治疗方案" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.status === 'completed' ? 'success' : 'warning'">
                  {{ scope.row.status === 'completed' ? '已完成' : '进行中' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>
    </div>

    <!-- 弹窗部分保持不变 -->
    <el-dialog v-model="followUpDialogVisible" title="新增随访记录" width="60%">
      <FollowUpTemplateForm
        ref="templateFormRef"
        :patient-id="patient.id"
        @submit="handleTemplateSubmit"
      />
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="followUpDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitFollowUp">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <el-dialog v-model="medicalDialogVisible" title="新增医疗记录" width="50%">
      <el-form
        ref="medicalFormRef"
        :model="medicalForm"
        :rules="medicalRules"
        label-width="100px"
      >
        <el-form-item label="就诊日期" prop="diagnosisDate">
          <el-date-picker
            v-model="medicalForm.diagnosisDate"
            type="datetime"
            placeholder="选择就诊日期"
          />
        </el-form-item>
        <el-form-item label="症状" prop="symptoms">
          <el-select
            v-model="medicalForm.symptoms"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="请输入症状"
          />
        </el-form-item>
        <el-form-item label="诊断结果" prop="diagnosis">
          <el-input v-model="medicalForm.diagnosis" placeholder="请输入诊断结果" />
        </el-form-item>
        <el-form-item label="治疗方案" prop="treatment">
          <el-input
            v-model="medicalForm.treatment"
            type="textarea"
            rows="3"
            placeholder="请输入治疗方案"
          />
        </el-form-item>
        <el-form-item label="处方" prop="prescription">
          <el-input
            v-model="medicalForm.prescription"
            type="textarea"
            rows="3"
            placeholder="请输入处方"
          />
        </el-form-item>
        <el-form-item label="就诊类型" prop="type">
          <el-select v-model="medicalForm.type" placeholder="请选择就诊类型">
            <el-option label="门诊" value="outpatient" />
            <el-option label="住院" value="inpatient" />
          </el-select>
        </el-form-item>
        <el-form-item label="就诊科室" prop="department">
          <el-input v-model="medicalForm.department" placeholder="请输入就诊科室" />
        </el-form-item>
        <el-form-item label="费用" prop="cost">
          <el-input-number v-model="medicalForm.cost" :precision="2" :step="0.01" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="medicalForm.status" placeholder="请选择状态">
            <el-option label="进行中" value="ongoing" />
            <el-option label="已完成" value="completed" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注" prop="notes">
          <el-input
            v-model="medicalForm.notes"
            type="textarea"
            rows="2"
            placeholder="请输入备注"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="medicalDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitMedical">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import type { Patient, FollowUpRecord, MedicalRecord } from '@/types'
import { medicalApi } from '@/api/medical'
import FollowUpTemplateForm from '@/components/FollowUpTemplateForm.vue'

const props = defineProps<{
  patient?: Patient
}>()

const followUpRecords = ref<FollowUpRecord[]>([])
const medicalRecords = ref<MedicalRecord[]>([])

// 加载随访记录
const loadFollowUpRecords = async () => {
  if (!props.patient) return
  try {
    followUpRecords.value = await medicalApi.getFollowUpRecords(props.patient.id)
  } catch (error) {
    console.error('Failed to load follow-up records:', error)
    ElMessage.error('加载随访记录失败')
  }
}

// 加载医疗记录
const loadMedicalRecords = async () => {
  if (!props.patient) return
  try {
    medicalRecords.value = await medicalApi.getMedicalRecords(props.patient.id)
  } catch (error) {
    console.error('Failed to load medical records:', error)
    ElMessage.error('加载医疗记录失败')
  }
}

// 监听患者变化
watch(() => props.patient?.id, () => {
  loadFollowUpRecords()
  loadMedicalRecords()
})

onMounted(() => {
  if (props.patient) {
    loadFollowUpRecords()
    loadMedicalRecords()
  }
})

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}

const formatGender = (gender: string) => {
  return gender === 'male' ? '男' : gender === 'female' ? '女' : '其他'
}

// 随访记录表单相关
const followUpDialogVisible = ref(false)
const templateFormRef = ref<InstanceType<typeof FollowUpTemplateForm> | null>(null)

// 显示新增随访记录弹窗
const showAddFollowUp = () => {
  followUpDialogVisible.value = true
}

// 处理模板表单提交
const handleTemplateSubmit = async (templateData: any) => {
  try {
    const record = {
      patientId: props.patient?.id,
      templateId: templateData.templateId,
      title: `随访记录 - ${new Date().toLocaleDateString()}`,
      content: JSON.stringify(templateData.data),
      followUpDate: new Date().toISOString(),
      nextFollowUp: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString(), // 默认下次随访时间为一周后
      status: 'completed',
      type: 'regular'
    }

    await medicalApi.createFollowUpRecord(record)
    ElMessage.success('随访记录创建成功')
    followUpDialogVisible.value = false
    loadFollowUpRecords()
  } catch (error) {
    console.error('Failed to create follow-up record:', error)
    ElMessage.error('创建随访记录失败')
  }
}

// 提交随访记录
const submitFollowUp = async () => {
  if (!templateFormRef.value) return
  
  const valid = await templateFormRef.value.validateAndSubmit()
  if (!valid) {
    ElMessage.warning('请完善表单信息')
  }
}

// 医疗记录表单相关
const medicalDialogVisible = ref(false)
const medicalFormRef = ref<FormInstance>()
const medicalForm = ref({
  diagnosisDate: '',
  symptoms: [],
  diagnosis: '',
  treatment: '',
  prescription: '',
  type: 'outpatient',
  department: '',
  cost: 0,
  status: 'ongoing',
  notes: ''
})

const medicalRules = {
  diagnosisDate: [{ required: true, message: '请选择就诊日期', trigger: 'change' }],
  diagnosis: [{ required: true, message: '请输入诊断结果', trigger: 'blur' }],
  treatment: [{ required: true, message: '请输入治疗方案', trigger: 'blur' }],
  type: [{ required: true, message: '请选择就诊类型', trigger: 'change' }],
  department: [{ required: true, message: '请输入就诊科室', trigger: 'blur' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }]
}

// 显示新增医疗记录弹窗
const showAddMedical = () => {
  medicalDialogVisible.value = true
  medicalForm.value = {
    diagnosisDate: '',
    symptoms: [],
    diagnosis: '',
    treatment: '',
    prescription: '',
    type: 'outpatient',
    department: '',
    cost: 0,
    status: 'ongoing',
    notes: ''
  }
}

// 提交医疗记录
const submitMedical = async () => {
  if (!medicalFormRef.value) return
  
  await medicalFormRef.value.validate(async (valid) => {
    if (valid && props.patient) {
      try {
        const record = {
          ...medicalForm.value,
          patientId: props.patient.id
        }
        await medicalApi.createMedicalRecord(record)
        ElMessage.success('医疗记录创建成功')
        medicalDialogVisible.value = false
        loadMedicalRecords()
      } catch (error) {
        console.error('Failed to create medical record:', error)
        ElMessage.error('创建医疗记录失败')
      }
    }
  })
}

// 添加tab点击处理函数
const handleTabClick = (tab: any) => {
  const element = document.getElementById(tab.props.name)
  if (element) {
    element.scrollIntoView({ behavior: 'smooth' })
  }
}
</script>

<style scoped>
.patient-profile {
  position: relative;
  padding: 24px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  height: calc(100vh - 120px);
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 32px;
  padding-bottom: 20px;
  border-bottom: 1px solid #ebeef5;
}

.profile-nav {
  position: sticky;
  top: 0;
  background: #fff;
  z-index: 10;
  border-bottom: 1px solid #ebeef5;
  margin: 0 -24px;
  padding: 0 24px;
}

.profile-content {
  margin-top: 20px;
  overflow-y: auto;
  height: calc(100% - 180px);
  scroll-behavior: smooth;
}

.section {
  padding: 24px 0;
  border-bottom: 1px solid #ebeef5;
}

.section:last-child {
  border-bottom: none;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-title {
  margin: 0 0 20px 0;
  font-size: 18px;
  color: #303133;
  font-weight: bold;
}

.patient-basic {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.patient-basic h2 {
  margin: 0;
  color: #303133;
  font-size: 24px;
}

.patient-tags {
  display: flex;
  gap: 8px;
}

.profile-info {
  margin-top: 16px;
}

.mx-1 {
  margin: 0 4px;
}

.follow-up-records,
.medical-records {
  margin-top: 20px;
}

.record-meta {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-top: 8px;
  color: #909399;
  font-size: 14px;
}

:deep(.el-form-item) {
  margin-bottom: 22px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* 滚动条样式优化 */
.profile-content {
  scrollbar-width: thin;
  scrollbar-color: #909399 #f4f4f5;
}

.profile-content::-webkit-scrollbar {
  width: 6px;
}

.profile-content::-webkit-scrollbar-track {
  background: #f4f4f5;
}

.profile-content::-webkit-scrollbar-thumb {
  background-color: #909399;
  border-radius: 3px;
}
</style> 
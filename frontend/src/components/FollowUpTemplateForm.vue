<template>
  <div class="template-form">
    <!-- 模板选择 -->
    <el-form-item label="随访模板" required>
      <el-select
        v-model="selectedTemplateId"
        placeholder="请选择随访模板"
        @change="handleTemplateChange"
      >
        <el-option-group
          v-for="group in groupedTemplates"
          :key="group.category"
          :label="formatCategory(group.category)"
        >
          <el-option
            v-for="template in group.templates"
            :key="template.id"
            :label="template.name"
            :value="template.id"
          >
            <span>{{ template.name }}</span>
            <span class="template-version">v{{ template.version }}</span>
          </el-option>
        </el-option-group>
      </el-select>
    </el-form-item>

    <!-- 动态表单 -->
    <template v-if="selectedTemplate">
      <el-form
        ref="dynamicFormRef"
        :model="formData"
        label-width="120px"
      >
        <template v-for="(field, key) in formFields" :key="key">
          <!-- 日期选择器 -->
          <el-form-item
            v-if="field.type === 'string' && field.format === 'date'"
            :label="field.title || key"
            :prop="key"
            :rules="generateRules(field)"
          >
            <el-date-picker
              v-model="formData[key]"
              type="date"
              :placeholder="'请选择' + (field.title || key)"
              value-format="YYYY-MM-DD"
            />
          </el-form-item>

          <!-- 文本输入 -->
          <el-form-item
            v-else-if="field.type === 'string' && !field.enum && field.format !== 'textarea'"
            :label="field.title || key"
            :prop="key"
            :rules="generateRules(field)"
          >
            <el-input
              v-model="formData[key]"
              :placeholder="'请输入' + (field.title || key)"
            />
          </el-form-item>

          <!-- 文本域 -->
          <el-form-item
            v-else-if="field.type === 'string' && field.format === 'textarea'"
            :label="field.title || key"
            :prop="key"
            :rules="generateRules(field)"
          >
            <el-input
              v-model="formData[key]"
              type="textarea"
              :rows="3"
              :placeholder="'请输入' + (field.title || key)"
            />
          </el-form-item>

          <!-- 数字输入 -->
          <el-form-item
            v-else-if="field.type === 'number'"
            :label="field.title || key"
            :prop="key"
            :rules="generateRules(field)"
          >
            <el-input-number
              v-model="formData[key]"
              :min="field.minimum"
              :max="field.maximum"
              :step="field.multipleOf || 1"
            />
            <span v-if="field.unit" class="unit-label">{{ field.unit }}</span>
          </el-form-item>

          <!-- 单选框 -->
          <el-form-item
            v-else-if="field.type === 'string' && field.enum"
            :label="field.title || key"
            :prop="key"
            :rules="generateRules(field)"
          >
            <el-select
              v-model="formData[key]"
              :placeholder="'请选择' + (field.title || key)"
            >
              <el-option
                v-for="option in field.enum"
                :key="option"
                :label="option"
                :value="option"
              />
            </el-select>
          </el-form-item>

          <!-- 多选框 -->
          <el-form-item
            v-else-if="field.type === 'array' && field.items?.type === 'string'"
            :label="field.title || key"
            :prop="key"
            :rules="generateRules(field)"
          >
            <el-select
              v-model="formData[key]"
              multiple
              :placeholder="'请选择' + (field.title || key)"
            >
              <el-option
                v-for="option in field.items.enum || []"
                :key="option"
                :label="option"
                :value="option"
              />
            </el-select>
          </el-form-item>

          <!-- 对象数组 -->
          <template v-else-if="field.type === 'array' && field.items?.type === 'object'">
            <el-divider>{{ field.title || key }}</el-divider>
            <div class="array-form">
              <div
                v-for="(item, index) in formData[key] || []"
                :key="index"
                class="array-item"
              >
                <div class="array-item-header">
                  <span>{{ field.title }} #{{ index + 1 }}</span>
                  <el-button
                    type="danger"
                    link
                    @click="removeArrayItem(key, index)"
                  >
                    删除
                  </el-button>
                </div>
                <template v-for="(subField, subKey) in field.items.properties" :key="subKey">
                  <el-form-item
                    :label="subField.title || subKey"
                    :prop="`${key}.${index}.${subKey}`"
                    :rules="generateRules(subField)"
                  >
                    <el-input
                      v-if="subField.type === 'string' && !subField.enum"
                      v-model="item[subKey]"
                      :placeholder="'请输入' + (subField.title || subKey)"
                    />
                    <el-select
                      v-else-if="subField.type === 'string' && subField.enum"
                      v-model="item[subKey]"
                      :placeholder="'请选择' + (subField.title || subKey)"
                    >
                      <el-option
                        v-for="option in subField.enum"
                        :key="option"
                        :label="option"
                        :value="option"
                      />
                    </el-select>
                  </el-form-item>
                </template>
              </div>
              <el-button
                type="primary"
                link
                @click="addArrayItem(key, field.items.properties)"
              >
                添加{{ field.title }}
              </el-button>
            </div>
          </template>

          <!-- 嵌套对象 -->
          <template v-else-if="field.type === 'object'">
            <el-divider>{{ field.title || key }}</el-divider>
            <div class="nested-form" v-if="field.properties">
              <template v-for="(subField, subKey) in field.properties" :key="subKey">
                <!-- 处理二级嵌套对象 -->
                <template v-if="subField.type === 'object'">
                  <el-divider content-position="left">{{ subField.title || subKey }}</el-divider>
                  <div class="nested-form" v-if="subField.properties">
                    <template v-for="(thirdField, thirdKey) in subField.properties" :key="thirdKey">
                      <el-form-item
                        :label="thirdField.title || thirdKey"
                        :prop="`${key}.${subKey}.${thirdKey}`"
                        :rules="generateRules(thirdField)"
                      >
                        <!-- 三级对象的各种类型处理 -->
                        <template v-if="thirdField.type === 'string' && !thirdField.enum">
                          <el-input
                            v-model="formData[key][subKey][thirdKey]"
                            :placeholder="'请输入' + (thirdField.title || thirdKey)"
                          />
                        </template>
                        <template v-else-if="thirdField.type === 'number'">
                          <el-input-number
                            v-model="formData[key][subKey][thirdKey]"
                            :min="thirdField.minimum"
                            :max="thirdField.maximum"
                            :step="thirdField.multipleOf || 1"
                          />
                          <span v-if="thirdField.unit" class="unit-label">{{ thirdField.unit }}</span>
                        </template>
                        <template v-else-if="thirdField.type === 'string' && thirdField.enum">
                          <el-select
                            v-model="formData[key][subKey][thirdKey]"
                            :placeholder="'请选择' + (thirdField.title || thirdKey)"
                          >
                            <el-option
                              v-for="option in thirdField.enum"
                              :key="option"
                              :label="option"
                              :value="option"
                            />
                          </el-select>
                        </template>
                      </el-form-item>
                    </template>
                  </div>
                </template>
                <!-- 处理普通字段 -->
                <el-form-item
                  v-else
                  :label="subField.title || subKey"
                  :prop="`${key}.${subKey}`"
                  :rules="generateRules(subField)"
                >
                  <!-- 子对象的各种类型处理 -->
                  <template v-if="subField.type === 'string' && !subField.enum">
                    <el-input
                      v-model="formData[key][subKey]"
                      :placeholder="'请输入' + (subField.title || subKey)"
                    />
                  </template>
                  <template v-else-if="subField.type === 'number'">
                    <el-input-number
                      v-model="formData[key][subKey]"
                      :min="subField.minimum"
                      :max="subField.maximum"
                      :step="subField.multipleOf || 1"
                    />
                    <span v-if="subField.unit" class="unit-label">{{ subField.unit }}</span>
                  </template>
                  <template v-else-if="subField.type === 'string' && subField.enum">
                    <el-select
                      v-model="formData[key][subKey]"
                      :placeholder="'请选择' + (subField.title || subKey)"
                    >
                      <el-option
                        v-for="option in subField.enum"
                        :key="option"
                        :label="option"
                        :value="option"
                      />
                    </el-select>
                  </template>
                </el-form-item>
              </template>
            </div>
          </template>
        </template>
      </el-form>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { request } from '@/utils/request'
import { medicalApi } from '@/api/medical'
import { FormField } from '@/types'

const props = defineProps<{
  patientId: string
}>()

const emit = defineEmits<{
  (e: 'submit', data: any): void
}>()

const templates = ref<any[]>([])
const selectedTemplateId = ref('')
const selectedTemplate = ref<any>(null)
const formData = ref<any>({})
const dynamicFormRef = ref<FormInstance>()

// 按分类分组的模板列表
const groupedTemplates = computed(() => {
  const groups: { [key: string]: any[] } = {}
  templates.value.forEach(template => {
    template.categories.forEach((category: string) => {
      if (!groups[category]) {
        groups[category] = []
      }
      groups[category].push(template)
    })
  })
  return Object.entries(groups).map(([category, templates]) => ({
    category,
    templates
  }))
})

// 从模板schema中提取表单字段
const formFields = computed(() => {
  if (!selectedTemplate.value?.schema) return {}
  try {
    const schema = JSON.parse(selectedTemplate.value.schema)
    return schema.properties || {}
  } catch (error) {
    console.error('Invalid schema:', error)
    return {}
  }
})

// 加载模板列表
const loadTemplates = async () => {
  try {
    const response = await request.get('/templates')
    templates.value = response.filter((t: any) => t.status === 'enabled')
  } catch (error) {
    console.error('Failed to load templates:', error)
    ElMessage.error('加载模板列表失败')
  }
}

// 处理模板选择变化
const handleTemplateChange = async (templateId: string) => {
  selectedTemplate.value = templates.value.find(t => t.id === templateId)
  if (selectedTemplate.value) {
    // 初始化表单数据
    formData.value = initFormData(formFields.value)
  }
}

// 初始化表单数据
const initFormData = (fields: any) => {
  const data: any = {}
  Object.entries(fields).forEach(([key, field]: [string, any]) => {
    if (field.type === 'object' && field.properties) {
      data[key] = initFormData(field.properties)
    } else if (field.type === 'array') {
      if (field.items?.type === 'object') {
        data[key] = []
      } else {
        data[key] = field.default || []
      }
    } else if (field.type === 'number') {
      data[key] = field.default !== undefined ? field.default : null
    } else if (field.type === 'string' && field.format === 'date') {
      data[key] = field.default || null
    } else {
      data[key] = field.default || ''
    }
  })
  return data
}

// 验证并提交表单
const validateAndSubmit = async () => {
  if (!dynamicFormRef.value) return false
  
  try {
    await dynamicFormRef.value.validate()
    // 验证数据是否符合模板要求
    const response = await request.post('/templates/validate', formData.value, {
      params: { templateId: selectedTemplate.value.id }
    })
    
    if (response.valid) {
      emit('submit', {
        templateId: selectedTemplate.value.id,
        data: formData.value
      })
      return true
    } else {
      ElMessage.error('表单数据不符合模板要求')
      return false
    }
  } catch (error) {
    console.error('Validation failed:', error)
    return false
  }
}

// 格式化分类显示
const formatCategory = (category: string) => {
  const categories: Record<string, string> = {
    general: '通用随访',
    diabetes: '糖尿病随访',
    cardiac: '心脏病随访',
    pediatric: '儿科随访',
    elderly: '老年人随访'
  }
  return categories[category] || category
}

// 使用默认模板
const useDefaultSchema = async () => {
  try {
    const response = await medicalApi.getDefaultTemplateSchema()
    form.value.schema = response
  } catch (error) {
    console.error('Failed to get default schema:', error)
    ElMessage.error('获取默认模板失败')
  }
}

// 暴露方法给父组件
defineExpose({
  validateAndSubmit
})

// 组件挂载时加载模板列表
onMounted(() => {
  loadTemplates()
})

// 添加数组项
const addArrayItem = (key: string, properties: any) => {
  if (!formData.value[key]) {
    formData.value[key] = []
  }
  const newItem: any = {}
  Object.keys(properties).forEach(prop => {
    newItem[prop] = properties[prop].default || ''
  })
  formData.value[key].push(newItem)
}

// 删除数组项
const removeArrayItem = (key: string, index: number) => {
  formData.value[key].splice(index, 1)
}

// 生成验证规则
const generateRules = (field: FormField): any => {
  const rules: any[] = []
  
  if (field.required) {
    rules.push({
      required: true,
      message: `请${field.type === 'array' ? '选择' : '输入'}${field.title}`,
      trigger: field.type === 'array' ? 'change' : 'blur'
    })
  }

  if (field.type === 'number') {
    if (field.minimum !== undefined) {
      rules.push({
        type: 'number',
        min: field.minimum,
        message: `${field.title}不能小于${field.minimum}${field.unit || ''}`,
        trigger: 'blur'
      })
    }
    if (field.maximum !== undefined) {
      rules.push({
        type: 'number',
        max: field.maximum,
        message: `${field.title}不能大于${field.maximum}${field.unit || ''}`,
        trigger: 'blur'
      })
    }
  }

  return rules
}

const submitForm = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    // 验证数据是否符合模板要求
    await templateApi.validateTemplateData(selectedTemplate.value.id, formData.value)
    
    // 提交表单数据
    const record = {
      patientId: props.patientId,
      templateId: selectedTemplate.value.id,
      content: JSON.stringify(formData.value),
      followUpDate: formData.value.basicInfo?.followUpDate,
      type: formData.value.basicInfo?.followUpType,
      status: 'completed'
    }
    
    await medicalApi.createFollowUpRecord(record)
    ElMessage.success('随访记录创建成功')
    emit('success')
  } catch (error) {
    console.error('Form submission failed:', error)
    ElMessage.error('提交失败，请检查表单数据')
  }
}
</script>

<style scoped>
.template-form {
  margin-top: 20px;
}

.template-version {
  float: right;
  color: #909399;
  font-size: 12px;
}

.nested-form {
  padding-left: 20px;
  border-left: 2px solid #ebeef5;
  margin: 16px 0;
}

.array-form {
  margin: 16px 0;
}

.array-item {
  border: 1px solid #ebeef5;
  border-radius: 4px;
  padding: 16px;
  margin-bottom: 16px;
}

.array-item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.unit-label {
  margin-left: 8px;
  color: #909399;
}

:deep(.el-divider__text) {
  font-size: 14px;
  color: #606266;
}

:deep(.el-divider--horizontal) {
  margin: 16px 0;
}

:deep(.el-divider[content-position=left]) {
  margin: 8px 0;
}

:deep(.el-form-item) {
  margin-bottom: 22px;
}
</style> 
<template>
  <div class="template-manager">
    <!-- 模板列表 -->
    <div class="template-list">
      <div class="list-header">
        <h3>随访模板管理</h3>
        <el-button type="primary" @click="showCreateDialog">
          <el-icon><Plus /></el-icon>新建模板
        </el-button>
      </div>

      <el-table :data="templates" style="width: 100%">
        <el-table-column prop="name" label="模板名称" />
        <el-table-column prop="description" label="描述" show-overflow-tooltip />
        <el-table-column prop="version" label="版本" width="100" />
        <el-table-column prop="categories" label="适用分类">
          <template #default="scope">
            <el-tag
              v-for="category in scope.row.categories"
              :key="category"
              size="small"
              class="mx-1"
            >
              {{ formatCategory(category) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'enabled' ? 'success' : 'info'">
              {{ scope.row.status === 'enabled' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <el-button-group>
              <el-button
                type="primary"
                link
                @click="editTemplate(scope.row)"
              >
                编辑
              </el-button>
              <el-button
                type="primary"
                link
                @click="previewTemplate(scope.row)"
              >
                预览
              </el-button>
              <el-button
                type="danger"
                link
                @click="deleteTemplate(scope.row)"
              >
                删除
              </el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 创建/编辑模板对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑模板' : '新建模板'"
      width="60%"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="模板名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入模板名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入模板描述"
          />
        </el-form-item>
        <el-form-item label="版本" prop="version">
          <el-input v-model="form.version" placeholder="请输入版本号，如: 1.0.0" />
        </el-form-item>
        <el-form-item label="适用分类" prop="categories">
          <el-select
            v-model="form.categories"
            multiple
            placeholder="请选择适用分类"
          >
            <el-option
              v-for="category in categoryOptions"
              :key="category.value"
              :label="category.label"
              :value="category.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio label="enabled">启用</el-radio>
            <el-radio label="disabled">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="模板结构" prop="schema">
          <el-input
            v-model="form.schema"
            type="textarea"
            :rows="10"
            placeholder="请输入JSON Schema格式的模板结构"
          />
          <div class="schema-tools">
            <el-button type="primary" link @click="formatSchema">格式化</el-button>
            <el-button type="primary" link @click="useDefaultSchema">使用默认模板</el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitForm">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 预览模板对话框 -->
    <el-dialog
      v-model="previewDialogVisible"
      title="预览模板"
      width="50%"
    >
      <div class="preview-content">
        <pre>{{ JSON.stringify(selectedTemplate?.schema, null, 2) }}</pre>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { request } from '@/utils/request'

const templates = ref([])
const dialogVisible = ref(false)
const previewDialogVisible = ref(false)
const isEdit = ref(false)
const selectedTemplate = ref(null)
const formRef = ref<FormInstance>()

const form = ref({
  name: '',
  description: '',
  version: '',
  categories: [],
  status: 'enabled',
  schema: ''
})

const categoryOptions = [
  { label: '通用随访', value: 'general' },
  { label: '糖尿病随访', value: 'diabetes' },
  { label: '心脏病随访', value: 'cardiac' },
  { label: '儿科随访', value: 'pediatric' },
  { label: '老年人随访', value: 'elderly' }
]

const rules = {
  name: [{ required: true, message: '请输入模板名称', trigger: 'blur' }],
  version: [{ required: true, message: '请输入版本号', trigger: 'blur' }],
  categories: [{ required: true, message: '请选择适用分类', trigger: 'change' }],
  schema: [{ required: true, message: '请输入模板结构', trigger: 'blur' }]
}

// 加载模板列表
const loadTemplates = async () => {
  try {
    const response = await request.get('/templates')
    templates.value = response
  } catch (error) {
    console.error('Failed to load templates:', error)
    ElMessage.error('加载模板列表失败')
  }
}

// 显示创建对话框
const showCreateDialog = () => {
  isEdit.value = false
  form.value = {
    name: '',
    description: '',
    version: '1.0.0',
    categories: [],
    status: 'enabled',
    schema: ''
  }
  dialogVisible.value = true
}

// 编辑模板
const editTemplate = (template: any) => {
  isEdit.value = true
  selectedTemplate.value = template
  form.value = { ...template }
  dialogVisible.value = true
}

// 预览模板
const previewTemplate = (template: any) => {
  selectedTemplate.value = template
  previewDialogVisible.value = true
}

// 删除模板
const deleteTemplate = async (template: any) => {
  try {
    await ElMessageBox.confirm('确定要删除该模板吗？', '提示', {
      type: 'warning'
    })
    
    await request.delete(`/templates/${template.id}`)
    ElMessage.success('删除成功')
    loadTemplates()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete template:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 格式化Schema
const formatSchema = () => {
  try {
    const parsed = JSON.parse(form.value.schema)
    form.value.schema = JSON.stringify(parsed, null, 2)
  } catch (error) {
    ElMessage.error('JSON格式不正确')
  }
}

// 使用默认模板
const useDefaultSchema = async () => {
  try {
    const response = await request.get('/templates/default-schema')
    console.log(response)
    form.value.schema = JSON.stringify(response, null, 2)
  } catch (error) {
    console.error('Failed to get default schema:', error)
    ElMessage.error('获取默认模板失败')
  }
}

// 提交表单
const submitForm = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        // 验证JSON Schema格式
        JSON.parse(form.value.schema)
        
        if (isEdit.value) {
          await request.put(`/templates/${selectedTemplate.value.id}`, form.value)
          ElMessage.success('更新成功')
        } else {
          await request.post('/templates', form.value)
          ElMessage.success('创建成功')
        }
        
        dialogVisible.value = false
        loadTemplates()
      } catch (error) {
        if (error instanceof SyntaxError) {
          ElMessage.error('JSON Schema格式不正确')
        } else {
          console.error('Failed to save template:', error)
          ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
        }
      }
    }
  })
}

// 格式化分类显示
const formatCategory = (category: string) => {
  const option = categoryOptions.find(opt => opt.value === category)
  return option ? option.label : category
}

onMounted(() => {
  loadTemplates()
})
</script>

<style scoped>
.template-manager {
  padding: 20px;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.list-header h3 {
  margin: 0;
  font-size: 18px;
  color: #303133;
}

.mx-1 {
  margin: 0 4px;
}

.schema-tools {
  margin-top: 8px;
  display: flex;
  gap: 12px;
}

.preview-content {
  background: #f5f7fa;
  padding: 16px;
  border-radius: 4px;
  max-height: 500px;
  overflow-y: auto;
}

.preview-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
}

:deep(.el-form-item) {
  margin-bottom: 22px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style> 
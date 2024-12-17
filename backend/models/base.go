package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// BaseModel 基础模型
type BaseModel struct {
	ID        string         `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
}

// Department 科室
type Department struct {
	BaseModel
	Name        string `json:"name" gorm:"uniqueIndex"` // 科室名称
	Description string `json:"description"`             // 科室描述
	Code        string `json:"code"`                    // 科室编码
}

// Doctor 医生
type Doctor struct {
	BaseModel
	Name         string     `json:"name"`                       // 姓名
	Username     string     `json:"username" gorm:"unique"`     // 登录用户名
	Password     string     `json:"-" gorm:"not null"`          // 密码（json中不返回）
	Salt         string     `json:"-"`                          // 密码盐
	Title        string     `json:"title"`                      // 职称
	DepartmentID string     `json:"departmentId"`               // 所属科室ID
	Department   Department `json:"department"`                 // 所属科室
	License      string     `json:"license"`                    // 执业证号
	Specialty    string     `json:"specialty"`                  // 专长
	Avatar       string     `json:"avatar"`                     // 头像
	Status       string     `json:"status"`                     // 状态（在职/离职等）
	Role         string     `json:"role" gorm:"default:doctor"` // 角色（医生/管理员）
	LastLoginAt  time.Time  `json:"lastLoginAt"`                // 最后登录时间
	Patients     []Patient  `json:"patients,omitempty" gorm:"foreignKey:DoctorID"`
}

// Patient 患者
type Patient struct {
	BaseModel
	Name            string         `json:"name" gorm:"index"`
	Gender          string         `json:"gender"`
	Age             int            `json:"age"` // 由身份证号计算
	Birthday        time.Time      `json:"birthday"`
	Phone           string         `json:"phone" gorm:"index"`
	EmergencyPhone  string         `json:"emergencyPhone"`
	Address         string         `json:"address"`
	IDCard          string         `json:"idCard" gorm:"uniqueIndex"`
	BloodType       string         `json:"bloodType"`
	Allergies       pq.StringArray `json:"allergies" gorm:"type:text[]"`
	ChronicDiseases pq.StringArray `json:"chronicDiseases" gorm:"type:text[]"`
	Avatar          string         `json:"avatar"`
	Messages        []Message      `json:"messages,omitempty" gorm:"foreignKey:PatientID"`

	// TODO主治医生id，后续需要改为多对多关系。使用签约关系管理
	DoctorID string `json:"doctorId" gorm:"index"` // 主治医生ID
	Doctor   Doctor `json:"doctor" gorm:"foreignKey:DoctorID"`
}

// MedicalRecord 病历记录
type MedicalRecord struct {
	BaseModel
	PatientID     string         `json:"patientId" gorm:"index"`         // 患者ID
	DoctorID      string         `json:"doctorId" gorm:"index"`          // 医生ID
	DiagnosisDate time.Time      `json:"diagnosisDate"`                  // 诊断日期
	Symptoms      pq.StringArray `json:"symptoms" gorm:"type:text[]"`    // 症状
	Diagnosis     string         `json:"diagnosis"`                      // 诊断结果
	Treatment     string         `json:"treatment"`                      // 治疗方案
	Prescription  string         `json:"prescription"`                   // 处方
	Notes         string         `json:"notes"`                          // 备注
	Status        string         `json:"status"`                         // 状态（进行中/已完成等）
	Type          string         `json:"type"`                           // 就诊类型（门诊/住院等）
	Department    string         `json:"department"`                     // 就诊科室
	Cost          float64        `json:"cost"`                           // 费用，非必填
	Attachments   []string       `json:"attachments" gorm:"type:text[]"` // 附件列表
}

// Message 消息
type Message struct {
	BaseModel
	PatientID string `json:"patientId"` // 患者ID
	DoctorID  string `json:"doctorId"`  // 医生ID
	RecordID  string `json:"recordId"`  // 关联的病历ID
	Content   string `json:"content"`   // 消息内容
	Type      string `json:"type"`      // 消息类型（文本/图片/语音等）
	Role      string `json:"role"`      // 发送者角色（医生/患者）
	Read      bool   `json:"read"`      // 是否已读
	ReplyTo   string `json:"replyTo"`   // 回复的消息ID
}

// AISuggestion AI 建议
type AISuggestion struct {
	BaseModel
	MessageID   string    `json:"messageId"`   // 关联的消息ID
	PatientID   string    `json:"patientId"`   // 患者ID
	Content     string    `json:"content"`     // 建议内容
	PromptUsed  string    `json:"-"`           // 使用的提示词
	ContextUsed string    `json:"-"`           // 使用的上下文
	ModelUsed   string    `json:"-"`           // 使用的模型
	Confidence  float64   `json:"confidence"`  // 置信度
	Category    string    `json:"category"`    // 建议类别（用药/就医/生活等）
	Priority    int       `json:"priority"`    // 优先级（1-5）
	Status      string    `json:"status"`      // 状态（待审核/已采纳/已拒绝等）
	ReviewedBy  string    `json:"reviewedBy"`  // 审核医生ID
	ReviewedAt  time.Time `json:"reviewedAt"`  // 审核时间
	ReviewNotes string    `json:"reviewNotes"` // 审核备注
	// Embedding   []float32 `json:"-" gorm:"type:vector(1536)"`
}

// Attachment 附件（检查报告、图片等）
type Attachment struct {
	BaseModel
	MessageID   string `json:"messageId"`   // 关联的消息ID
	RecordID    string `json:"recordId"`    // 关联的病历ID
	Type        string `json:"type"`        // 类型（检查报告/图片/语音等）
	URL         string `json:"url"`         // 文件URL
	Name        string `json:"name"`        // 文件名
	Size        int64  `json:"size"`        // 文件大小
	ContentType string `json:"contentType"` // 文件类型
	UploadedBy  string `json:"uploadedBy"`  // 上传者ID
}

// FollowUpTemplate 随访模板
type FollowUpTemplate struct {
	BaseModel
	Name        string         `json:"name" gorm:"index"`             // 模板名称
	Description string         `json:"description"`                   // 模板描述
	Schema      string         `json:"schema" gorm:"type:text"`       // 模板格式（JSON Schema）
	Version     string         `json:"version"`                       // 版本号
	Status      string         `json:"status"`                        // 状态（启用/禁用）
	CreatedBy   string         `json:"createdBy"`                     // 创建人ID
	UpdatedBy   string         `json:"updatedBy"`                     // 最后修改人ID
	Categories  pq.StringArray `json:"categories" gorm:"type:text[]"` // 适用分类
}

// FollowUpRecord 随访记录
type FollowUpRecord struct {
	BaseModel
	PatientID    string    `json:"patientId" gorm:"index"`         // 患者ID
	DoctorID     string    `json:"doctorId" gorm:"index"`          // 医生ID
	TemplateID   string    `json:"templateId"`                     // 使用的模板ID
	Title        string    `json:"title"`                          // 随访标题
	Content      string    `json:"content"`                        // 随访内容（根据模板填写的JSON数据）
	FollowUpDate time.Time `json:"followUpDate"`                   // 随访日期
	NextFollowUp time.Time `json:"nextFollowUp"`                   // 下次随访日期
	Status       string    `json:"status"`                         // 状态(completed/pending)
	Type         string    `json:"type"`                           // 随访类型(常规/特殊)
	Attachments  []string  `json:"attachments" gorm:"type:text[]"` // 附件列表
}

// AIAgentTemplate AI代理模板
type AIAgentTemplate struct {
	BaseModel
	Name        string         `json:"name" gorm:"index"`             // 模板名称
	Version     string         `json:"version"`                       // 版本号
	Description string         `json:"description"`                   // 模板描述
	Content     string         `json:"content" gorm:"type:text"`      // 模板内容（JSON格式）
	Categories  pq.StringArray `json:"categories" gorm:"type:text[]"` // 适用分类
	Status      string         `json:"status"`                        // 状态（启用/禁用）
	CreatedBy   string         `json:"createdBy"`                     // 创建人ID
	UpdatedBy   string         `json:"updatedBy"`                     // 最后修改人ID
	LastAuditBy string         `json:"lastAuditBy"`                   // 最后审核人ID
	LastAuditAt time.Time      `json:"lastAuditAt"`                   // 最后审核时间
	AuditStatus string         `json:"auditStatus"`                   // 审核状态
	AuditNotes  string         `json:"auditNotes"`                    // 审核备注
}

// AISuggestionFeedback AI建议评价
type AISuggestionFeedback struct {
	BaseModel
	SuggestionID string    `json:"suggestionId" gorm:"index"` // 关联的AI建议ID
	PatientID    string    `json:"patientId" gorm:"index"`    // 患者ID
	Rating       int       `json:"rating"`                    // 评价 (1: 点赞, -1: 踩)
	Comment      string    `json:"comment"`                   // 评论内容
	Tags         []string  `json:"tags" gorm:"type:text[]"`   // 标签（例如：有帮助、专业、易懂等）
	CreatedBy    string    `json:"createdBy"`                 // 评价人ID
	UpdatedBy    string    `json:"updatedBy"`                 // 最后修改人ID
	ReviewedBy   string    `json:"reviewedBy"`                // 审核人ID
	ReviewedAt   time.Time `json:"reviewedAt"`                // 审核时间
	Status       string    `json:"status"`                    // 状态（待审核/已通过/已拒绝）
}

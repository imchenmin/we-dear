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
	Name         string     `json:"name"`         // 姓名
	Title        string     `json:"title"`        // 职称
	DepartmentID string     `json:"departmentId"` // 所属科室ID
	Department   Department `json:"department"`   // 所属科室
	License      string     `json:"license"`      // 执业证号
	Specialty    string     `json:"specialty"`    // 专长
	Avatar       string     `json:"avatar"`       // 头像
	Status       string     `json:"status"`       // 状态（在职/离职等）
	Patients     []Patient  `json:"patients,omitempty" gorm:"foreignKey:DoctorID"`
}

// Patient 患者
type Patient struct {
	BaseModel
	Name            string         `json:"name" gorm:"index"`
	Gender          string         `json:"gender"`
	Age             int            `json:"age"`
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
	DoctorID        string         `json:"doctorId" gorm:"index"` // 主治医生ID
	Doctor          Doctor         `json:"doctor" gorm:"foreignKey:DoctorID"`
}

// MedicalRecord 病历记录
type MedicalRecord struct {
	BaseModel
	PatientID     string    `json:"patientId"`                   // 患者ID
	DoctorID      string    `json:"doctorId"`                    // 医生ID
	DiagnosisDate time.Time `json:"diagnosisDate"`               // 诊断日期
	Symptoms      []string  `json:"symptoms" gorm:"type:text[]"` // 症状
	Diagnosis     string    `json:"diagnosis"`                   // 诊断结果
	Treatment     string    `json:"treatment"`                   // 治疗方案
	Prescription  string    `json:"prescription"`                // 处方
	Notes         string    `json:"notes"`                       // 备注
	Status        string    `json:"status"`                      // 状态（进行中/已完成等）
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
	Status    string `json:"status"`    // 状态（未读/已读等）
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

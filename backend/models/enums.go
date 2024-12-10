package models

// 性别
const (
	GenderMale   = "male"
	GenderFemale = "female"
	GenderOther  = "other"
)

// 血型
const (
	BloodTypeA  = "A"
	BloodTypeB  = "B"
	BloodTypeAB = "AB"
	BloodTypeO  = "O"
)

// 消息类型
const (
	MessageTypeText  = "text"
	MessageTypeImage = "image"
	MessageTypeVoice = "voice"
	MessageTypeFile  = "file"
)

// 消息状态
const (
	MessageStatusUnread = "unread"
	MessageStatusRead   = "read"
)

// 消息角色
const (
	MessageRoleDoctor  = "doctor"
	MessageRolePatient = "patient"
	MessageRoleSystem  = "system"
)

// AI建议类别
const (
	AISuggestionCategoryMedication = "medication" // 用药建议
	AISuggestionCategoryVisit      = "visit"      // 就医建议
	AISuggestionCategoryLifestyle  = "lifestyle"  // 生活建议
	AISuggestionCategoryUrgent     = "urgent"     // 紧急建议
)

// AI建议状态
const (
	AISuggestionStatusPending  = "pending"  // 待审核
	AISuggestionStatusApproved = "approved" // 已采纳
	AISuggestionStatusRejected = "rejected" // 已拒绝
)

// AI建议优先级
const (
	AISuggestionPriorityLow      = 1 // 低优先级
	AISuggestionPriorityNormal   = 2 // 普通优先级
	AISuggestionPriorityHigh     = 3 // 高优先级
	AISuggestionPriorityUrgent   = 4 // 紧急
	AISuggestionPriorityCritical = 5 // 危急
)

// 病历状态
const (
	MedicalRecordStatusInProgress = "in_progress" // 进行中
	MedicalRecordStatusCompleted  = "completed"   // 已完成
	MedicalRecordStatusCancelled  = "cancelled"   // 已取消
)

// 医生状态
const (
	DoctorStatusActive   = "active"   // 在职
	DoctorStatusInactive = "inactive" // 离职
	DoctorStatusVacation = "vacation" // 休假
)

// 附件类型
const (
	AttachmentTypeReport = "report" // 检查报告
	AttachmentTypeImage  = "image"  // 图片
	AttachmentTypeVoice  = "voice"  // 语音
	AttachmentTypeVideo  = "video"  // 视频
	AttachmentTypeOther  = "other"  // 其他
)

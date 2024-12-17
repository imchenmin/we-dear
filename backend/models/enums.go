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

// 模板状态
const (
	TemplateStatusEnabled  = "enabled"  // 启用
	TemplateStatusDisabled = "disabled" // 禁用
)

// 随访模板分类
const (
	TemplateCategoryGeneral   = "general"   // 通用随访
	TemplateCategoryDiabetes  = "diabetes"  // 糖尿病随访
	TemplateCategoryCardiac   = "cardiac"   // 心脏病随访
	TemplateCategoryPediatric = "pediatric" // 儿科随访
	TemplateCategoryElderly   = "elderly"   // 老年人随访
)

// AI代理模板分类
const (
	AIAgentCategoryGeneral     = "general"     // 通用模板
	AIAgentCategoryDiabetes    = "diabetes"    // 糖尿病
	AIAgentCategoryCardiac     = "cardiac"     // 心脏病
	AIAgentCategoryPediatric   = "pediatric"   // 儿科
	AIAgentCategoryGeriatric   = "geriatric"   // 老年科
	AIAgentCategoryOncology    = "oncology"    // 肿瘤科
	AIAgentCategoryPsychiatric = "psychiatric" // 精神科
)

// AI代理模板审核状态
const (
	AIAgentAuditStatusPending  = "pending"  // 待审核
	AIAgentAuditStatusApproved = "approved" // 已通过
	AIAgentAuditStatusRejected = "rejected" // 已拒绝
)

// AI代理模板状态
const (
	AIAgentStatusDraft    = "draft"    // 草稿
	AIAgentStatusEnabled  = "enabled"  // 已启用
	AIAgentStatusDisabled = "disabled" // 已禁用
	AIAgentStatusArchived = "archived" // 已归档
)

// AI建议评价状态
const (
	AISuggestionFeedbackStatusPending  = "pending"  // 待审核
	AISuggestionFeedbackStatusApproved = "approved" // 已通过
	AISuggestionFeedbackStatusRejected = "rejected" // 已拒绝
)

// AI建议评价标签
const (
	AISuggestionFeedbackTagHelpful          = "helpful"            // 有帮助
	AISuggestionFeedbackTagProfessional     = "professional"       // 专业
	AISuggestionFeedbackTagEasyToUnderstand = "easy_to_understand" // 易懂
	AISuggestionFeedbackTagAccurate         = "accurate"           // 准确
	AISuggestionFeedbackTagInnovative       = "innovative"         // 创新
	AISuggestionFeedbackTagTimely           = "timely"             // 及时
)

// AI建议评价类型
const (
	AISuggestionFeedbackRatingLike    = 1  // 点赞
	AISuggestionFeedbackRatingDislike = -1 // 踩
)

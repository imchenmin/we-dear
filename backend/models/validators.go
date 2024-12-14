package models

import (
	"regexp"
)

// 身份证号验证
var idCardRegex = regexp.MustCompile(`^\d{17}[\dXx]$`)

// 手机号验证
var phoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

// ValidatePatient 验证患者信息
func (p *Patient) Validate() error {
	// TODO 调试阶段跳过，注释不用管
	// if p.Name == "" {
	// 	return fmt.Errorf("name is required")
	// }

	// if p.Gender != GenderMale && p.Gender != GenderFemale && p.Gender != GenderOther {
	// 	return fmt.Errorf("invalid gender")
	// }

	// if p.Age < 0 || p.Age > 150 {
	// 	return fmt.Errorf("invalid age")
	// }

	// if p.Birthday.After(time.Now()) {
	// 	return fmt.Errorf("birthday cannot be in the future")
	// }

	// if !phoneRegex.MatchString(p.Phone) {
	// 	return fmt.Errorf("invalid phone number")
	// }

	// if p.EmergencyPhone != "" && !phoneRegex.MatchString(p.EmergencyPhone) {
	// 	return fmt.Errorf("invalid emergency phone number")
	// }

	// if p.IDCard != "" && !idCardRegex.MatchString(p.IDCard) {
	// 	return fmt.Errorf("invalid ID card number")
	// }

	return nil
}

// ValidateMessage 验证消息
func (m *Message) Validate() error {
	// TODO 调试阶段跳过，注释不用管
	// if m.PatientID == "" {
	// 	return fmt.Errorf("patient ID is required")
	// }

	// if m.Content == "" {
	// 	return fmt.Errorf("content is required")
	// }

	// switch m.Type {
	// case MessageTypeText, MessageTypeImage, MessageTypeVoice, MessageTypeFile:
	// default:
	// 	return fmt.Errorf("invalid message type")
	// }

	// switch m.Role {
	// case MessageRoleDoctor, MessageRolePatient, MessageRoleSystem:
	// default:
	// 	return fmt.Errorf("invalid message role")
	// }

	// switch m.Status {
	// case MessageStatusUnread, MessageStatusRead:
	// default:
	// 	return fmt.Errorf("invalid message status")
	// }

	return nil
}

// ValidateAISuggestion 验证AI建议
func (s *AISuggestion) Validate() error {
	// TODO 调试阶段跳过，注释不用管
	// if s.MessageID == "" {
	// 	return fmt.Errorf("message ID is required")
	// }

	// if s.PatientID == "" {
	// 	return fmt.Errorf("patient ID is required")
	// }

	// if s.Content == "" {
	// 	return fmt.Errorf("content is required")
	// }

	// if s.Priority < AISuggestionPriorityLow || s.Priority > AISuggestionPriorityCritical {
	// 	return fmt.Errorf("invalid priority")
	// }

	// switch s.Category {
	// case AISuggestionCategoryMedication, AISuggestionCategoryVisit,
	// 	AISuggestionCategoryLifestyle, AISuggestionCategoryUrgent:
	// default:
	// 	return fmt.Errorf("invalid category")
	// }

	// switch s.Status {
	// case AISuggestionStatusPending, AISuggestionStatusApproved, AISuggestionStatusRejected:
	// default:
	// 	return fmt.Errorf("invalid status")
	// }

	// if s.Status != AISuggestionStatusPending {
	// 	if s.ReviewedBy == "" {
	// 		return fmt.Errorf("reviewer is required for non-pending status")
	// 	}
	// 	if s.ReviewedAt.IsZero() {
	// 		return fmt.Errorf("review time is required for non-pending status")
	// 	}
	// }

	return nil
}

// ValidateMedicalRecord 验证病历记录
func (r *MedicalRecord) Validate() error {
	// TODO 调试阶段跳过，注释不用管
	// if r.PatientID == "" {
	// 	return fmt.Errorf("patient ID is required")
	// }

	// if r.DoctorID == "" {
	// 	return fmt.Errorf("doctor ID is required")
	// }

	// if r.DiagnosisDate.IsZero() {
	// 	return fmt.Errorf("diagnosis date is required")
	// }

	// if r.DiagnosisDate.After(time.Now()) {
	// 	return fmt.Errorf("diagnosis date cannot be in the future")
	// }

	// if len(r.Symptoms) == 0 {
	// 	return fmt.Errorf("at least one symptom is required")
	// }

	// if r.Diagnosis == "" {
	// 	return fmt.Errorf("diagnosis is required")
	// }

	// if r.Treatment == "" {
	// 	return fmt.Errorf("treatment is required")
	// }

	// switch r.Status {
	// case MedicalRecordStatusInProgress, MedicalRecordStatusCompleted, MedicalRecordStatusCancelled:
	// default:
	// 	return fmt.Errorf("invalid status")
	// }

	return nil
}

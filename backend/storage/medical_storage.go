package storage

import (
	"we-dear/config"
	"we-dear/models"

	"gorm.io/gorm"
)

type MedicalStorage struct {
	db *gorm.DB
}

var medicalInstance *MedicalStorage

func GetMedicalStorage() *MedicalStorage {
	if medicalInstance == nil {
		medicalInstance = &MedicalStorage{
			db: config.DB,
		}
	}
	return medicalInstance
}

// 获取患者的所有随访记录
func (s *MedicalStorage) GetFollowUpRecords(patientID string) ([]models.FollowUpRecord, error) {
	var records []models.FollowUpRecord
	err := s.db.Where("patient_id = ?", patientID).Order("follow_up_date desc").Find(&records).Error
	return records, err
}

// 创建随访记录
func (s *MedicalStorage) CreateFollowUpRecord(record *models.FollowUpRecord) error {
	return s.db.Create(record).Error
}

// 更新随访记录
func (s *MedicalStorage) UpdateFollowUpRecord(record *models.FollowUpRecord) error {
	return s.db.Save(record).Error
}

// 删除随访记录
func (s *MedicalStorage) DeleteFollowUpRecord(id string) error {
	return s.db.Delete(&models.FollowUpRecord{}, "id = ?", id).Error
}

// 获取患者的所有医疗记录
func (s *MedicalStorage) GetMedicalRecords(patientID string) ([]models.MedicalRecord, error) {
	var records []models.MedicalRecord
	err := s.db.Where("patient_id = ?", patientID).Order("diagnosis_date desc").Find(&records).Error
	return records, err
}

// 创建医疗记录
func (s *MedicalStorage) CreateMedicalRecord(record *models.MedicalRecord) error {
	return s.db.Create(record).Error
}

// 更新医疗记录
func (s *MedicalStorage) UpdateMedicalRecord(record *models.MedicalRecord) error {
	return s.db.Save(record).Error
}

// 删除医疗记录
func (s *MedicalStorage) DeleteMedicalRecord(id string) error {
	return s.db.Delete(&models.MedicalRecord{}, "id = ?", id).Error
}

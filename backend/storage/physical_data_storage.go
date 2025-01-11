package storage

import (
	"we-dear/config"
	"we-dear/models"

	"gorm.io/gorm"
)

type PhysiologicalDataStorage struct {
	db *gorm.DB
}

var physiologicalDataInstance *PhysiologicalDataStorage

func GetPhysiologicalDataStorage() *PhysiologicalDataStorage {
	if physiologicalDataInstance == nil {
		physiologicalDataInstance = &PhysiologicalDataStorage{
			db: config.DB,
		}
	}
	return physiologicalDataInstance
}

// Create 创建生理数据记录
func (s *PhysiologicalDataStorage) Create(data *models.PhysiologicalData) error {
	return s.db.Create(data).Error
}

// GetByPatientID 获取患者的生理数据记录
func (s *PhysiologicalDataStorage) GetByPatientID(patientID string) ([]models.PhysiologicalData, error) {
	var records []models.PhysiologicalData
	err := s.db.Where("patient_id = ?", patientID).Order("measured_at desc").Find(&records).Error
	return records, err
}

// GetByPatientIDAndType 获取患者特定类型的生理数据记录
func (s *PhysiologicalDataStorage) GetByPatientIDAndType(patientID string, dataType string) ([]models.PhysiologicalData, error) {
	var records []models.PhysiologicalData
	err := s.db.Where("patient_id = ? AND type = ?", patientID, dataType).Order("measured_at desc").Find(&records).Error
	return records, err
}

// Update 更新生理数据记录
func (s *PhysiologicalDataStorage) Update(data *models.PhysiologicalData) error {
	return s.db.Save(data).Error
}

// Delete 删除生理数据记录
func (s *PhysiologicalDataStorage) Delete(id string) error {
	return s.db.Delete(&models.PhysiologicalData{}, "id = ?", id).Error
}
package storage

import (
	"errors"
	"sync"
	"we-dear/config"
	"we-dear/models"

	"gorm.io/gorm"
)

type DoctorStorage struct {
	db *gorm.DB
}

var (
	doctorInstance *DoctorStorage
	doctorOnce     sync.Once
)

func GetDoctorStorage() *DoctorStorage {
	doctorOnce.Do(func() {
		doctorInstance = &DoctorStorage{
			db: config.DB,
		}
	})
	return doctorInstance
}

func (s *DoctorStorage) CreateDoctor(doctor *models.Doctor) error {
	return s.db.Create(doctor).Error
}

func (s *DoctorStorage) GetAllDoctors() ([]models.Doctor, error) {
	var doctors []models.Doctor
	err := s.db.Preload("Department").Find(&doctors).Error
	return doctors, err
}

func (s *DoctorStorage) GetDoctorByID(id string) (*models.Doctor, error) {
	var doctor models.Doctor
	err := s.db.Preload("Department").First(&doctor, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("doctor not found")
		}
		return nil, err
	}
	return &doctor, nil
}

func (s *DoctorStorage) UpdateDoctor(doctor *models.Doctor) error {
	return s.db.Save(doctor).Error
}

func (s *DoctorStorage) DeleteDoctor(id string) error {
	return s.db.Delete(&models.Doctor{}, "id = ?", id).Error
}

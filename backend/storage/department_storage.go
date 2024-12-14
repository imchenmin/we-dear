package storage

import (
	"errors"
	"sync"
	"we-dear/config"
	"we-dear/models"

	"gorm.io/gorm"
)

type DepartmentStorage struct {
	db *gorm.DB
}

var (
	departmentInstance *DepartmentStorage
	departmentOnce     sync.Once
)

func GetDepartmentStorage() *DepartmentStorage {
	departmentOnce.Do(func() {
		departmentInstance = &DepartmentStorage{
			db: config.DB,
		}
	})
	return departmentInstance
}

func (s *DepartmentStorage) CreateDepartment(department *models.Department) error {
	return s.db.Create(department).Error
}

func (s *DepartmentStorage) GetAllDepartments() ([]models.Department, error) {
	var departments []models.Department
	err := s.db.Find(&departments).Error
	return departments, err
}

func (s *DepartmentStorage) GetDepartmentByID(id string) (*models.Department, error) {
	var department models.Department
	err := s.db.First(&department, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("department not found")
		}
		return nil, err
	}
	return &department, nil
}

func (s *DepartmentStorage) UpdateDepartment(department *models.Department) error {
	return s.db.Save(department).Error
}

func (s *DepartmentStorage) DeleteDepartment(id string) error {
	return s.db.Delete(&models.Department{}, "id = ?", id).Error
}

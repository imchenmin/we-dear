package storage

import (
	"errors"
	"we-dear/config"
	"we-dear/models"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type PatientStorage struct {
	db *gorm.DB
}

var patientInstance *PatientStorage

func GetPatientStorage() *PatientStorage {
	if patientInstance == nil {
		patientInstance = &PatientStorage{
			db: config.DB,
		}
	}
	return patientInstance
}

func (s *PatientStorage) CreatePatient(patient *models.Patient) error {
	if patient.Allergies == nil {
		patient.Allergies = pq.StringArray{}
	}
	if patient.ChronicDiseases == nil {
		patient.ChronicDiseases = pq.StringArray{}
	}

	return s.db.Create(patient).Error
}

func (s *PatientStorage) GetAllPatients() []models.Patient {
	var patients []models.Patient
	s.db.Preload("Doctor").Find(&patients)
	return patients
}

func (s *PatientStorage) GetPatientByID(id string) (*models.Patient, error) {
	var patient models.Patient
	err := s.db.Preload("Doctor").First(&patient, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("patient not found")
		}
		return nil, err
	}
	return &patient, nil
}

func (s *PatientStorage) GetChatHistory(patientID string) ([]models.Message, error) {
	var messages []models.Message
	err := s.db.Where("patient_id = ?", patientID).Order("created_at asc").Find(&messages).Error
	return messages, err
}

func (s *PatientStorage) AddMessage(patientID string, message models.Message) error {
	return s.db.Create(&message).Error
}

func (s *PatientStorage) GetAISuggestions(patientID string, messageID string) ([]models.AISuggestion, error) {
	var suggestions []models.AISuggestion
	query := s.db.Where("patient_id = ?", patientID)
	if messageID != "" {
		query = query.Where("message_id = ?", messageID)
	}
	err := query.Find(&suggestions).Error
	return suggestions, err
}

func (s *PatientStorage) SaveAISuggestion(suggestion *models.AISuggestion) error {
	return s.db.Create(suggestion).Error
}

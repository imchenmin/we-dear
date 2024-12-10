package storage

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"we-dear/models"
)

type PatientStorage struct {
	mu            sync.RWMutex
	patients      []models.Patient
	messages      map[string][]models.Message      // patientId -> messages
	aiSuggestions map[string][]models.AISuggestion // messageId -> suggestions
}

var (
	instance *PatientStorage
	once     sync.Once
)

func GetPatientStorage() *PatientStorage {
	once.Do(func() {
		instance = &PatientStorage{
			messages:      make(map[string][]models.Message),
			aiSuggestions: make(map[string][]models.AISuggestion),
		}
		instance.loadData()
	})
	return instance
}

func (s *PatientStorage) loadData() {
	s.loadPatients()
	s.loadMessages()
	s.loadAISuggestions()
}

func (s *PatientStorage) loadPatients() {
	file, err := os.Open("data/patients.csv")
	if err != nil {
		log.Printf("Warning: Could not open patients.csv: %v", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Error reading patients.csv: %v", err)
		return
	}

	for _, record := range records[1:] { // Skip header
		age, _ := strconv.Atoi(record[3])
		allergies := strings.Split(record[10], "|")
		if len(allergies) == 1 && allergies[0] == "" {
			allergies = []string{}
		}
		chronicDiseases := strings.Split(record[11], "|")
		if len(chronicDiseases) == 1 && chronicDiseases[0] == "" {
			chronicDiseases = []string{}
		}

		birthday, _ := time.Parse("2006-01-02", record[4])
		createdAt, _ := time.Parse(time.RFC3339, record[1])
		updatedAt, _ := time.Parse(time.RFC3339, record[2])

		patient := models.Patient{
			BaseModel: models.BaseModel{
				ID:        record[0],
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			Name:            record[3],
			Gender:          record[4],
			Age:             age,
			Birthday:        birthday,
			Phone:           record[5],
			EmergencyPhone:  record[6],
			Address:         record[7],
			IDCard:          record[8],
			BloodType:       record[9],
			Allergies:       allergies,
			ChronicDiseases: chronicDiseases,
			Avatar:          record[12],
		}
		s.patients = append(s.patients, patient)
	}
}

func (s *PatientStorage) loadMessages() {
	file, err := os.Open("data/messages.csv")
	if err != nil {
		log.Printf("Warning: Could not open messages.csv: %v", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Error reading messages.csv: %v", err)
		return
	}

	for _, record := range records[1:] { // Skip header
		createdAt, _ := time.Parse(time.RFC3339, record[1])
		updatedAt, _ := time.Parse(time.RFC3339, record[2])

		message := models.Message{
			BaseModel: models.BaseModel{
				ID:        record[0],
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			PatientID: record[3],
			DoctorID:  record[4],
			RecordID:  record[5],
			Content:   record[6],
			Type:      record[7],
			Role:      record[8],
			Status:    record[9],
			ReplyTo:   record[10],
		}

		patientID := message.PatientID
		s.messages[patientID] = append(s.messages[patientID], message)
	}
}

func (s *PatientStorage) loadAISuggestions() {
	file, err := os.Open("data/ai_suggestions.csv")
	if err != nil {
		log.Printf("Warning: Could not open ai_suggestions.csv: %v", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Error reading ai_suggestions.csv: %v", err)
		return
	}

	for _, record := range records[1:] { // Skip header
		createdAt, _ := time.Parse(time.RFC3339, record[1])
		updatedAt, _ := time.Parse(time.RFC3339, record[2])
		reviewedAt, _ := time.Parse(time.RFC3339, record[11])
		confidence, _ := strconv.ParseFloat(record[7], 64)
		priority, _ := strconv.Atoi(record[9])

		suggestion := models.AISuggestion{
			BaseModel: models.BaseModel{
				ID:        record[0],
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			MessageID:   record[3],
			PatientID:   record[4],
			Content:     record[5],
			ModelUsed:   record[6],
			Confidence:  confidence,
			Category:    record[8],
			Priority:    priority,
			Status:      record[10],
			ReviewedBy:  record[12],
			ReviewedAt:  reviewedAt,
			ReviewNotes: record[13],
		}

		messageID := suggestion.MessageID
		s.aiSuggestions[messageID] = append(s.aiSuggestions[messageID], suggestion)
	}
}

func (s *PatientStorage) GetAllPatients() []models.Patient {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.patients
}

func (s *PatientStorage) GetPatientByID(id string) (*models.Patient, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, p := range s.patients {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("patient not found")
}

func (s *PatientStorage) GetChatHistory(patientID string) ([]models.Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	messages, ok := s.messages[patientID]
	if !ok {
		return nil, fmt.Errorf("no messages found for patient")
	}
	return messages, nil
}

func (s *PatientStorage) AddMessage(patientID string, message models.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.messages[patientID] = append(s.messages[patientID], message)
	return s.saveMessages()
}

func (s *PatientStorage) GetAISuggestions(patientID string, messageID string) ([]models.AISuggestion, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if suggestions, ok := s.aiSuggestions[messageID]; ok {
		return suggestions, nil
	}
	return nil, nil
}

func (s *PatientStorage) SaveAISuggestion(suggestion *models.AISuggestion) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	messageID := suggestion.MessageID
	s.aiSuggestions[messageID] = append(s.aiSuggestions[messageID], *suggestion)
	return s.saveAISuggestions()
}

// Save data to CSV files
func (s *PatientStorage) savePatients() error {
	return s.saveToCSV("data/patients.csv", []string{
		"ID", "CreatedAt", "UpdatedAt", "Name", "Gender", "Age", "Birthday",
		"Phone", "EmergencyPhone", "Address", "IDCard", "BloodType",
		"Allergies", "ChronicDiseases", "Avatar",
	}, func(w *csv.Writer) error {
		for _, p := range s.patients {
			record := []string{
				p.ID,
				p.CreatedAt.Format(time.RFC3339),
				p.UpdatedAt.Format(time.RFC3339),
				p.Name,
				p.Gender,
				strconv.Itoa(p.Age),
				p.Birthday.Format("2006-01-02"),
				p.Phone,
				p.EmergencyPhone,
				p.Address,
				p.IDCard,
				p.BloodType,
				strings.Join(p.Allergies, "|"),
				strings.Join(p.ChronicDiseases, "|"),
				p.Avatar,
			}
			if err := w.Write(record); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *PatientStorage) saveMessages() error {
	return s.saveToCSV("data/messages.csv", []string{
		"ID", "CreatedAt", "UpdatedAt", "PatientID", "DoctorID",
		"RecordID", "Content", "Type", "Role", "Status", "ReplyTo",
	}, func(w *csv.Writer) error {
		for _, messages := range s.messages {
			for _, m := range messages {
				record := []string{
					m.ID,
					m.CreatedAt.Format(time.RFC3339),
					m.UpdatedAt.Format(time.RFC3339),
					m.PatientID,
					m.DoctorID,
					m.RecordID,
					m.Content,
					m.Type,
					m.Role,
					m.Status,
					m.ReplyTo,
				}
				if err := w.Write(record); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *PatientStorage) saveAISuggestions() error {
	return s.saveToCSV("data/ai_suggestions.csv", []string{
		"ID", "CreatedAt", "UpdatedAt", "MessageID", "PatientID",
		"Content", "ModelUsed", "Confidence", "Category", "Priority",
		"Status", "ReviewedAt", "ReviewedBy", "ReviewNotes",
	}, func(w *csv.Writer) error {
		for _, suggestions := range s.aiSuggestions {
			for _, s := range suggestions {
				record := []string{
					s.ID,
					s.CreatedAt.Format(time.RFC3339),
					s.UpdatedAt.Format(time.RFC3339),
					s.MessageID,
					s.PatientID,
					s.Content,
					s.ModelUsed,
					strconv.FormatFloat(s.Confidence, 'f', 2, 64),
					s.Category,
					strconv.Itoa(s.Priority),
					s.Status,
					s.ReviewedAt.Format(time.RFC3339),
					s.ReviewedBy,
					s.ReviewNotes,
				}
				if err := w.Write(record); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *PatientStorage) saveToCSV(filename string, header []string, writeFunc func(*csv.Writer) error) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(header); err != nil {
		return err
	}

	return writeFunc(writer)
}

package storage

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"we-dear/models"
)

type PatientStorage struct {
	patients        []models.Patient
	mu              sync.RWMutex
	messages        map[string][]models.Message
	suggestions     map[string][]models.AISuggestion // patientID -> suggestions
	messagesPath    string
	suggestionsPath string
}

var instance *PatientStorage
var once sync.Once

func GetPatientStorage() *PatientStorage {
	once.Do(func() {
		instance = &PatientStorage{
			messages:        make(map[string][]models.Message),
			suggestions:     make(map[string][]models.AISuggestion),
			messagesPath:    filepath.Join("data", "messages.json"),
			suggestionsPath: filepath.Join("data", "ai_suggestions.json"),
		}
		instance.loadPatientsFromCSV()
		instance.loadAllMessages()
		instance.loadAllSuggestions()
	})
	return instance
}

func (s *PatientStorage) loadPatientsFromCSV() error {
	file, err := os.Open("../data/patients.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Skip header
	_, err = reader.Read()
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		age, _ := strconv.Atoi(record[3])
		patient := models.Patient{
			ID:        record[0],
			Name:      record[1],
			Gender:    record[2],
			Age:       age,
			Phone:     record[4],
			Diagnosis: record[5],
			Doctor:    record[6],
			Avatar:    record[7],
			Messages:  []models.Message{},
		}
		s.patients = append(s.patients, patient)
	}
	return nil
}

func (s *PatientStorage) GetAllPatients() []models.Patient {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.patients
}

func (s *PatientStorage) GetPatientByID(id string) (*models.Patient, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := range s.patients {
		if s.patients[i].ID == id {
			return &s.patients[i], nil
		}
	}
	return nil, fmt.Errorf("patient not found")
}

func (s *PatientStorage) AddMessage(patientID string, message models.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.patients {
		if s.patients[i].ID == patientID {
			if message.Timestamp.IsZero() {
				message.Timestamp = time.Now()
			}
			message.UnixTime = message.Timestamp.UnixMilli()

			s.patients[i].Messages = append(s.patients[i].Messages, message)
			s.messages[patientID] = s.patients[i].Messages
			go s.saveMessages()
			return nil
		}
	}
	return fmt.Errorf("patient not found")
}

func (s *PatientStorage) GetMessages(id string) []models.Message {
	s.mu.RLock()
	defer s.mu.RUnlock()
	fmt.Println(s.messages)
	return s.messages[id]
}

func (s *PatientStorage) loadAllMessages() error {
	data, err := os.ReadFile(s.messagesPath)
	if err != nil {
		if os.IsNotExist(err) {
			s.messages = make(map[string][]models.Message)
			return nil
		}
		return err
	}

	if err := json.Unmarshal(data, &s.messages); err != nil {
		return err
	}

	for i := range s.patients {
		if messages, ok := s.messages[s.patients[i].ID]; ok {
			for j := range messages {
				if messages[j].Timestamp.IsZero() {
					messages[j].Timestamp = time.Unix(0, messages[j].UnixTime*int64(time.Millisecond))
				}
			}
			s.patients[i].Messages = messages
		}
	}
	return nil
}

func (s *PatientStorage) saveMessages() error {
	data, err := json.Marshal(s.messages)
	if err != nil {
		return err
	}
	return os.WriteFile(s.messagesPath, data, 0644)
}

func (s *PatientStorage) GetChatHistory(patientID string) ([]models.Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := range s.patients {
		if s.patients[i].ID == patientID {
			return s.patients[i].Messages, nil
		}
	}
	return nil, fmt.Errorf("patient not found")
}

func (s *PatientStorage) SaveAISuggestion(suggestion *models.AISuggestion) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.suggestions == nil {
		s.suggestions = make(map[string][]models.AISuggestion)
	}

	s.suggestions[suggestion.PatientID] = append(s.suggestions[suggestion.PatientID], *suggestion)
	return s.saveSuggestions()
}

func (s *PatientStorage) GetAISuggestions(patientID string, messageID string) ([]models.AISuggestion, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	suggestions := s.suggestions[patientID]
	if messageID != "" {
		// 如果指定了消息ID，只返回该消息的建议
		var filtered []models.AISuggestion
		for _, sug := range suggestions {
			if sug.MessageID == messageID {
				filtered = append(filtered, sug)
			}
		}
		return filtered, nil
	}
	return suggestions, nil
}

func (s *PatientStorage) loadAllSuggestions() error {
	data, err := os.ReadFile(s.suggestionsPath)
	if err != nil {
		if os.IsNotExist(err) {
			s.suggestions = make(map[string][]models.AISuggestion)
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &s.suggestions)
}

func (s *PatientStorage) saveSuggestions() error {
	data, err := json.Marshal(s.suggestions)
	if err != nil {
		return err
	}
	return os.WriteFile(s.suggestionsPath, data, 0644)
}

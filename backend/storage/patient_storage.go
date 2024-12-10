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
	patients     []models.Patient
	mu           sync.RWMutex
	messages     map[string][]models.Message
	messagesPath string
}

var instance *PatientStorage
var once sync.Once

func GetPatientStorage() *PatientStorage {
	once.Do(func() {
		instance = &PatientStorage{
			messages:     make(map[string][]models.Message),
			messagesPath: filepath.Join("../data", "messages.json"),
		}
		instance.loadPatientsFromCSV()
		instance.loadAllMessages()
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

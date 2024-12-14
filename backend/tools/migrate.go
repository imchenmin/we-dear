package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"we-dear/config"
	"we-dear/models"
	"we-dear/utils"

	"github.com/lib/pq"
)

func main() {
	// 初始化数据库连接
	config.Init()
	config.InitDB()

	// 迁移数据
	if err := migratePatients(); err != nil {
		log.Printf("Failed to migrate patients: %v", err)
	}
	if err := migrateMessages(); err != nil {
		log.Printf("Failed to migrate messages: %v", err)
	}
	if err := migrateAISuggestions(); err != nil {
		log.Printf("Failed to migrate AI suggestions: %v", err)
	}
}

func migratePatients() error {
	file, err := os.Open("../data/patients.csv")
	if err != nil {
		return fmt.Errorf("failed to open patients.csv: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %v", err)
	}

	// Skip header
	for _, record := range records[1:] {
		age, _ := strconv.Atoi(record[3])
		birthday, _ := time.Parse("2006-01-02", record[4])
		createdAt, _ := time.Parse(time.RFC3339, record[1])
		updatedAt, _ := time.Parse(time.RFC3339, record[2])

		patient := models.Patient{
			BaseModel: models.BaseModel{
				ID:        utils.GenerateID(),
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
			Allergies:       pq.StringArray(strings.Split(record[10], "|")),
			ChronicDiseases: pq.StringArray(strings.Split(record[11], "|")),
			Avatar:          record[12],
		}

		if err := config.DB.Create(&patient).Error; err != nil {
			return fmt.Errorf("failed to create patient: %v", err)
		}
	}

	return nil
}

func migrateMessages() error {
	file, err := os.Open("../data/messages.csv")
	if err != nil {
		return fmt.Errorf("failed to open messages.csv: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %v", err)
	}

	for _, record := range records[1:] {
		createdAt, _ := time.Parse(time.RFC3339, record[1])
		updatedAt, _ := time.Parse(time.RFC3339, record[2])

		message := models.Message{
			BaseModel: models.BaseModel{
				ID:        utils.GenerateID(),
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

		if err := config.DB.Create(&message).Error; err != nil {
			return fmt.Errorf("failed to create message: %v", err)
		}
	}

	return nil
}

func migrateAISuggestions() error {
	file, err := os.Open("../data/ai_suggestions.csv")
	if err != nil {
		return fmt.Errorf("failed to open ai_suggestions.csv: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %v", err)
	}

	for _, record := range records[1:] {
		createdAt, _ := time.Parse(time.RFC3339, record[1])
		updatedAt, _ := time.Parse(time.RFC3339, record[2])
		reviewedAt, _ := time.Parse(time.RFC3339, record[11])
		confidence, _ := strconv.ParseFloat(record[7], 64)
		priority, _ := strconv.Atoi(record[9])

		suggestion := models.AISuggestion{
			BaseModel: models.BaseModel{
				ID:        utils.GenerateID(),
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

		if err := config.DB.Create(&suggestion).Error; err != nil {
			return fmt.Errorf("failed to create AI suggestion: %v", err)
		}
	}

	return nil
}

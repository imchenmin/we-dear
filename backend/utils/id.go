package utils

import (
	"github.com/google/uuid"
)

// GenerateID 生成唯一ID
func GenerateID() string {
	uuid := uuid.New()
	return uuid.String()
}

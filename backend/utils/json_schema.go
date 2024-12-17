package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/xeipuuv/gojsonschema"
)

// ValidateJSONSchema 验证数据是否符合JSON Schema
func ValidateJSONSchema(schema string, data string) (bool, error) {
	// 解析schema确保是有效的JSON
	var schemaObj interface{}
	if err := json.Unmarshal([]byte(schema), &schemaObj); err != nil {
		return false, fmt.Errorf("无效的schema格式: %v", err)
	}

	// 创建schema加载器
	schemaLoader := gojsonschema.NewGoLoader(schemaObj)
	documentLoader := gojsonschema.NewStringLoader(data)

	// 验证数据
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return false, fmt.Errorf("验证过程出错: %v", err)
	}

	// 如果验证失败，收集错误信息
	if !result.Valid() {
		var errorMessages []string
		for _, err := range result.Errors() {
			errorMessages = append(errorMessages, err.String())
		}
		return false, fmt.Errorf("验证失败: %v", errorMessages)
	}

	return true, nil
}

// GetDefaultSchema 获取默认的随访模板Schema, 从default_template.json中获取
func GetDefaultSchema() string {
	schema, err := os.ReadFile("config/default_template.json")
	if err != nil {
		return ""
	}
	return string(schema)
}

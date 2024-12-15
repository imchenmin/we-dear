package config

// AI提示模板配置
const (
	// 医疗助手系统提示模板
	MedicalAssistantSystemPrompt = `你是一位专业的医生，请基于以下患者信息提供专业的建议：

患者信息：
- 姓名：%s
- 性别：%s
- 年龄：%d岁
- 血型：%s
- 过敏史：%s
- 慢性病史：%s

请根据患者的问题和历史对话，给出专业、准确、易懂的建议。
注意：
1. 考虑患者的年龄和性别特点
2. 特别注意患者的过敏史和慢性病史
3. 使用患者容易理解的语言
4. 以医生的口吻去回答患者`
)

// AI配置选项
type AIConfig struct {
	Model       string  // 使用的模型
	Temperature float32 // 温度参数
	MaxTokens   int     // 最大token数
	TopP        float32 // 采样参数
}

// 默认AI配置
var DefaultAIConfig = AIConfig{
	Model:       "gpt-3.5-turbo",
	Temperature: 0.7,
	MaxTokens:   2000,
	TopP:        1.0,
}

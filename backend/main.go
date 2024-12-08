package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

type Patient struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	Age       int       `json:"age"`
	Phone     string    `json:"phone"`
	Diagnosis string    `json:"diagnosis"`
	Doctor    string    `json:"doctor"`
	Avatar    string    `json:"avatar"`
	Messages  []Message `json:"messages"`
}

type Message struct {
	Type         string    `json:"type"`
	ContentType  string    `json:"contentType"`
	Content      string    `json:"content"`
	Timestamp    time.Time `json:"timestamp"`
	AISuggestion string    `json:"aiSuggestion,omitempty"`
}

var (
	patients     []Patient
	openAIClient *openai.Client
	mu           sync.Mutex
)

// 消息持久化相关函数
func loadMessagesFromFile() error {
	// 确保messages目录存在
	if err := os.MkdirAll("messages", 0755); err != nil {
		return fmt.Errorf("failed to create messages directory: %v", err)
	}

	// 遍历所有患者
	for i := range patients {
		filename := fmt.Sprintf("messages/patient_%s.json", patients[i].ID)

		// 如果消息文件存在，则加载消息
		if _, err := os.Stat(filename); err == nil {
			data, err := os.ReadFile(filename)
			if err != nil {
				log.Printf("Error reading messages for patient %s: %v", patients[i].ID, err)
				continue
			}

			var messages []Message
			if err := json.Unmarshal(data, &messages); err != nil {
				log.Printf("Error parsing messages for patient %s: %v", patients[i].ID, err)
				continue
			}

			patients[i].Messages = messages
		}
	}
	return nil
}

func saveMessagesToFile(patientID string) error {
	mu.Lock()
	defer mu.Unlock()

	// 查找患者
	var patient *Patient
	for i := range patients {
		if patients[i].ID == patientID {
			patient = &patients[i]
			break
		}
	}

	if patient == nil {
		return fmt.Errorf("patient not found: %s", patientID)
	}

	// 将消息保存到JSON文件
	filename := fmt.Sprintf("messages/patient_%s.json", patientID)
	data, err := json.MarshalIndent(patient.Messages, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling messages: %v", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("error writing messages file: %v", err)
	}

	return nil
}

// 更新患者消息并持久化
func updatePatientMessages(id string, message Message) error {
	mu.Lock()
	defer mu.Unlock()

	for i := range patients {
		if patients[i].ID == id {
			patients[i].Messages = append(patients[i].Messages, message)
			// 异步保存消息
			go func(pid string) {
				if err := saveMessagesToFile(pid); err != nil {
					log.Printf("Error saving messages for patient %s: %v", pid, err)
				}
			}(id)
			return nil
		}
	}
	return fmt.Errorf("patient not found")
}

func loadPatientsFromCSV() error {
	file, err := os.Open("../data/patients.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// 跳过标题行
	_, err = reader.Read()
	if err != nil {
		return err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		age, _ := strconv.Atoi(record[3])
		patient := Patient{
			ID:        record[0],
			Name:      record[1],
			Gender:    record[2],
			Age:       age,
			Phone:     record[4],
			Diagnosis: record[5],
			Doctor:    record[6],
			Avatar:    record[7],
			Messages:  []Message{},
		}
		patients = append(patients, patient)
	}

	// 加载历史消息
	if err := loadMessagesFromFile(); err != nil {
		log.Printf("Warning: failed to load messages: %v", err)
	}

	return nil
}

func initOpenAI() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is not set")
	}
	openAIClient = openai.NewClient(apiKey)
}

func main() {
	// 初始化OpenAI客户端
	initOpenAI()

	if err := loadPatientsFromCSV(); err != nil {
		log.Fatal("Failed to load patients data:", err)
	}

	r := gin.Default()

	// 配置CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	r.Use(cors.New(config))

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	// API路由
	api := r.Group("/api")
	{
		// 医生端API
		api.GET("/patients", getPatients)
		api.GET("/patients/:id", getPatient)
		api.POST("/patients/:id/messages", addMessage)
		api.POST("/upload", handleUpload)

		// 患者端API
		api.GET("/patient/:id/messages", getPatientMessages)
		api.POST("/patient/:id/question", addPatientQuestion)
		api.POST("/patient/upload", handlePatientUpload)
	}

	r.Run(":8080")
}

// 医生端API实现
func getPatients(c *gin.Context) {
	c.JSON(http.StatusOK, patients)
}

func getPatient(c *gin.Context) {
	id := c.Param("id")
	for _, p := range patients {
		if p.ID == id {
			c.JSON(http.StatusOK, p)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
}

func addMessage(c *gin.Context) {
	id := c.Param("id")
	var message Message
	if err := c.BindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message.Timestamp = time.Now()
	message.Type = "doctor" // 确保消息类型为医生

	if err := updatePatientMessages(id, message); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	c.JSON(http.StatusOK, message)
}

// 患者端API实现
func getPatientMessages(c *gin.Context) {
	id := c.Param("id")
	for _, p := range patients {
		if p.ID == id {
			c.JSON(http.StatusOK, gin.H{
				"messages": p.Messages,
				"patient":  p,
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
}

func addPatientQuestion(c *gin.Context) {
	id := c.Param("id")
	var message Message
	if err := c.BindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message.Type = "patient" // 强制设置为患者消息
	message.Timestamp = time.Now()
	message.AISuggestion = generateAISuggestion(message.Content)

	if err := updatePatientMessages(id, message); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      message,
		"aiSuggestion": message.AISuggestion,
	})
}

func handleUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 确保上传目录存在
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	filepath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": fmt.Sprintf("/uploads/%s", filename),
	})
}

// 患者文件上传处理（与医生端共用同一个处理函数）
func handlePatientUpload(c *gin.Context) {
	handleUpload(c)
}

func generateAISuggestion(content string) string {
	if openAIClient == nil {
		return "AI服务未初始化，请检查配置"
	}

	// 构建提示信息
	prompt := fmt.Sprintf(`作为一名专业的医生助手，请针对以下患者的问题提供专业的建议回复。
问题内容：%s

请提供具体的、有针对性的回复建议，包括：
1. 需要询问的具体症状
2. 可能需要的检查项目
3. 初步建议和注意事项

请以"建议回复："开头，使用分点方式回答。`, content)

	resp, err := openAIClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "你是一个专业的医生助手，负责为医生提供回复患者问题的建议。请使用专业、关心的语气，并确保建议的准确性和实用性。",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.7,
			MaxTokens:   500,
		},
	)

	if err != nil {
		log.Printf("OpenAI API error: %v", err)
		return "AI服务暂时不可用，请使用默认建议模板"
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content
	}

	// 如果AI服务失败，使用基础模板
	return getBasicSuggestion(content)
}

// 基础建议模板
func getBasicSuggestion(content string) string {
	suggestions := map[string]string{
		"血糖": "建议回复：1. 请告知您最近的血糖监测值\n2. 您是否按时服用降糖药物？\n3. 最近的饮食情况如何？\n4. 有进行运动锻炼吗？",
		"血压": "建议回复：1. 请提供您最近的血压测量值\n2. 是否规律服用降压药？\n3. 有注意限制盐分摄入吗？\n4. 作息时间是否规律？",
		"胸痛": "建议回复：1. 疼痛的具体位置在哪里？\n2. 是否伴有出汗、呼吸困难等症状？\n3. 疼痛持续多长时间？\n4. 是否服用了硝酸甘油？如有不适请及时就医。",
		"头痛": "建议回复：1. 头痛的部位在哪里？\n2. 是持续性还是间歇性疼痛？\n3. 是否测量了血压？\n4. 有其他不适症状吗？",
	}

	for keyword, suggestion := range suggestions {
		if strings.Contains(content, keyword) {
			return suggestion
		}
	}

	return "建议回复：1. 请详细描述您的具体症状\n2. 症状持续多长时间了？\n3. 是否进行过相关检查？\n4. 目前是否在服用任何药物？"
}

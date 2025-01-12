package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	// 写入超时时间
	writeWait = 10 * time.Second

	// 读取超时时间
	pongWait = 60 * time.Second

	// 发送 ping 的时间间隔，必须小于 pongWait
	pingPeriod = (pongWait * 9) / 10

	// 最大消息大小
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源，生产环境应该配置具体的域名
	},
}

// WSMessage WebSocket消息结构
type WSMessage struct {
	Type    string      `json:"type"`    // 消息类型：chat/notification/physiological
	Action  string      `json:"action"`   // 动作：created/updated/deleted
	Payload interface{} `json:"payload"`  // 消息内容
}

// HandleWebSocket 处理WebSocket连接
func (m *Manager) HandleWebSocket(c *gin.Context) {
	userID := c.Query("userId")
	role := c.Query("role")
	if userID == "" || role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing userId or role"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := &Client{
		ID:      userID + "_" + time.Now().String(),
		Role:    role,
		UserID:  userID,
		Conn:    conn,
		Send:    make(chan []byte, 256),
		Manager: m,
	}

	m.register <- client

	// 启动goroutine处理读写
	go client.writePump()
	go client.readPump()
}

// readPump 处理WebSocket读取
func (c *Client) readPump() {
	defer func() {
		c.Manager.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// 处理接收到的消息
		var wsMessage WSMessage
		if err := json.Unmarshal(message, &wsMessage); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		// 根据消息类型处理
		switch wsMessage.Type {
		case "chat":
			// 处理聊天消息
			handleChatMessage(c, &wsMessage)
		case "physiological":
			// 处理生理数据消息
			handlePhysiologicalMessage(c, &wsMessage)
		}
	}
}

// writePump 处理WebSocket写入
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// 通道已关闭
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 添加队列中的其他消息
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleChatMessage 处理聊天消息
func handleChatMessage(c *Client, msg *WSMessage) {
	// 根据消息动作处理
	switch msg.Action {
	case "created":
		// 新消息创建时，通知相关用户
		if data, err := json.Marshal(msg); err == nil {
			if c.Role == "doctor" {
				// 如果是医生发送的消息，通知相关患者
				if payload, ok := msg.Payload.(map[string]interface{}); ok {
					if patientID, ok := payload["patientId"].(string); ok {
						c.Manager.SendToUser(patientID, data)
					}
				}
			} else {
				// 如果是患者发送的消息，通知相关医生
				if payload, ok := msg.Payload.(map[string]interface{}); ok {
					if doctorID, ok := payload["doctorId"].(string); ok {
						c.Manager.SendToUser(doctorID, data)
					}
				}
			}
		}
	}
}

// handlePhysiologicalMessage 处理生理数据消息
func handlePhysiologicalMessage(c *Client, msg *WSMessage) {
	// 根据消息动作处理
	switch msg.Action {
	case "created":
		// 新生理数据创建时，通知相关医生
		if data, err := json.Marshal(msg); err == nil {
			if payload, ok := msg.Payload.(map[string]interface{}); ok {
				if doctorID, ok := payload["doctorId"].(string); ok {
					c.Manager.SendToUser(doctorID, data)
				}
			}
		}
	}
} 
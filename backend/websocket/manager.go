package websocket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Client 表示一个WebSocket客户端连接
type Client struct {
	ID       string
	Role     string          // doctor/patient
	UserID   string          // 用户ID（医生ID或患者ID）
	Conn     *websocket.Conn // WebSocket连接
	Send     chan []byte     // 发送消息的通道
	Manager  *Manager        // 所属的Manager
}

// Manager 管理所有WebSocket连接
type Manager struct {
	clients    map[*Client]bool      // 所有连接的客户端
	broadcast  chan []byte           // 广播消息通道
	register   chan *Client          // 注册客户端通道
	unregister chan *Client          // 注销客户端通道
	mu         sync.RWMutex         // 读写锁
	userConns  map[string][]*Client // 用户ID到客户端连接的映射
}

// NewManager 创建一个新的Manager
func NewManager() *Manager {
	return &Manager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		userConns:  make(map[string][]*Client),
	}
}

// Start 启动WebSocket管理器
func (m *Manager) Start() {
	for {
		select {
		case client := <-m.register:
			m.mu.Lock()
			m.clients[client] = true
			m.userConns[client.UserID] = append(m.userConns[client.UserID], client)
			m.mu.Unlock()
			log.Printf("Client registered: %s (%s)", client.UserID, client.Role)

		case client := <-m.unregister:
			m.mu.Lock()
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				close(client.Send)
				// 从userConns中移除
				conns := m.userConns[client.UserID]
				for i, c := range conns {
					if c == client {
						m.userConns[client.UserID] = append(conns[:i], conns[i+1:]...)
						break
					}
				}
				if len(m.userConns[client.UserID]) == 0 {
					delete(m.userConns, client.UserID)
				}
			}
			m.mu.Unlock()
			log.Printf("Client unregistered: %s (%s)", client.UserID, client.Role)

		case message := <-m.broadcast:
			m.mu.RLock()
			for client := range m.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(m.clients, client)
				}
			}
			m.mu.RUnlock()
		}
	}
}

// SendToUser 向指定用户发送消息
func (m *Manager) SendToUser(userID string, message []byte) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if clients, ok := m.userConns[userID]; ok {
		for _, client := range clients {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(m.clients, client)
			}
		}
	}
}

// SendToRole 向指定角色的所有用户发送消息
func (m *Manager) SendToRole(role string, message []byte) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for client := range m.clients {
		if client.Role == role {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(m.clients, client)
			}
		}
	}
}

// BroadcastMessage 广播消息给所有客户端
func (m *Manager) BroadcastMessage(message []byte) {
	m.broadcast <- message
} 
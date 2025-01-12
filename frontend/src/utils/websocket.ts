import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

export interface WSMessage {
  type: 'chat' | 'ai_suggestion' | 'physiological'
  action: 'created' | 'updated' | 'deleted'
  payload: any
}

export class WebSocketService {
  private static instance: WebSocketService
  private ws: WebSocket | null = null
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private reconnectTimeout = 3000
  private messageHandlers: Map<string, Set<(payload: any) => void>> = new Map()

  private constructor() {
    // 私有构造函数，确保单例
  }

  public static getInstance(): WebSocketService {
    if (!WebSocketService.instance) {
      WebSocketService.instance = new WebSocketService()
    }
    return WebSocketService.instance
  }

  public connect() {
    const userStore = useUserStore()
    if (!userStore.isLoggedIn) {
      console.warn('[WebSocket] User not logged in, cannot connect to WebSocket')
      return
    }

    const wsUrl = `ws://${window.location.hostname}:8080/ws?userId=${userStore.userId}&role=${userStore.role}`
    console.log('[WebSocket] Attempting to connect to:', wsUrl)
    this.ws = new WebSocket(wsUrl)

    this.ws.onopen = () => {
      console.log('[WebSocket] Connection established successfully')
      this.reconnectAttempts = 0
    }

    this.ws.onmessage = (event) => {
      try {
        const messages = event.data.split('\n')
        console.log(`[WebSocket] Received ${messages.length} message(s)`)
        messages.forEach(message => {
          if (!message) return
          const wsMessage: WSMessage = JSON.parse(message)
          console.log('[WebSocket] Processing message:', {
            type: wsMessage.type,
            action: wsMessage.action,
            payload: wsMessage.payload
          })
          this.handleMessage(wsMessage)
        })
      } catch (error) {
        console.error('[WebSocket] Failed to parse message:', error, '\nRaw message:', event.data)
      }
    }

    this.ws.onclose = (event) => {
      console.log('[WebSocket] Connection closed', {
        code: event.code,
        reason: event.reason,
        wasClean: event.wasClean
      })
      this.attemptReconnect()
    }

    this.ws.onerror = (error) => {
      console.error('[WebSocket] Connection error:', error)
      ElMessage.error('WebSocket连接错误')
    }
  }

  private attemptReconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error('[WebSocket] Max reconnection attempts reached')
      ElMessage.error('WebSocket连接失败，请刷新页面重试')
      return
    }

    this.reconnectAttempts++
    console.log(`[WebSocket] Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts}) in ${this.reconnectTimeout}ms`)

    setTimeout(() => {
      this.connect()
    }, this.reconnectTimeout)
  }

  public disconnect() {
    if (this.ws) {
      console.log('[WebSocket] Manually disconnecting')
      this.ws.close()
      this.ws = null
    }
  }

  public addMessageHandler(type: string, handler: (payload: any) => void) {
    if (!this.messageHandlers.has(type)) {
      this.messageHandlers.set(type, new Set())
    }
    this.messageHandlers.get(type)?.add(handler)
    console.log(`[WebSocket] Added message handler for type: ${type}`)
  }

  public removeMessageHandler(type: string, handler: (payload: any) => void) {
    const removed = this.messageHandlers.get(type)?.delete(handler)
    console.log(`[WebSocket] Removed message handler for type: ${type}`, { success: removed })
  }

  private handleMessage(message: WSMessage) {
    const handlers = this.messageHandlers.get(message.type)
    console.log(`[WebSocket] Handling message of type: ${message.type}`, {
      handlersCount: handlers?.size || 0
    })
    
    if (handlers) {
      handlers.forEach(handler => {
        try {
          handler(message.payload)
        } catch (error) {
          console.error('[WebSocket] Error in message handler:', error)
        }
      })
    }
  }

  public isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN
  }
} 
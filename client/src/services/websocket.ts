import { useAuthStore } from '@/stores/authStore'
import { useGameStore } from '@/stores/gameStore'
import type { Message, Entity, Region } from '@/types'

type MessageHandler = (message: Message) => void

class WebSocketService {
  private ws: WebSocket | null = null
  private handlers: Map<string, MessageHandler[]> = new Map()
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private reconnectDelay = 1000

  connect(token: string) {
    const wsUrl = `ws://localhost:8080/ws?token=${token}`
    this.ws = new WebSocket(wsUrl)

    this.ws.onopen = () => {
      console.log('WebSocket connected')
      this.reconnectAttempts = 0
      useGameStore.getState().setConnected(true)
    }

    this.ws.onclose = () => {
      console.log('WebSocket disconnected')
      useGameStore.getState().setConnected(false)
      this.attemptReconnect(token)
    }

    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error)
    }

    this.ws.onmessage = (event) => {
      try {
        const message: Message = JSON.parse(event.data)
        this.handleMessage(message)
      } catch (e) {
        console.error('Failed to parse message:', e)
      }
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
  }

  send(type: string, payload: Record<string, unknown>) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      const message: Message = {
        type,
        payload,
        timestamp: Date.now(),
      }
      this.ws.send(JSON.stringify(message))
    }
  }

  on(type: string, handler: MessageHandler) {
    if (!this.handlers.has(type)) {
      this.handlers.set(type, [])
    }
    this.handlers.get(type)!.push(handler)
  }

  off(type: string, handler: MessageHandler) {
    const handlers = this.handlers.get(type)
    if (handlers) {
      const index = handlers.indexOf(handler)
      if (index > -1) {
        handlers.splice(index, 1)
      }
    }
  }

  private handleMessage(message: Message) {
    const handlers = this.handlers.get(message.type)
    if (handlers) {
      handlers.forEach((handler) => handler(message))
    }

    switch (message.type) {
      case 'state_sync':
        this.handleStateSync(message)
        break
      case 'entity_update':
        this.handleEntityUpdate(message)
        break
      case 'chat':
        this.handleChat(message)
        break
      case 'world_event':
        this.handleWorldEvent(message)
        break
    }
  }

  private handleStateSync(message: Message) {
    const { entity, region, world_time } = message.payload as {
      entity: Entity
      region: Region
      world_time: number
    }

    useAuthStore.getState().updateEntity(entity)
    useGameStore.getState().setCurrentRegion(region)
    useGameStore.getState().setWorldTime(world_time)
  }

  private handleEntityUpdate(message: Message) {
    const { changes } = message.payload as { changes: Partial<Entity> }
    useAuthStore.getState().updateEntity(changes)
  }

  private handleChat(message: Message) {
    const { sender_id, sender_name, channel, content } = message.payload as {
      sender_id: string
      sender_name: string
      channel: string
      content: string
    }

    useGameStore.getState().addChatMessage({
      id: `${Date.now()}-${Math.random()}`,
      sender_id,
      sender_name,
      channel,
      content,
      timestamp: message.timestamp,
    })
  }

  private handleWorldEvent(message: Message) {
    console.log('World event:', message.payload)
  }

  private attemptReconnect(token: string) {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++
      setTimeout(() => {
        console.log(`Reconnecting... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`)
        this.connect(token)
      }, this.reconnectDelay * this.reconnectAttempts)
    }
  }
}

export const wsService = new WebSocketService()

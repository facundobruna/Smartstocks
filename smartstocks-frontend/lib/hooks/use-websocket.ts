'use client'

import { useRef, useCallback, useState } from 'react'
import { useAuthStore } from '@/lib/stores/auth-store'

const WS_URL = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8081'

export type WebSocketMessageType =
  | 'match_found'
  | 'round_start'
  | 'round_result'
  | 'match_result'
  | 'error'
  | 'opponent_left'
  | 'ping'
  | 'pong'

export interface WebSocketMessage {
  type: WebSocketMessageType
  data: any
  timestamp?: string
}

interface UseWebSocketOptions {
  onMessage?: (message: WebSocketMessage) => void
  onConnect?: () => void
  onDisconnect?: () => void
  onError?: (error: Event) => void
}

export function useWebSocket(options: UseWebSocketOptions = {}) {
  const { accessToken } = useAuthStore()
  const socketRef = useRef<WebSocket | null>(null)
  const [isConnected, setIsConnected] = useState(false)
  const [isConnecting, setIsConnecting] = useState(false)

  const connect = useCallback((): Promise<void> => {
    return new Promise((resolve, reject) => {
      if (socketRef.current?.readyState === WebSocket.OPEN) {
        resolve()
        return
      }

      if (isConnecting) {
        reject(new Error('Already connecting'))
        return
      }

      if (!accessToken) {
        reject(new Error('No access token available for WebSocket connection'))
        return
      }

      setIsConnecting(true)

      // Usar WebSocket nativo con token en query param
      const wsUrl = `${WS_URL}/api/v1/pvp/ws?token=${accessToken}`
      const socket = new WebSocket(wsUrl)

      socket.onopen = () => {
        setIsConnected(true)
        setIsConnecting(false)
        options.onConnect?.()
        resolve()
      }

      socket.onclose = () => {
        setIsConnected(false)
        setIsConnecting(false)
        options.onDisconnect?.()
      }

      socket.onerror = (error) => {
        setIsConnecting(false)
        setIsConnected(false)
        options.onError?.(error)
        reject(error)
      }

      socket.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data) as WebSocketMessage
          console.log('📨 WS Message received:', message)
          options.onMessage?.(message)
        } catch (e) {
          console.error('Failed to parse WebSocket message:', e)
        }
      }

      socketRef.current = socket
    })
  }, [accessToken, isConnecting, options])

  const disconnect = useCallback(() => {
    if (socketRef.current) {
      socketRef.current.close()
      socketRef.current = null
      setIsConnected(false)
    }
  }, [])

  const send = useCallback((type: string, data?: any) => {
    if (socketRef.current?.readyState === WebSocket.OPEN) {
      const message = { type, data, timestamp: new Date().toISOString() }
      console.log('📤 WS Message sending:', message)
      socketRef.current.send(JSON.stringify(message))
    }
  }, [])

  const submitDecision = useCallback((decision: string) => {
    send('submit_decision', { decision })
  }, [send])

  return {
    socket: socketRef.current,
    isConnected,
    isConnecting,
    connect,
    disconnect,
    send,
    submitDecision,
  }
}

export default useWebSocket

'use client'

import { useEffect, useCallback } from 'react'
import { Swords, Play, History, Trophy, Target, Loader2 } from 'lucide-react'
import { toast } from 'sonner'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { QueueCard } from '@/components/pvp/queue-card'
import { MatchArena } from '@/components/pvp/match-arena'
import { RoundResultCard } from '@/components/pvp/round-result'
import { MatchResultCard } from '@/components/pvp/match-result'
import { useWebSocket, type WebSocketMessage } from '@/lib/hooks/use-websocket'
import { usePvPStore } from '@/lib/stores/pvp-store'
import { useAuthStore } from '@/lib/stores/auth-store'
import { pvpApi } from '@/lib/api/pvp'
import type { SimulatorDecision } from '@/types'

export default function PvPPage() {
  const { updateStats } = useAuthStore()
  const {
    status,
    queuePosition,
    matchId,
    currentRound,
    roundResult,
    matchResult,
    setQueuing,
    setMatchFound,
    setRoundStart,
    setSelectedDecision,
    setRoundResult,
    setMatchEnd,
    reset,
  } = usePvPStore()

  const handleMessage = useCallback((message: WebSocketMessage) => {
    console.log('📨 Handling message:', message.type, message.data)

    switch (message.type) {
      case 'match_found':
        setMatchFound({
          match_id: message.data.match_id,
          opponent: message.data.opponent,
          total_rounds: message.data.total_rounds,
        })
        toast.success('Oponente encontrado!')
        break
      case 'round_start':
        setRoundStart({
          round_number: message.data.round_number,
          scenario: message.data.scenario,
          time_limit_seconds: message.data.time_limit_seconds,
        })
        break
      case 'round_result':
        setRoundResult({
          round_number: message.data.round_number,
          your_decision: message.data.your_decision,
          opponent_decision: message.data.opponent_decision,
          correct_decision: message.data.correct_decision,
          your_points: message.data.your_points,
          opponent_points: message.data.opponent_points,
          explanation: message.data.explanation,
        })
        break
      case 'match_result':
        setMatchEnd({
          winner_id: message.data.winner, // "you", "opponent", "tie"
          your_final_score: message.data.your_final_score,
          opponent_final_score: message.data.opponent_final_score,
          points_earned: message.data.points_gained,
          new_total_points: message.data.new_total_points,
          new_rank_tier: message.data.new_rank_tier,
        })
        updateStats({
          smartpoints: message.data.new_total_points,
          rank_tier: message.data.new_rank_tier,
        })
        break
      case 'opponent_left':
        toast.error('Tu oponente se desconecto')
        reset()
        break
      case 'error':
        toast.error(message.data?.error || message.data?.message || 'Error en el juego')
        break
    }
  }, [
    setMatchFound,
    setRoundStart,
    setRoundResult,
    setMatchEnd,
    updateStats,
    reset,
  ])

  const {
    isConnected,
    isConnecting,
    connect,
    disconnect,
  } = useWebSocket({
    onMessage: handleMessage,
    onConnect: () => toast.success('Conectado al servidor'),
    onDisconnect: () => {
      if (status !== 'idle') {
        toast.error('Desconectado del servidor')
        reset()
      }
    },
    onError: (error) => {
      toast.error('Error de conexion')
      console.error('WebSocket error:', error)
    },
  })

  const handleStartQueue = async () => {
    try {
      // Primero conectar al WebSocket y esperar la conexión
      await connect()
      // Una vez conectado, unirse a la cola
      const response = await pvpApi.joinQueue()
      setQueuing(response.position ?? 1)
    } catch (error: any) {
      disconnect()
      toast.error(error.response?.data?.message || error.message || 'Error al unirse a la cola')
    }
  }

  const handleCancelQueue = async () => {
    try {
      await pvpApi.leaveQueue()
      disconnect()
      reset()
    } catch (error) {
      console.error('Error leaving queue:', error)
    }
  }

  const handleDecision = async (decision: SimulatorDecision, timeElapsed: number) => {
    if (!matchId) return
    setSelectedDecision(decision)
    try {
      await pvpApi.submitDecision(matchId, currentRound, decision, timeElapsed)
    } catch (error: any) {
      console.error('Error submitting decision:', error)
      // El resultado llegará via WebSocket
    }
  }

  const handleContinue = () => {
    // El servidor enviará round_start automáticamente
  }

  const handlePlayAgain = () => {
    reset()
  }

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      disconnect()
      reset()
    }
  }, [disconnect, reset])

  // Render based on status
  if (status === 'queuing') {
    return (
      <div className="space-y-6">
        <QueueCard
          queuePosition={queuePosition}
          onCancel={handleCancelQueue}
          isConnecting={isConnecting}
        />
      </div>
    )
  }

  if (status === 'matched' || status === 'playing' || status === 'waiting_opponent') {
    return (
      <div className="max-w-4xl mx-auto">
        <MatchArena onDecision={handleDecision} />
      </div>
    )
  }

  if (status === 'round_result' && roundResult) {
    return (
      <div className="max-w-4xl mx-auto">
        <RoundResultCard result={roundResult} onContinue={handleContinue} />
      </div>
    )
  }

  if (status === 'match_end' && matchResult) {
    return (
      <div className="max-w-4xl mx-auto">
        <MatchResultCard result={matchResult} onPlayAgain={handlePlayAgain} />
      </div>
    )
  }

  // Idle state - show lobby
  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="text-center">
        <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-purple-100 mb-4">
          <Swords className="w-8 h-8 text-purple-600" />
        </div>
        <h1 className="text-3xl font-bold text-foreground">Modo PvP</h1>
        <p className="text-muted-foreground mt-2 max-w-xl mx-auto">
          Compite en tiempo real contra otros jugadores. Analiza los mismos escenarios
          y demuestra quien toma mejores decisiones.
        </p>
      </div>

      {/* Play Button */}
      <div className="max-w-md mx-auto">
        <Card className="border-2 border-purple-200 bg-purple-50">
          <CardContent className="pt-6 text-center">
            <Button
              size="lg"
              className="w-full bg-purple-600 hover:bg-purple-700 text-lg py-6"
              onClick={handleStartQueue}
              disabled={isConnecting}
            >
              {isConnecting ? (
                <>
                  <Loader2 className="mr-2 h-5 w-5 animate-spin" />
                  Conectando...
                </>
              ) : (
                <>
                  <Play className="mr-2 h-5 w-5" />
                  Buscar Partida
                </>
              )}
            </Button>
            <p className="text-sm text-muted-foreground mt-3">
              5 rondas • 15 segundos por ronda
            </p>
          </CardContent>
        </Card>
      </div>

      {/* How to Play */}
      <div className="max-w-2xl mx-auto">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Como Jugar</CardTitle>
          </CardHeader>
          <CardContent>
            <ol className="list-decimal list-inside space-y-2 text-sm text-muted-foreground">
              <li>Presiona "Buscar Partida" para entrar en la cola</li>
              <li>Seras emparejado con otro jugador de nivel similar</li>
              <li>Ambos veran el mismo escenario y tendran 15 segundos para decidir</li>
              <li>Gana puntos si tu decision es correcta</li>
              <li>El jugador con mas puntos al final de 5 rondas gana</li>
            </ol>
          </CardContent>
        </Card>
      </div>

      {/* Stats Preview */}
      <div className="max-w-2xl mx-auto grid grid-cols-3 gap-4">
        <Card>
          <CardContent className="pt-6 text-center">
            <Trophy className="w-8 h-8 text-yellow-500 mx-auto mb-2" />
            <p className="text-2xl font-bold">0</p>
            <p className="text-sm text-muted-foreground">Victorias</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6 text-center">
            <Target className="w-8 h-8 text-blue-500 mx-auto mb-2" />
            <p className="text-2xl font-bold">0</p>
            <p className="text-sm text-muted-foreground">Partidas</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="pt-6 text-center">
            <History className="w-8 h-8 text-green-500 mx-auto mb-2" />
            <p className="text-2xl font-bold">0%</p>
            <p className="text-sm text-muted-foreground">Win Rate</p>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

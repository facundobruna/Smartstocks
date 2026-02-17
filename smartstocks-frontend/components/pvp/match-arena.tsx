'use client'

import { useEffect, useState } from 'react'
import { Clock } from 'lucide-react'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { PlayerCard } from './player-card'
import { ScenarioCard } from '@/components/simulator/scenario-card'
import { StockChart } from '@/components/simulator/stock-chart'
import { DecisionButtons } from '@/components/simulator/decision-buttons'
import { useAuthStore } from '@/lib/stores/auth-store'
import { usePvPStore } from '@/lib/stores/pvp-store'
import { cn } from '@/lib/utils'
import type { SimulatorDecision } from '@/types'

interface MatchArenaProps {
  onDecision: (decision: SimulatorDecision, timeElapsed: number) => void
}

export function MatchArena({ onDecision }: MatchArenaProps) {
  const { user, stats } = useAuthStore()
  const {
    opponent,
    currentRound,
    totalRounds,
    yourScore,
    opponentScore,
    scenario,
    timeLimit,
    opponentDecided,
    selectedDecision,
    status,
  } = usePvPStore()

  const [timeRemaining, setTimeRemaining] = useState(timeLimit)

  useEffect(() => {
    setTimeRemaining(timeLimit)
  }, [timeLimit, currentRound])

  useEffect(() => {
    if (status !== 'playing' || timeRemaining <= 0) return

    const interval = setInterval(() => {
      setTimeRemaining((prev) => {
        if (prev <= 1) {
          clearInterval(interval)
          // Auto-submit hold if time runs out
          if (!selectedDecision) {
            onDecision('hold', timeLimit)
          }
          return 0
        }
        return prev - 1
      })
    }, 1000)

    return () => clearInterval(interval)
  }, [status, timeRemaining, selectedDecision, onDecision, timeLimit])

  const handleDecision = (decision: SimulatorDecision) => {
    if (selectedDecision) return
    const timeElapsed = timeLimit - timeRemaining
    onDecision(decision, timeElapsed)
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <Badge variant="outline" className="text-lg px-4 py-2">
          Ronda {currentRound}/{totalRounds}
        </Badge>

        <Card className={cn(
          'border-0 shadow-none',
          timeRemaining < 5 ? 'bg-red-100' : 'bg-gray-100'
        )}>
          <CardContent className="py-2 px-4 flex items-center gap-2">
            <Clock className={cn(
              'w-5 h-5',
              timeRemaining < 5 ? 'text-red-500' : 'text-muted-foreground'
            )} />
            <span className={cn(
              'font-mono text-2xl font-bold',
              timeRemaining < 5 && 'text-red-500'
            )}>
              {timeRemaining}s
            </span>
          </CardContent>
        </Card>
      </div>

      {/* Players */}
      <div className="grid grid-cols-2 gap-6">
        <PlayerCard
          user={user}
          score={yourScore}
          isCurrentUser
          rankTier={stats?.rank_tier}
          hasDecided={!!selectedDecision}
        />
        <PlayerCard
          user={opponent}
          score={opponentScore}
          rankTier="bronze_1" // TODO: get from opponent data
          hasDecided={opponentDecided}
        />
      </div>

      {/* Score Progress */}
      <div className="space-y-2">
        <div className="flex justify-between text-sm">
          <span className="text-primary font-medium">Tu: {yourScore}</span>
          <span className="text-muted-foreground font-medium">Oponente: {opponentScore}</span>
        </div>
        <div className="flex gap-1 h-3">
          <Progress
            value={yourScore > 0 ? (yourScore / (yourScore + opponentScore)) * 100 : 50}
            className="flex-1"
          />
        </div>
      </div>

      {/* Scenario */}
      {scenario && (
        <>
          <ScenarioCard newsContent={scenario.news_content} />
          <StockChart data={scenario.chart_data} />
        </>
      )}

      {/* Decision Buttons */}
      <DecisionButtons
        onDecision={handleDecision}
        disabled={!!selectedDecision || status !== 'playing'}
        loading={status === 'waiting_opponent'}
        selectedDecision={selectedDecision}
      />

      {status === 'waiting_opponent' && (
        <p className="text-center text-muted-foreground animate-pulse">
          Esperando decision del oponente...
        </p>
      )}
    </div>
  )
}

export default MatchArena

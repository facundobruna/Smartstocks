'use client'

import { CheckCircle, XCircle, ArrowRight } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { cn } from '@/lib/utils'
import type { RoundResult } from '@/lib/stores/pvp-store'
import type { SimulatorDecision } from '@/types'

interface RoundResultProps {
  result: RoundResult
  onContinue: () => void
}

const decisionLabels: Record<SimulatorDecision, string> = {
  buy: 'Comprar',
  hold: 'Mantener',
  sell: 'Vender',
}

export function RoundResultCard({ result, onContinue }: RoundResultProps) {
  const youWon = result.your_points > result.opponent_points
  const draw = result.your_points === result.opponent_points

  return (
    <Card className="max-w-lg mx-auto">
      <CardHeader className="text-center">
        <CardTitle className="flex items-center justify-center gap-2 text-2xl">
          {youWon ? (
            <>
              <CheckCircle className="w-8 h-8 text-bullish" />
              <span className="text-bullish">Ganaste la Ronda!</span>
            </>
          ) : draw ? (
            <>
              <span className="text-yellow-500">Empate!</span>
            </>
          ) : (
            <>
              <XCircle className="w-8 h-8 text-bearish" />
              <span className="text-bearish">Perdiste la Ronda</span>
            </>
          )}
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* Decisions */}
        <div className="grid grid-cols-3 gap-4 text-center">
          <div className="p-3 bg-gray-50 rounded-lg">
            <p className="text-xs text-muted-foreground mb-1">Tu decision</p>
            <Badge variant={result.your_decision === result.correct_decision ? 'default' : 'secondary'}>
              {decisionLabels[result.your_decision]}
            </Badge>
          </div>
          <div className="p-3 bg-green-50 rounded-lg">
            <p className="text-xs text-muted-foreground mb-1">Correcta</p>
            <Badge className="bg-bullish">
              {decisionLabels[result.correct_decision]}
            </Badge>
          </div>
          <div className="p-3 bg-gray-50 rounded-lg">
            <p className="text-xs text-muted-foreground mb-1">Oponente</p>
            <Badge variant={result.opponent_decision === result.correct_decision ? 'default' : 'secondary'}>
              {decisionLabels[result.opponent_decision]}
            </Badge>
          </div>
        </div>

        {/* Points */}
        <div className="grid grid-cols-2 gap-4 text-center">
          <div className={cn(
            'p-4 rounded-lg',
            youWon ? 'bg-green-100' : draw ? 'bg-yellow-50' : 'bg-gray-100'
          )}>
            <p className="text-2xl font-bold">+{result.your_points}</p>
            <p className="text-sm text-muted-foreground">Tu puntos</p>
          </div>
          <div className={cn(
            'p-4 rounded-lg',
            !youWon && !draw ? 'bg-red-100' : 'bg-gray-100'
          )}>
            <p className="text-2xl font-bold">+{result.opponent_points}</p>
            <p className="text-sm text-muted-foreground">Puntos oponente</p>
          </div>
        </div>

        {/* Explanation */}
        <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
          <p className="text-sm text-blue-700">{result.explanation}</p>
        </div>

        <Button className="w-full" onClick={onContinue}>
          Siguiente Ronda
          <ArrowRight className="ml-2 w-4 h-4" />
        </Button>
      </CardContent>
    </Card>
  )
}

export default RoundResultCard

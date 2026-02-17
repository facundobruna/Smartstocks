'use client'

import { useRouter } from 'next/navigation'
import { Trophy, Medal, ArrowRight, Home } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { RankBadge } from '@/components/rankings/rank-badge'
import { useAuthStore } from '@/lib/stores/auth-store'
import { usePvPStore, type MatchResult } from '@/lib/stores/pvp-store'
import { cn } from '@/lib/utils'

interface MatchResultProps {
  result: MatchResult
  onPlayAgain: () => void
}

export function MatchResultCard({ result, onPlayAgain }: MatchResultProps) {
  const router = useRouter()
  const { user, updateStats } = useAuthStore()
  const { opponent } = usePvPStore()

  const isWinner = result.winner_id === user?.id
  const isDraw = result.your_final_score === result.opponent_final_score

  return (
    <Card className="max-w-lg mx-auto">
      <CardHeader className={cn(
        'text-center rounded-t-lg',
        isWinner ? 'bg-gradient-to-r from-yellow-400 to-yellow-600' :
        isDraw ? 'bg-gradient-to-r from-gray-400 to-gray-600' :
        'bg-gradient-to-r from-gray-600 to-gray-800'
      )}>
        <CardTitle className="text-white">
          <div className="flex items-center justify-center gap-3 text-3xl">
            {isWinner ? (
              <>
                <Trophy className="w-10 h-10" />
                Victoria!
              </>
            ) : isDraw ? (
              <>
                <Medal className="w-10 h-10" />
                Empate!
              </>
            ) : (
              <>
                <Medal className="w-10 h-10" />
                Derrota
              </>
            )}
          </div>
        </CardTitle>
      </CardHeader>
      <CardContent className="pt-6 space-y-6">
        {/* Final Score */}
        <div className="grid grid-cols-3 gap-4 items-center text-center">
          <div>
            <p className="text-4xl font-bold text-primary">{result.your_final_score}</p>
            <p className="text-sm text-muted-foreground">Tu puntuacion</p>
          </div>
          <div className="text-2xl font-bold text-muted-foreground">VS</div>
          <div>
            <p className="text-4xl font-bold">{result.opponent_final_score}</p>
            <p className="text-sm text-muted-foreground">{opponent?.username}</p>
          </div>
        </div>

        {/* Rewards */}
        <div className="bg-gray-50 rounded-lg p-4 space-y-3">
          <div className="flex justify-between items-center">
            <span className="text-muted-foreground">Puntos ganados</span>
            <span className={cn(
              'font-bold text-lg',
              result.points_earned > 0 ? 'text-bullish' : 'text-bearish'
            )}>
              {result.points_earned > 0 ? '+' : ''}{result.points_earned}
            </span>
          </div>
          <div className="flex justify-between items-center">
            <span className="text-muted-foreground">Total de puntos</span>
            <span className="font-bold text-lg">{result.new_total_points.toLocaleString()}</span>
          </div>
          <div className="flex justify-between items-center">
            <span className="text-muted-foreground">Rango</span>
            <RankBadge tier={result.new_rank_tier} size="sm" />
          </div>
        </div>

        {/* Actions */}
        <div className="grid grid-cols-2 gap-3">
          <Button variant="outline" onClick={() => router.push('/home')}>
            <Home className="w-4 h-4 mr-2" />
            Inicio
          </Button>
          <Button onClick={onPlayAgain}>
            Jugar de Nuevo
            <ArrowRight className="ml-2 w-4 h-4" />
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}

export default MatchResultCard

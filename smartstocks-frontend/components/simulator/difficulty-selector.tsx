'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { Clock, Zap, Trophy, Lock, Loader2 } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { cn } from '@/lib/utils'
import { simulatorApi } from '@/lib/api/simulator'
import type { SimulatorDifficulty, CooldownStatus } from '@/types'

interface DifficultyOption {
  id: SimulatorDifficulty
  name: string
  description: string
  points: number
  color: string
  bgColor: string
}

const difficulties: DifficultyOption[] = [
  {
    id: 'easy',
    name: 'Facil',
    description: 'Escenarios simples para principiantes',
    points: 25,
    color: 'text-green-600',
    bgColor: 'bg-green-50 border-green-200 hover:border-green-400',
  },
  {
    id: 'medium',
    name: 'Medio',
    description: 'Escenarios con mayor complejidad',
    points: 50,
    color: 'text-yellow-600',
    bgColor: 'bg-yellow-50 border-yellow-200 hover:border-yellow-400',
  },
  {
    id: 'hard',
    name: 'Dificil',
    description: 'Escenarios desafiantes para expertos',
    points: 100,
    color: 'text-red-600',
    bgColor: 'bg-red-50 border-red-200 hover:border-red-400',
  },
]

export function DifficultySelector() {
  const router = useRouter()
  const [cooldowns, setCooldowns] = useState<Record<SimulatorDifficulty, CooldownStatus | null>>({
    easy: null,
    medium: null,
    hard: null,
  })
  const [loading, setLoading] = useState(true)
  const [starting, setStarting] = useState<SimulatorDifficulty | null>(null)

  useEffect(() => {
    const fetchCooldowns = async () => {
      try {
        const [easy, medium, hard] = await Promise.all([
          simulatorApi.getCooldownStatus('easy'),
          simulatorApi.getCooldownStatus('medium'),
          simulatorApi.getCooldownStatus('hard'),
        ])
        setCooldowns({ easy, medium, hard })
      } catch (error) {
        console.error('Error fetching cooldowns:', error)
      } finally {
        setLoading(false)
      }
    }
    fetchCooldowns()
  }, [])

  const handlePlay = async (difficulty: SimulatorDifficulty) => {
    setStarting(difficulty)
    router.push(`/simulator/${difficulty}`)
  }

  const formatTimeRemaining = (hours: number): string => {
    if (hours < 1) {
      return `${Math.ceil(hours * 60)}m`
    }
    return `${Math.floor(hours)}h ${Math.ceil((hours % 1) * 60)}m`
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
      {difficulties.map((diff) => {
        const cooldown = cooldowns[diff.id]
        const canPlay = cooldown?.can_attempt ?? true
        const isStarting = starting === diff.id

        return (
          <Card
            key={diff.id}
            className={cn(
              'relative transition-all duration-200 border-2',
              diff.bgColor,
              !canPlay && 'opacity-70'
            )}
          >
            <CardHeader className="text-center pb-2">
              <div className="flex justify-center mb-2">
                {diff.id === 'easy' && <Zap className={cn('w-12 h-12', diff.color)} />}
                {diff.id === 'medium' && <Trophy className={cn('w-12 h-12', diff.color)} />}
                {diff.id === 'hard' && <Trophy className={cn('w-12 h-12', diff.color)} />}
              </div>
              <CardTitle className={cn('text-2xl', diff.color)}>
                {diff.name}
              </CardTitle>
            </CardHeader>
            <CardContent className="text-center space-y-4">
              <p className="text-sm text-muted-foreground">{diff.description}</p>

              <Badge variant="secondary" className="text-lg px-4 py-1">
                +{diff.points} pts
              </Badge>

              {canPlay ? (
                <Button
                  className="w-full"
                  size="lg"
                  onClick={() => handlePlay(diff.id)}
                  disabled={isStarting}
                >
                  {isStarting ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      Cargando...
                    </>
                  ) : (
                    'Jugar'
                  )}
                </Button>
              ) : (
                <div className="space-y-2">
                  <Button className="w-full" size="lg" disabled>
                    <Lock className="mr-2 h-4 w-4" />
                    En Cooldown
                  </Button>
                  <div className="flex items-center justify-center gap-1 text-sm text-muted-foreground">
                    <Clock className="w-4 h-4" />
                    <span>
                      Disponible en {formatTimeRemaining(cooldown?.hours_remaining || 0)}
                    </span>
                  </div>
                </div>
              )}

              <p className="text-xs text-muted-foreground">
                {canPlay ? 'Disponible ahora' : 'Vuelve manana'}
              </p>
            </CardContent>
          </Card>
        )
      })}
    </div>
  )
}

export default DifficultySelector

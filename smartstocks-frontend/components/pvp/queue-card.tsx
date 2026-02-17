'use client'

import { Swords, Loader2, X } from 'lucide-react'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { PlayerCard } from './player-card'
import { useAuthStore } from '@/lib/stores/auth-store'

interface QueueCardProps {
  queuePosition: number | null
  onCancel: () => void
  isConnecting?: boolean
}

export function QueueCard({ queuePosition, onCancel, isConnecting }: QueueCardProps) {
  const { user, stats } = useAuthStore()

  return (
    <Card className="max-w-2xl mx-auto">
      <CardContent className="pt-8 pb-6">
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-primary/10 mb-4">
            <Swords className="w-8 h-8 text-primary animate-pulse" />
          </div>
          <h2 className="text-2xl font-bold">Buscando Oponente...</h2>
          <p className="text-muted-foreground mt-1">
            Preparate para la batalla
          </p>
        </div>

        {/* VS Display */}
        <div className="grid grid-cols-3 gap-4 items-center mb-8">
          <PlayerCard
            user={user}
            score={0}
            isCurrentUser
            rankTier={stats?.rank_tier}
          />

          <div className="text-center">
            <div className="text-4xl font-bold text-muted-foreground">VS</div>
            <Loader2 className="w-6 h-6 animate-spin mx-auto mt-2 text-primary" />
          </div>

          <PlayerCard
            user={null}
            score={0}
            isSearching
          />
        </div>

        {/* Queue Info */}
        <div className="text-center space-y-2 mb-6">
          {queuePosition && (
            <p className="text-sm text-muted-foreground">
              Posicion en cola: <span className="font-semibold">#{queuePosition}</span>
            </p>
          )}
          <p className="text-sm text-muted-foreground">
            Tiempo estimado: <span className="font-semibold">~30s</span>
          </p>
        </div>

        {/* Cancel Button */}
        <div className="text-center">
          <Button
            variant="outline"
            onClick={onCancel}
            disabled={isConnecting}
          >
            <X className="w-4 h-4 mr-2" />
            Cancelar Busqueda
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}

export default QueueCard

'use client'

import { useRouter } from 'next/navigation'
import { Calendar, Users, Coins, Trophy, Clock, Lock, CheckCircle } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { cn } from '@/lib/utils'
import type { Tournament } from '@/types'
import { format, formatDistanceToNow } from 'date-fns'
import { es } from 'date-fns/locale'

interface TournamentCardProps {
  tournament: Tournament
  isRegistered?: boolean
  onJoin?: () => void
  loading?: boolean
}

const statusLabels: Record<string, { label: string; color: string }> = {
  upcoming: { label: 'Proximo', color: 'bg-blue-500' },
  registration: { label: 'Inscripcion Abierta', color: 'bg-green-500' },
  in_progress: { label: 'En Progreso', color: 'bg-yellow-500' },
  completed: { label: 'Finalizado', color: 'bg-gray-500' },
  cancelled: { label: 'Cancelado', color: 'bg-red-500' },
}

const typeLabels: Record<string, string> = {
  weekly: 'Semanal',
  monthly: 'Mensual',
  special: 'Especial',
}

export function TournamentCard({
  tournament,
  isRegistered = false,
  onJoin,
  loading = false,
}: TournamentCardProps) {
  const router = useRouter()
  const statusInfo = statusLabels[tournament.status] || statusLabels.upcoming

  const participantPercentage = (tournament.current_participants / tournament.max_participants) * 100
  const isFull = tournament.current_participants >= tournament.max_participants
  const canJoin = tournament.status === 'registration' && !isFull && !isRegistered

  const getTimeLabel = () => {
    const now = new Date()
    const startTime = new Date(tournament.start_time)

    if (tournament.status === 'upcoming' || tournament.status === 'registration') {
      return `Inicia ${formatDistanceToNow(startTime, { addSuffix: true, locale: es })}`
    }
    if (tournament.status === 'in_progress') {
      return 'En curso'
    }
    return format(startTime, 'dd MMM yyyy', { locale: es })
  }

  return (
    <Card className="overflow-hidden hover:shadow-lg transition-shadow">
      <CardHeader className="pb-2">
        <div className="flex items-start justify-between">
          <div>
            <Badge className={cn('mb-2', statusInfo.color)}>
              {statusInfo.label}
            </Badge>
            <CardTitle className="text-xl">{tournament.name}</CardTitle>
          </div>
          <Badge variant="outline">{typeLabels[tournament.tournament_type]}</Badge>
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        <p className="text-sm text-muted-foreground line-clamp-2">
          {tournament.description}
        </p>

        {/* Info Grid */}
        <div className="grid grid-cols-2 gap-3 text-sm">
          <div className="flex items-center gap-2">
            <Calendar className="w-4 h-4 text-muted-foreground" />
            <span>{getTimeLabel()}</span>
          </div>
          <div className="flex items-center gap-2">
            <Users className="w-4 h-4 text-muted-foreground" />
            <span>{tournament.current_participants}/{tournament.max_participants}</span>
          </div>
          <div className="flex items-center gap-2">
            <Coins className="w-4 h-4 text-yellow-500" />
            <span>Entry: {tournament.entry_fee} tokens</span>
          </div>
          <div className="flex items-center gap-2">
            <Trophy className="w-4 h-4 text-yellow-500" />
            <span>Premio: {tournament.prize_pool.toLocaleString()}</span>
          </div>
        </div>

        {/* Participants Progress */}
        <div>
          <div className="flex justify-between text-xs mb-1">
            <span className="text-muted-foreground">Participantes</span>
            <span>{participantPercentage.toFixed(0)}%</span>
          </div>
          <Progress value={participantPercentage} className="h-2" />
        </div>

        {/* Requirements */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          <Lock className="w-4 h-4" />
          <span>Rango minimo: {tournament.min_rank_required}</span>
        </div>

        {/* Actions */}
        <div className="flex gap-2">
          <Button
            variant="outline"
            className="flex-1"
            onClick={() => router.push(`/tournaments/${tournament.id}`)}
          >
            Ver Detalles
          </Button>
          {canJoin && onJoin && (
            <Button className="flex-1" onClick={onJoin} disabled={loading}>
              Inscribirse
            </Button>
          )}
          {isRegistered && (
            <Button className="flex-1" variant="secondary" disabled>
              <CheckCircle className="w-4 h-4 mr-2" />
              Inscrito
            </Button>
          )}
          {isFull && !isRegistered && (
            <Button className="flex-1" variant="secondary" disabled>
              Torneo Lleno
            </Button>
          )}
        </div>
      </CardContent>
    </Card>
  )
}

export default TournamentCard

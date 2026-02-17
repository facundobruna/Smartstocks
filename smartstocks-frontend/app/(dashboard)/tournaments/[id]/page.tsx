'use client'

import { useEffect, useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import {
  ArrowLeft,
  Trophy,
  Users,
  Coins,
  Calendar,
  Clock,
  Lock,
  CheckCircle,
  Loader2,
} from 'lucide-react'
import { toast } from 'sonner'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { RankBadge } from '@/components/rankings/rank-badge'
import { tournamentsApi, type TournamentDetailsResponse, type TournamentStandingsResponse } from '@/lib/api/tournaments'
import { cn } from '@/lib/utils'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'

export default function TournamentDetailPage() {
  const params = useParams()
  const router = useRouter()
  const tournamentId = params.id as string

  const [tournament, setTournament] = useState<TournamentDetailsResponse | null>(null)
  const [standings, setStandings] = useState<TournamentStandingsResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [joining, setJoining] = useState(false)

  useEffect(() => {
    fetchData()
  }, [tournamentId])

  const fetchData = async () => {
    try {
      const [details, standingsData] = await Promise.all([
        tournamentsApi.getTournamentDetails(tournamentId),
        tournamentsApi.getTournamentStandings(tournamentId).catch(() => null),
      ])
      setTournament(details)
      setStandings(standingsData)
    } catch (error) {
      toast.error('Error al cargar el torneo')
      router.push('/tournaments')
    } finally {
      setLoading(false)
    }
  }

  const handleJoin = async () => {
    setJoining(true)
    try {
      await tournamentsApi.joinTournament(tournamentId)
      toast.success('Inscripcion exitosa!')
      fetchData()
    } catch (error: any) {
      toast.error(error.response?.data?.message || 'Error al inscribirse')
    } finally {
      setJoining(false)
    }
  }

  if (loading || !tournament) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  const canJoin = tournament.status === 'registration' &&
    tournament.current_participants < tournament.max_participants &&
    !tournament.is_registered

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Button variant="ghost" onClick={() => router.push('/tournaments')}>
          <ArrowLeft className="w-4 h-4 mr-2" />
          Volver
        </Button>
      </div>

      {/* Tournament Info */}
      <Card>
        <CardHeader>
          <div className="flex items-start justify-between">
            <div>
              <Badge className="mb-2">
                {tournament.status === 'registration' ? 'Inscripcion Abierta' : tournament.status}
              </Badge>
              <CardTitle className="text-2xl">{tournament.name}</CardTitle>
              <p className="text-muted-foreground mt-2">{tournament.description}</p>
            </div>
            {canJoin && (
              <Button size="lg" onClick={handleJoin} disabled={joining}>
                {joining ? (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                ) : (
                  <Trophy className="mr-2 h-4 w-4" />
                )}
                Inscribirse
              </Button>
            )}
            {tournament.is_registered && (
              <Badge variant="secondary" className="text-lg py-2 px-4">
                <CheckCircle className="w-4 h-4 mr-2" />
                Inscrito
              </Badge>
            )}
          </div>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-full bg-blue-100 flex items-center justify-center">
                <Calendar className="w-5 h-5 text-blue-600" />
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Inicio</p>
                <p className="font-medium">
                  {format(new Date(tournament.start_time), 'dd MMM, HH:mm', { locale: es })}
                </p>
              </div>
            </div>
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-full bg-green-100 flex items-center justify-center">
                <Users className="w-5 h-5 text-green-600" />
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Participantes</p>
                <p className="font-medium">
                  {tournament.current_participants}/{tournament.max_participants}
                </p>
              </div>
            </div>
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-full bg-yellow-100 flex items-center justify-center">
                <Coins className="w-5 h-5 text-yellow-600" />
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Entry Fee</p>
                <p className="font-medium">{tournament.entry_fee} tokens</p>
              </div>
            </div>
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-full bg-purple-100 flex items-center justify-center">
                <Trophy className="w-5 h-5 text-purple-600" />
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Premio Total</p>
                <p className="font-medium">{tournament.prize_pool.toLocaleString()} tokens</p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <Tabs defaultValue="prizes">
        <TabsList>
          <TabsTrigger value="prizes">Premios</TabsTrigger>
          <TabsTrigger value="standings">Clasificacion</TabsTrigger>
          <TabsTrigger value="rules">Reglas</TabsTrigger>
        </TabsList>

        <TabsContent value="prizes" className="mt-4">
          <Card>
            <CardHeader>
              <CardTitle>Distribucion de Premios</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {tournament.prizes?.map((prize, index) => (
                  <div
                    key={index}
                    className={cn(
                      'flex items-center justify-between p-4 rounded-lg',
                      index === 0 && 'bg-yellow-50 border border-yellow-200',
                      index === 1 && 'bg-gray-100',
                      index === 2 && 'bg-amber-50',
                      index > 2 && 'bg-gray-50'
                    )}
                  >
                    <div className="flex items-center gap-3">
                      <span className="text-2xl">
                        {index === 0 && '🥇'}
                        {index === 1 && '🥈'}
                        {index === 2 && '🥉'}
                        {index > 2 && `#${prize.position_from}`}
                      </span>
                      <span className="font-medium">
                        {prize.position_from === prize.position_to
                          ? `Posicion ${prize.position_from}`
                          : `Posiciones ${prize.position_from}-${prize.position_to}`}
                      </span>
                    </div>
                    <div className="text-right">
                      <p className="font-bold text-lg">{prize.token_reward.toLocaleString()} tokens</p>
                      {prize.special_reward && (
                        <p className="text-sm text-muted-foreground">{prize.special_reward}</p>
                      )}
                    </div>
                  </div>
                )) || (
                  <p className="text-center text-muted-foreground py-4">
                    Premios por definir
                  </p>
                )}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="standings" className="mt-4">
          <Card>
            <CardHeader>
              <CardTitle>Clasificacion Actual</CardTitle>
            </CardHeader>
            <CardContent>
              {standings?.participants && standings.participants.length > 0 ? (
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead className="w-16">Pos.</TableHead>
                      <TableHead>Jugador</TableHead>
                      <TableHead>Rango</TableHead>
                      <TableHead className="text-right">Puntos</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {standings.participants.map((participant) => (
                      <TableRow key={participant.user_id}>
                        <TableCell className="font-medium">
                          {participant.position <= 3 && (
                            <span className="mr-1">
                              {participant.position === 1 && '🥇'}
                              {participant.position === 2 && '🥈'}
                              {participant.position === 3 && '🥉'}
                            </span>
                          )}
                          #{participant.position}
                        </TableCell>
                        <TableCell>
                          <div className="flex items-center gap-2">
                            <Avatar className="h-8 w-8">
                              <AvatarImage src={participant.profile_picture_url} />
                              <AvatarFallback>
                                {participant.username.charAt(0).toUpperCase()}
                              </AvatarFallback>
                            </Avatar>
                            <span>{participant.username}</span>
                          </div>
                        </TableCell>
                        <TableCell>
                          <RankBadge tier={participant.rank_tier} size="sm" showName={false} />
                        </TableCell>
                        <TableCell className="text-right font-medium">
                          {participant.score}
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              ) : (
                <p className="text-center text-muted-foreground py-8">
                  El torneo aun no ha comenzado
                </p>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="rules" className="mt-4">
          <Card>
            <CardHeader>
              <CardTitle>Reglas del Torneo</CardTitle>
            </CardHeader>
            <CardContent className="prose prose-sm max-w-none">
              <ul className="space-y-2 text-muted-foreground">
                <li>El torneo consiste en multiples rondas de simulador</li>
                <li>Cada ronda tiene un tiempo limite de decision</li>
                <li>Los puntos se acumulan segun las respuestas correctas</li>
                <li>En caso de empate, se considera el tiempo de respuesta</li>
                <li>Rango minimo requerido: {tournament.min_rank_required}</li>
                <li>Entry fee: {tournament.entry_fee} tokens (no reembolsable)</li>
              </ul>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}

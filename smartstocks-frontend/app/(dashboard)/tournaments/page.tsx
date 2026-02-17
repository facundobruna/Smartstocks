'use client'

import { useEffect, useState } from 'react'
import { Gamepad2, Trophy, Calendar, Loader2 } from 'lucide-react'
import { toast } from 'sonner'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { TournamentCard } from '@/components/tournaments/tournament-card'
import { tournamentsApi } from '@/lib/api/tournaments'
import type { Tournament } from '@/types'

export default function TournamentsPage() {
  const [activeTournaments, setActiveTournaments] = useState<Tournament[]>([])
  const [myTournaments, setMyTournaments] = useState<Tournament[]>([])
  const [loading, setLoading] = useState(true)
  const [joiningId, setJoiningId] = useState<string | null>(null)

  useEffect(() => {
    fetchTournaments()
  }, [])

  const fetchTournaments = async () => {
    try {
      const [active, mine] = await Promise.all([
        tournamentsApi.getActiveTournaments(),
        tournamentsApi.getMyTournaments(),
      ])
      setActiveTournaments(active)
      setMyTournaments(mine)
    } catch (error) {
      toast.error('Error al cargar torneos')
    } finally {
      setLoading(false)
    }
  }

  const handleJoin = async (tournamentId: string) => {
    setJoiningId(tournamentId)
    try {
      await tournamentsApi.joinTournament(tournamentId)
      toast.success('Inscripcion exitosa!')
      fetchTournaments() // Refresh
    } catch (error: any) {
      toast.error(error.response?.data?.message || 'Error al inscribirse')
    } finally {
      setJoiningId(null)
    }
  }

  const isRegistered = (tournamentId: string) => {
    return myTournaments.some((t) => t.id === tournamentId)
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-foreground flex items-center gap-3">
          <Gamepad2 className="w-8 h-8 text-primary" />
          Torneos
        </h1>
        <p className="text-muted-foreground mt-1">
          Compite en torneos y gana tokens y premios exclusivos
        </p>
      </div>

      <Tabs defaultValue="active">
        <TabsList>
          <TabsTrigger value="active" className="flex items-center gap-2">
            <Trophy className="w-4 h-4" />
            Torneos Activos
          </TabsTrigger>
          <TabsTrigger value="my" className="flex items-center gap-2">
            <Calendar className="w-4 h-4" />
            Mis Torneos
          </TabsTrigger>
        </TabsList>

        <TabsContent value="active" className="mt-6">
          {activeTournaments.length > 0 ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {activeTournaments.map((tournament) => (
                <TournamentCard
                  key={tournament.id}
                  tournament={tournament}
                  isRegistered={isRegistered(tournament.id)}
                  onJoin={() => handleJoin(tournament.id)}
                  loading={joiningId === tournament.id}
                />
              ))}
            </div>
          ) : (
            <div className="text-center py-12">
              <Gamepad2 className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
              <h3 className="text-lg font-medium">No hay torneos activos</h3>
              <p className="text-muted-foreground">
                Vuelve pronto para ver nuevos torneos
              </p>
            </div>
          )}
        </TabsContent>

        <TabsContent value="my" className="mt-6">
          {myTournaments.length > 0 ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {myTournaments.map((tournament) => (
                <TournamentCard
                  key={tournament.id}
                  tournament={tournament}
                  isRegistered={true}
                />
              ))}
            </div>
          ) : (
            <div className="text-center py-12">
              <Calendar className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
              <h3 className="text-lg font-medium">No estas inscrito en ningun torneo</h3>
              <p className="text-muted-foreground">
                Inscribete en un torneo activo para comenzar
              </p>
            </div>
          )}
        </TabsContent>
      </Tabs>
    </div>
  )
}

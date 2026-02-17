'use client'

import { useEffect, useState } from 'react'
import { Globe, School, Trophy, Users, Loader2 } from 'lucide-react'
import { toast } from 'sonner'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { LeaderboardTable } from '@/components/rankings/leaderboard-table'
import { RankBadge } from '@/components/rankings/rank-badge'
import { rankingsApi } from '@/lib/api/rankings'
import { useAuthStore } from '@/lib/stores/auth-store'
import type { LeaderboardResponse, UserPositionResponse } from '@/types'

export default function RankingsPage() {
  const { user } = useAuthStore()
  const [globalData, setGlobalData] = useState<LeaderboardResponse | null>(null)
  const [schoolData, setSchoolData] = useState<LeaderboardResponse | null>(null)
  const [myPosition, setMyPosition] = useState<UserPositionResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [activeTab, setActiveTab] = useState('global')

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [global, position] = await Promise.all([
          rankingsApi.getGlobalLeaderboard(50),
          rankingsApi.getMyPosition(),
        ])
        setGlobalData(global)
        setMyPosition(position)

        // Try to fetch school leaderboard if user has a school
        if (user?.school_id) {
          try {
            const school = await rankingsApi.getMySchoolLeaderboard(50)
            setSchoolData(school)
          } catch (error) {
            // School leaderboard not available
          }
        }
      } catch (error: any) {
        toast.error('Error al cargar los rankings')
      } finally {
        setLoading(false)
      }
    }

    fetchData()
  }, [user?.school_id])

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
          <Trophy className="w-8 h-8 text-primary" />
          Rankings
        </h1>
        <p className="text-muted-foreground mt-1">
          Compara tu progreso con otros estudiantes
        </p>
      </div>

      {/* Position Card */}
      {myPosition && (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <Card>
            <CardContent className="pt-6">
              <div className="flex items-center gap-4">
                <div className="w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center">
                  <Globe className="w-6 h-6 text-primary" />
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Posicion Global</p>
                  <p className="text-2xl font-bold">#{myPosition.global_position}</p>
                  <p className="text-xs text-muted-foreground">
                    de {myPosition.total_players.toLocaleString()} jugadores
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>

          {myPosition.school_position && (
            <Card>
              <CardContent className="pt-6">
                <div className="flex items-center gap-4">
                  <div className="w-12 h-12 rounded-full bg-green-100 flex items-center justify-center">
                    <School className="w-6 h-6 text-green-600" />
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">En tu Colegio</p>
                    <p className="text-2xl font-bold">#{myPosition.school_position}</p>
                  </div>
                </div>
              </CardContent>
            </Card>
          )}

          <Card>
            <CardContent className="pt-6">
              <div className="flex items-center gap-4">
                <div className="w-12 h-12 rounded-full bg-yellow-100 flex items-center justify-center">
                  <Users className="w-6 h-6 text-yellow-600" />
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Top Percentil</p>
                  <p className="text-2xl font-bold">
                    {((myPosition.global_position / myPosition.total_players) * 100).toFixed(1)}%
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Leaderboard */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Trophy className="w-5 h-5 text-yellow-500" />
            Top Jugadores
          </CardTitle>
        </CardHeader>
        <CardContent>
          <Tabs value={activeTab} onValueChange={setActiveTab}>
            <TabsList className="mb-4">
              <TabsTrigger value="global" className="flex items-center gap-2">
                <Globe className="w-4 h-4" />
                Global
              </TabsTrigger>
              {schoolData && (
                <TabsTrigger value="school" className="flex items-center gap-2">
                  <School className="w-4 h-4" />
                  Mi Colegio
                </TabsTrigger>
              )}
            </TabsList>

            <TabsContent value="global">
              {globalData && globalData.top_players.length > 0 ? (
                <LeaderboardTable entries={globalData.top_players} showSchool />
              ) : (
                <div className="text-center py-8 text-muted-foreground">
                  No hay datos disponibles
                </div>
              )}
            </TabsContent>

            {schoolData && (
              <TabsContent value="school">
                {schoolData.top_players.length > 0 ? (
                  <LeaderboardTable entries={schoolData.top_players} />
                ) : (
                  <div className="text-center py-8 text-muted-foreground">
                    No hay datos disponibles para tu colegio
                  </div>
                )}
              </TabsContent>
            )}
          </Tabs>
        </CardContent>
      </Card>
    </div>
  )
}

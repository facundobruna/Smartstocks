'use client'

import { useEffect, useState } from 'react'
import {
  User,
  Trophy,
  Target,
  Flame,
  Medal,
  TrendingUp,
  Award,
  Lock,
  Loader2,
  Mail,
  School,
  Calendar,
} from 'lucide-react'
import { toast } from 'sonner'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { Separator } from '@/components/ui/separator'
import { RankBadge } from '@/components/rankings/rank-badge'
import { rankingsApi } from '@/lib/api/rankings'
import { simulatorApi } from '@/lib/api/simulator'
import { useAuthStore } from '@/lib/stores/auth-store'
import { cn } from '@/lib/utils'
import type { AllAchievementsResponse, SimulatorHistoryResponse } from '@/types'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'

export default function ProfilePage() {
  const { user, stats } = useAuthStore()
  const [achievements, setAchievements] = useState<AllAchievementsResponse | null>(null)
  const [history, setHistory] = useState<SimulatorHistoryResponse | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [achievementsData, historyData] = await Promise.all([
          rankingsApi.getMyAchievements(),
          simulatorApi.getHistory(10),
        ])
        setAchievements(achievementsData)
        setHistory(historyData)
      } catch (error) {
        toast.error('Error al cargar datos del perfil')
      } finally {
        setLoading(false)
      }
    }

    fetchData()
  }, [])

  if (!user || !stats) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  const winRate = stats.total_wins + stats.total_losses > 0
    ? (stats.total_wins / (stats.total_wins + stats.total_losses)) * 100
    : 0

  return (
    <div className="space-y-6">
      {/* Profile Header */}
      <Card>
        <CardContent className="pt-6">
          <div className="flex flex-col md:flex-row gap-6 items-start md:items-center">
            {/* Avatar */}
            <Avatar className="w-24 h-24">
              <AvatarImage src={user.profile_picture_url || undefined} />
              <AvatarFallback className="bg-primary text-white text-3xl">
                {user.username.charAt(0).toUpperCase()}
              </AvatarFallback>
            </Avatar>

            {/* Info */}
            <div className="flex-1">
              <div className="flex items-center gap-3 flex-wrap">
                <h1 className="text-2xl font-bold">{user.username}</h1>
                <RankBadge tier={stats.rank_tier} size="md" />
              </div>

              <div className="flex flex-wrap gap-4 mt-3 text-sm text-muted-foreground">
                <span className="flex items-center gap-1">
                  <Mail className="w-4 h-4" />
                  {user.email}
                </span>
                {user.school_id && (
                  <span className="flex items-center gap-1">
                    <School className="w-4 h-4" />
                    Colegio registrado
                  </span>
                )}
                <span className="flex items-center gap-1">
                  <Calendar className="w-4 h-4" />
                  Miembro desde {format(new Date(user.created_at), 'MMM yyyy', { locale: es })}
                </span>
              </div>
            </div>

            {/* Quick Stats */}
            <div className="flex gap-6 text-center">
              <div>
                <p className="text-3xl font-bold text-primary">{stats.smartpoints.toLocaleString()}</p>
                <p className="text-sm text-muted-foreground">Puntos</p>
              </div>
              <div>
                <p className="text-3xl font-bold">{stats.total_wins}</p>
                <p className="text-sm text-muted-foreground">Victorias</p>
              </div>
              <div>
                <p className="text-3xl font-bold">{stats.win_streak}</p>
                <p className="text-sm text-muted-foreground">Racha</p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Stats */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <TrendingUp className="w-5 h-5 text-primary" />
              Estadisticas
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="p-4 bg-gray-50 rounded-lg">
                <div className="flex items-center gap-2 text-muted-foreground mb-1">
                  <Medal className="w-4 h-4" />
                  <span className="text-sm">Victorias / Derrotas</span>
                </div>
                <p className="text-xl font-bold">
                  <span className="text-bullish">{stats.total_wins}</span>
                  <span className="text-muted-foreground"> / </span>
                  <span className="text-bearish">{stats.total_losses}</span>
                </p>
              </div>
              <div className="p-4 bg-gray-50 rounded-lg">
                <div className="flex items-center gap-2 text-muted-foreground mb-1">
                  <Target className="w-4 h-4" />
                  <span className="text-sm">Win Rate</span>
                </div>
                <p className="text-xl font-bold">{winRate.toFixed(1)}%</p>
              </div>
              <div className="p-4 bg-gray-50 rounded-lg">
                <div className="flex items-center gap-2 text-muted-foreground mb-1">
                  <Flame className="w-4 h-4" />
                  <span className="text-sm">Racha Actual</span>
                </div>
                <p className="text-xl font-bold">{stats.win_streak}</p>
              </div>
              <div className="p-4 bg-gray-50 rounded-lg">
                <div className="flex items-center gap-2 text-muted-foreground mb-1">
                  <TrendingUp className="w-4 h-4" />
                  <span className="text-sm">Partidas Simulador</span>
                </div>
                <p className="text-xl font-bold">{stats.total_simulator_games}</p>
              </div>
            </div>

            {/* Simulator Stats */}
            {history?.stats && (
              <>
                <Separator />
                <div>
                  <h4 className="font-medium mb-3">Precision por Dificultad</h4>
                  <div className="space-y-3">
                    {Object.entries(history.stats.by_difficulty || {}).map(([difficulty, data]) => (
                      <div key={difficulty}>
                        <div className="flex justify-between text-sm mb-1">
                          <span className="capitalize">{difficulty}</span>
                          <span>{data.accuracy_rate.toFixed(1)}%</span>
                        </div>
                        <Progress value={data.accuracy_rate} className="h-2" />
                      </div>
                    ))}
                  </div>
                </div>
              </>
            )}
          </CardContent>
        </Card>

        {/* Achievements */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Award className="w-5 h-5 text-yellow-500" />
              Logros
              {achievements && (
                <Badge variant="secondary" className="ml-2">
                  {achievements.unlocked.length}/{achievements.total_count}
                </Badge>
              )}
            </CardTitle>
          </CardHeader>
          <CardContent>
            {loading ? (
              <div className="flex justify-center py-8">
                <Loader2 className="h-6 w-6 animate-spin text-primary" />
              </div>
            ) : achievements ? (
              <div className="grid grid-cols-4 gap-3">
                {/* Unlocked */}
                {achievements.unlocked.map((achievement) => (
                  <div
                    key={achievement.id}
                    className="p-3 bg-yellow-50 border border-yellow-200 rounded-lg text-center cursor-pointer hover:bg-yellow-100 transition-colors"
                    title={achievement.achievement_description}
                  >
                    <div className="text-2xl mb-1">
                      {achievement.icon_url || '🏆'}
                    </div>
                    <p className="text-xs font-medium truncate">
                      {achievement.achievement_name}
                    </p>
                  </div>
                ))}

                {/* Locked */}
                {achievements.locked.map((achievement) => (
                  <div
                    key={achievement.achievement_type}
                    className="p-3 bg-gray-100 border border-gray-200 rounded-lg text-center opacity-60"
                    title={`${achievement.description} (${achievement.current}/${achievement.required})`}
                  >
                    <div className="text-2xl mb-1 grayscale">
                      <Lock className="w-6 h-6 mx-auto text-gray-400" />
                    </div>
                    <p className="text-xs font-medium truncate text-gray-500">
                      ???
                    </p>
                    <Progress
                      value={achievement.progress}
                      className="h-1 mt-1"
                    />
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-center text-muted-foreground py-4">
                No hay logros disponibles
              </p>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Recent Activity */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <TrendingUp className="w-5 h-5" />
            Actividad Reciente
          </CardTitle>
        </CardHeader>
        <CardContent>
          {history && history.attempts.length > 0 ? (
            <div className="space-y-3">
              {history.attempts.slice(0, 5).map((attempt) => (
                <div
                  key={attempt.id}
                  className="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
                >
                  <div className="flex items-center gap-3">
                    <div className={cn(
                      'w-10 h-10 rounded-full flex items-center justify-center',
                      attempt.was_correct ? 'bg-green-100' : 'bg-red-100'
                    )}>
                      {attempt.was_correct ? (
                        <Trophy className="w-5 h-5 text-bullish" />
                      ) : (
                        <Target className="w-5 h-5 text-bearish" />
                      )}
                    </div>
                    <div>
                      <p className="font-medium">
                        Simulador - <span className="capitalize">{attempt.difficulty}</span>
                      </p>
                      <p className="text-sm text-muted-foreground">
                        Decision: {attempt.user_decision}
                      </p>
                    </div>
                  </div>
                  <div className="text-right">
                    <p className={cn(
                      'font-bold',
                      attempt.was_correct ? 'text-bullish' : 'text-bearish'
                    )}>
                      {attempt.was_correct ? '+' : ''}{attempt.points_earned} pts
                    </p>
                    <p className="text-xs text-muted-foreground">
                      {format(new Date(attempt.created_at), 'dd MMM', { locale: es })}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-center text-muted-foreground py-8">
              No hay actividad reciente. ¡Juega al simulador para comenzar!
            </p>
          )}
        </CardContent>
      </Card>
    </div>
  )
}

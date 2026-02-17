'use client'

import { Trophy, Target, Flame, TrendingUp, Medal, Gamepad2 } from 'lucide-react'
import { WelcomeBanner } from '@/components/dashboard/welcome-banner'
import { StatsCard } from '@/components/dashboard/stats-card'
import { QuickActions } from '@/components/dashboard/quick-actions'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useAuthStore } from '@/lib/stores/auth-store'

export default function HomePage() {
  const { stats } = useAuthStore()

  return (
    <div className="space-y-6">
      {/* Welcome Banner */}
      <WelcomeBanner />

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatsCard
          title="Puntos Totales"
          value={stats?.smartpoints || 0}
          icon={Trophy}
          description="Tu puntuacion acumulada"
          iconClassName="bg-primary"
        />
        <StatsCard
          title="Victorias"
          value={stats?.total_wins || 0}
          icon={Medal}
          description={`${stats?.total_losses || 0} derrotas`}
          iconClassName="bg-green-500"
        />
        <StatsCard
          title="Racha Actual"
          value={stats?.win_streak || 0}
          icon={Flame}
          description="Victorias consecutivas"
          iconClassName="bg-orange-500"
        />
        <StatsCard
          title="Partidas Simulador"
          value={stats?.total_simulator_games || 0}
          icon={TrendingUp}
          description="Escenarios completados"
          iconClassName="bg-blue-500"
        />
      </div>

      {/* Quick Actions */}
      <QuickActions />

      {/* Activity & Progress */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Recent Activity */}
        <Card>
          <CardHeader>
            <CardTitle className="text-lg flex items-center gap-2">
              <Gamepad2 className="w-5 h-5" />
              Actividad Reciente
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-center py-8 text-muted-foreground">
              <p>Comienza a jugar para ver tu actividad aqui</p>
            </div>
          </CardContent>
        </Card>

        {/* Daily Challenges */}
        <Card>
          <CardHeader>
            <CardTitle className="text-lg flex items-center gap-2">
              <Target className="w-5 h-5" />
              Desafios Diarios
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <div>
                  <p className="font-medium text-sm">Completa 1 escenario facil</p>
                  <p className="text-xs text-muted-foreground">+25 puntos</p>
                </div>
                <div className="text-sm text-muted-foreground">0/1</div>
              </div>
              <div className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <div>
                  <p className="font-medium text-sm">Gana 1 partida PvP</p>
                  <p className="text-xs text-muted-foreground">+50 puntos</p>
                </div>
                <div className="text-sm text-muted-foreground">0/1</div>
              </div>
              <div className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <div>
                  <p className="font-medium text-sm">Completa 1 leccion</p>
                  <p className="text-xs text-muted-foreground">+15 puntos</p>
                </div>
                <div className="text-sm text-muted-foreground">0/1</div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

'use client'

import { Trophy, Coins, Zap } from 'lucide-react'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { useAuthStore } from '@/lib/stores/auth-store'
import { RANK_TIERS } from '@/types'

export function WelcomeBanner() {
  const { user, stats } = useAuthStore()

  if (!user || !stats) return null

  const rankInfo = RANK_TIERS[stats.rank_tier as keyof typeof RANK_TIERS] || {
    name: stats.rank_tier,
    color: '#60A5FA',
  }

  return (
    <Card className="bg-gradient-to-r from-primary to-blue-600 text-white border-0">
      <CardContent className="pt-6">
        <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
          <div>
            <h1 className="text-2xl font-bold mb-1">
              Hola, {user.username}!
            </h1>
            <p className="text-white/80">
              Bienvenido de vuelta a SmartStocks
            </p>
          </div>

          <div className="flex flex-wrap gap-3">
            <Badge
              variant="secondary"
              className="bg-white/20 text-white hover:bg-white/30 px-3 py-1.5"
            >
              <Trophy className="w-4 h-4 mr-1.5" />
              {rankInfo.name}
            </Badge>
            <Badge
              variant="secondary"
              className="bg-white/20 text-white hover:bg-white/30 px-3 py-1.5"
            >
              <Coins className="w-4 h-4 mr-1.5" />
              {stats.smartpoints.toLocaleString()} pts
            </Badge>
            <Badge
              variant="secondary"
              className="bg-white/20 text-white hover:bg-white/30 px-3 py-1.5"
            >
              <Zap className="w-4 h-4 mr-1.5" />
              Racha: {stats.win_streak}
            </Badge>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

export default WelcomeBanner

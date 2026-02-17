import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Card, CardContent } from '@/components/ui/card'
import { RankBadge } from '@/components/rankings/rank-badge'
import { cn } from '@/lib/utils'
import type { User } from '@/types'

interface PlayerCardProps {
  user: User | null
  score: number
  isCurrentUser?: boolean
  rankTier?: string
  isSearching?: boolean
  hasDecided?: boolean
}

export function PlayerCard({
  user,
  score,
  isCurrentUser = false,
  rankTier = 'bronze_1',
  isSearching = false,
  hasDecided = false,
}: PlayerCardProps) {
  return (
    <Card className={cn(
      'transition-all',
      isCurrentUser && 'border-primary border-2',
      hasDecided && 'ring-2 ring-green-500'
    )}>
      <CardContent className="pt-6 text-center">
        <Avatar className="w-20 h-20 mx-auto mb-3">
          {isSearching ? (
            <AvatarFallback className="bg-gray-200 animate-pulse">
              ?
            </AvatarFallback>
          ) : (
            <>
              <AvatarImage src={user?.profile_picture_url || undefined} />
              <AvatarFallback className={cn(
                'text-2xl',
                isCurrentUser ? 'bg-primary text-white' : 'bg-gray-200'
              )}>
                {user?.username?.charAt(0).toUpperCase() || '?'}
              </AvatarFallback>
            </>
          )}
        </Avatar>

        <h3 className="font-semibold text-lg mb-1">
          {isSearching ? (
            <span className="text-muted-foreground">Buscando...</span>
          ) : (
            <>
              {user?.username}
              {isCurrentUser && <span className="text-primary ml-1">(Tu)</span>}
            </>
          )}
        </h3>

        {!isSearching && (
          <div className="flex justify-center mb-3">
            <RankBadge tier={rankTier} size="sm" />
          </div>
        )}

        <div className="mt-4 pt-4 border-t">
          <p className="text-3xl font-bold">{score}</p>
          <p className="text-sm text-muted-foreground">Puntos</p>
        </div>

        {hasDecided && (
          <p className="text-xs text-green-600 mt-2">Ya decidio</p>
        )}
      </CardContent>
    </Card>
  )
}

export default PlayerCard

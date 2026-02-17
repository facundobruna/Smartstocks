'use client'

import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { ScrollArea } from '@/components/ui/scroll-area'
import { RankBadge } from './rank-badge'
import { cn } from '@/lib/utils'
import type { LeaderboardEntry } from '@/types'

interface LeaderboardTableProps {
  entries: LeaderboardEntry[]
  showSchool?: boolean
}

function getMedalEmoji(position: number): string {
  switch (position) {
    case 1:
      return '🥇'
    case 2:
      return '🥈'
    case 3:
      return '🥉'
    default:
      return ''
  }
}

export function LeaderboardTable({ entries, showSchool = false }: LeaderboardTableProps) {
  return (
    <ScrollArea className="h-[500px]">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="w-16">Pos.</TableHead>
            <TableHead>Usuario</TableHead>
            <TableHead className="text-center">Rango</TableHead>
            <TableHead className="text-right">Puntos</TableHead>
            <TableHead className="text-right">Victorias</TableHead>
            <TableHead className="text-right">Win Rate</TableHead>
            {showSchool && <TableHead>Colegio</TableHead>}
          </TableRow>
        </TableHeader>
        <TableBody>
          {entries.map((entry) => (
            <TableRow
              key={entry.user_id}
              className={cn(
                entry.is_current_user && 'bg-primary/10 border-primary/20'
              )}
            >
              <TableCell className="font-medium">
                <span className="flex items-center gap-1">
                  {getMedalEmoji(entry.rank_position)}
                  <span className={cn(
                    entry.rank_position <= 3 && 'font-bold',
                    entry.rank_position === 1 && 'text-yellow-600',
                    entry.rank_position === 2 && 'text-gray-500',
                    entry.rank_position === 3 && 'text-amber-700'
                  )}>
                    #{entry.rank_position}
                  </span>
                </span>
              </TableCell>
              <TableCell>
                <div className="flex items-center gap-3">
                  <Avatar className="h-8 w-8">
                    <AvatarImage src={entry.profile_picture_url} />
                    <AvatarFallback className="bg-primary text-white text-xs">
                      {entry.username.charAt(0).toUpperCase()}
                    </AvatarFallback>
                  </Avatar>
                  <span className={cn(
                    'font-medium',
                    entry.is_current_user && 'text-primary'
                  )}>
                    {entry.username}
                    {entry.is_current_user && ' (Tu)'}
                  </span>
                </div>
              </TableCell>
              <TableCell>
                <div className="flex justify-center">
                  <RankBadge tier={entry.rank_tier} size="sm" showName={false} />
                </div>
              </TableCell>
              <TableCell className="text-right font-medium">
                {entry.smartpoints.toLocaleString()}
              </TableCell>
              <TableCell className="text-right">
                <span className="text-bullish">{entry.total_wins}</span>
                <span className="text-muted-foreground"> / </span>
                <span className="text-bearish">{entry.total_losses}</span>
              </TableCell>
              <TableCell className="text-right">
                {(entry.win_rate * 100).toFixed(1)}%
              </TableCell>
              {showSchool && (
                <TableCell className="text-muted-foreground text-sm">
                  {entry.school_name || '-'}
                </TableCell>
              )}
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </ScrollArea>
  )
}

export default LeaderboardTable

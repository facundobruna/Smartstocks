import { cn } from '@/lib/utils'
import { RANK_TIERS } from '@/types'

interface RankBadgeProps {
  tier: string
  size?: 'sm' | 'md' | 'lg'
  showName?: boolean
  className?: string
}

export function RankBadge({ tier, size = 'md', showName = true, className }: RankBadgeProps) {
  const rankInfo = RANK_TIERS[tier as keyof typeof RANK_TIERS] || {
    name: tier,
    color: '#60A5FA',
    level: 0,
  }

  const sizeClasses = {
    sm: 'w-6 h-6 text-xs',
    md: 'w-8 h-8 text-sm',
    lg: 'w-12 h-12 text-base',
  }

  const stars = tier.includes('_') ? parseInt(tier.split('_')[1]) || 0 : 0

  return (
    <div className={cn('flex items-center gap-2', className)}>
      {/* Hexagon Badge */}
      <div
        className={cn(
          'relative flex items-center justify-center font-bold text-white',
          sizeClasses[size]
        )}
        style={{
          backgroundColor: rankInfo.color,
          clipPath: 'polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%)',
        }}
      >
        {stars > 0 && (
          <span className="text-yellow-300">
            {'★'.repeat(stars)}
          </span>
        )}
        {tier === 'master' && <span>M</span>}
      </div>

      {showName && (
        <span
          className="font-medium"
          style={{ color: rankInfo.color }}
        >
          {rankInfo.name}
        </span>
      )}
    </div>
  )
}

export default RankBadge

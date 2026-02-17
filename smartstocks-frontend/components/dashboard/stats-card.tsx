import { LucideIcon } from 'lucide-react'
import { Card, CardContent } from '@/components/ui/card'
import { cn } from '@/lib/utils'

interface StatsCardProps {
  title: string
  value: string | number
  icon: LucideIcon
  description?: string
  trend?: {
    value: number
    isPositive: boolean
  }
  className?: string
  iconClassName?: string
}

export function StatsCard({
  title,
  value,
  icon: Icon,
  description,
  trend,
  className,
  iconClassName,
}: StatsCardProps) {
  return (
    <Card className={cn('transition-all hover:shadow-md', className)}>
      <CardContent className="pt-6">
        <div className="flex items-start justify-between">
          <div>
            <p className="text-sm font-medium text-muted-foreground">{title}</p>
            <p className="text-2xl font-bold text-foreground mt-1">
              {typeof value === 'number' ? value.toLocaleString() : value}
            </p>
            {description && (
              <p className="text-xs text-muted-foreground mt-1">{description}</p>
            )}
            {trend && (
              <p
                className={cn(
                  'text-xs font-medium mt-1',
                  trend.isPositive ? 'text-bullish' : 'text-bearish'
                )}
              >
                {trend.isPositive ? '+' : '-'}
                {Math.abs(trend.value)}% desde ayer
              </p>
            )}
          </div>
          <div
            className={cn(
              'w-12 h-12 rounded-lg flex items-center justify-center',
              iconClassName || 'bg-primary/10'
            )}
          >
            <Icon
              className={cn(
                'w-6 h-6',
                iconClassName ? 'text-white' : 'text-primary'
              )}
            />
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

export default StatsCard

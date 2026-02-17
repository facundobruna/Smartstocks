'use client'

import { TrendingUp, Minus, TrendingDown, Loader2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'
import type { SimulatorDecision } from '@/types'

interface DecisionButtonsProps {
  onDecision: (decision: SimulatorDecision) => void
  disabled?: boolean
  loading?: boolean
  selectedDecision?: SimulatorDecision | null
}

const decisions: {
  id: SimulatorDecision
  label: string
  description: string
  icon: typeof TrendingUp
  color: string
  bgColor: string
  hoverColor: string
}[] = [
  {
    id: 'buy',
    label: 'COMPRAR',
    description: 'El precio subira',
    icon: TrendingUp,
    color: 'text-white',
    bgColor: 'bg-bullish',
    hoverColor: 'hover:bg-green-600',
  },
  {
    id: 'hold',
    label: 'MANTENER',
    description: 'El precio se mantendra',
    icon: Minus,
    color: 'text-white',
    bgColor: 'bg-yellow-500',
    hoverColor: 'hover:bg-yellow-600',
  },
  {
    id: 'sell',
    label: 'VENDER',
    description: 'El precio bajara',
    icon: TrendingDown,
    color: 'text-white',
    bgColor: 'bg-bearish',
    hoverColor: 'hover:bg-red-600',
  },
]

export function DecisionButtons({
  onDecision,
  disabled = false,
  loading = false,
  selectedDecision = null,
}: DecisionButtonsProps) {
  return (
    <div className="space-y-4">
      <h3 className="text-lg font-semibold text-center">¿Que harias?</h3>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {decisions.map((decision) => {
          const Icon = decision.icon
          const isSelected = selectedDecision === decision.id
          const isLoading = loading && isSelected

          return (
            <Button
              key={decision.id}
              onClick={() => onDecision(decision.id)}
              disabled={disabled || loading}
              className={cn(
                'h-auto py-6 flex-col gap-2 transition-all',
                decision.bgColor,
                decision.hoverColor,
                decision.color,
                isSelected && 'ring-4 ring-offset-2',
                disabled && !isSelected && 'opacity-50'
              )}
            >
              {isLoading ? (
                <Loader2 className="w-8 h-8 animate-spin" />
              ) : (
                <Icon className="w-8 h-8" />
              )}
              <span className="text-lg font-bold">{decision.label}</span>
              <span className="text-xs opacity-80">{decision.description}</span>
            </Button>
          )
        })}
      </div>
    </div>
  )
}

export default DecisionButtons

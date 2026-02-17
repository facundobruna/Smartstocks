'use client'

import { useRouter } from 'next/navigation'
import { CheckCircle, XCircle, Trophy, TrendingUp, ArrowRight } from 'lucide-react'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { StockChart } from './stock-chart'
import { cn } from '@/lib/utils'
import type { SimulatorResult, SimulatorDecision } from '@/types'
import { RANK_TIERS } from '@/types'

interface ResultModalProps {
  result: SimulatorResult
  open: boolean
  onClose: () => void
}

const decisionLabels: Record<SimulatorDecision, string> = {
  buy: 'Comprar',
  hold: 'Mantener',
  sell: 'Vender',
}

export function ResultModal({ result, open, onClose }: ResultModalProps) {
  const router = useRouter()

  const rankInfo = RANK_TIERS[result.new_rank_tier as keyof typeof RANK_TIERS] || {
    name: result.new_rank_tier,
    color: '#60A5FA',
  }

  const handlePlayAgain = () => {
    onClose()
    router.push('/simulator')
  }

  const handleViewStats = () => {
    onClose()
    router.push('/profile')
  }

  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-3 text-2xl">
            {result.was_correct ? (
              <>
                <CheckCircle className="w-8 h-8 text-bullish" />
                <span className="text-bullish">Correcto!</span>
              </>
            ) : (
              <>
                <XCircle className="w-8 h-8 text-bearish" />
                <span className="text-bearish">Incorrecto</span>
              </>
            )}
          </DialogTitle>
        </DialogHeader>

        <div className="space-y-6 py-4">
          {/* Decision Summary */}
          <div className="grid grid-cols-2 gap-4">
            <div className="text-center p-4 bg-gray-50 rounded-lg">
              <p className="text-sm text-muted-foreground mb-1">Tu decision</p>
              <p className={cn(
                'font-bold text-lg',
                result.user_decision === result.correct_decision ? 'text-bullish' : 'text-bearish'
              )}>
                {decisionLabels[result.user_decision]}
              </p>
            </div>
            <div className="text-center p-4 bg-gray-50 rounded-lg">
              <p className="text-sm text-muted-foreground mb-1">Decision correcta</p>
              <p className="font-bold text-lg text-bullish">
                {decisionLabels[result.correct_decision]}
              </p>
            </div>
          </div>

          {/* Chart with full data */}
          <StockChart
            data={result.full_chart_data}
            showFullData={true}
            title="Grafico Completo"
          />

          {/* Explanation */}
          <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <h4 className="font-semibold text-blue-800 mb-2 flex items-center gap-2">
              <TrendingUp className="w-5 h-5" />
              Explicacion
            </h4>
            <p className="text-blue-700 text-sm leading-relaxed">
              {result.explanation}
            </p>
          </div>

          <Separator />

          {/* Stats */}
          <div className="grid grid-cols-3 gap-4 text-center">
            <div>
              <p className="text-3xl font-bold text-primary">
                +{result.points_earned}
              </p>
              <p className="text-sm text-muted-foreground">Puntos ganados</p>
            </div>
            <div>
              <p className="text-3xl font-bold">
                {result.new_total_points.toLocaleString()}
              </p>
              <p className="text-sm text-muted-foreground">Total de puntos</p>
            </div>
            <div>
              <Badge
                className="text-lg px-4 py-2"
                style={{ backgroundColor: rankInfo.color, color: 'white' }}
              >
                <Trophy className="w-4 h-4 mr-1" />
                {rankInfo.name}
              </Badge>
              <p className="text-sm text-muted-foreground mt-1">Rango</p>
            </div>
          </div>

          {/* Actions */}
          <div className="flex gap-3">
            <Button
              variant="outline"
              className="flex-1"
              onClick={handleViewStats}
            >
              Ver Estadisticas
            </Button>
            <Button
              className="flex-1"
              onClick={handlePlayAgain}
            >
              Jugar de Nuevo
              <ArrowRight className="ml-2 w-4 h-4" />
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}

export default ResultModal

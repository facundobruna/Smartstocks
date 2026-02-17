'use client'

import { useEffect, useState, useRef } from 'react'
import { useParams, useRouter } from 'next/navigation'
import { Clock, ArrowLeft, Loader2 } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent } from '@/components/ui/card'
import { ScenarioCard } from '@/components/simulator/scenario-card'
import { StockChart } from '@/components/simulator/stock-chart'
import { DecisionButtons } from '@/components/simulator/decision-buttons'
import { ResultModal } from '@/components/simulator/result-modal'
import { simulatorApi } from '@/lib/api/simulator'
import { useAuthStore } from '@/lib/stores/auth-store'
import type { SimulatorScenario, SimulatorResult, SimulatorDecision, SimulatorDifficulty } from '@/types'

const difficultyLabels: Record<SimulatorDifficulty, { name: string; color: string }> = {
  easy: { name: 'Facil', color: 'bg-green-500' },
  medium: { name: 'Medio', color: 'bg-yellow-500' },
  hard: { name: 'Dificil', color: 'bg-red-500' },
}

export default function SimulatorGamePage() {
  const params = useParams()
  const router = useRouter()
  const { updateStats } = useAuthStore()
  const difficulty = params.difficulty as SimulatorDifficulty

  const [scenario, setScenario] = useState<SimulatorScenario | null>(null)
  const [result, setResult] = useState<SimulatorResult | null>(null)
  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [selectedDecision, setSelectedDecision] = useState<SimulatorDecision | null>(null)
  const [showResult, setShowResult] = useState(false)
  const [timeRemaining, setTimeRemaining] = useState<number>(0)
  const startTimeRef = useRef<number>(Date.now())

  // Validate difficulty
  useEffect(() => {
    if (!['easy', 'medium', 'hard'].includes(difficulty)) {
      router.push('/simulator')
    }
  }, [difficulty, router])

  // Fetch scenario
  useEffect(() => {
    const fetchScenario = async () => {
      try {
        const data = await simulatorApi.getScenario(difficulty)
        setScenario(data)
        startTimeRef.current = Date.now()

        // Calculate time remaining
        const expiresAt = new Date(data.expires_at).getTime()
        const remaining = Math.max(0, Math.floor((expiresAt - Date.now()) / 1000))
        setTimeRemaining(remaining)
      } catch (error: any) {
        const message = error.response?.data?.message || 'Error al cargar el escenario'
        toast.error(message)
        router.push('/simulator')
      } finally {
        setLoading(false)
      }
    }

    fetchScenario()
  }, [difficulty, router])

  // Timer countdown
  useEffect(() => {
    if (timeRemaining <= 0 || result) return

    const interval = setInterval(() => {
      setTimeRemaining((prev) => {
        if (prev <= 1) {
          clearInterval(interval)
          // Auto-submit if time runs out
          if (!result && scenario) {
            handleDecision('hold') // Default to hold if no decision
          }
          return 0
        }
        return prev - 1
      })
    }, 1000)

    return () => clearInterval(interval)
  }, [timeRemaining, result, scenario])

  const handleDecision = async (decision: SimulatorDecision) => {
    if (!scenario || submitting) return

    setSelectedDecision(decision)
    setSubmitting(true)

    const timeTaken = Math.floor((Date.now() - startTimeRef.current) / 1000)

    try {
      const resultData = await simulatorApi.submitDecision({
        scenario_id: scenario.scenario_id,
        decision,
        time_taken_seconds: timeTaken,
      })

      setResult(resultData)
      setShowResult(true)

      // Update user stats
      updateStats({
        smartpoints: resultData.new_total_points,
        rank_tier: resultData.new_rank_tier,
      })

      if (resultData.was_correct) {
        toast.success(`Correcto! +${resultData.points_earned} puntos`)
      } else {
        toast.error('Incorrecto. La respuesta era: ' + resultData.correct_decision)
      }
    } catch (error: any) {
      const message = error.response?.data?.message || 'Error al enviar la decision'
      toast.error(message)
      setSelectedDecision(null)
    } finally {
      setSubmitting(false)
    }
  }

  const formatTime = (seconds: number): string => {
    const mins = Math.floor(seconds / 60)
    const secs = seconds % 60
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <div className="text-center">
          <Loader2 className="h-12 w-12 animate-spin text-primary mx-auto mb-4" />
          <p className="text-muted-foreground">Generando escenario...</p>
        </div>
      </div>
    )
  }

  if (!scenario) {
    return null
  }

  const diffInfo = difficultyLabels[difficulty]

  return (
    <div className="space-y-6 max-w-4xl mx-auto">
      {/* Header */}
      <div className="flex items-center justify-between">
        <Button variant="ghost" onClick={() => router.push('/simulator')}>
          <ArrowLeft className="w-4 h-4 mr-2" />
          Volver
        </Button>

        <div className="flex items-center gap-4">
          <Badge className={diffInfo.color}>{diffInfo.name}</Badge>

          <Card className="border-0 shadow-none bg-gray-100">
            <CardContent className="py-2 px-4 flex items-center gap-2">
              <Clock className={timeRemaining < 30 ? 'text-red-500' : 'text-muted-foreground'} />
              <span className={`font-mono text-lg font-bold ${timeRemaining < 30 ? 'text-red-500' : ''}`}>
                {formatTime(timeRemaining)}
              </span>
            </CardContent>
          </Card>
        </div>
      </div>

      {/* Scenario Content */}
      <ScenarioCard newsContent={scenario.news_content} />

      {/* Chart */}
      <StockChart data={scenario.chart_data} />

      {/* Decision Buttons */}
      <DecisionButtons
        onDecision={handleDecision}
        disabled={!!result}
        loading={submitting}
        selectedDecision={selectedDecision}
      />

      {/* Result Modal */}
      {result && (
        <ResultModal
          result={result}
          open={showResult}
          onClose={() => setShowResult(false)}
        />
      )}
    </div>
  )
}

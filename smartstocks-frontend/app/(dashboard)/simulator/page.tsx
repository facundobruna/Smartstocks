import { TrendingUp } from 'lucide-react'
import { DifficultySelector } from '@/components/simulator/difficulty-selector'

export default function SimulatorPage() {
  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="text-center">
        <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-primary/10 mb-4">
          <TrendingUp className="w-8 h-8 text-primary" />
        </div>
        <h1 className="text-3xl font-bold text-foreground">Simulador de Trading</h1>
        <p className="text-muted-foreground mt-2 max-w-xl mx-auto">
          Analiza noticias del mercado y graficos de precios para tomar decisiones de
          inversion. Gana puntos por cada respuesta correcta.
        </p>
      </div>

      {/* Difficulty Selector */}
      <div className="max-w-4xl mx-auto">
        <h2 className="text-xl font-semibold mb-4 text-center">
          Elige tu nivel de desafio
        </h2>
        <DifficultySelector />
      </div>

      {/* Instructions */}
      <div className="max-w-2xl mx-auto bg-blue-50 border border-blue-200 rounded-lg p-6">
        <h3 className="font-semibold text-blue-800 mb-3">Como jugar</h3>
        <ol className="list-decimal list-inside space-y-2 text-blue-700 text-sm">
          <li>Lee la noticia del mercado cuidadosamente</li>
          <li>Analiza el grafico de precios historicos</li>
          <li>Decide si el precio subira (Comprar), bajara (Vender) o se mantendra (Mantener)</li>
          <li>Recibe puntos si tu prediccion es correcta</li>
          <li>Solo puedes jugar una vez por dia en cada dificultad</li>
        </ol>
      </div>
    </div>
  )
}

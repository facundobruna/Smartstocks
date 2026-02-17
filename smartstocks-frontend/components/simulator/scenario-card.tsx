import { Newspaper } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

interface ScenarioCardProps {
  newsContent: string
}

export function ScenarioCard({ newsContent }: ScenarioCardProps) {
  return (
    <Card className="bg-gray-50 border-gray-200">
      <CardHeader className="pb-2">
        <CardTitle className="text-lg flex items-center gap-2">
          <Newspaper className="w-5 h-5 text-primary" />
          Noticia del Mercado
        </CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-foreground leading-relaxed whitespace-pre-wrap">
          {newsContent}
        </p>
      </CardContent>
    </Card>
  )
}

export default ScenarioCard

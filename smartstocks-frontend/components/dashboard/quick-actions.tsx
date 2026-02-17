import Link from 'next/link'
import { TrendingUp, Swords, BookOpen, Trophy } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'

const actions = [
  {
    title: 'Jugar Simulador',
    description: 'Practica tus habilidades de trading',
    icon: TrendingUp,
    href: '/simulator',
    color: 'bg-blue-500',
  },
  {
    title: 'Batalla PvP',
    description: 'Compite contra otros jugadores',
    icon: Swords,
    href: '/pvp',
    color: 'bg-purple-500',
  },
  {
    title: 'Aprender',
    description: 'Continua tu educacion financiera',
    icon: BookOpen,
    href: '/education',
    color: 'bg-green-500',
  },
  {
    title: 'Ver Rankings',
    description: 'Compara tu progreso',
    icon: Trophy,
    href: '/rankings',
    color: 'bg-yellow-500',
  },
]

export function QuickActions() {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-lg">Acciones Rapidas</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
          {actions.map((action) => {
            const Icon = action.icon
            return (
              <Link key={action.title} href={action.href}>
                <Button
                  variant="outline"
                  className="w-full h-auto flex-col gap-2 py-4 hover:border-primary transition-colors"
                >
                  <div
                    className={`w-10 h-10 rounded-lg ${action.color} flex items-center justify-center`}
                  >
                    <Icon className="w-5 h-5 text-white" />
                  </div>
                  <span className="text-sm font-medium">{action.title}</span>
                </Button>
              </Link>
            )
          })}
        </div>
      </CardContent>
    </Card>
  )
}

export default QuickActions

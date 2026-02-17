'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import {
  BookOpen,
  Play,
  Lock,
  CheckCircle,
  Clock,
  Star,
  TrendingUp,
  BarChart3,
  Target,
  Trophy,
  Loader2,
} from 'lucide-react'
import { toast } from 'sonner'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { cn } from '@/lib/utils'
import { coursesApi, type Course, type CoursesListResponse } from '@/lib/api/courses'

const difficultyColors = {
  principiante: 'bg-green-100 text-green-700',
  intermedio: 'bg-yellow-100 text-yellow-700',
  avanzado: 'bg-red-100 text-red-700',
}

const difficultyLabels = {
  principiante: 'Principiante',
  intermedio: 'Intermedio',
  avanzado: 'Avanzado',
}

const categoryIcons: Record<string, React.ReactNode> = {
  fundamentos: <BookOpen className="w-6 h-6" />,
  analisis: <BarChart3 className="w-6 h-6" />,
  estrategia: <Target className="w-6 h-6" />,
  avanzado: <TrendingUp className="w-6 h-6" />,
}

function CourseCard({ course, onClick }: { course: Course; onClick: () => void }) {
  const progress = course.total_lessons > 0 ? (course.completed_lessons / course.total_lessons) * 100 : 0
  const isCompleted = course.is_completed
  const isLocked = course.is_premium // Por ahora, solo premium esta bloqueado

  return (
    <Card className={cn(
      'relative overflow-hidden transition-all hover:shadow-lg cursor-pointer',
      isLocked && 'opacity-75'
    )}
    onClick={!isLocked ? onClick : undefined}
    >
      {course.is_premium && (
        <div className="absolute top-3 right-3">
          <Badge className="bg-gradient-to-r from-yellow-400 to-orange-500 text-white">
            <Star className="w-3 h-3 mr-1" />
            Premium
          </Badge>
        </div>
      )}
      <CardHeader className="pb-3">
        <div className="flex items-start gap-4">
          <div className={cn(
            'w-12 h-12 rounded-xl flex items-center justify-center',
            isCompleted ? 'bg-green-100 text-green-600' : 'bg-blue-100 text-blue-600'
          )}>
            {isCompleted ? <CheckCircle className="w-6 h-6" /> : categoryIcons[course.category] || <BookOpen className="w-6 h-6" />}
          </div>
          <div className="flex-1">
            <CardTitle className="text-lg leading-tight">{course.title}</CardTitle>
            <div className="flex items-center gap-2 mt-1">
              <Badge variant="secondary" className={difficultyColors[course.difficulty]}>
                {difficultyLabels[course.difficulty]}
              </Badge>
            </div>
          </div>
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        <p className="text-sm text-muted-foreground line-clamp-2">
          {course.description}
        </p>

        <div className="flex items-center justify-between text-sm text-muted-foreground">
          <div className="flex items-center gap-1">
            <BookOpen className="w-4 h-4" />
            <span>{course.total_lessons} lecciones</span>
          </div>
          <div className="flex items-center gap-1">
            <Clock className="w-4 h-4" />
            <span>{course.duration_minutes} min</span>
          </div>
          <div className="flex items-center gap-1">
            <Trophy className="w-4 h-4" />
            <span>{course.points_reward} pts</span>
          </div>
        </div>

        {!isLocked && (
          <div>
            <div className="flex justify-between text-xs mb-1">
              <span className="text-muted-foreground">Progreso</span>
              <span className="font-medium">{Math.round(progress)}%</span>
            </div>
            <Progress value={progress} className="h-2" />
          </div>
        )}

        <Button
          className="w-full"
          variant={isLocked ? 'secondary' : 'default'}
          disabled={isLocked}
          onClick={(e) => {
            e.stopPropagation()
            if (!isLocked) onClick()
          }}
        >
          {isLocked ? (
            <>
              <Lock className="w-4 h-4 mr-2" />
              Requiere Premium
            </>
          ) : isCompleted ? (
            <>
              <CheckCircle className="w-4 h-4 mr-2" />
              Revisar
            </>
          ) : course.completed_lessons > 0 ? (
            <>
              <Play className="w-4 h-4 mr-2" />
              Continuar
            </>
          ) : (
            <>
              <Play className="w-4 h-4 mr-2" />
              Comenzar
            </>
          )}
        </Button>
      </CardContent>
    </Card>
  )
}

export default function EducationPage() {
  const router = useRouter()
  const [selectedCategory, setSelectedCategory] = useState('all')
  const [data, setData] = useState<CoursesListResponse | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchCourses()
  }, [])

  const fetchCourses = async () => {
    try {
      const response = await coursesApi.getAllCourses()
      setData(response)
    } catch (error) {
      toast.error('Error al cargar los cursos')
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  const courses = data?.courses || []
  const totalLessons = data?.total_lessons || 0
  const completedLessons = data?.completed_lessons || 0
  const overallProgress = data?.overall_progress || 0

  const filteredCourses = selectedCategory === 'all'
    ? courses
    : courses.filter((c) => c.category === selectedCategory)

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-foreground flex items-center gap-3">
          <BookOpen className="w-8 h-8 text-primary" />
          Educacion Financiera
        </h1>
        <p className="text-muted-foreground mt-1">
          Aprende sobre finanzas, inversiones y el mercado de valores
        </p>
      </div>

      {/* Progress Overview */}
      <Card className="bg-gradient-to-r from-blue-50 to-purple-50 border-blue-100">
        <CardContent className="py-6">
          <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
            <div>
              <h2 className="text-lg font-semibold">Tu Progreso General</h2>
              <p className="text-sm text-muted-foreground">
                Has completado {completedLessons} de {totalLessons} lecciones
              </p>
            </div>
            <div className="flex items-center gap-4">
              <div className="w-48">
                <Progress value={overallProgress} className="h-3" />
              </div>
              <span className="font-bold text-lg text-primary">
                {Math.round(overallProgress)}%
              </span>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Categories */}
      <Tabs defaultValue="all" onValueChange={setSelectedCategory}>
        <TabsList>
          <TabsTrigger value="all">Todos</TabsTrigger>
          <TabsTrigger value="fundamentos">Fundamentos</TabsTrigger>
          <TabsTrigger value="analisis">Analisis</TabsTrigger>
          <TabsTrigger value="estrategia">Estrategia</TabsTrigger>
          <TabsTrigger value="avanzado">Avanzado</TabsTrigger>
        </TabsList>

        <TabsContent value={selectedCategory} className="mt-6">
          {filteredCourses.length === 0 ? (
            <div className="text-center py-12 text-muted-foreground">
              <BookOpen className="w-12 h-12 mx-auto mb-4 opacity-50" />
              <p>No hay cursos disponibles en esta categoria</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {filteredCourses.map((course) => (
                <CourseCard
                  key={course.id}
                  course={course}
                  onClick={() => router.push(`/education/courses/${course.id}`)}
                />
              ))}
            </div>
          )}
        </TabsContent>
      </Tabs>

      {/* Premium Upsell */}
      <Card className="bg-gradient-to-r from-yellow-50 to-orange-50 border-yellow-200">
        <CardContent className="py-6">
          <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 rounded-full bg-gradient-to-r from-yellow-400 to-orange-500 flex items-center justify-center">
                <Star className="w-6 h-6 text-white" />
              </div>
              <div>
                <h3 className="font-semibold text-lg">Desbloquea todos los cursos</h3>
                <p className="text-sm text-muted-foreground">
                  Accede a contenido exclusivo y cursos avanzados con Premium
                </p>
              </div>
            </div>
            <Button className="bg-gradient-to-r from-yellow-400 to-orange-500 hover:from-yellow-500 hover:to-orange-600 text-white">
              Ver Premium
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

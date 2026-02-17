'use client'

import { useEffect, useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import {
  ArrowLeft,
  BookOpen,
  Play,
  CheckCircle,
  Clock,
  Trophy,
  Loader2,
  Lock,
  FileText,
  Video,
  HelpCircle,
} from 'lucide-react'
import { toast } from 'sonner'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { Separator } from '@/components/ui/separator'
import { coursesApi, type CourseDetailResponse } from '@/lib/api/courses'
import { cn } from '@/lib/utils'

const difficultyColors = {
  principiante: 'bg-green-100 text-green-700',
  intermedio: 'bg-yellow-100 text-yellow-700',
  avanzado: 'bg-red-100 text-red-700',
}

const contentTypeIcons = {
  text: FileText,
  video: Video,
  quiz: HelpCircle,
}

export default function CoursePage() {
  const params = useParams()
  const router = useRouter()
  const courseId = params.courseId as string

  const [data, setData] = useState<CourseDetailResponse | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchCourse()
  }, [courseId])

  const fetchCourse = async () => {
    try {
      const response = await coursesApi.getCourseById(courseId)
      setData(response)
    } catch (error) {
      toast.error('Error al cargar el curso')
      router.push('/education')
    } finally {
      setLoading(false)
    }
  }

  if (loading || !data) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  const { course, lessons } = data
  const progress = course.total_lessons > 0
    ? (course.completed_lessons / course.total_lessons) * 100
    : 0

  // Find first incomplete lesson or first lesson
  const nextLesson = lessons.find(l => !l.is_completed) || lessons[0]

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Button variant="ghost" onClick={() => router.push('/education')}>
          <ArrowLeft className="w-4 h-4 mr-2" />
          Volver
        </Button>
      </div>

      {/* Course Info */}
      <Card>
        <CardHeader>
          <div className="flex flex-col md:flex-row md:items-start md:justify-between gap-4">
            <div className="flex-1">
              <div className="flex items-center gap-2 mb-2">
                <Badge className={difficultyColors[course.difficulty]}>
                  {course.difficulty.charAt(0).toUpperCase() + course.difficulty.slice(1)}
                </Badge>
                <Badge variant="outline">{course.category}</Badge>
              </div>
              <CardTitle className="text-2xl mb-2">{course.title}</CardTitle>
              <p className="text-muted-foreground">{course.description}</p>
            </div>

            {nextLesson && (
              <Button
                size="lg"
                onClick={() => router.push(`/education/lessons/${nextLesson.id}`)}
              >
                {course.completed_lessons > 0 ? (
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
            )}
          </div>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-full bg-blue-100 flex items-center justify-center">
                <BookOpen className="w-5 h-5 text-blue-600" />
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Lecciones</p>
                <p className="font-medium">{course.total_lessons}</p>
              </div>
            </div>
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-full bg-green-100 flex items-center justify-center">
                <Clock className="w-5 h-5 text-green-600" />
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Duracion</p>
                <p className="font-medium">{course.duration_minutes} min</p>
              </div>
            </div>
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-full bg-yellow-100 flex items-center justify-center">
                <Trophy className="w-5 h-5 text-yellow-600" />
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Recompensa</p>
                <p className="font-medium">{course.points_reward} pts</p>
              </div>
            </div>
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-full bg-purple-100 flex items-center justify-center">
                <CheckCircle className="w-5 h-5 text-purple-600" />
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Completado</p>
                <p className="font-medium">{course.completed_lessons}/{course.total_lessons}</p>
              </div>
            </div>
          </div>

          {/* Progress Bar */}
          <div>
            <div className="flex justify-between text-sm mb-2">
              <span className="text-muted-foreground">Progreso del curso</span>
              <span className="font-medium">{Math.round(progress)}%</span>
            </div>
            <Progress value={progress} className="h-3" />
          </div>
        </CardContent>
      </Card>

      {/* Lessons List */}
      <Card>
        <CardHeader>
          <CardTitle>Contenido del Curso</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            {lessons.map((lesson, index) => {
              const Icon = contentTypeIcons[lesson.content_type]
              const isAccessible = index === 0 || lessons[index - 1]?.is_completed

              return (
                <div
                  key={lesson.id}
                  className={cn(
                    'flex items-center gap-4 p-4 rounded-lg border transition-colors',
                    lesson.is_completed && 'bg-green-50 border-green-200',
                    !lesson.is_completed && isAccessible && 'hover:bg-gray-50 cursor-pointer',
                    !isAccessible && 'opacity-50'
                  )}
                  onClick={() => {
                    if (isAccessible) {
                      router.push(`/education/lessons/${lesson.id}`)
                    }
                  }}
                >
                  {/* Lesson Number / Status */}
                  <div className={cn(
                    'w-10 h-10 rounded-full flex items-center justify-center text-sm font-medium',
                    lesson.is_completed ? 'bg-green-500 text-white' : 'bg-gray-100'
                  )}>
                    {lesson.is_completed ? (
                      <CheckCircle className="w-5 h-5" />
                    ) : !isAccessible ? (
                      <Lock className="w-4 h-4" />
                    ) : (
                      index + 1
                    )}
                  </div>

                  {/* Lesson Info */}
                  <div className="flex-1">
                    <p className="font-medium">{lesson.title}</p>
                    <div className="flex items-center gap-3 text-sm text-muted-foreground">
                      <span className="flex items-center gap-1">
                        <Icon className="w-4 h-4" />
                        {lesson.content_type === 'text' && 'Lectura'}
                        {lesson.content_type === 'video' && 'Video'}
                        {lesson.content_type === 'quiz' && 'Quiz'}
                      </span>
                      <span className="flex items-center gap-1">
                        <Clock className="w-4 h-4" />
                        {lesson.duration_minutes} min
                      </span>
                    </div>
                  </div>

                  {/* Action */}
                  {isAccessible && !lesson.is_completed && (
                    <Button variant="ghost" size="sm">
                      <Play className="w-4 h-4" />
                    </Button>
                  )}
                </div>
              )
            })}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

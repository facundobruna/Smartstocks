'use client'

import { useEffect, useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import {
  ArrowLeft,
  ArrowRight,
  CheckCircle,
  Loader2,
  BookOpen,
  Trophy,
  HelpCircle,
} from 'lucide-react'
import { toast } from 'sonner'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group'
import { Label } from '@/components/ui/label'
import {
  coursesApi,
  type LessonDetailResponse,
  type QuizAnswer,
  type CompleteLessonResponse,
} from '@/lib/api/courses'
import { cn } from '@/lib/utils'
import { useAuthStore } from '@/lib/stores/auth-store'

export default function LessonPage() {
  const params = useParams()
  const router = useRouter()
  const lessonId = params.lessonId as string
  const { updateStats } = useAuthStore()

  const [data, setData] = useState<LessonDetailResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [completing, setCompleting] = useState(false)
  const [showResult, setShowResult] = useState(false)
  const [result, setResult] = useState<CompleteLessonResponse | null>(null)

  // Quiz state
  const [quizAnswers, setQuizAnswers] = useState<Record<string, string>>({})
  const [quizSubmitted, setQuizSubmitted] = useState(false)

  useEffect(() => {
    fetchLesson()
  }, [lessonId])

  const fetchLesson = async () => {
    try {
      const response = await coursesApi.getLessonById(lessonId)
      setData(response)
    } catch (error) {
      toast.error('Error al cargar la leccion')
      router.push('/education')
    } finally {
      setLoading(false)
    }
  }

  const handleCompleteLesson = async () => {
    setCompleting(true)
    try {
      let request = undefined

      // If quiz, send answers
      if (data?.lesson.content_type === 'quiz' && data.quiz_questions) {
        const answers: QuizAnswer[] = Object.entries(quizAnswers).map(([questionId, answer]) => ({
          question_id: questionId,
          answer,
        }))
        request = { quiz_answers: answers }
      }

      const response = await coursesApi.completeLesson(lessonId, request)
      setResult(response)
      setShowResult(true)
      setQuizSubmitted(true)

      // Update user stats
      if (response.points_earned > 0) {
        updateStats({ smartpoints: response.new_total_points })
      }

      if (response.course_completed) {
        toast.success('Felicitaciones! Completaste el curso!')
      } else {
        toast.success('Leccion completada!')
      }
    } catch (error) {
      toast.error('Error al completar la leccion')
    } finally {
      setCompleting(false)
    }
  }

  const handleQuizAnswerChange = (questionId: string, answer: string) => {
    setQuizAnswers(prev => ({ ...prev, [questionId]: answer }))
  }

  const canSubmitQuiz = () => {
    if (!data?.quiz_questions) return true
    return data.quiz_questions.every(q => quizAnswers[q.id])
  }

  if (loading || !data) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  const { lesson, quiz_questions, next_lesson_id, prev_lesson_id } = data

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <Button variant="ghost" onClick={() => router.push(`/education/courses/${lesson.course_id}`)}>
          <ArrowLeft className="w-4 h-4 mr-2" />
          Volver al curso
        </Button>

        <div className="flex items-center gap-2">
          {prev_lesson_id && (
            <Button
              variant="outline"
              size="sm"
              onClick={() => router.push(`/education/lessons/${prev_lesson_id}`)}
            >
              <ArrowLeft className="w-4 h-4" />
            </Button>
          )}
          {next_lesson_id && (
            <Button
              variant="outline"
              size="sm"
              onClick={() => router.push(`/education/lessons/${next_lesson_id}`)}
            >
              <ArrowRight className="w-4 h-4" />
            </Button>
          )}
        </div>
      </div>

      {/* Lesson Content */}
      <Card>
        <CardHeader>
          <div className="flex items-center gap-2 mb-2">
            <Badge variant="outline">
              {lesson.content_type === 'text' && 'Lectura'}
              {lesson.content_type === 'video' && 'Video'}
              {lesson.content_type === 'quiz' && 'Quiz'}
            </Badge>
            {lesson.is_completed && (
              <Badge className="bg-green-500">
                <CheckCircle className="w-3 h-3 mr-1" />
                Completado
              </Badge>
            )}
          </div>
          <CardTitle className="text-2xl">{lesson.title}</CardTitle>
        </CardHeader>
        <CardContent>
          {lesson.content_type === 'quiz' && quiz_questions ? (
            // Quiz Content
            <div className="space-y-8">
              {quiz_questions.map((question, index) => (
                <div key={question.id} className="space-y-4">
                  <h3 className="font-medium text-lg">
                    {index + 1}. {question.question_text}
                  </h3>

                  <RadioGroup
                    value={quizAnswers[question.id] || ''}
                    onValueChange={(value) => handleQuizAnswerChange(question.id, value)}
                    disabled={quizSubmitted}
                  >
                    {['A', 'B', 'C', 'D'].map((option) => {
                      const optionText = question[`option_${option.toLowerCase()}` as keyof typeof question] as string
                      const isCorrect = question.correct_option === option
                      const isSelected = quizAnswers[question.id] === option

                      return (
                        <div
                          key={option}
                          className={cn(
                            'flex items-center space-x-3 p-3 rounded-lg border',
                            quizSubmitted && isCorrect && 'bg-green-50 border-green-300',
                            quizSubmitted && isSelected && !isCorrect && 'bg-red-50 border-red-300',
                            !quizSubmitted && 'hover:bg-gray-50'
                          )}
                        >
                          <RadioGroupItem value={option} id={`${question.id}-${option}`} />
                          <Label
                            htmlFor={`${question.id}-${option}`}
                            className="flex-1 cursor-pointer"
                          >
                            <span className="font-medium mr-2">{option}.</span>
                            {optionText}
                          </Label>
                          {quizSubmitted && isCorrect && (
                            <CheckCircle className="w-5 h-5 text-green-500" />
                          )}
                        </div>
                      )
                    })}
                  </RadioGroup>

                  {quizSubmitted && (
                    <div className="p-4 bg-blue-50 rounded-lg">
                      <p className="text-sm text-blue-800">
                        <strong>Explicacion:</strong> {question.explanation}
                      </p>
                    </div>
                  )}
                </div>
              ))}
            </div>
          ) : (
            // Text/Video Content
            <div className="prose prose-sm max-w-none">
              {lesson.content.split('\n').map((paragraph, index) => {
                // Handle markdown-like headers
                if (paragraph.startsWith('# ')) {
                  return <h1 key={index} className="text-2xl font-bold mt-6 mb-4">{paragraph.slice(2)}</h1>
                }
                if (paragraph.startsWith('## ')) {
                  return <h2 key={index} className="text-xl font-semibold mt-5 mb-3">{paragraph.slice(3)}</h2>
                }
                if (paragraph.startsWith('### ')) {
                  return <h3 key={index} className="text-lg font-medium mt-4 mb-2">{paragraph.slice(4)}</h3>
                }
                if (paragraph.startsWith('- ')) {
                  return <li key={index} className="ml-4">{paragraph.slice(2)}</li>
                }
                if (paragraph.trim() === '') {
                  return <br key={index} />
                }
                // Handle bold text
                const formattedParagraph = paragraph.replace(
                  /\*\*(.*?)\*\*/g,
                  '<strong>$1</strong>'
                )
                return (
                  <p
                    key={index}
                    className="mb-3 leading-relaxed"
                    dangerouslySetInnerHTML={{ __html: formattedParagraph }}
                  />
                )
              })}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Result Card */}
      {showResult && result && (
        <Card className="bg-gradient-to-r from-green-50 to-emerald-50 border-green-200">
          <CardContent className="py-6">
            <div className="flex items-center gap-4">
              <div className="w-16 h-16 rounded-full bg-green-500 flex items-center justify-center">
                <Trophy className="w-8 h-8 text-white" />
              </div>
              <div className="flex-1">
                <h3 className="text-xl font-bold text-green-800">
                  {result.course_completed ? 'Curso Completado!' : 'Leccion Completada!'}
                </h3>
                {result.quiz_score !== undefined && (
                  <p className="text-green-700">
                    Quiz: {result.quiz_score}/{result.quiz_total} respuestas correctas
                  </p>
                )}
                {result.points_earned > 0 && (
                  <p className="text-green-700 font-medium">
                    +{result.points_earned} puntos ganados!
                  </p>
                )}
                <p className="text-sm text-green-600">
                  Progreso: {result.completed_lessons}/{result.total_lessons} lecciones
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Actions */}
      <div className="flex justify-between">
        <div>
          {prev_lesson_id && (
            <Button
              variant="outline"
              onClick={() => router.push(`/education/lessons/${prev_lesson_id}`)}
            >
              <ArrowLeft className="w-4 h-4 mr-2" />
              Leccion anterior
            </Button>
          )}
        </div>

        <div className="flex gap-2">
          {!lesson.is_completed && !showResult && (
            <Button
              onClick={handleCompleteLesson}
              disabled={completing || (lesson.content_type === 'quiz' && !canSubmitQuiz())}
            >
              {completing ? (
                <Loader2 className="w-4 h-4 mr-2 animate-spin" />
              ) : (
                <CheckCircle className="w-4 h-4 mr-2" />
              )}
              {lesson.content_type === 'quiz' ? 'Enviar respuestas' : 'Marcar como completada'}
            </Button>
          )}

          {(lesson.is_completed || showResult) && next_lesson_id && (
            <Button onClick={() => router.push(`/education/lessons/${next_lesson_id}`)}>
              Siguiente leccion
              <ArrowRight className="w-4 h-4 ml-2" />
            </Button>
          )}

          {(lesson.is_completed || showResult) && !next_lesson_id && (
            <Button onClick={() => router.push(`/education/courses/${lesson.course_id}`)}>
              Volver al curso
              <ArrowRight className="w-4 h-4 ml-2" />
            </Button>
          )}
        </div>
      </div>
    </div>
  )
}

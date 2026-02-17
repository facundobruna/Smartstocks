import apiClient from './client'
import type { ApiResponse } from '@/types'

// Types
export interface Course {
  id: string
  title: string
  description: string
  icon: string
  category: 'fundamentos' | 'analisis' | 'estrategia' | 'avanzado'
  difficulty: 'principiante' | 'intermedio' | 'avanzado'
  duration_minutes: number
  points_reward: number
  is_premium: boolean
  is_active: boolean
  order_index: number
  total_lessons: number
  completed_lessons: number
  is_started: boolean
  is_completed: boolean
}

export interface Lesson {
  id: string
  course_id: string
  title: string
  content: string
  content_type: 'text' | 'video' | 'quiz'
  video_url?: string
  duration_minutes: number
  order_index: number
  is_completed: boolean
}

export interface QuizQuestion {
  id: string
  lesson_id: string
  question_text: string
  option_a: string
  option_b: string
  option_c: string
  option_d: string
  correct_option: string
  explanation: string
  order_index: number
}

export interface CoursesListResponse {
  courses: Course[]
  total_lessons: number
  completed_lessons: number
  overall_progress: number
}

export interface CourseDetailResponse {
  course: Course
  lessons: Lesson[]
}

export interface LessonDetailResponse {
  lesson: Lesson
  quiz_questions?: QuizQuestion[]
  next_lesson_id?: string
  prev_lesson_id?: string
}

export interface QuizAnswer {
  question_id: string
  answer: string
}

export interface CompleteLessonRequest {
  quiz_answers?: QuizAnswer[]
}

export interface CompleteLessonResponse {
  lesson_completed: boolean
  course_completed: boolean
  quiz_score?: number
  quiz_total?: number
  points_earned: number
  new_total_points: number
  completed_lessons: number
  total_lessons: number
}

export const coursesApi = {
  // Get all courses with progress
  getAllCourses: async (): Promise<CoursesListResponse> => {
    const response = await apiClient.get<ApiResponse<CoursesListResponse>>('/courses')
    return response.data.data!
  },

  // Get course details with lessons
  getCourseById: async (courseId: string): Promise<CourseDetailResponse> => {
    const response = await apiClient.get<ApiResponse<CourseDetailResponse>>(`/courses/${courseId}`)
    return response.data.data!
  },

  // Get lesson details
  getLessonById: async (lessonId: string): Promise<LessonDetailResponse> => {
    const response = await apiClient.get<ApiResponse<LessonDetailResponse>>(`/courses/lessons/${lessonId}`)
    return response.data.data!
  },

  // Complete a lesson
  completeLesson: async (lessonId: string, request?: CompleteLessonRequest): Promise<CompleteLessonResponse> => {
    const response = await apiClient.post<ApiResponse<CompleteLessonResponse>>(
      `/courses/lessons/${lessonId}/complete`,
      request || {}
    )
    return response.data.data!
  },
}

export default coursesApi

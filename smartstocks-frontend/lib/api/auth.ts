import apiClient from './client'
import type {
  ApiResponse,
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  User,
  UserStats,
  School
} from '@/types'

export const authApi = {
  login: async (data: LoginRequest): Promise<AuthResponse> => {
    const response = await apiClient.post<ApiResponse<AuthResponse>>('/auth/login', data)
    return response.data.data!
  },

  register: async (data: RegisterRequest): Promise<AuthResponse> => {
    const response = await apiClient.post<ApiResponse<AuthResponse>>('/auth/register', data)
    return response.data.data!
  },

  logout: async (refreshToken: string): Promise<void> => {
    await apiClient.post('/auth/logout', { refresh_token: refreshToken })
  },

  refreshToken: async (refreshToken: string): Promise<AuthResponse> => {
    const response = await apiClient.post<ApiResponse<AuthResponse>>('/auth/refresh', {
      refresh_token: refreshToken,
    })
    return response.data.data!
  },

  verifyEmail: async (token: string): Promise<void> => {
    await apiClient.post('/auth/verify-email', { token })
  },

  getSchools: async (): Promise<School[]> => {
    const response = await apiClient.get<ApiResponse<School[]>>('/schools')
    return response.data.data || []
  },
}

export const userApi = {
  getProfile: async (): Promise<{ user: User; stats: UserStats }> => {
    const response = await apiClient.get<ApiResponse<{ user: User; stats: UserStats }>>('/user/profile')
    return response.data.data!
  },

  updateProfile: async (data: Partial<User>): Promise<User> => {
    const response = await apiClient.put<ApiResponse<User>>('/user/profile', data)
    return response.data.data!
  },

  getStats: async (): Promise<UserStats> => {
    const response = await apiClient.get<ApiResponse<UserStats>>('/user/stats')
    return response.data.data!
  },
}

export default authApi

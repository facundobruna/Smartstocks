import apiClient from './client'
import type {
  ApiResponse,
  LeaderboardResponse,
  UserPositionResponse,
  UserProfilePublic,
  AllAchievementsResponse,
} from '@/types'

export const rankingsApi = {
  getGlobalLeaderboard: async (limit: number = 50, offset: number = 0): Promise<LeaderboardResponse> => {
    const response = await apiClient.get<ApiResponse<LeaderboardResponse>>(
      `/rankings/global?limit=${limit}&offset=${offset}`
    )
    return response.data.data!
  },

  getSchoolLeaderboard: async (schoolId: string, limit: number = 50): Promise<LeaderboardResponse> => {
    const response = await apiClient.get<ApiResponse<LeaderboardResponse>>(
      `/rankings/school/${schoolId}?limit=${limit}`
    )
    return response.data.data!
  },

  getMySchoolLeaderboard: async (limit: number = 50): Promise<LeaderboardResponse> => {
    const response = await apiClient.get<ApiResponse<LeaderboardResponse>>(
      `/rankings/my-school?limit=${limit}`
    )
    return response.data.data!
  },

  getMyPosition: async (): Promise<UserPositionResponse> => {
    const response = await apiClient.get<ApiResponse<UserPositionResponse>>(
      '/rankings/my-position'
    )
    return response.data.data!
  },

  getPublicProfile: async (userId: string): Promise<UserProfilePublic> => {
    const response = await apiClient.get<ApiResponse<UserProfilePublic>>(
      `/rankings/profile/${userId}`
    )
    return response.data.data!
  },

  getMyAchievements: async (): Promise<AllAchievementsResponse> => {
    const response = await apiClient.get<ApiResponse<AllAchievementsResponse>>(
      '/rankings/achievements'
    )
    return response.data.data!
  },
}

export default rankingsApi

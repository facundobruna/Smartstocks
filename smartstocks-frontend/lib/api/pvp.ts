import apiClient from './client'
import type {
  ApiResponse,
  PvPMatch,
  SimulatorDecision,
} from '@/types'

export interface PvPHistoryEntry {
  match_id: string
  opponent_username: string
  opponent_profile_picture?: string
  your_score: number
  opponent_score: number
  result: 'win' | 'loss' | 'draw'
  points_earned: number
  played_at: string
}

export interface PvPHistoryResponse {
  matches: PvPHistoryEntry[]
  total_matches: number
  wins: number
  losses: number
  draws: number
}

export const pvpApi = {
  joinQueue: async (): Promise<{ message: string; position?: number }> => {
    const response = await apiClient.post<ApiResponse<{ message: string; position?: number }>>(
      '/pvp/queue/join'
    )
    return response.data.data!
  },

  leaveQueue: async (): Promise<void> => {
    await apiClient.post('/pvp/queue/leave')
  },

  submitDecision: async (
    matchId: string,
    roundNumber: number,
    decision: SimulatorDecision,
    timeElapsed: number
  ): Promise<void> => {
    await apiClient.post('/pvp/submit', {
      match_id: matchId,
      round_number: roundNumber,
      decision,
      time_elapsed: timeElapsed,
    })
  },

  getHistory: async (limit: number = 20): Promise<PvPHistoryResponse> => {
    const response = await apiClient.get<ApiResponse<PvPHistoryResponse>>(
      `/pvp/history?limit=${limit}`
    )
    return response.data.data!
  },
}

export default pvpApi

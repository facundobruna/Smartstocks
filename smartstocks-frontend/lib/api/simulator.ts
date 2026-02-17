import apiClient from './client'
import type {
  ApiResponse,
  SimulatorScenario,
  SimulatorResult,
  SimulatorSubmitRequest,
  SimulatorHistoryResponse,
  CooldownStatus,
  SimulatorStats,
  SimulatorDifficulty,
} from '@/types'

export const simulatorApi = {
  getScenario: async (difficulty: SimulatorDifficulty): Promise<SimulatorScenario> => {
    const response = await apiClient.get<ApiResponse<SimulatorScenario>>(
      `/simulator/${difficulty}`
    )
    return response.data.data!
  },

  submitDecision: async (data: SimulatorSubmitRequest): Promise<SimulatorResult> => {
    const response = await apiClient.post<ApiResponse<SimulatorResult>>(
      '/simulator/submit',
      data
    )
    return response.data.data!
  },

  getHistory: async (limit: number = 20): Promise<SimulatorHistoryResponse> => {
    const response = await apiClient.get<ApiResponse<SimulatorHistoryResponse>>(
      `/simulator/history?limit=${limit}`
    )
    return response.data.data!
  },

  getCooldownStatus: async (difficulty: SimulatorDifficulty): Promise<CooldownStatus> => {
    const response = await apiClient.get<ApiResponse<CooldownStatus>>(
      `/simulator/cooldown/${difficulty}`
    )
    return response.data.data!
  },

  getStats: async (): Promise<SimulatorStats> => {
    const response = await apiClient.get<ApiResponse<SimulatorStats>>(
      '/simulator/stats'
    )
    return response.data.data!
  },
}

export default simulatorApi

import apiClient from './client'
import type {
  ApiResponse,
  Tournament,
  TournamentParticipant,
  TournamentBracketMatch,
} from '@/types'

export interface TournamentDetailsResponse extends Tournament {
  prizes: {
    position_from: number
    position_to: number
    token_reward: number
    special_reward?: string
  }[]
  is_registered: boolean
  my_position?: number
}

export interface TournamentStandingsResponse {
  tournament_id: string
  participants: TournamentParticipant[]
  total_participants: number
}

export interface TournamentBracketResponse {
  tournament_id: string
  rounds: {
    round_number: number
    matches: TournamentBracketMatch[]
  }[]
}

interface TournamentListResponse {
  tournaments: Tournament[]
  total: number
}

export const tournamentsApi = {
  getActiveTournaments: async (): Promise<Tournament[]> => {
    const response = await apiClient.get<ApiResponse<TournamentListResponse>>('/tournaments')
    return response.data.data?.tournaments || []
  },

  getMyTournaments: async (): Promise<Tournament[]> => {
    const response = await apiClient.get<ApiResponse<TournamentListResponse>>('/tournaments/my-tournaments')
    return response.data.data?.tournaments || []
  },

  getTournamentDetails: async (tournamentId: string): Promise<TournamentDetailsResponse> => {
    const response = await apiClient.get<ApiResponse<TournamentDetailsResponse>>(
      `/tournaments/${tournamentId}`
    )
    return response.data.data!
  },

  joinTournament: async (tournamentId: string): Promise<{ message: string }> => {
    const response = await apiClient.post<ApiResponse<{ message: string }>>(
      '/tournaments/join',
      { tournament_id: tournamentId }
    )
    return response.data.data!
  },

  getTournamentStandings: async (tournamentId: string): Promise<TournamentStandingsResponse> => {
    const response = await apiClient.get<ApiResponse<TournamentStandingsResponse>>(
      `/tournaments/${tournamentId}/standings`
    )
    return response.data.data!
  },

  getTournamentBracket: async (tournamentId: string): Promise<TournamentBracketResponse> => {
    const response = await apiClient.get<ApiResponse<TournamentBracketResponse>>(
      `/tournaments/${tournamentId}/bracket`
    )
    return response.data.data!
  },
}

export default tournamentsApi

import { create } from 'zustand'
import type {
  User,
  SimulatorScenario,
  SimulatorDecision,
} from '@/types'

export type PvPStatus =
  | 'idle'
  | 'queuing'
  | 'matched'
  | 'playing'
  | 'waiting_opponent'
  | 'round_result'
  | 'match_end'

export interface RoundResult {
  round_number: number
  your_decision: SimulatorDecision
  opponent_decision: SimulatorDecision
  correct_decision: SimulatorDecision
  your_points: number
  opponent_points: number
  explanation: string
}

export interface MatchResult {
  winner_id: string
  your_final_score: number
  opponent_final_score: number
  points_earned: number
  new_total_points: number
  new_rank_tier: string
}

interface PvPState {
  status: PvPStatus
  queuePosition: number | null
  matchId: string | null
  opponent: User | null
  currentRound: number
  totalRounds: number
  yourScore: number
  opponentScore: number
  scenario: SimulatorScenario | null
  timeLimit: number
  opponentDecided: boolean
  roundResult: RoundResult | null
  matchResult: MatchResult | null
  selectedDecision: SimulatorDecision | null

  // Actions
  setQueuing: (position?: number) => void
  updateQueuePosition: (position: number) => void
  setMatchFound: (data: {
    match_id: string
    opponent: User
    total_rounds: number
  }) => void
  setRoundStart: (data: {
    round_number: number
    scenario: SimulatorScenario
    time_limit_seconds: number
  }) => void
  setOpponentDecided: () => void
  setSelectedDecision: (decision: SimulatorDecision) => void
  setRoundResult: (result: RoundResult) => void
  setMatchEnd: (result: MatchResult) => void
  reset: () => void
}

const initialState = {
  status: 'idle' as PvPStatus,
  queuePosition: null,
  matchId: null,
  opponent: null,
  currentRound: 0,
  totalRounds: 5,
  yourScore: 0,
  opponentScore: 0,
  scenario: null,
  timeLimit: 15,
  opponentDecided: false,
  roundResult: null,
  matchResult: null,
  selectedDecision: null,
}

export const usePvPStore = create<PvPState>((set) => ({
  ...initialState,

  setQueuing: (position) => {
    set({
      status: 'queuing',
      queuePosition: position || null,
    })
  },

  updateQueuePosition: (position) => {
    set({ queuePosition: position })
  },

  setMatchFound: ({ match_id, opponent, total_rounds }) => {
    set({
      status: 'matched',
      matchId: match_id,
      opponent,
      totalRounds: total_rounds,
      yourScore: 0,
      opponentScore: 0,
    })
  },

  setRoundStart: ({ round_number, scenario, time_limit_seconds }) => {
    set({
      status: 'playing',
      currentRound: round_number,
      scenario,
      timeLimit: time_limit_seconds,
      opponentDecided: false,
      roundResult: null,
      selectedDecision: null,
    })
  },

  setOpponentDecided: () => {
    set({ opponentDecided: true })
  },

  setSelectedDecision: (decision) => {
    set({
      selectedDecision: decision,
      status: 'waiting_opponent',
    })
  },

  setRoundResult: (result) => {
    set((state) => ({
      status: 'round_result',
      roundResult: result,
      yourScore: state.yourScore + result.your_points,
      opponentScore: state.opponentScore + result.opponent_points,
    }))
  },

  setMatchEnd: (result) => {
    set({
      status: 'match_end',
      matchResult: result,
    })
  },

  reset: () => {
    set(initialState)
  },
}))

export default usePvPStore

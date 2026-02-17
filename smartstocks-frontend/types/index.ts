// ============================================
// API Response Wrapper
// ============================================
export interface ApiResponse<T> {
  success: boolean
  message?: string
  data?: T
  error?: string
}

// ============================================
// User Types
// ============================================
export interface User {
  id: string
  username: string
  email: string
  profile_picture_url?: string | null
  school_id?: string | null
  email_verified: boolean
  created_at: string
}

export interface UserStats {
  user_id: string
  smartpoints: number
  rank_tier: string
  total_quizzes_completed: number
  total_simulator_games: number
  win_streak: number
  total_wins: number
  total_losses: number
  updated_at: string
}

export interface School {
  id: string
  name: string
  location: string
  is_active: boolean
  created_at: string
}

// ============================================
// Auth Types
// ============================================
export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
  school_id?: string
  profile_picture_url?: string
}

export interface AuthResponse {
  access_token: string
  refresh_token: string
  user: User
  stats: UserStats
}

export interface RefreshTokenRequest {
  refresh_token: string
}

export interface UpdateProfileRequest {
  username?: string
  profile_picture_url?: string
  school_id?: string
}

// ============================================
// Simulator Types
// ============================================
export type SimulatorDifficulty = 'easy' | 'medium' | 'hard'
export type SimulatorDecision = 'buy' | 'sell' | 'hold'

export interface ChartData {
  labels: string[]
  prices: number[]
  full_prices?: number[]
  ticker: string
  asset_name: string
}

export interface SimulatorScenario {
  scenario_id: string
  difficulty: SimulatorDifficulty
  news_content: string
  chart_data: ChartData
  expires_at: string
}

export interface SimulatorSubmitRequest {
  scenario_id: string
  decision: SimulatorDecision
  time_taken_seconds?: number
}

export interface SimulatorResult {
  was_correct: boolean
  correct_decision: SimulatorDecision
  user_decision: SimulatorDecision
  points_earned: number
  explanation: string
  full_chart_data: ChartData
  new_total_points: number
  new_rank_tier: string
}

export interface SimulatorAttempt {
  id: string
  user_id: string
  scenario_id: string
  difficulty: SimulatorDifficulty
  user_decision: SimulatorDecision
  was_correct: boolean
  points_earned: number
  time_taken_seconds?: number
  created_at: string
  news_content?: string
  explanation?: string
}

export interface SimulatorStats {
  total_attempts: number
  correct_attempts: number
  accuracy_rate: number
  total_points: number
  by_difficulty: Record<string, {
    attempts: number
    correct: number
    accuracy_rate: number
    points_earned: number
  }>
}

export interface SimulatorHistoryResponse {
  attempts: SimulatorAttempt[]
  stats: SimulatorStats
}

export interface CooldownStatus {
  can_attempt: boolean
  last_attempt_date?: string
  next_available?: string
  hours_remaining?: number
}

// ============================================
// PvP Types
// ============================================
export interface PvPMatch {
  id: string
  player1_id: string
  player2_id: string
  player1_score: number
  player2_score: number
  status: string
  current_round: number
  total_rounds: number
}

export interface MatchFoundResponse {
  match_id: string
  opponent_id: string
  opponent: User
  total_rounds: number
  message: string
}

export interface RoundStartResponse {
  match_id: string
  round_number: number
  total_rounds: number
  scenario: SimulatorScenario
  time_limit_seconds: number
}

export interface RoundResultResponse {
  match_id: string
  round_number: number
  your_decision: SimulatorDecision
  opponent_decision: SimulatorDecision
  correct_decision: SimulatorDecision
  your_points: number
  opponent_points: number
  total_your_score: number
  total_opponent_score: number
  explanation: string
}

export interface MatchEndResponse {
  match_id: string
  winner_id: string
  your_final_score: number
  opponent_final_score: number
  points_earned: number
  new_total_points: number
  new_rank_tier: string
}

// ============================================
// Rankings Types
// ============================================
export interface LeaderboardEntry {
  rank_position: number
  user_id: string
  username: string
  smartpoints: number
  rank_tier: string
  total_wins: number
  total_losses: number
  win_rate: number
  profile_picture_url?: string
  school_name?: string
  school_id?: string
  is_current_user: boolean
}

export interface LeaderboardResponse {
  type: 'global' | 'school'
  top_players: LeaderboardEntry[]
  user_position?: LeaderboardEntry
  total_players: number
  last_updated: string
}

export interface UserPositionResponse {
  global_position: number
  school_position?: number
  total_players: number
}

export interface Achievement {
  id: string
  user_id: string
  achievement_type: string
  achievement_name: string
  achievement_description: string
  icon_url?: string
  unlocked_at: string
  is_unlocked: boolean
}

export interface AchievementProgress {
  achievement_type: string
  name: string
  description: string
  current: number
  required: number
  progress: number
  is_unlocked: boolean
}

export interface AllAchievementsResponse {
  unlocked: Achievement[]
  locked: AchievementProgress[]
  total_count: number
}

export interface UserProfilePublic {
  id: string
  username: string
  email: string
  profile_picture_url?: string
  school_id?: string
  email_verified: boolean
  created_at: string
  stats: UserStats
  achievements: Achievement[]
  global_rank: number
  school_rank?: number
}

// ============================================
// Tournament Types
// ============================================
export type TournamentType = 'weekly' | 'monthly' | 'special'
export type TournamentFormat = 'bracket' | 'league' | 'battle_royale'
export type TournamentStatus = 'upcoming' | 'registration' | 'in_progress' | 'completed' | 'cancelled'

export interface Tournament {
  id: string
  name: string
  description: string
  tournament_type: TournamentType
  format: TournamentFormat
  entry_fee: number
  prize_pool: number
  min_rank_required: string
  max_participants: number
  current_participants: number
  status: TournamentStatus
  start_time: string
  end_time: string
  registration_start: string
  registration_end: string
}

export interface TournamentPrize {
  position_from: number
  position_to: number
  token_reward: number
  special_reward?: string
}

export interface TournamentParticipant {
  user_id: string
  username: string
  profile_picture_url?: string
  rank_tier: string
  score: number
  position: number
  is_eliminated: boolean
}

export interface TournamentBracketMatch {
  match_id: string
  round: number
  position: number
  player1?: TournamentParticipant
  player2?: TournamentParticipant
  winner_id?: string
  status: string
}

// ============================================
// Tokens Types
// ============================================
export interface TokenBalance {
  balance: number
  total_earned: number
  total_spent: number
}

export interface TokenTransaction {
  id: string
  transaction_type: string
  amount: number
  balance_after: number
  description: string
  created_at: string
}

// ============================================
// Rank Tiers
// ============================================
export const RANK_TIERS = {
  'bronze_1': { name: 'Bronce 1', color: '#CD7F32', level: 1 },
  'bronze_2': { name: 'Bronce 2', color: '#CD7F32', level: 2 },
  'bronze_3': { name: 'Bronce 3', color: '#CD7F32', level: 3 },
  'silver_1': { name: 'Plata 1', color: '#C0C0C0', level: 4 },
  'silver_2': { name: 'Plata 2', color: '#C0C0C0', level: 5 },
  'silver_3': { name: 'Plata 3', color: '#C0C0C0', level: 6 },
  'gold_1': { name: 'Oro 1', color: '#FFD700', level: 7 },
  'gold_2': { name: 'Oro 2', color: '#FFD700', level: 8 },
  'gold_3': { name: 'Oro 3', color: '#FFD700', level: 9 },
  'master': { name: 'Maestro', color: '#9333EA', level: 10 },
} as const

export type RankTier = keyof typeof RANK_TIERS

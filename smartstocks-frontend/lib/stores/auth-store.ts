import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User, UserStats } from '@/types'

interface AuthState {
  user: User | null
  stats: UserStats | null
  accessToken: string | null
  refreshToken: string | null
  isAuthenticated: boolean
  isLoading: boolean

  // Actions
  setAuth: (data: {
    user: User
    stats: UserStats
    accessToken: string
    refreshToken: string
  }) => void
  updateUser: (user: Partial<User>) => void
  updateStats: (stats: Partial<UserStats>) => void
  logout: () => void
  setLoading: (loading: boolean) => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      stats: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,
      isLoading: true,

      setAuth: ({ user, stats, accessToken, refreshToken }) => {
        // Also store tokens in localStorage for the API client interceptor
        if (typeof window !== 'undefined') {
          localStorage.setItem('access_token', accessToken)
          localStorage.setItem('refresh_token', refreshToken)
        }

        set({
          user,
          stats,
          accessToken,
          refreshToken,
          isAuthenticated: true,
          isLoading: false,
        })
      },

      updateUser: (userData) => {
        set((state) => ({
          user: state.user ? { ...state.user, ...userData } : null,
        }))
      },

      updateStats: (statsData) => {
        set((state) => ({
          stats: state.stats ? { ...state.stats, ...statsData } : null,
        }))
      },

      logout: () => {
        if (typeof window !== 'undefined') {
          localStorage.removeItem('access_token')
          localStorage.removeItem('refresh_token')
        }

        set({
          user: null,
          stats: null,
          accessToken: null,
          refreshToken: null,
          isAuthenticated: false,
          isLoading: false,
        })
      },

      setLoading: (loading) => {
        set({ isLoading: loading })
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        user: state.user,
        stats: state.stats,
        accessToken: state.accessToken,
        refreshToken: state.refreshToken,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
)

export default useAuthStore

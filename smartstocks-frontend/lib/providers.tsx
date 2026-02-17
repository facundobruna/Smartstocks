'use client'

import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { useState, useEffect } from 'react'
import { Toaster } from 'sonner'
import { useAuthStore } from './stores/auth-store'

function AuthInitializer({ children }: { children: React.ReactNode }) {
  const { setLoading, isAuthenticated } = useAuthStore()

  useEffect(() => {
    // Check if we have stored auth on mount
    const accessToken = localStorage.getItem('access_token')
    if (!accessToken) {
      setLoading(false)
    } else {
      // Token exists, auth store will be rehydrated by zustand persist
      setLoading(false)
    }
  }, [setLoading])

  return <>{children}</>
}

export function Providers({ children }: { children: React.ReactNode }) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            staleTime: 60 * 1000, // 1 minute
            retry: 1,
            refetchOnWindowFocus: false,
          },
        },
      })
  )

  return (
    <QueryClientProvider client={queryClient}>
      <AuthInitializer>
        {children}
        <Toaster
          position="top-right"
          toastOptions={{
            style: {
              background: 'white',
              border: '1px solid #e5e7eb',
            },
          }}
        />
      </AuthInitializer>
    </QueryClientProvider>
  )
}

export default Providers

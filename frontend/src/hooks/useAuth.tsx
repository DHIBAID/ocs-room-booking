import { createContext, useContext, useEffect, useMemo, useState } from 'react'
import { apiFetch } from '../api/client'
import type { User } from '../api/types'

type AuthContextValue = {
  token: string | null
  user: User | null
  loading: boolean
  login: (username: string, password: string) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined)

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const [token, setToken] = useState<string | null>(() => {
    return localStorage.getItem('ocs_token')
  })
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!token) {
      setUser(null)
      setLoading(false)
      return
    }

    setLoading(true)
    apiFetch('/api/me', { token })
      .then((data) => setUser(data))
      .catch(() => {
        setToken(null)
        setUser(null)
        localStorage.removeItem('ocs_token')
      })
      .finally(() => setLoading(false))
  }, [token])

  const login = async (username: string, password: string) => {
    const data = await apiFetch('/api/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    })
    setToken(data.token)
    setUser(data.user)
    localStorage.setItem('ocs_token', data.token)
  }

  const logout = () => {
    setToken(null)
    setUser(null)
    localStorage.removeItem('ocs_token')
  }

  const value = useMemo(
    () => ({ token, user, loading, login, logout }),
    [token, user, loading],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export const useAuth = () => {
  const ctx = useContext(AuthContext)
  if (!ctx) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return ctx
}

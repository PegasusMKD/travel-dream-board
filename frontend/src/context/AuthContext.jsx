import { createContext, useContext, useState, useEffect } from 'react'
import { api, AuthError, setUnauthorizedHandler } from '../services/api'

const AuthContext = createContext(null)

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    setUnauthorizedHandler(() => setUser(null))

    api.auth.me()
      .then(setUser)
      .catch((err) => {
        if (!(err instanceof AuthError)) {
          console.error('Failed to check auth status:', err)
        }
        setUser(null)
      })
      .finally(() => setLoading(false))
  }, [])

  const logout = async () => {
    try {
      await api.auth.logout()
    } catch {
      // ignore errors during logout
    }
    setUser(null)
  }

  return (
    <AuthContext.Provider value={{ user, loading, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return ctx
}

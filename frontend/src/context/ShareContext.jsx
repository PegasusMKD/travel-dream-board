import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { setShareTokenProvider } from '../services/api'

const ShareContext = createContext(null)

const TOKEN_KEY = 'tdb_share_token'
const GUEST_UUID_KEY = 'tdb_guest_uuid'
const GUEST_NAME_KEY = 'tdb_guest_name'

function readToken() {
  try { return sessionStorage.getItem(TOKEN_KEY) || null } catch { return null }
}

function readGuestUuid() {
  try { return localStorage.getItem(GUEST_UUID_KEY) || null } catch { return null }
}

function readGuestName() {
  try { return localStorage.getItem(GUEST_NAME_KEY) || null } catch { return null }
}

export function ShareProvider({ children }) {
  const navigate = useNavigate()
  const location = useLocation()
  const [shareToken, setShareToken] = useState(readToken)
  const [guestUuid, setGuestUuid] = useState(readGuestUuid)
  const [guestName, setGuestName] = useState(readGuestName)

  // Capture ?token= from URL on first render and stash in sessionStorage,
  // then strip the query string from the visible URL.
  useEffect(() => {
    const params = new URLSearchParams(location.search)
    const urlToken = params.get('token')
    if (urlToken) {
      try { sessionStorage.setItem(TOKEN_KEY, urlToken) } catch {}
      setShareToken(urlToken)
      params.delete('token')
      const cleanSearch = params.toString()
      navigate(
        { pathname: location.pathname, search: cleanSearch ? `?${cleanSearch}` : '' },
        { replace: true }
      )
    }
  }, [location.pathname, location.search, navigate])

  // Provide the token to api.js for header injection.
  useEffect(() => {
    setShareTokenProvider(() => shareToken)
  }, [shareToken])

  const persistGuest = useCallback((uuid, name) => {
    try {
      localStorage.setItem(GUEST_UUID_KEY, uuid)
      localStorage.setItem(GUEST_NAME_KEY, name)
    } catch {}
    setGuestUuid(uuid)
    setGuestName(name)
  }, [])

  const clearShare = useCallback(() => {
    try { sessionStorage.removeItem(TOKEN_KEY) } catch {}
    setShareToken(null)
  }, [])

  const value = useMemo(() => ({
    shareToken,
    guestUuid,
    guestName,
    persistGuest,
    clearShare,
  }), [shareToken, guestUuid, guestName, persistGuest, clearShare])

  return (
    <ShareContext.Provider value={value}>
      {children}
    </ShareContext.Provider>
  )
}

export function useShare() {
  const ctx = useContext(ShareContext)
  if (!ctx) throw new Error('useShare must be used within a ShareProvider')
  return ctx
}

import { useEffect, useState } from 'react'
import { X, Copy, Check, Users, Loader2, Trash2 } from 'lucide-react'
import { useLang } from '../context/LanguageContext'
import { api } from '../services/api'

export default function ShareModal({ boardUuid, boardName, onClose }) {
  const { t } = useLang()
  const [copied, setCopied] = useState(false)
  const [token, setToken] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [revoking, setRevoking] = useState(false)
  const [revoked, setRevoked] = useState(false)

  useEffect(() => {
    let cancelled = false
    async function loadOrCreate() {
      try {
        const existing = await api.boards.shareTokens.list(boardUuid)
        if (cancelled) return
        if (existing && existing.length > 0) {
          setToken(existing[0].token)
        } else {
          const created = await api.boards.shareTokens.create(boardUuid)
          if (cancelled) return
          setToken(created.token)
        }
      } catch (err) {
        if (!cancelled) setError(err.message)
      } finally {
        if (!cancelled) setLoading(false)
      }
    }
    loadOrCreate()
    return () => { cancelled = true }
  }, [boardUuid])

  const shareUrl = token
    ? `${window.location.origin}/board/${boardUuid}?token=${token}`
    : ''

  const handleCopy = () => {
    if (!shareUrl) return
    navigator.clipboard.writeText(shareUrl)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  const handleRevoke = async () => {
    if (!token || revoking) return
    if (!window.confirm(t.confirmRevokeLink)) return
    setRevoking(true)
    setError(null)
    try {
      await api.boards.shareTokens.delete(boardUuid, token)
      setToken(null)
      setRevoked(true)
    } catch (err) {
      setError(err.message)
    } finally {
      setRevoking(false)
    }
  }

  const handleRegenerate = async () => {
    setLoading(true)
    setError(null)
    setRevoked(false)
    try {
      const created = await api.boards.shareTokens.create(boardUuid)
      setToken(created.token)
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-black/30 backdrop-blur-sm" onClick={onClose} />
      <div className="relative bg-surface-0 rounded-2xl shadow-2xl w-full max-w-md overflow-hidden border border-surface-200">
        <div className="px-6 pt-5 pb-3 flex items-center justify-between">
          <h3 className="text-base font-bold text-text-primary">{t.shareBoard}</h3>
          <button
            onClick={onClose}
            className="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-surface-100 transition-colors text-text-tertiary cursor-pointer"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <div className="px-6 pb-6">
          <p className="text-sm text-text-secondary mb-4">
            {t.shareDesc} <strong>"{boardName}"</strong>.
          </p>

          <div className="flex items-center gap-2 mb-3">
            <div className="flex-1 bg-surface-50 border border-surface-200 rounded-xl px-3 py-2.5 text-xs text-text-tertiary font-mono truncate min-h-[2.5rem] flex items-center">
              {loading ? (
                <Loader2 className="w-3.5 h-3.5 animate-spin" />
              ) : revoked && !token ? (
                <span className="italic">{t.linkRevoked}</span>
              ) : error ? (
                <span className="text-red-500">{error}</span>
              ) : (
                shareUrl
              )}
            </div>
            <button
              onClick={handleCopy}
              disabled={!shareUrl}
              className="shrink-0 w-10 h-10 flex items-center justify-center bg-accent-50 hover:bg-accent-100 text-accent-500 rounded-xl transition-colors cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {copied ? <Check className="w-4 h-4" /> : <Copy className="w-4 h-4" />}
            </button>
          </div>

          <div className="flex items-center justify-end gap-2 mb-4">
            {token && !loading && (
              <button
                onClick={handleRevoke}
                disabled={revoking}
                className="flex items-center gap-1.5 text-xs font-semibold text-red-500 hover:text-red-600 hover:bg-red-50 dark:hover:bg-red-950/30 px-2.5 py-1.5 rounded-lg transition-colors cursor-pointer disabled:opacity-50"
              >
                {revoking ? <Loader2 className="w-3.5 h-3.5 animate-spin" /> : <Trash2 className="w-3.5 h-3.5" />}
                {t.revokeLink}
              </button>
            )}
            {!token && !loading && (
              <button
                onClick={handleRegenerate}
                className="flex items-center gap-1.5 text-xs font-semibold text-accent-500 hover:text-accent-600 bg-accent-50 hover:bg-accent-100 px-3 py-1.5 rounded-lg transition-colors cursor-pointer"
              >
                {t.share}
              </button>
            )}
          </div>

          <div className="bg-surface-50 rounded-xl p-4">
            <div className="flex items-center gap-2 mb-2">
              <Users className="w-4 h-4 text-text-secondary" />
              <span className="text-xs font-semibold text-text-primary">{t.friendAccess}</span>
            </div>
            <ul className="text-xs text-text-secondary space-y-1 ml-6 list-disc">
              <li>{t.canVote}</li>
              <li>{t.canComment}</li>
              <li>{t.canView}</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  )
}

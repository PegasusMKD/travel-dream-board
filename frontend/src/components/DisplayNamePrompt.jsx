import { useState } from 'react'
import { Loader2, Users } from 'lucide-react'
import { useLang } from '../context/LanguageContext'
import { api } from '../services/api'

export default function DisplayNamePrompt({ onSubmitted }) {
  const { t } = useLang()
  const [name, setName] = useState('')
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState(null)

  const handleSubmit = async (e) => {
    e.preventDefault()
    const trimmed = name.trim()
    if (!trimmed) return
    setBusy(true)
    setError(null)
    try {
      const guest = await api.guests.create(trimmed)
      onSubmitted({ uuid: guest.uuid, name: guest.name })
    } catch (err) {
      setError(err.message || t.guestNameError)
      setBusy(false)
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-black/40 backdrop-blur-sm" />
      <div className="relative bg-surface-0 rounded-2xl shadow-2xl w-full max-w-sm overflow-hidden border border-surface-200">
        <div className="px-6 pt-6 pb-2">
          <div className="flex items-center gap-2 mb-3">
            <div className="w-9 h-9 rounded-xl bg-accent-100 flex items-center justify-center">
              <Users className="w-5 h-5 text-accent-600" />
            </div>
            <h3 className="text-base font-bold text-text-primary">{t.guestNameTitle}</h3>
          </div>
          <p className="text-sm text-text-secondary mb-4">{t.guestNameDesc}</p>
        </div>

        <form onSubmit={handleSubmit} className="px-6 pb-6 space-y-3">
          <input
            type="text"
            autoFocus
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder={t.guestNamePlaceholder}
            maxLength={60}
            disabled={busy}
            className="w-full text-sm text-text-primary bg-surface-50 border border-surface-200 rounded-xl px-3 py-2.5 focus:outline-none focus:border-accent-400 focus:ring-2 focus:ring-accent-100 placeholder:text-text-muted transition-colors disabled:opacity-50"
          />
          {error && (
            <p className="text-xs text-red-500 bg-red-50 dark:bg-red-950/30 rounded-lg p-2.5">{error}</p>
          )}
          <button
            type="submit"
            disabled={busy || !name.trim()}
            className="w-full flex items-center justify-center gap-2 bg-accent-500 hover:bg-accent-600 text-white text-sm font-semibold px-3 py-2.5 rounded-xl transition-colors cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {busy ? <Loader2 className="w-4 h-4 animate-spin" /> : null}
            {t.guestNameContinue}
          </button>
        </form>
      </div>
    </div>
  )
}

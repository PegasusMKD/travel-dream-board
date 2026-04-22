import { useState } from 'react'
import { X, Link as LinkIcon, Loader2 } from 'lucide-react'
import { useLang } from '../context/LanguageContext'
import { api } from '../services/api'
import { sectionToItemApi } from '../services/mappers'

export default function AddItemModal({ sectionType, boardUuid, onClose, onAdded }) {
  const { t } = useLang()
  const [url, setUrl] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)

  const addLabels = {
    accommodation: t.addAccommodation,
    transport: t.addTransport,
    activities: t.addActivity,
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError(null)

    try {
      const itemApi = sectionToItemApi(api, sectionType)
      await itemApi.create(url, boardUuid)
      onAdded?.()
      onClose()
    } catch (err) {
      setError(err.message)
      setLoading(false)
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-black/30 backdrop-blur-sm" onClick={onClose} />
      <div className="relative bg-surface-0 rounded-2xl shadow-2xl w-full max-w-md overflow-hidden border border-surface-200">
        <div className="flex items-center justify-between px-6 pt-5 pb-3">
          <h3 className="text-base font-bold text-text-primary">
            {addLabels[sectionType]}
          </h3>
          <button
            onClick={onClose}
            className="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-surface-100 transition-colors text-text-tertiary cursor-pointer"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="px-6 pb-6">
          <div className="relative mb-4">
            <LinkIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-text-muted" />
            <input
              type="url"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              placeholder={t.pasteLink}
              className="w-full pl-10 pr-4 py-3 bg-surface-50 border border-surface-200 rounded-xl text-sm focus:outline-none focus:border-accent-400 focus:ring-2 focus:ring-accent-100 placeholder:text-text-muted text-text-primary transition-colors"
              autoFocus
              required
            />
          </div>

          <p className="text-xs text-text-tertiary mb-4 leading-relaxed">
            {t.addLinkHelp}
          </p>

          {error && <p className="text-xs text-red-500 mb-3">{error}</p>}

          <button
            type="submit"
            disabled={!url || loading}
            className="w-full bg-accent-500 hover:bg-accent-600 text-white py-3 rounded-xl text-sm font-semibold transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 cursor-pointer"
          >
            {loading ? (
              <>
                <Loader2 className="w-4 h-4 animate-spin" />
                {t.fetching}
              </>
            ) : (
              t.addLinkBtn
            )}
          </button>
        </form>
      </div>
    </div>
  )
}

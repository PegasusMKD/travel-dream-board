import { useState } from 'react'
import { X, Loader2 } from 'lucide-react'
import { useLang } from '../context/LanguageContext'

export default function EditBoardModal({ board, onClose }) {
  const { t } = useLang()
  const [name, setName] = useState(board.name)
  const [destination, setDestination] = useState(board.destination)
  const [startDate, setStartDate] = useState(board.dateRange?.start || '')
  const [endDate, setEndDate] = useState(board.dateRange?.end || '')
  const [coverImage, setCoverImage] = useState(board.coverImage)
  const [saving, setSaving] = useState(false)

  const handleSubmit = (e) => {
    e.preventDefault()
    setSaving(true)
    // TODO: connect to backend API
    setTimeout(() => {
      setSaving(false)
      onClose()
    }, 1500)
  }

  const inputClass =
    'w-full px-4 py-3 bg-surface-50 border border-surface-200 rounded-xl text-sm focus:outline-none focus:border-accent-400 focus:ring-2 focus:ring-accent-100 placeholder:text-text-muted text-text-primary transition-colors'

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-black/30 backdrop-blur-sm" onClick={onClose} />
      <div className="relative bg-surface-0 rounded-2xl shadow-2xl w-full max-w-md overflow-hidden border border-surface-200">
        <div className="flex items-center justify-between px-6 pt-5 pb-3">
          <h3 className="text-base font-bold text-text-primary">{t.editBoard}</h3>
          <button
            onClick={onClose}
            className="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-surface-100 transition-colors text-text-tertiary cursor-pointer"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="px-6 pb-6 space-y-4">
          {/* Name */}
          <div>
            <label className="block text-xs font-semibold text-text-secondary mb-1.5">
              {t.editBoardName}
            </label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder={t.editBoardNamePlaceholder}
              className={inputClass}
              autoFocus
            />
          </div>

          {/* Destination */}
          <div>
            <label className="block text-xs font-semibold text-text-secondary mb-1.5">
              {t.editBoardDestination}
            </label>
            <input
              type="text"
              value={destination}
              onChange={(e) => setDestination(e.target.value)}
              placeholder={t.editBoardDestinationPlaceholder}
              className={inputClass}
            />
          </div>

          {/* Date range */}
          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="block text-xs font-semibold text-text-secondary mb-1.5">
                {t.editBoardStartDate}
              </label>
              <input
                type="date"
                value={startDate}
                onChange={(e) => setStartDate(e.target.value)}
                className={inputClass}
              />
            </div>
            <div>
              <label className="block text-xs font-semibold text-text-secondary mb-1.5">
                {t.editBoardEndDate}
              </label>
              <input
                type="date"
                value={endDate}
                onChange={(e) => setEndDate(e.target.value)}
                className={inputClass}
              />
            </div>
          </div>

          {/* Cover image */}
          <div>
            <label className="block text-xs font-semibold text-text-secondary mb-1.5">
              {t.editBoardCoverImage}
            </label>
            <input
              type="url"
              value={coverImage}
              onChange={(e) => setCoverImage(e.target.value)}
              placeholder={t.editBoardCoverImagePlaceholder}
              className={inputClass}
            />
            {coverImage && (
              <img
                src={coverImage}
                alt=""
                className="mt-2 w-full h-32 object-cover rounded-xl border border-surface-200"
              />
            )}
          </div>

          {/* Submit */}
          <button
            type="submit"
            disabled={!name || saving}
            className="w-full bg-accent-500 hover:bg-accent-600 text-white py-3 rounded-xl text-sm font-semibold transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 cursor-pointer"
          >
            {saving ? (
              <>
                <Loader2 className="w-4 h-4 animate-spin" />
                {t.saving}
              </>
            ) : (
              t.save
            )}
          </button>
        </form>
      </div>
    </div>
  )
}

import { useState } from 'react'
import {
  X,
  ExternalLink,
  ThumbsUp,
  ThumbsDown,
  Bookmark,
  BookmarkX,
  Bed,
  Plane,
  MapPinned,
  Pencil,
  Save,
  XCircle,
  ImagePlus,
  Trash2,
} from 'lucide-react'
import StatusBadge from './StatusBadge'
import { useLang } from '../context/LanguageContext'

const sectionIcons = {
  accommodation: Bed,
  transport: Plane,
  activities: MapPinned,
}

const ALL_STATUSES = ['considering', 'finalist', 'rejected', 'booked', 'completed']

export default function ItemDetailSidebar({ item, sectionType, onClose }) {
  const { lang, t } = useLang()
  const [editing, setEditing] = useState(false)

  // Edit state — local draft
  const [draft, setDraft] = useState({
    title: item.title,
    note: item.note || '',
    url: item.url,
    image: item.image || '',
    status: item.status,
    isFinal: item.isFinal,
    bookingRef: item.bookingRef || '',
  })

  if (!item) return null

  const FallbackIcon = sectionIcons[sectionType] || MapPinned
  const upVotes = item.votes.filter((v) => v.value === 'up')
  const downVotes = item.votes.filter((v) => v.value === 'down')
  const locale = lang === 'pl' ? 'pl-PL' : 'en-US'

  const startEditing = () => {
    setDraft({
      title: item.title,
      note: item.note || '',
      url: item.url,
      image: item.image || '',
      status: item.status,
      isFinal: item.isFinal,
      bookingRef: item.bookingRef || '',
    })
    setEditing(true)
  }

  const cancelEditing = () => setEditing(false)

  const handleSave = () => {
    // In a real app this would call an API — for now just close edit mode
    setEditing(false)
  }

  const updateDraft = (field, value) => setDraft((prev) => ({ ...prev, [field]: value }))

  // Determine which image/title to show (draft when editing, item otherwise)
  const displayImage = editing ? draft.image : item.image
  const displayTitle = editing ? draft.title : item.title

  return (
    <>
      {/* Backdrop */}
      <div
        className="fixed inset-0 bg-black/30 backdrop-blur-sm z-40"
        onClick={onClose}
      />

      {/* Sidebar */}
      <div className="fixed top-0 right-0 h-full w-full max-w-md bg-surface-0 z-50 shadow-2xl flex flex-col animate-slide-in">
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b border-surface-200">
          <h2 className="text-base font-bold text-text-primary">{t.details}</h2>
          <div className="flex items-center gap-1">
            {!editing ? (
              <button
                onClick={startEditing}
                className="flex items-center gap-1.5 text-xs font-semibold text-accent-500 hover:text-accent-600 bg-accent-50 hover:bg-accent-100 px-3 py-1.5 rounded-lg transition-colors cursor-pointer"
              >
                <Pencil className="w-3.5 h-3.5" />
                {t.edit}
              </button>
            ) : (
              <>
                <button
                  onClick={cancelEditing}
                  className="flex items-center gap-1.5 text-xs font-semibold text-text-tertiary hover:text-text-secondary bg-surface-100 hover:bg-surface-200 px-3 py-1.5 rounded-lg transition-colors cursor-pointer"
                >
                  <XCircle className="w-3.5 h-3.5" />
                  {t.cancel}
                </button>
                <button
                  onClick={handleSave}
                  className="flex items-center gap-1.5 text-xs font-semibold text-white bg-accent-500 hover:bg-accent-600 px-3 py-1.5 rounded-lg transition-colors cursor-pointer"
                >
                  <Save className="w-3.5 h-3.5" />
                  {t.save}
                </button>
              </>
            )}
            <button
              onClick={onClose}
              className="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-surface-100 transition-colors text-text-tertiary cursor-pointer ml-1"
            >
              <X className="w-5 h-5" />
            </button>
          </div>
        </div>

        {/* Scrollable content */}
        <div className="flex-1 overflow-y-auto">
          {/* Image area */}
          <div className="relative group">
            {displayImage ? (
              <>
                <img src={displayImage} alt={displayTitle} className="w-full h-52 object-cover" />
                {editing && (
                  <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-3">
                    <button
                      onClick={() => {
                        const url = prompt(t.editImagePlaceholder)
                        if (url) updateDraft('image', url)
                      }}
                      className="flex items-center gap-1.5 bg-white text-gray-800 text-xs font-semibold px-3 py-2 rounded-lg cursor-pointer"
                    >
                      <ImagePlus className="w-3.5 h-3.5" />
                      {t.changeImage}
                    </button>
                    <button
                      onClick={() => updateDraft('image', '')}
                      className="flex items-center gap-1.5 bg-red-500 text-white text-xs font-semibold px-3 py-2 rounded-lg cursor-pointer"
                    >
                      <Trash2 className="w-3.5 h-3.5" />
                      {t.removeImage}
                    </button>
                  </div>
                )}
              </>
            ) : (
              <div className="w-full h-36 bg-surface-100 flex flex-col items-center justify-center gap-2">
                <FallbackIcon className="w-12 h-12 text-surface-300" />
                {editing && (
                  <button
                    onClick={() => {
                      const url = prompt(t.editImagePlaceholder)
                      if (url) updateDraft('image', url)
                    }}
                    className="flex items-center gap-1.5 text-xs font-semibold text-accent-500 bg-accent-50 hover:bg-accent-100 px-3 py-1.5 rounded-lg transition-colors cursor-pointer"
                  >
                    <ImagePlus className="w-3.5 h-3.5" />
                    {t.changeImage}
                  </button>
                )}
              </div>
            )}
          </div>

          <div className="p-6 space-y-5">
            {/* Final badge + Title */}
            <div>
              {editing ? (
                <>
                  {/* Final toggle */}
                  <button
                    onClick={() => updateDraft('isFinal', !draft.isFinal)}
                    className={`inline-flex items-center gap-1.5 text-[11px] font-bold px-2.5 py-1 rounded-lg mb-3 cursor-pointer transition-colors ${
                      draft.isFinal
                        ? 'bg-accent-500 text-white'
                        : 'bg-surface-100 text-text-tertiary hover:bg-surface-200'
                    }`}
                  >
                    {draft.isFinal ? (
                      <Bookmark className="w-3 h-3" fill="currentColor" />
                    ) : (
                      <BookmarkX className="w-3 h-3" />
                    )}
                    {draft.isFinal ? t.unmarkAsFinal : t.markAsFinal}
                  </button>

                  <label className="text-xs font-semibold text-text-tertiary uppercase tracking-wider mb-1.5 block">
                    {t.editTitle}
                  </label>
                  <input
                    type="text"
                    value={draft.title}
                    onChange={(e) => updateDraft('title', e.target.value)}
                    className="w-full text-lg font-bold text-text-primary bg-surface-50 border border-surface-200 rounded-xl px-3 py-2.5 focus:outline-none focus:border-accent-400 focus:ring-2 focus:ring-accent-100 transition-colors"
                  />
                </>
              ) : (
                <>
                  {item.isFinal && (
                    <span className="inline-flex items-center gap-1 bg-accent-500 text-white text-[11px] font-bold px-2.5 py-1 rounded-lg mb-2">
                      <Bookmark className="w-3 h-3" fill="currentColor" />
                      {t.selected}
                    </span>
                  )}
                  <h3 className="text-lg font-bold text-text-primary leading-snug">
                    {item.title}
                  </h3>
                </>
              )}
            </div>

            {/* Status */}
            <div>
              <label className="text-xs font-semibold text-text-tertiary uppercase tracking-wider mb-1.5 block">
                {t.status}
              </label>
              {editing ? (
                <div className="flex flex-wrap gap-2">
                  {ALL_STATUSES.map((s) => (
                    <button
                      key={s}
                      onClick={() => updateDraft('status', s)}
                      className={`cursor-pointer rounded-full transition-all ${
                        draft.status === s
                          ? 'ring-2 ring-accent-400 ring-offset-2 ring-offset-surface-0'
                          : 'opacity-60 hover:opacity-100'
                      }`}
                    >
                      <StatusBadge status={s} />
                    </button>
                  ))}
                </div>
              ) : (
                <StatusBadge status={item.status} />
              )}
            </div>

            {/* Booking ref */}
            {editing ? (
              <div>
                <label className="text-xs font-semibold text-text-tertiary uppercase tracking-wider mb-1.5 block">
                  {t.editBookingRef}
                </label>
                <input
                  type="text"
                  value={draft.bookingRef}
                  onChange={(e) => updateDraft('bookingRef', e.target.value)}
                  placeholder={t.editBookingRefPlaceholder}
                  className="w-full text-sm font-mono text-text-primary bg-surface-50 border border-surface-200 rounded-xl px-3 py-2.5 focus:outline-none focus:border-accent-400 focus:ring-2 focus:ring-accent-100 placeholder:text-text-muted transition-colors"
                />
              </div>
            ) : (
              item.bookingRef && (
                <div className="bg-emerald-50 dark:bg-emerald-950/30 text-emerald-700 dark:text-emerald-400 rounded-xl p-3 text-sm font-mono">
                  {t.bookingRef}: <strong>{item.bookingRef}</strong>
                </div>
              )
            )}

            {/* Notes */}
            <div>
              <label className="text-xs font-semibold text-text-tertiary uppercase tracking-wider mb-1.5 block">
                {t.notes}
              </label>
              {editing ? (
                <textarea
                  value={draft.note}
                  onChange={(e) => updateDraft('note', e.target.value)}
                  placeholder={t.editNotesPlaceholder}
                  rows={3}
                  className="w-full text-sm text-text-primary bg-surface-50 border border-surface-200 rounded-xl px-3 py-2.5 focus:outline-none focus:border-accent-400 focus:ring-2 focus:ring-accent-100 placeholder:text-text-muted resize-y leading-relaxed transition-colors"
                />
              ) : item.note ? (
                <p className="text-sm text-text-secondary bg-surface-50 rounded-xl p-3 leading-relaxed">
                  {item.note}
                </p>
              ) : (
                <p className="text-sm text-text-muted italic">{t.noNotes}</p>
              )}
            </div>

            {/* URL */}
            {editing ? (
              <div>
                <label className="text-xs font-semibold text-text-tertiary uppercase tracking-wider mb-1.5 block">
                  {t.editUrl}
                </label>
                <input
                  type="url"
                  value={draft.url}
                  onChange={(e) => updateDraft('url', e.target.value)}
                  className="w-full text-sm text-text-primary bg-surface-50 border border-surface-200 rounded-xl px-3 py-2.5 focus:outline-none focus:border-accent-400 focus:ring-2 focus:ring-accent-100 placeholder:text-text-muted transition-colors"
                />
              </div>
            ) : (
              <a
                href={item.url}
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-2 text-sm font-semibold text-accent-500 hover:text-accent-600 transition-colors no-underline"
              >
                <ExternalLink className="w-4 h-4" />
                {t.openLink}
              </a>
            )}

            {/* Image URL (edit mode only) */}
            {editing && (
              <div>
                <label className="text-xs font-semibold text-text-tertiary uppercase tracking-wider mb-1.5 block">
                  {t.editImageUrl}
                </label>
                <input
                  type="url"
                  value={draft.image}
                  onChange={(e) => updateDraft('image', e.target.value)}
                  placeholder={t.editImagePlaceholder}
                  className="w-full text-sm text-text-primary bg-surface-50 border border-surface-200 rounded-xl px-3 py-2.5 focus:outline-none focus:border-accent-400 focus:ring-2 focus:ring-accent-100 placeholder:text-text-muted transition-colors"
                />
              </div>
            )}

            {/* Divider between editable fields and social section */}
            <div className="border-t border-surface-200" />

            {/* Votes */}
            <div>
              <label className="text-xs font-semibold text-text-tertiary uppercase tracking-wider mb-2 block">
                {t.votes}
              </label>
              <div className="space-y-2">
                {upVotes.length > 0 && (
                  <div className="flex flex-wrap items-center gap-2">
                    <ThumbsUp className="w-4 h-4 text-accent-500" />
                    {upVotes.map((v, i) => (
                      <span
                        key={i}
                        className="inline-flex items-center gap-1.5 bg-accent-50 text-accent-600 text-xs font-semibold px-2.5 py-1 rounded-full"
                      >
                        <span className="w-5 h-5 rounded-full bg-accent-200 flex items-center justify-center text-[10px] font-bold text-accent-700">
                          {v.displayName[0]}
                        </span>
                        {v.displayName}
                      </span>
                    ))}
                  </div>
                )}
                {downVotes.length > 0 && (
                  <div className="flex flex-wrap items-center gap-2">
                    <ThumbsDown className="w-4 h-4 text-text-muted" />
                    {downVotes.map((v, i) => (
                      <span
                        key={i}
                        className="inline-flex items-center gap-1.5 bg-surface-100 text-text-secondary text-xs font-semibold px-2.5 py-1 rounded-full"
                      >
                        <span className="w-5 h-5 rounded-full bg-surface-200 flex items-center justify-center text-[10px] font-bold text-text-tertiary">
                          {v.displayName[0]}
                        </span>
                        {v.displayName}
                      </span>
                    ))}
                  </div>
                )}
                {upVotes.length === 0 && downVotes.length === 0 && (
                  <p className="text-sm text-text-muted italic">—</p>
                )}
              </div>
            </div>

            {/* Comments */}
            <div>
              <label className="text-xs font-semibold text-text-tertiary uppercase tracking-wider mb-2 block">
                {t.comments} ({item.comments.length})
              </label>
              {item.comments.length > 0 ? (
                <div className="space-y-3">
                  {item.comments.map((c, i) => (
                    <div key={i} className="bg-surface-50 rounded-xl p-3">
                      <div className="flex items-center gap-2 mb-1.5">
                        <span className="w-6 h-6 rounded-full bg-accent-100 flex items-center justify-center text-[10px] font-bold text-accent-600">
                          {c.displayName[0]}
                        </span>
                        <span className="text-xs font-semibold text-text-primary">
                          {c.displayName}
                        </span>
                        <span className="text-[11px] text-text-muted ml-auto">
                          {new Date(c.createdAt).toLocaleDateString(locale, {
                            day: 'numeric',
                            month: 'short',
                            hour: '2-digit',
                            minute: '2-digit',
                          })}
                        </span>
                      </div>
                      <p className="text-sm text-text-secondary leading-relaxed pl-8">
                        {c.text}
                      </p>
                    </div>
                  ))}
                </div>
              ) : (
                <p className="text-sm text-text-muted italic">{t.noComments}</p>
              )}

              <div className="flex items-center gap-2 mt-3">
                <input
                  type="text"
                  placeholder={t.addComment}
                  className="flex-1 text-sm bg-surface-50 border border-surface-200 rounded-xl px-3 py-2.5 focus:outline-none focus:border-accent-400 focus:ring-2 focus:ring-accent-100 placeholder:text-text-muted text-text-primary transition-colors"
                />
              </div>
            </div>
          </div>
        </div>
      </div>

      <style>{`
        @keyframes slideIn {
          from { transform: translateX(100%); }
          to { transform: translateX(0); }
        }
        .animate-slide-in {
          animation: slideIn 0.25s ease-out;
        }
      `}</style>
    </>
  )
}

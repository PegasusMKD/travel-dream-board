import {
  ThumbsUp,
  ThumbsDown,
  MessageCircle,
  ExternalLink,
  Bookmark,
  Bed,
  Plane,
  MapPinned,
  Loader2,
} from 'lucide-react'
import StatusBadge from './StatusBadge'
import { useLang } from '../context/LanguageContext'
import { useAuth } from '../context/AuthContext'

const sectionIcons = {
  accommodation: Bed,
  transport: Plane,
  activities: MapPinned,
}

export default function ItemCard({ item, sectionType, onClick, onVote }) {
  const { t } = useLang()
  const { user } = useAuth()
  const upVotes = item.likes ?? 0
  const downVotes = item.dislikes ?? 0
  const myVote = user ? item.votes.find((v) => v.userUuid === user.uuid) : null
  const FallbackIcon = sectionIcons[sectionType] || MapPinned
  const isPending = !!item.pending

  const handleVote = (e, direction) => {
    e.stopPropagation()
    if (isPending) return
    onVote?.(item, direction)
  }

  const handleClick = () => {
    if (isPending) return
    onClick?.()
  }

  return (
    <div
      onClick={handleClick}
      className={`relative bg-surface-0 rounded-2xl overflow-hidden shadow-sm border transition-all duration-200 ${
        isPending
          ? 'border-surface-200 cursor-wait'
          : item.isFinal
            ? 'border-accent-400 ring-2 ring-accent-100 hover:shadow-md cursor-pointer'
            : item.status === 'rejected'
              ? 'border-surface-200 opacity-50 hover:shadow-md cursor-pointer'
              : 'border-surface-200 hover:border-accent-300 hover:shadow-md cursor-pointer'
      }`}
    >
      {isPending && (
        <div className="absolute inset-0 bg-surface-0/60 backdrop-blur-[1px] z-10 flex flex-col items-center justify-center gap-2">
          <Loader2 className="w-6 h-6 text-accent-500 animate-spin" />
          <span className="text-xs font-semibold text-text-secondary">{t.fetching}</span>
        </div>
      )}
      {/* Image */}
      <div className="relative">
        {item.image ? (
          <img src={item.image} alt={item.title} className="w-full h-40 object-cover" />
        ) : (
          <div className="w-full h-28 bg-surface-100 flex items-center justify-center">
            <FallbackIcon className="w-10 h-10 text-surface-300" />
          </div>
        )}
        {item.isFinal && (
          <div className="absolute top-2 left-2 bg-accent-500 text-white text-[11px] font-bold px-2.5 py-1 rounded-lg flex items-center gap-1 shadow-sm">
            <Bookmark className="w-3 h-3" fill="currentColor" />
            {t.selected}
          </div>
        )}
      </div>

      {/* Content */}
      <div className="p-4">
        <div className="flex items-start justify-between gap-2 mb-2">
          <h4 className="text-sm font-semibold text-text-primary leading-snug line-clamp-2">
            {item.title}
          </h4>
          <a
            href={item.url}
            target="_blank"
            rel="noopener noreferrer"
            className="text-text-muted hover:text-accent-500 transition-colors shrink-0 mt-0.5"
            onClick={(e) => e.stopPropagation()}
          >
            <ExternalLink className="w-4 h-4" />
          </a>
        </div>

        <StatusBadge status={item.status} />

        {item.note && (
          <p className="mt-2.5 text-xs text-text-secondary bg-surface-50 rounded-lg p-2.5 leading-relaxed">
            {item.note}
          </p>
        )}

        {item.bookingRef && (
          <div className="mt-2 text-xs bg-emerald-50 dark:bg-emerald-950/30 text-emerald-700 dark:text-emerald-400 rounded-lg p-2 font-mono">
            {t.bookingRef}: {item.bookingRef}
          </div>
        )}

        {/* Footer */}
        <div className="flex items-center justify-between mt-3 pt-3 border-t border-surface-200">
          <div className="flex items-center gap-3">
            <button
              onClick={(e) => handleVote(e, 'up')}
              className={`flex items-center gap-1 text-xs transition-colors cursor-pointer ${
                myVote?.value === 'up'
                  ? 'text-accent-500'
                  : 'text-text-tertiary hover:text-accent-500'
              }`}
            >
              <ThumbsUp className="w-3.5 h-3.5" fill={myVote?.value === 'up' ? 'currentColor' : 'none'} />
              <span className="font-semibold">{upVotes}</span>
            </button>
            <button
              onClick={(e) => handleVote(e, 'down')}
              className={`flex items-center gap-1 text-xs transition-colors cursor-pointer ${
                myVote?.value === 'down'
                  ? 'text-text-primary'
                  : 'text-text-tertiary hover:text-text-secondary'
              }`}
            >
              <ThumbsDown className="w-3.5 h-3.5" fill={myVote?.value === 'down' ? 'currentColor' : 'none'} />
              <span className="font-semibold">{downVotes}</span>
            </button>
          </div>

          {item.comments.length > 0 && (
            <span className="flex items-center gap-1 text-xs text-text-tertiary">
              <MessageCircle className="w-3.5 h-3.5" />
              {item.comments.length}
            </span>
          )}
        </div>

        {/* Vote avatars */}
        {item.votes.length > 0 && (
          <div className="flex items-center gap-1 mt-2">
            {item.votes
              .filter((v) => v.value === 'up')
              .map((v) => (
                <span
                  key={v.id}
                  className="w-6 h-6 rounded-full bg-accent-100 flex items-center justify-center text-[10px] font-bold text-accent-600"
                  title={v.displayName}
                >
                  {v.displayName[0]}
                </span>
              ))}
          </div>
        )}
      </div>
    </div>
  )
}

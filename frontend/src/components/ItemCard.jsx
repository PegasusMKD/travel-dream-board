import {
  ThumbsUp,
  ThumbsDown,
  MessageCircle,
  ExternalLink,
  Bookmark,
  Bed,
  Plane,
  MapPinned,
} from 'lucide-react'
import StatusBadge from './StatusBadge'
import { useLang } from '../context/LanguageContext'

const sectionIcons = {
  accommodation: Bed,
  transport: Plane,
  activities: MapPinned,
}

export default function ItemCard({ item, sectionType, onClick }) {
  const { t } = useLang()
  const upVotes = item.votes.filter((v) => v.value === 'up').length
  const downVotes = item.votes.filter((v) => v.value === 'down').length
  const FallbackIcon = sectionIcons[sectionType] || MapPinned

  return (
    <div
      onClick={onClick}
      className={`bg-surface-0 rounded-2xl overflow-hidden shadow-sm border transition-all duration-200 hover:shadow-md cursor-pointer ${
        item.isFinal
          ? 'border-accent-400 ring-2 ring-accent-100'
          : item.status === 'rejected'
            ? 'border-surface-200 opacity-50'
            : 'border-surface-200 hover:border-accent-300'
      }`}
    >
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
              className="flex items-center gap-1 text-xs text-text-tertiary hover:text-accent-500 transition-colors cursor-pointer"
              onClick={(e) => e.stopPropagation()}
            >
              <ThumbsUp className="w-3.5 h-3.5" />
              <span className="font-semibold">{upVotes}</span>
            </button>
            <button
              className="flex items-center gap-1 text-xs text-text-tertiary hover:text-text-secondary transition-colors cursor-pointer"
              onClick={(e) => e.stopPropagation()}
            >
              <ThumbsDown className="w-3.5 h-3.5" />
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
              .map((v, i) => (
                <span
                  key={i}
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

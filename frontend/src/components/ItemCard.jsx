import {
  Star,
  MessageCircle,
  ExternalLink,
  Bookmark,
  Bed,
  Plane,
  MapPinned,
  Loader2,
  CornerUpLeft,
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
  const avgRating = item.avgRating ?? 0
  const ratingCount = item.ratingCount ?? 0
  const FallbackIcon = sectionIcons[sectionType] || MapPinned
  const isPending = !!item.pending

  const handleClick = () => {
    if (isPending) return
    onClick?.()
  }

  const isRejected = item.status === 'rejected'

  return (
    <div
      onClick={handleClick}
      className={`relative bg-surface-0 rounded-2xl overflow-hidden shadow-sm border transition-all duration-200 ${
        isPending
          ? 'border-surface-200 cursor-wait'
          : item.isFinal
            ? 'border-accent-400 ring-2 ring-accent-100 hover:shadow-md cursor-pointer'
            : isRejected
              ? 'border-surface-200 opacity-60 hover:shadow-md cursor-pointer'
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
          <img
            src={item.image}
            alt={item.title}
            className={`w-full h-40 object-cover transition-[filter] duration-200 ${
              isRejected ? 'grayscale' : ''
            }`}
          />
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
          <h4 className={`text-sm font-semibold text-text-primary leading-snug line-clamp-2 ${
            isRejected ? 'line-through text-text-tertiary' : ''
          }`}>
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

        {sectionType === 'transport' && (item.outboundDepartingAt || item.inboundDepartingAt) && (
          <div className="mt-2.5 space-y-1 text-xs">
            {item.outboundDepartingAt && (
              <div className="flex items-center gap-1.5 text-text-secondary">
                <Plane className="w-3.5 h-3.5 text-accent-500 shrink-0" />
                <span className="font-semibold text-text-primary">{formatDepartureShort(item.outboundDepartingAt)}</span>
                {item.outboundDepartingLocation && (
                  <span className="text-text-tertiary truncate">· {item.outboundDepartingLocation}</span>
                )}
              </div>
            )}
            {item.inboundDepartingAt && (
              <div className="flex items-center gap-1.5 text-text-secondary">
                <CornerUpLeft className="w-3.5 h-3.5 text-accent-500 shrink-0" />
                <span className="font-semibold text-text-primary">{formatDepartureShort(item.inboundDepartingAt)}</span>
                {item.inboundDepartingLocation && (
                  <span className="text-text-tertiary truncate">· {item.inboundDepartingLocation}</span>
                )}
              </div>
            )}
          </div>
        )}

        {item.note && (
          <p className="mt-2.5 text-xs text-text-secondary bg-surface-50 rounded-lg p-2.5 leading-relaxed">
            {item.note}
          </p>
        )}

        {item.status === 'booked' && item.bookingRef && (
          <div className="mt-2 text-xs bg-emerald-50 dark:bg-emerald-950/30 text-emerald-700 dark:text-emerald-400 rounded-lg p-2 font-mono">
            {t.bookingRef}: {item.bookingRef}
          </div>
        )}

        {/* Footer */}
        <div className="flex items-center justify-between mt-3 pt-3 border-t border-surface-200">
          {ratingCount > 0 ? (
            <div className="flex items-center gap-1.5 text-xs">
              <div className="flex items-center gap-0.5">
                {[1, 2, 3, 4, 5].map((n) => (
                  <Star
                    key={n}
                    className={`w-3.5 h-3.5 ${
                      n <= Math.round(avgRating)
                        ? 'text-amber-400'
                        : 'text-surface-300'
                    }`}
                    fill={n <= Math.round(avgRating) ? 'currentColor' : 'none'}
                  />
                ))}
              </div>
              <span className="font-semibold text-text-secondary">
                {avgRating.toFixed(1)}
              </span>
              <span className="text-text-tertiary">({ratingCount})</span>
            </div>
          ) : (
            <span className="text-xs text-text-muted italic">{t.noRatingsYet}</span>
          )}

          {item.comments.length > 0 && (
            <span className="flex items-center gap-1 text-xs text-text-tertiary">
              <MessageCircle className="w-3.5 h-3.5" />
              {item.comments.length}
            </span>
          )}
        </div>
      </div>
    </div>
  )
}

// Times are wall-clock (UTC-marked) — display in UTC so they match the airport clock.
function formatDepartureShort(iso) {
  if (!iso) return ''
  const d = new Date(iso)
  if (isNaN(d.getTime())) return ''
  return d.toLocaleString(undefined, {
    day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit', timeZone: 'UTC',
  })
}

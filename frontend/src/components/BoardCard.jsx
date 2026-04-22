import { Link } from 'react-router-dom'
import { MapPin, Calendar, Bed, Plane, MapPinned } from 'lucide-react'
import { useLang } from '../context/LanguageContext'

function formatDateRange(range, lang) {
  if (!range) return null
  const locale = lang === 'pl' ? 'pl-PL' : 'en-US'
  const fmt = (d) => new Date(d).toLocaleDateString(locale, { day: 'numeric', month: 'short' })
  return `${fmt(range.start)} — ${fmt(range.end)}`
}

export default function BoardCard({ board }) {
  const { lang, t } = useLang()
  const counts = board.counts || { accommodation: 0, transport: 0, activities: 0 }
  const total = counts.accommodation + counts.transport + counts.activities
  const dateLabel = formatDateRange(board.dateRange, lang)

  return (
    <Link to={`/board/${board.id}`} className="group block no-underline">
      <div className="bg-surface-0 rounded-2xl overflow-hidden shadow-sm hover:shadow-lg transition-all duration-300 border border-surface-200 hover:border-accent-300 hover:-translate-y-1">
        {/* Cover */}
        <div className="relative h-48 overflow-hidden bg-surface-100">
          {board.coverImage ? (
            <img
              src={board.coverImage}
              alt={board.destination}
              className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
            />
          ) : (
            <div className="w-full h-full flex items-center justify-center">
              <MapPin className="w-10 h-10 text-surface-300" />
            </div>
          )}
          <div className="absolute inset-0 bg-gradient-to-t from-black/50 via-transparent to-transparent" />
          <div className="absolute bottom-3 left-4 right-4">
            <h3 className="text-white text-lg font-bold drop-shadow-md">
              {board.name}
            </h3>
          </div>
        </div>

        {/* Info */}
        <div className="p-4">
          <div className="flex items-center gap-1.5 text-text-secondary text-sm mb-2">
            <MapPin className="w-3.5 h-3.5 shrink-0" />
            <span>{board.destination}</span>
          </div>
          <div className="flex items-center gap-1.5 text-text-tertiary text-xs mb-4">
            <Calendar className="w-3.5 h-3.5 shrink-0" />
            <span>{dateLabel || t.someday}</span>
          </div>

          <div className="flex items-center gap-4 text-xs text-text-tertiary">
            <span className="flex items-center gap-1">
              <Bed className="w-3.5 h-3.5" />
              {counts.accommodation}
            </span>
            <span className="flex items-center gap-1">
              <Plane className="w-3.5 h-3.5" />
              {counts.transport}
            </span>
            <span className="flex items-center gap-1">
              <MapPinned className="w-3.5 h-3.5" />
              {counts.activities}
            </span>
            {total > 0 && (
              <span className="ml-auto text-text-secondary font-semibold">
                {total} {total === 1 ? t.link : t.links}
              </span>
            )}
          </div>
        </div>
      </div>
    </Link>
  )
}

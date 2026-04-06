import { useEffect, useCallback } from 'react'
import {
  X,
  ChevronLeft,
  ChevronRight,
  MapPin,
  Clock,
  Camera,
  HardDrive,
} from 'lucide-react'
import { useLang } from '../context/LanguageContext'

export default function MemoryLightbox({ memories, currentIndex, onClose, onNavigate }) {
  const { lang, t } = useLang()
  const memory = memories[currentIndex]
  const locale = lang === 'pl' ? 'pl-PL' : 'en-US'
  const hasPrev = currentIndex > 0
  const hasNext = currentIndex < memories.length - 1

  const goNext = useCallback(() => {
    if (hasNext) onNavigate(currentIndex + 1)
  }, [hasNext, currentIndex, onNavigate])

  const goPrev = useCallback(() => {
    if (hasPrev) onNavigate(currentIndex - 1)
  }, [hasPrev, currentIndex, onNavigate])

  useEffect(() => {
    const handleKey = (e) => {
      if (e.key === 'Escape') onClose()
      if (e.key === 'ArrowRight') goNext()
      if (e.key === 'ArrowLeft') goPrev()
    }
    window.addEventListener('keydown', handleKey)
    return () => window.removeEventListener('keydown', handleKey)
  }, [onClose, goNext, goPrev])

  const formatDate = (iso) =>
    new Date(iso).toLocaleDateString(locale, {
      weekday: 'long',
      day: 'numeric',
      month: 'long',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })

  return (
    <div className="fixed inset-0 z-50 bg-black/90 flex">
      {/* Close */}
      <button
        onClick={onClose}
        className="absolute top-4 right-4 z-10 w-10 h-10 flex items-center justify-center rounded-full bg-white/10 hover:bg-white/20 text-white transition-colors cursor-pointer"
      >
        <X className="w-5 h-5" />
      </button>

      {/* Counter */}
      <div className="absolute top-4 left-4 z-10 text-white/60 text-sm font-medium">
        {currentIndex + 1} {t.memoryOf} {memories.length}
      </div>

      {/* Main area */}
      <div className="flex-1 flex items-center justify-center relative min-w-0">
        {/* Prev */}
        {hasPrev && (
          <button
            onClick={goPrev}
            className="absolute left-4 z-10 w-12 h-12 flex items-center justify-center rounded-full bg-white/10 hover:bg-white/20 text-white transition-colors cursor-pointer"
          >
            <ChevronLeft className="w-6 h-6" />
          </button>
        )}

        {/* Image */}
        <img
          key={memory.id}
          src={memory.src}
          alt={memory.caption}
          className="max-h-[85vh] max-w-full object-contain rounded-lg animate-fade-in"
        />

        {/* Next */}
        {hasNext && (
          <button
            onClick={goNext}
            className="absolute right-4 z-10 w-12 h-12 flex items-center justify-center rounded-full bg-white/10 hover:bg-white/20 text-white transition-colors cursor-pointer"
          >
            <ChevronRight className="w-6 h-6" />
          </button>
        )}
      </div>

      {/* Info panel — right side */}
      <div className="w-80 bg-surface-0 border-l border-surface-200 flex flex-col overflow-y-auto shrink-0 hidden lg:flex">
        {/* Thumbnail strip */}
        <div className="p-4 border-b border-surface-200">
          <div className="flex gap-2 overflow-x-auto pb-1">
            {memories.map((m, i) => (
              <button
                key={m.id}
                onClick={() => onNavigate(i)}
                className={`shrink-0 w-14 h-14 rounded-lg overflow-hidden cursor-pointer transition-all ${
                  i === currentIndex
                    ? 'ring-2 ring-accent-500 ring-offset-2 ring-offset-surface-0'
                    : 'opacity-50 hover:opacity-80'
                }`}
              >
                <img src={m.src} alt="" className="w-full h-full object-cover" />
              </button>
            ))}
          </div>
        </div>

        <div className="p-5 space-y-4 flex-1">
          {/* Caption */}
          {memory.caption && (
            <div>
              <label className="text-xs font-semibold text-text-tertiary uppercase tracking-wider mb-1 block">
                {t.memoryCaption}
              </label>
              <p className="text-sm text-text-primary font-medium leading-relaxed">
                {memory.caption}
              </p>
            </div>
          )}

          {/* Metadata */}
          <div className="space-y-3">
            {memory.takenAt && (
              <div className="flex items-start gap-3">
                <Clock className="w-4 h-4 text-text-muted mt-0.5 shrink-0" />
                <div>
                  <label className="text-[11px] font-semibold text-text-muted uppercase tracking-wider block">
                    {t.memoryTakenAt}
                  </label>
                  <p className="text-sm text-text-primary">{formatDate(memory.takenAt)}</p>
                </div>
              </div>
            )}

            {memory.location && (
              <div className="flex items-start gap-3">
                <MapPin className="w-4 h-4 text-text-muted mt-0.5 shrink-0" />
                <div>
                  <label className="text-[11px] font-semibold text-text-muted uppercase tracking-wider block">
                    {t.memoryLocation}
                  </label>
                  <p className="text-sm text-text-primary">{memory.location}</p>
                </div>
              </div>
            )}

            {memory.camera && (
              <div className="flex items-start gap-3">
                <Camera className="w-4 h-4 text-text-muted mt-0.5 shrink-0" />
                <div>
                  <label className="text-[11px] font-semibold text-text-muted uppercase tracking-wider block">
                    {t.memoryCamera}
                  </label>
                  <p className="text-sm text-text-primary">{memory.camera}</p>
                </div>
              </div>
            )}

            {memory.size && (
              <div className="flex items-start gap-3">
                <HardDrive className="w-4 h-4 text-text-muted mt-0.5 shrink-0" />
                <div>
                  <label className="text-[11px] font-semibold text-text-muted uppercase tracking-wider block">
                    {t.memorySize}
                  </label>
                  <p className="text-sm text-text-primary">{memory.size}</p>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      <style>{`
        @keyframes fadeIn {
          from { opacity: 0; transform: scale(0.97); }
          to { opacity: 1; transform: scale(1); }
        }
        .animate-fade-in {
          animation: fadeIn 0.2s ease-out;
        }
      `}</style>
    </div>
  )
}

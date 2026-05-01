import { useRef, useState } from 'react'
import { Camera, MapPin, Clock, Plus, Upload, Trash2, Loader2 } from 'lucide-react'
import { useLang } from '../context/LanguageContext'
import MemoryLightbox from './MemoryLightbox'

export default function MemoryGallery({ memories, onUpload, onDelete, canDelete = false, uploading = false }) {
  const { lang, t } = useLang()
  const [lightboxIndex, setLightboxIndex] = useState(null)
  const [busyDeleteId, setBusyDeleteId] = useState(null)
  const fileInputRef = useRef(null)
  const locale = lang === 'pl' ? 'pl-PL' : 'en-US'

  const triggerUpload = () => fileInputRef.current?.click()

  const handleFiles = async (e) => {
    const files = Array.from(e.target.files || [])
    if (!files.length || !onUpload) return
    e.target.value = ''
    for (const file of files) {
      try {
        await onUpload(file)
      } catch {}
    }
  }

  const handleDelete = async (memoryId) => {
    if (!onDelete) return
    if (!window.confirm(t.confirmDeleteMemory || 'Delete this memory?')) return
    setBusyDeleteId(memoryId)
    try {
      await onDelete(memoryId)
    } finally {
      setBusyDeleteId(null)
    }
  }

  const fileInput = (
    <input
      ref={fileInputRef}
      type="file"
      accept="image/*"
      multiple
      hidden
      onChange={handleFiles}
    />
  )

  if (memories.length === 0) {
    return (
      <div className="border-2 border-dashed border-surface-200 rounded-2xl py-16 text-center">
        <div className="w-16 h-16 rounded-2xl bg-surface-100 flex items-center justify-center mx-auto mb-4">
          <Camera className="w-8 h-8 text-surface-300" />
        </div>
        <p className="text-sm text-text-muted mb-4">{t.noMemories}</p>
        {onUpload && (
          <button
            onClick={triggerUpload}
            disabled={uploading}
            className="text-xs font-semibold text-accent-500 hover:text-accent-600 bg-accent-50 hover:bg-accent-100 px-4 py-2 rounded-xl transition-colors cursor-pointer disabled:opacity-50"
          >
            {uploading ? <Loader2 className="w-3.5 h-3.5 inline mr-1 animate-spin" /> : <Plus className="w-3.5 h-3.5 inline mr-1" />}
            {t.addMemory}
          </button>
        )}
        {fileInput}
      </div>
    )
  }

  return (
    <>
      {/* Upload zone */}
      {onUpload && (
        <div className="mb-6">
          <button
            onClick={triggerUpload}
            disabled={uploading}
            className="w-full border-2 border-dashed border-surface-200 hover:border-accent-300 rounded-2xl py-6 flex flex-col items-center gap-2 transition-colors cursor-pointer group hover:bg-accent-50/50 disabled:opacity-50"
          >
            {uploading ? (
              <Loader2 className="w-6 h-6 text-accent-500 animate-spin" />
            ) : (
              <Upload className="w-6 h-6 text-text-muted group-hover:text-accent-500 transition-colors" />
            )}
            <span className="text-sm text-text-muted group-hover:text-accent-500 font-medium transition-colors">
              {t.memoryUploadDesc}
            </span>
          </button>
        </div>
      )}
      {fileInput}

      {/* Gallery grid — variable row heights for visual interest */}
      <div className="columns-2 sm:columns-3 gap-3 space-y-3">
        {memories.map((memory, index) => (
          <div
            key={memory.id}
            className="break-inside-avoid group cursor-pointer"
            onClick={() => setLightboxIndex(index)}
          >
            <div className="relative bg-surface-0 rounded-xl overflow-hidden border border-surface-200 hover:border-accent-300 shadow-sm hover:shadow-md transition-all duration-200">
              <img
                src={memory.src}
                alt={memory.caption}
                className="w-full object-cover group-hover:scale-[1.03] transition-transform duration-300"
              />
              {canDelete && (
                <button
                  onClick={(e) => {
                    e.stopPropagation()
                    handleDelete(memory.id)
                  }}
                  disabled={busyDeleteId === memory.id}
                  title={t.confirmDeleteMemory || 'Delete'}
                  className="absolute top-2 right-2 w-8 h-8 flex items-center justify-center rounded-lg bg-black/50 text-white opacity-0 group-hover:opacity-100 hover:bg-red-500 transition-all disabled:opacity-50"
                >
                  {busyDeleteId === memory.id ? (
                    <Loader2 className="w-4 h-4 animate-spin" />
                  ) : (
                    <Trash2 className="w-4 h-4" />
                  )}
                </button>
              )}
              {/* Hover overlay */}
              <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none">
                <div className="absolute bottom-0 left-0 right-0 p-3">
                  {memory.caption && (
                    <p className="text-white text-xs font-semibold mb-1.5 leading-snug">
                      {memory.caption}
                    </p>
                  )}
                  <div className="flex items-center gap-3 text-white/70 text-[11px]">
                    {memory.location && (
                      <span className="flex items-center gap-1">
                        <MapPin className="w-3 h-3" />
                        {memory.location}
                      </span>
                    )}
                    {memory.takenAt && (
                      <span className="flex items-center gap-1">
                        <Clock className="w-3 h-3" />
                        {new Date(memory.takenAt).toLocaleDateString(locale, {
                          day: 'numeric',
                          month: 'short',
                        })}
                      </span>
                    )}
                  </div>
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Lightbox */}
      {lightboxIndex !== null && (
        <MemoryLightbox
          memories={memories}
          currentIndex={lightboxIndex}
          onClose={() => setLightboxIndex(null)}
          onNavigate={setLightboxIndex}
        />
      )}
    </>
  )
}

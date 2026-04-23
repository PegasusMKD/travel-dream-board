import { useEffect, useRef, useState } from 'react'
import { X, Link as LinkIcon, ImagePlus } from 'lucide-react'
import { useLang } from '../context/LanguageContext'

const MAX_IMAGE_BYTES = 10 * 1024 * 1024

export default function AddItemModal({ sectionType, onClose, onSubmit }) {
  const { t } = useLang()
  const [url, setUrl] = useState('')
  const [fileError, setFileError] = useState(null)
  const [isDragging, setIsDragging] = useState(false)
  const fileInputRef = useRef(null)

  const addLabels = {
    accommodation: t.addAccommodation,
    transport: t.addTransport,
    activities: t.addActivity,
  }

  const acceptFile = (candidate) => {
    if (!candidate) return
    if (!candidate.type || !candidate.type.startsWith('image/')) {
      setFileError(t.imageWrongType)
      return
    }
    if (candidate.size > MAX_IMAGE_BYTES) {
      setFileError(t.imageTooLarge)
      return
    }
    setFileError(null)
    onSubmit({ file: candidate })
    onClose()
  }

  useEffect(() => {
    const onPaste = (e) => {
      const items = e.clipboardData?.items
      if (!items) return
      for (const item of items) {
        if (item.kind === 'file' && item.type.startsWith('image/')) {
          const f = item.getAsFile()
          if (f) {
            e.preventDefault()
            acceptFile(f)
            return
          }
        }
      }
    }
    window.addEventListener('paste', onPaste)
    return () => window.removeEventListener('paste', onPaste)
  }, [])

  const handleDrop = (e) => {
    e.preventDefault()
    setIsDragging(false)
    const dropped = e.dataTransfer?.files?.[0]
    if (dropped) acceptFile(dropped)
  }

  const handleDragOver = (e) => {
    e.preventDefault()
    if (!isDragging) setIsDragging(true)
  }

  const handleDragLeave = (e) => {
    e.preventDefault()
    setIsDragging(false)
  }

  const handleFileInput = (e) => {
    const picked = e.target.files?.[0]
    if (picked) acceptFile(picked)
    e.target.value = ''
  }

  const handleSubmit = (e) => {
    e.preventDefault()
    const trimmed = url.trim()
    if (!trimmed) return
    onSubmit({ url: trimmed })
    onClose()
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
          <div className="relative mb-3">
            <LinkIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-text-muted" />
            <input
              type="url"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              placeholder={t.pasteLink}
              className="w-full pl-10 pr-4 py-3 bg-surface-50 border border-surface-200 rounded-xl text-sm focus:outline-none focus:border-accent-400 focus:ring-2 focus:ring-accent-100 placeholder:text-text-muted text-text-primary transition-colors"
              autoFocus
            />
          </div>

          <button
            type="submit"
            disabled={!url.trim()}
            className="w-full bg-accent-500 hover:bg-accent-600 text-white py-3 rounded-xl text-sm font-semibold transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 cursor-pointer mb-4"
          >
            {t.addLinkBtn}
          </button>

          <div className="flex items-center gap-3 mb-3">
            <div className="flex-1 h-px bg-surface-200" />
            <span className="text-xs text-text-tertiary uppercase tracking-wider">
              {t.orDivider}
            </span>
            <div className="flex-1 h-px bg-surface-200" />
          </div>

          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            className="hidden"
            onChange={handleFileInput}
          />

          <button
            type="button"
            onClick={() => fileInputRef.current?.click()}
            onDrop={handleDrop}
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            className={`w-full rounded-xl border-2 border-dashed transition-colors px-4 py-6 flex flex-col items-center justify-center gap-2 cursor-pointer ${
              isDragging
                ? 'border-accent-400 bg-accent-50/50'
                : 'border-surface-200 bg-surface-50 hover:border-accent-300 hover:bg-surface-100'
            }`}
          >
            <ImagePlus className="w-6 h-6 text-text-tertiary" />
            <span className="text-sm text-text-secondary text-center">{t.dropImage}</span>
            <span className="text-xs text-text-tertiary text-center">{t.pasteImageHint}</span>
          </button>

          {fileError && (
            <p className="text-xs text-red-500 mt-2">{fileError}</p>
          )}

          <p className="text-xs text-text-tertiary mt-4 leading-relaxed">
            {t.addLinkHelp}
          </p>
        </form>
      </div>
    </div>
  )
}

import { useState } from 'react'
import { X, Copy, Check, Users } from 'lucide-react'
import { useLang } from '../context/LanguageContext'

export default function ShareModal({ boardName, onClose }) {
  const { t } = useLang()
  const [copied, setCopied] = useState(false)
  const shareUrl = `${window.location.origin}/boards/abc123?token=demo-share-token`

  const handleCopy = () => {
    navigator.clipboard.writeText(shareUrl)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-black/30 backdrop-blur-sm" onClick={onClose} />
      <div className="relative bg-surface-0 rounded-2xl shadow-2xl w-full max-w-md overflow-hidden border border-surface-200">
        <div className="px-6 pt-5 pb-3 flex items-center justify-between">
          <h3 className="text-base font-bold text-text-primary">{t.shareBoard}</h3>
          <button
            onClick={onClose}
            className="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-surface-100 transition-colors text-text-tertiary cursor-pointer"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <div className="px-6 pb-6">
          <p className="text-sm text-text-secondary mb-4">
            {t.shareDesc} <strong>"{boardName}"</strong>.
          </p>

          <div className="flex items-center gap-2 mb-4">
            <div className="flex-1 bg-surface-50 border border-surface-200 rounded-xl px-3 py-2.5 text-xs text-text-tertiary font-mono truncate">
              {shareUrl}
            </div>
            <button
              onClick={handleCopy}
              className="shrink-0 w-10 h-10 flex items-center justify-center bg-accent-50 hover:bg-accent-100 text-accent-500 rounded-xl transition-colors cursor-pointer"
            >
              {copied ? <Check className="w-4 h-4" /> : <Copy className="w-4 h-4" />}
            </button>
          </div>

          <div className="bg-surface-50 rounded-xl p-4">
            <div className="flex items-center gap-2 mb-2">
              <Users className="w-4 h-4 text-text-secondary" />
              <span className="text-xs font-semibold text-text-primary">{t.friendAccess}</span>
            </div>
            <ul className="text-xs text-text-secondary space-y-1 ml-6 list-disc">
              <li>{t.canVote}</li>
              <li>{t.canComment}</li>
              <li>{t.canView}</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  )
}

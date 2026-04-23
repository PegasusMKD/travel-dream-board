import { AlertTriangle, WifiOff, Link2Off, RotateCw } from 'lucide-react'
import { useLang } from '../context/LanguageContext'

const variantConfig = {
  network: { icon: WifiOff, titleKey: 'networkErrorTitle', descKey: 'networkErrorDesc' },
  scrape: { icon: Link2Off, titleKey: 'scrapeErrorTitle', descKey: 'scrapeErrorDesc' },
  generic: { icon: AlertTriangle, titleKey: null, descKey: null },
}

export default function ErrorState({ variant = 'generic', message, onRetry, className = '' }) {
  const { t } = useLang()
  const config = variantConfig[variant] || variantConfig.generic
  const Icon = config.icon
  const title = config.titleKey ? t[config.titleKey] : null
  const desc = config.descKey ? t[config.descKey] : null

  return (
    <div className={`border border-red-200 dark:border-red-900/40 bg-red-50 dark:bg-red-950/20 rounded-2xl py-6 px-5 text-center ${className}`}>
      <div className="w-10 h-10 mx-auto mb-2 rounded-xl bg-red-100 dark:bg-red-900/40 flex items-center justify-center">
        <Icon className="w-5 h-5 text-red-500" />
      </div>
      {title && <h3 className="text-sm font-semibold text-text-primary mb-1">{title}</h3>}
      {(message || desc) && (
        <p className="text-xs text-text-secondary">{message || desc}</p>
      )}
      {onRetry && (
        <button
          onClick={onRetry}
          className="mt-3 inline-flex items-center gap-1.5 text-xs font-semibold text-accent-500 hover:text-accent-600 bg-surface-0 hover:bg-surface-50 px-3 py-1.5 rounded-lg transition-colors cursor-pointer border border-surface-200"
        >
          <RotateCw className="w-3 h-3" />
          {t.retry}
        </button>
      )}
    </div>
  )
}

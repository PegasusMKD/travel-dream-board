import { Circle, Star, X, CheckCircle2, Trophy } from 'lucide-react'
import { useLang } from '../context/LanguageContext'

const statusConfig = {
  considering: {
    key: 'considering',
    bg: 'bg-surface-100',
    text: 'text-text-tertiary',
    icon: Circle,
  },
  finalist: {
    key: 'finalist',
    bg: 'bg-amber-50 dark:bg-amber-950/40',
    text: 'text-amber-600 dark:text-amber-400',
    icon: Star,
  },
  rejected: {
    key: 'rejected',
    bg: 'bg-surface-100',
    text: 'text-text-muted line-through',
    icon: X,
  },
  booked: {
    key: 'booked',
    bg: 'bg-emerald-50 dark:bg-emerald-950/40',
    text: 'text-emerald-600 dark:text-emerald-400',
    icon: CheckCircle2,
  },
  completed: {
    key: 'completed',
    bg: 'bg-blue-50 dark:bg-blue-950/40',
    text: 'text-blue-600 dark:text-blue-400',
    icon: Trophy,
  },
}

export default function StatusBadge({ status }) {
  const { t } = useLang()
  const config = statusConfig[status] || statusConfig.considering
  const Icon = config.icon

  return (
    <span className={`inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-semibold ${config.bg} ${config.text}`}>
      <Icon className="w-3 h-3" />
      {t[config.key]}
    </span>
  )
}

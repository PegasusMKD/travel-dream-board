export default function EmptyState({ icon, emoji, title, description, action, className = '' }) {
  const Icon = icon
  return (
    <div className={`border-2 border-dashed border-surface-200 rounded-2xl py-10 px-6 text-center ${className}`}>
      {Icon ? (
        <div className="w-12 h-12 mx-auto mb-3 rounded-2xl bg-surface-100 flex items-center justify-center">
          <Icon className="w-6 h-6 text-text-muted" />
        </div>
      ) : emoji ? (
        <div className="text-3xl mb-2">{emoji}</div>
      ) : null}
      {title && (
        <h3 className="text-sm font-semibold text-text-secondary mb-1">{title}</h3>
      )}
      {description && (
        <p className="text-xs text-text-muted">{description}</p>
      )}
      {action && (
        <div className="mt-4">{action}</div>
      )}
    </div>
  )
}

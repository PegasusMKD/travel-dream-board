export function Skeleton({ className = '' }) {
  return <div className={`bg-surface-100 dark:bg-surface-200 rounded animate-pulse ${className}`} />
}

export function ItemCardSkeleton() {
  return (
    <div className="bg-surface-0 rounded-2xl overflow-hidden shadow-sm border border-surface-200">
      <Skeleton className="w-full h-40 rounded-none" />
      <div className="p-4 space-y-3">
        <Skeleton className="h-4 w-3/4" />
        <Skeleton className="h-5 w-20 rounded-full" />
        <div className="pt-3 border-t border-surface-200 flex items-center justify-between">
          <Skeleton className="h-3 w-16" />
          <Skeleton className="h-3 w-8" />
        </div>
      </div>
    </div>
  )
}

export function BoardCardSkeleton() {
  return (
    <div className="bg-surface-0 rounded-2xl overflow-hidden shadow-sm border border-surface-200">
      <Skeleton className="w-full h-40 rounded-none" />
      <div className="p-4 space-y-3">
        <Skeleton className="h-5 w-2/3" />
        <Skeleton className="h-3 w-1/2" />
        <div className="flex items-center gap-3 pt-2">
          <Skeleton className="h-3 w-12" />
          <Skeleton className="h-3 w-12" />
          <Skeleton className="h-3 w-12" />
        </div>
      </div>
    </div>
  )
}

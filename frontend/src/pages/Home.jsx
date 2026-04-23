import { useCallback, useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Plus, Sparkles, MapPin } from 'lucide-react'
import BoardCard from '../components/BoardCard'
import EditBoardModal from '../components/EditBoardModal'
import EmptyState from '../components/EmptyState'
import ErrorState from '../components/ErrorState'
import { BoardCardSkeleton } from '../components/Skeleton'
import { useLang } from '../context/LanguageContext'
import { api, NetworkError } from '../services/api'
import { mapBoardSummary } from '../services/mappers'

export default function Home() {
  const { t } = useLang()
  const navigate = useNavigate()
  const [boards, setBoards] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [errorVariant, setErrorVariant] = useState('generic')
  const [creating, setCreating] = useState(false)

  const loadBoards = useCallback(() => {
    setLoading(true)
    setError(null)
    api.boards.getAll()
      .then((data) => setBoards((data || []).map(mapBoardSummary)))
      .catch((err) => {
        setErrorVariant(err instanceof NetworkError ? 'network' : 'generic')
        setError(err.message)
      })
      .finally(() => setLoading(false))
  }, [])

  useEffect(() => { loadBoards() }, [loadBoards])

  const handleCreated = (createdRaw) => {
    if (createdRaw?.uuid) {
      navigate(`/board/${createdRaw.uuid}`)
    }
  }

  return (
    <div className="pt-8">
      {/* Hero */}
      <div className="text-center mb-10">
        <div className="flex items-center justify-center gap-2 mb-3">
          <Sparkles className="w-5 h-5 text-accent-400" />
          <span className="text-sm text-text-tertiary font-semibold">{t.heroTagline}</span>
        </div>
        <h1 className="text-3xl sm:text-4xl font-extrabold text-text-primary mb-2">
          {t.heroTitle}
        </h1>
        <p className="text-text-secondary text-sm max-w-md mx-auto">
          {t.heroSubtitle}
        </p>
      </div>

      {loading ? (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          <BoardCardSkeleton />
          <BoardCardSkeleton />
          <BoardCardSkeleton />
        </div>
      ) : error ? (
        <ErrorState variant={errorVariant} message={errorVariant === 'generic' ? error : null} onRetry={loadBoards} />
      ) : boards.length === 0 ? (
        <EmptyState
          icon={MapPin}
          title={t.emptyBoardsTitle}
          description={t.emptyBoardsDesc}
          action={(
            <button
              onClick={() => setCreating(true)}
              className="inline-flex items-center gap-1.5 bg-accent-500 hover:bg-accent-600 text-white text-sm font-semibold px-4 py-2 rounded-xl transition-colors cursor-pointer"
            >
              <Plus className="w-4 h-4" />
              {t.newTrip}
            </button>
          )}
        />
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {boards.map((board) => (
            <BoardCard key={board.id} board={board} />
          ))}

          <button
            onClick={() => setCreating(true)}
            className="group border-2 border-dashed border-surface-200 hover:border-accent-400 rounded-2xl flex flex-col items-center justify-center min-h-[280px] transition-all duration-300 hover:bg-accent-50 cursor-pointer"
          >
            <div className="w-14 h-14 rounded-2xl bg-surface-100 group-hover:bg-accent-100 flex items-center justify-center transition-colors mb-3">
              <Plus className="w-7 h-7 text-text-muted group-hover:text-accent-500 transition-colors" />
            </div>
            <span className="text-sm font-semibold text-text-muted group-hover:text-accent-500 transition-colors">
              {t.newTrip}
            </span>
          </button>
        </div>
      )}

      {creating && (
        <EditBoardModal
          onClose={() => setCreating(false)}
          onSaved={handleCreated}
        />
      )}
    </div>
  )
}

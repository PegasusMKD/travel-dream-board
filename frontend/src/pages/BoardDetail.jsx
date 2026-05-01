import { useCallback, useEffect, useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import {
  ArrowLeft,
  Calendar,
  MapPin,
  Share2,
  Plus,
  Bed,
  Plane,
  MapPinned,
  Settings,
  ClipboardList,
  Camera,
} from 'lucide-react'
import { useLang } from '../context/LanguageContext'
import { useAuth } from '../context/AuthContext'
import { useShare } from '../context/ShareContext'
import ItemCard from '../components/ItemCard'
import AddItemModal from '../components/AddItemModal'
import ShareModal from '../components/ShareModal'
import ItemDetailSidebar from '../components/ItemDetailSidebar'
import EditBoardModal from '../components/EditBoardModal'
import MemoryGallery from '../components/MemoryGallery'
import DisplayNamePrompt from '../components/DisplayNamePrompt'
import ErrorState from '../components/ErrorState'
import { ItemCardSkeleton } from '../components/Skeleton'
import { api, NetworkError } from '../services/api'
import { mapAggregatedBoard, mapMemory, backendVoteTarget, sectionToItemApi } from '../services/mappers'

const sectionConfig = {
  accommodation: { key: 'accommodation', icon: Bed, emoji: '\u{1F3E8}' },
  transport: { key: 'transport', icon: Plane, emoji: '\u{2708}\u{FE0F}' },
  activities: { key: 'activities', icon: MapPinned, emoji: '\u{1F4CD}' },
}

function formatDateRange(range, lang) {
  if (!range) return null
  const locale = lang === 'pl' ? 'pl-PL' : 'en-US'
  const fmt = (d) =>
    new Date(d).toLocaleDateString(locale, { day: 'numeric', month: 'long', year: 'numeric' })
  return `${fmt(range.start)} — ${fmt(range.end)}`
}

export default function BoardDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { lang, t } = useLang()
  const { user } = useAuth()
  const { shareToken, guestUuid, guestName, persistGuest } = useShare()
  const isGuest = !user && !!shareToken
  const voterUuid = user?.uuid || guestUuid
  const needsGuestName = isGuest && (!guestUuid || !guestName)
  const [board, setBoard] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [activeTab, setActiveTab] = useState('planning')
  const [addingTo, setAddingTo] = useState(null)
  const [showShare, setShowShare] = useState(false)
  const [selectedItemId, setSelectedItemId] = useState(null)
  const [selectedSection, setSelectedSection] = useState(null)
  const [showEditBoard, setShowEditBoard] = useState(false)
  const [pending, setPending] = useState({ accommodation: [], transport: [], activities: [] })
  const [memories, setMemories] = useState([])
  const [uploadingMemory, setUploadingMemory] = useState(false)

  const loadBoard = useCallback(async () => {
    try {
      const data = await api.boards.getById(id)
      setBoard(mapAggregatedBoard(data))
      setError(null)
    } catch (err) {
      setError(err)
    }
  }, [id])

  const loadMemories = useCallback(async () => {
    try {
      const list = await api.memories.list(id)
      setMemories((list || []).map((m) => mapMemory(m, shareToken)))
    } catch {
      setMemories([])
    }
  }, [id, shareToken])

  useEffect(() => {
    loadMemories()
  }, [loadMemories])

  const handleUploadMemory = useCallback(async (file) => {
    setUploadingMemory(true)
    try {
      await api.memories.create({ boardUuid: id, file, uploadedBy: !user ? voterUuid : undefined })
      await loadMemories()
    } finally {
      setUploadingMemory(false)
    }
  }, [id, loadMemories, user, voterUuid])

  const handleDeleteMemory = useCallback(async (memoryId) => {
    await api.memories.delete(memoryId)
    await loadMemories()
  }, [loadMemories])

  useEffect(() => {
    setLoading(true)
    loadBoard().finally(() => setLoading(false))
  }, [loadBoard])

  const closeSidebar = () => {
    setSelectedItemId(null)
    setSelectedSection(null)
  }

  const handleAddItem = (sectionType, { url, file }) => {
    const tempId = `pending-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`
    const previewUrl = file ? URL.createObjectURL(file) : null
    const placeholder = buildPlaceholder(tempId, {
      url,
      title: url || file?.name || '',
      image: previewUrl,
    })
    setPending((prev) => ({
      ...prev,
      [sectionType]: [...prev[sectionType], placeholder],
    }))

    const itemApi = sectionToItemApi(api, sectionType)
    itemApi.create({ url, file, boardUuid: id })
      .then(() => loadBoard())
      .catch((err) => setError(err.message))
      .finally(() => {
        if (previewUrl) URL.revokeObjectURL(previewUrl)
        setPending((prev) => ({
          ...prev,
          [sectionType]: prev[sectionType].filter((p) => p.id !== tempId),
        }))
      })
  }

  const handleVote = async (item, rank) => {
    if (!voterUuid) return
    const myVote = item.votes.find((v) => v.userUuid === voterUuid)
    try {
      if (!myVote) {
        await api.votes.create({
          user_uuid: voterUuid,
          rank,
          voted_on: backendVoteTarget(selectedSectionOf(item, board)),
          voted_on_uuid: item.id,
        })
      } else if (myVote.rank === rank) {
        await api.votes.delete(myVote.id)
      } else {
        await api.votes.update(myVote.id, rank)
      }
      await loadBoard()
    } catch (err) {
      setError(err)
    }
  }

  const handleClearVote = async (item) => {
    if (!voterUuid) return
    const myVote = item.votes.find((v) => v.userUuid === voterUuid)
    if (!myVote) return
    try {
      await api.votes.delete(myVote.id)
      await loadBoard()
    } catch (err) {
      setError(err)
    }
  }

  if (loading) {
    return (
      <div className="pt-8 space-y-6">
        <div className="w-full h-56 sm:h-72 rounded-2xl bg-surface-100 animate-pulse" />
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          <ItemCardSkeleton />
          <ItemCardSkeleton />
          <ItemCardSkeleton />
        </div>
      </div>
    )
  }

  if (error || !board) {
    const variant = error instanceof NetworkError ? 'network' : 'generic'
    return (
      <div className="pt-12 max-w-md mx-auto">
        <ErrorState
          variant={variant}
          message={variant === 'generic' ? (error?.message || error || t.boardNotFound) : null}
          onRetry={loadBoard}
        />
        {!isGuest && (
          <div className="text-center mt-4">
            <Link to="/" className="text-accent-500 text-sm hover:underline">
              {t.backToList}
            </Link>
          </div>
        )}
      </div>
    )
  }

  const memoryCount = memories.length

  return (
    <div className="pt-6">
      {!isGuest && (
        <Link
          to="/"
          className="inline-flex items-center gap-1.5 text-sm text-text-tertiary hover:text-accent-500 transition-colors mb-4 no-underline"
        >
          <ArrowLeft className="w-4 h-4" />
          {t.allBoards}
        </Link>
      )}

      {/* Hero */}
      <div className="relative rounded-2xl overflow-hidden mb-6 bg-surface-100">
        {board.coverImage ? (
          <img
            src={board.coverImage}
            alt={board.destination}
            className="w-full h-56 sm:h-72 object-cover"
          />
        ) : (
          <div className="w-full h-56 sm:h-72 flex items-center justify-center">
            <MapPin className="w-16 h-16 text-surface-300" />
          </div>
        )}
        <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-black/20 to-transparent" />
        <div className="absolute bottom-0 left-0 right-0 p-6 sm:p-8">
          <h1 className="text-2xl sm:text-3xl font-extrabold text-white mb-2 drop-shadow-md">
            {board.name}
          </h1>
          <div className="flex flex-wrap items-center gap-4 text-white/90 text-sm">
            <span className="flex items-center gap-1.5">
              <MapPin className="w-4 h-4" />
              {board.destination}
            </span>
            {board.dateRange && (
              <span className="flex items-center gap-1.5">
                <Calendar className="w-4 h-4" />
                {formatDateRange(board.dateRange, lang)}
              </span>
            )}
          </div>
        </div>

        {!isGuest && (
          <div className="absolute top-4 right-4 flex items-center gap-2">
            <button
              onClick={() => setShowShare(true)}
              className="flex items-center gap-1.5 bg-white/90 backdrop-blur-sm text-gray-800 px-3 py-2 rounded-xl text-xs font-semibold hover:bg-white transition-colors shadow-sm cursor-pointer"
            >
              <Share2 className="w-3.5 h-3.5" />
              {t.share}
            </button>
            <button
              onClick={() => setShowEditBoard(true)}
              className="w-9 h-9 flex items-center justify-center bg-white/90 backdrop-blur-sm text-gray-800 rounded-xl hover:bg-white transition-colors shadow-sm cursor-pointer"
            >
              <Settings className="w-4 h-4" />
            </button>
          </div>
        )}
        {isGuest && (
          <div className="absolute top-4 right-4">
            <span className="inline-flex items-center gap-1 bg-white/90 backdrop-blur-sm text-gray-800 px-3 py-2 rounded-xl text-xs font-semibold shadow-sm">
              {t.guestBadge}{guestName ? ` · ${guestName}` : ''}
            </span>
          </div>
        )}
      </div>

      {/* Tabs */}
      <div className="flex items-center gap-1 mb-8 border-b border-surface-200 overflow-x-auto">
        <button
          onClick={() => setActiveTab('planning')}
          className={`flex items-center gap-2 px-4 py-3 text-sm font-semibold border-b-2 transition-colors cursor-pointer ${
            activeTab === 'planning'
              ? 'border-accent-500 text-accent-600'
              : 'border-transparent text-text-tertiary hover:text-text-secondary'
          }`}
        >
          <ClipboardList className="w-4 h-4" />
          {t.tabPlanning}
        </button>
        <button
          onClick={() => setActiveTab('memories')}
          className={`flex items-center gap-2 px-4 py-3 text-sm font-semibold border-b-2 transition-colors cursor-pointer ${
            activeTab === 'memories'
              ? 'border-accent-500 text-accent-600'
              : 'border-transparent text-text-tertiary hover:text-text-secondary'
          }`}
        >
          <Camera className="w-4 h-4" />
          {t.tabMemories}
          {memoryCount > 0 && (
            <span className="text-[11px] bg-surface-100 text-text-muted font-semibold px-1.5 py-0.5 rounded-full">
              {memoryCount}
            </span>
          )}
        </button>
      </div>

      {activeTab === 'planning' && (
        <div className="space-y-10">
          {Object.entries(sectionConfig).map(([key, config]) => {
            const items = [...(board.sections[key] || []), ...(pending[key] || [])]
            const Icon = config.icon

            return (
              <section key={key}>
                <div className="flex items-center justify-between mb-4">
                  <div className="flex items-center gap-2.5">
                    <div className="w-8 h-8 rounded-lg bg-surface-100 flex items-center justify-center">
                      <Icon className="w-4 h-4 text-text-secondary" />
                    </div>
                    <h2 className="text-xl font-bold text-text-primary">
                      {t[config.key]}
                    </h2>
                    <span className="text-xs text-text-muted bg-surface-100 font-semibold px-2 py-0.5 rounded-full">
                      {items.length}
                    </span>
                  </div>
                  {!isGuest && (
                    <button
                      onClick={() => setAddingTo(key)}
                      className="flex items-center gap-1.5 text-xs font-semibold text-accent-500 hover:text-accent-600 bg-accent-50 hover:bg-accent-100 px-3 py-1.5 rounded-xl transition-colors cursor-pointer"
                    >
                      <Plus className="w-3.5 h-3.5" />
                      {t.addLink}
                    </button>
                  )}
                </div>

                {items.length > 0 ? (
                  <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                    {items.map((item) => (
                      <ItemCard
                        key={item.id}
                        item={item}
                        sectionType={key}
                        onClick={() => {
                          setSelectedItemId(item.id)
                          setSelectedSection(key)
                        }}
                      />
                    ))}
                  </div>
                ) : (
                  <div className="border-2 border-dashed border-surface-200 rounded-2xl py-10 text-center">
                    <div className="text-3xl mb-2">{config.emoji}</div>
                    <p className="text-sm text-text-muted">{t.noLinks}</p>
                    {!isGuest && (
                      <button
                        onClick={() => setAddingTo(key)}
                        className="mt-3 text-xs font-semibold text-accent-500 hover:text-accent-600 bg-accent-50 hover:bg-accent-100 px-4 py-2 rounded-xl transition-colors cursor-pointer"
                      >
                        <Plus className="w-3.5 h-3.5 inline mr-1" />
                        {t.addLink}
                      </button>
                    )}
                  </div>
                )}
              </section>
            )
          })}
        </div>
      )}

      {activeTab === 'memories' && (
        <MemoryGallery
          memories={memories}
          onUpload={handleUploadMemory}
          onDelete={!isGuest ? handleDeleteMemory : undefined}
          canDelete={!isGuest}
          uploading={uploadingMemory}
        />
      )}

      {/* Modals */}
      {addingTo && (
        <AddItemModal
          sectionType={addingTo}
          onClose={() => setAddingTo(null)}
          onSubmit={(payload) => handleAddItem(addingTo, payload)}
        />
      )}
      {showShare && (
        <ShareModal
          boardUuid={board.id}
          boardName={board.name}
          onClose={() => setShowShare(false)}
        />
      )}
      {showEditBoard && (
        <EditBoardModal
          board={board}
          onClose={() => setShowEditBoard(false)}
          onSaved={loadBoard}
          onDeleted={() => navigate('/')}
        />
      )}

      {/* Detail sidebar */}
      {selectedItemId && selectedSection && (
        <ItemDetailSidebar
          itemId={selectedItemId}
          sectionType={selectedSection}
          onClose={closeSidebar}
          onRefresh={loadBoard}
          onVote={handleVote}
          onClearVote={handleClearVote}
          isGuest={isGuest}
          voterUuid={voterUuid}
          guestName={guestName}
          board={board}
        />
      )}

      {needsGuestName && (
        <DisplayNamePrompt
          onSubmitted={({ uuid, name }) => persistGuest(uuid, name)}
        />
      )}
    </div>
  )
}

function selectedSectionOf(item, board) {
  for (const [section, items] of Object.entries(board.sections)) {
    if (items.some((i) => i.id === item.id)) return section
  }
  return 'accommodation'
}

function buildPlaceholder(id, { url, title, image }) {
  return {
    id,
    url: url || '',
    title: title || url || '',
    image: image || null,
    note: '',
    status: 'considering',
    isFinal: false,
    bookingRef: null,
    avgRating: 0,
    ratingCount: 0,
    votes: [],
    comments: [],
    pending: true,
  }
}

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
import ItemCard from '../components/ItemCard'
import AddItemModal from '../components/AddItemModal'
import ShareModal from '../components/ShareModal'
import ItemDetailSidebar from '../components/ItemDetailSidebar'
import EditBoardModal from '../components/EditBoardModal'
import MemoryGallery from '../components/MemoryGallery'
import { api } from '../services/api'
import { mapAggregatedBoard, backendVoteTarget, sectionToItemApi } from '../services/mappers'

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

  const loadBoard = useCallback(async () => {
    try {
      const data = await api.boards.getById(id)
      setBoard(mapAggregatedBoard(data))
      setError(null)
    } catch (err) {
      setError(err.message)
    }
  }, [id])

  useEffect(() => {
    setLoading(true)
    loadBoard().finally(() => setLoading(false))
  }, [loadBoard])

  const closeSidebar = () => {
    setSelectedItemId(null)
    setSelectedSection(null)
  }

  const handleAddItem = (sectionType, url) => {
    const tempId = `pending-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`
    const placeholder = buildPlaceholder(tempId, url)
    setPending((prev) => ({
      ...prev,
      [sectionType]: [...prev[sectionType], placeholder],
    }))

    const itemApi = sectionToItemApi(api, sectionType)
    itemApi.create(url, id)
      .then(() => loadBoard())
      .catch((err) => setError(err.message))
      .finally(() => {
        setPending((prev) => ({
          ...prev,
          [sectionType]: prev[sectionType].filter((p) => p.id !== tempId),
        }))
      })
  }

  const handleVote = async (item, direction) => {
    if (!user) return
    const myVote = item.votes.find((v) => v.userUuid === user.uuid)
    const rank = direction === 'up' ? 1 : -1
    try {
      if (!myVote) {
        await api.votes.create({
          user_uuid: user.uuid,
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
      setError(err.message)
    }
  }

  if (loading) {
    return (
      <div className="flex justify-center py-20">
        <div className="w-8 h-8 border-3 border-accent-200 border-t-accent-500 rounded-full animate-spin" />
      </div>
    )
  }

  if (error || !board) {
    return (
      <div className="pt-20 text-center">
        <p className="text-text-secondary">{error || t.boardNotFound}</p>
        <Link to="/" className="text-accent-500 text-sm mt-2 inline-block hover:underline">
          {t.backToList}
        </Link>
      </div>
    )
  }

  const memoryCount = board.memories?.length || 0

  return (
    <div className="pt-6">
      <Link
        to="/"
        className="inline-flex items-center gap-1.5 text-sm text-text-tertiary hover:text-accent-500 transition-colors mb-4 no-underline"
      >
        <ArrowLeft className="w-4 h-4" />
        {t.allBoards}
      </Link>

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
      </div>

      {/* Tabs */}
      <div className="flex items-center gap-1 mb-8 border-b border-surface-200">
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
                  <button
                    onClick={() => setAddingTo(key)}
                    className="flex items-center gap-1.5 text-xs font-semibold text-accent-500 hover:text-accent-600 bg-accent-50 hover:bg-accent-100 px-3 py-1.5 rounded-xl transition-colors cursor-pointer"
                  >
                    <Plus className="w-3.5 h-3.5" />
                    {t.addLink}
                  </button>
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
                        onVote={handleVote}
                      />
                    ))}
                  </div>
                ) : (
                  <div className="border-2 border-dashed border-surface-200 rounded-2xl py-10 text-center">
                    <div className="text-3xl mb-2">{config.emoji}</div>
                    <p className="text-sm text-text-muted">{t.noLinks}</p>
                    <button
                      onClick={() => setAddingTo(key)}
                      className="mt-3 text-xs font-semibold text-accent-500 hover:text-accent-600 bg-accent-50 hover:bg-accent-100 px-4 py-2 rounded-xl transition-colors cursor-pointer"
                    >
                      <Plus className="w-3.5 h-3.5 inline mr-1" />
                      {t.addLink}
                    </button>
                  </div>
                )}
              </section>
            )
          })}
        </div>
      )}

      {activeTab === 'memories' && (
        <MemoryGallery memories={board.memories || []} />
      )}

      {/* Modals */}
      {addingTo && (
        <AddItemModal
          sectionType={addingTo}
          onClose={() => setAddingTo(null)}
          onSubmit={(url) => handleAddItem(addingTo, url)}
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

function buildPlaceholder(id, url) {
  return {
    id,
    url,
    title: url,
    image: null,
    note: '',
    status: 'considering',
    isFinal: false,
    bookingRef: null,
    likes: 0,
    dislikes: 0,
    votes: [],
    comments: [],
    pending: true,
  }
}

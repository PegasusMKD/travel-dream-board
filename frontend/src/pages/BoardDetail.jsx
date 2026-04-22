import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
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
import ItemCard from '../components/ItemCard'
import AddItemModal from '../components/AddItemModal'
import ShareModal from '../components/ShareModal'
import ItemDetailSidebar from '../components/ItemDetailSidebar'
import EditBoardModal from '../components/EditBoardModal'
import MemoryGallery from '../components/MemoryGallery'
import { api } from '../services/api'
import { mapAggregatedBoard } from '../services/mappers'

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
  const { lang, t } = useLang()
  const [board, setBoard] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [activeTab, setActiveTab] = useState('planning')
  const [addingTo, setAddingTo] = useState(null)
  const [showShare, setShowShare] = useState(false)
  const [selectedItem, setSelectedItem] = useState(null)
  const [selectedSection, setSelectedSection] = useState(null)
  const [showEditBoard, setShowEditBoard] = useState(false)

  useEffect(() => {
    setLoading(true)
    setError(null)
    api.boards.getById(id)
      .then((data) => setBoard(mapAggregatedBoard(data)))
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false))
  }, [id])

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
      <div className="relative rounded-2xl overflow-hidden mb-6">
        <img
          src={board.coverImage}
          alt={board.destination}
          className="w-full h-56 sm:h-72 object-cover"
        />
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

      {/* Tab content */}
      {activeTab === 'planning' && (
        <div className="space-y-10">
          {Object.entries(sectionConfig).map(([key, config]) => {
            const items = board.sections[key] || []
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
                          setSelectedItem(item)
                          setSelectedSection(key)
                        }}
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
      {addingTo && <AddItemModal sectionType={addingTo} onClose={() => setAddingTo(null)} />}
      {showShare && <ShareModal boardName={board.name} onClose={() => setShowShare(false)} />}
      {showEditBoard && <EditBoardModal board={board} onClose={() => setShowEditBoard(false)} />}

      {/* Detail sidebar */}
      {selectedItem && (
        <ItemDetailSidebar
          item={selectedItem}
          sectionType={selectedSection}
          onClose={() => {
            setSelectedItem(null)
            setSelectedSection(null)
          }}
        />
      )}
    </div>
  )
}

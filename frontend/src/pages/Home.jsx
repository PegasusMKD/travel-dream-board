import { Plus, Sparkles } from 'lucide-react'
import BoardCard from '../components/BoardCard'
import { mockBoards } from '../data/mockData'
import { useLang } from '../context/LanguageContext'

export default function Home() {
  const { t } = useLang()

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

      {/* Board grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
        {mockBoards.map((board) => (
          <BoardCard key={board.id} board={board} />
        ))}

        {/* Add new board */}
        <button className="group border-2 border-dashed border-surface-200 hover:border-accent-400 rounded-2xl flex flex-col items-center justify-center min-h-[280px] transition-all duration-300 hover:bg-accent-50 cursor-pointer">
          <div className="w-14 h-14 rounded-2xl bg-surface-100 group-hover:bg-accent-100 flex items-center justify-center transition-colors mb-3">
            <Plus className="w-7 h-7 text-text-muted group-hover:text-accent-500 transition-colors" />
          </div>
          <span className="text-sm font-semibold text-text-muted group-hover:text-accent-500 transition-colors">
            {t.newTrip}
          </span>
        </button>
      </div>
    </div>
  )
}

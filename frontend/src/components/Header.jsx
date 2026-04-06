import { Link } from 'react-router-dom'
import { Plane, Plus, Sun, Moon, Globe } from 'lucide-react'
import { useLang } from '../context/LanguageContext'
import { useTheme } from '../context/ThemeContext'

export default function Header() {
  const { lang, toggle: toggleLang, t } = useLang()
  const { theme, toggle: toggleTheme } = useTheme()

  return (
    <header className="bg-surface-0/80 backdrop-blur-md border-b border-surface-200 sticky top-0 z-50 transition-colors">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          <Link to="/" className="flex items-center gap-2.5 no-underline">
            <div className="w-9 h-9 bg-accent-500 rounded-xl flex items-center justify-center shadow-sm">
              <Plane className="w-5 h-5 text-white" strokeWidth={2} />
            </div>
            <span className="text-lg font-bold text-text-primary tracking-tight">
              {t.appName}
            </span>
          </Link>

          <div className="flex items-center gap-2">
            {/* Language toggle */}
            <button
              onClick={toggleLang}
              className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold text-text-secondary hover:bg-surface-100 transition-colors cursor-pointer"
              title={lang === 'pl' ? 'Switch to English' : 'Przełącz na polski'}
            >
              <Globe className="w-3.5 h-3.5" />
              {lang === 'pl' ? 'EN' : 'PL'}
            </button>

            {/* Theme toggle */}
            <button
              onClick={toggleTheme}
              className="w-9 h-9 flex items-center justify-center rounded-lg text-text-secondary hover:bg-surface-100 transition-colors cursor-pointer"
              title={theme === 'dark' ? t.lightMode : t.darkMode}
            >
              {theme === 'dark' ? <Sun className="w-4 h-4" /> : <Moon className="w-4 h-4" />}
            </button>

            {/* New trip button */}
            <button className="flex items-center gap-2 bg-accent-500 hover:bg-accent-600 text-white px-4 py-2 rounded-xl text-sm font-semibold transition-colors shadow-sm cursor-pointer">
              <Plus className="w-4 h-4" />
              <span className="hidden sm:inline">{t.newTrip}</span>
            </button>

            {/* Avatar */}
            <div className="w-8 h-8 rounded-full bg-accent-200 flex items-center justify-center text-accent-700 text-xs font-bold">
              Z
            </div>
          </div>
        </div>
      </div>
    </header>
  )
}

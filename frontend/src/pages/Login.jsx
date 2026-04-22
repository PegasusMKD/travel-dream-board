import { Plane } from 'lucide-react'
import { useLang } from '../context/LanguageContext'
import { useTheme } from '../context/ThemeContext'
import { Sun, Moon, Globe } from 'lucide-react'
import { api } from '../services/api'

export default function Login() {
  const { lang, toggle: toggleLang, t } = useLang()
  const { theme, toggle: toggleTheme } = useTheme()

  return (
    <div className="min-h-screen bg-surface-50 flex flex-col items-center justify-center px-4 transition-colors duration-200">
      {/* Top-right controls */}
      <div className="fixed top-4 right-4 flex items-center gap-2">
        <button
          onClick={toggleLang}
          className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold text-text-secondary hover:bg-surface-100 transition-colors cursor-pointer"
        >
          <Globe className="w-3.5 h-3.5" />
          {lang === 'pl' ? 'EN' : 'PL'}
        </button>
        <button
          onClick={toggleTheme}
          className="w-9 h-9 flex items-center justify-center rounded-lg text-text-secondary hover:bg-surface-100 transition-colors cursor-pointer"
        >
          {theme === 'dark' ? <Sun className="w-4 h-4" /> : <Moon className="w-4 h-4" />}
        </button>
      </div>

      {/* Login card */}
      <div className="w-full max-w-sm">
        <div className="flex flex-col items-center mb-8">
          <div className="w-14 h-14 bg-accent-500 rounded-2xl flex items-center justify-center shadow-sm mb-4">
            <Plane className="w-7 h-7 text-white" strokeWidth={2} />
          </div>
          <h1 className="text-2xl font-extrabold text-text-primary">{t.appName}</h1>
          <p className="text-sm text-text-secondary mt-1">{t.heroTagline}</p>
        </div>

        <div className="bg-surface-0 border border-surface-200 rounded-2xl p-6 shadow-sm">
          <h2 className="text-lg font-bold text-text-primary text-center mb-2">
            {t.loginTitle}
          </h2>
          <p className="text-sm text-text-secondary text-center mb-6">
            {t.loginSubtitle}
          </p>

          <a
            href={api.auth.googleLoginUrl}
            className="flex items-center justify-center gap-3 w-full bg-surface-50 border border-surface-200 hover:border-accent-300 hover:bg-surface-100 rounded-xl px-4 py-3 text-sm font-semibold text-text-primary transition-colors cursor-pointer no-underline"
          >
            <svg className="w-5 h-5" viewBox="0 0 24 24">
              <path
                d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 0 1-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z"
                fill="#4285F4"
              />
              <path
                d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                fill="#34A853"
              />
              <path
                d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18A10.96 10.96 0 0 0 1 12c0 1.77.42 3.45 1.18 4.93l2.85-2.22.81-.62z"
                fill="#FBBC05"
              />
              <path
                d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                fill="#EA4335"
              />
            </svg>
            {t.loginWithGoogle}
          </a>
        </div>
      </div>
    </div>
  )
}

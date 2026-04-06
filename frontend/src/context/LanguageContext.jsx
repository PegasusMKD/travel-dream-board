import { createContext, useContext, useState, useCallback } from 'react'
import { translations } from '../data/translations'

const LanguageContext = createContext()

export function LanguageProvider({ children }) {
  const [lang, setLang] = useState(() => localStorage.getItem('lang') || 'pl')

  const toggle = useCallback(() => {
    setLang((prev) => {
      const next = prev === 'pl' ? 'en' : 'pl'
      localStorage.setItem('lang', next)
      document.documentElement.lang = next
      return next
    })
  }, [])

  const t = translations[lang]

  return (
    <LanguageContext.Provider value={{ lang, toggle, t }}>
      {children}
    </LanguageContext.Provider>
  )
}

export const useLang = () => useContext(LanguageContext)

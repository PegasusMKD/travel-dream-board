const KNOWN_CURRENCIES = new Set(['PLN', 'EUR', 'MKD'])

export function formatPrice({ price, currency, lang }) {
  if (price === null || price === undefined || price === '') return ''

  const num = Number(price)
  if (Number.isNaN(num)) return String(price)

  const locale = lang === 'pl' ? 'pl-PL' : 'en-US'
  const formatted = num.toLocaleString(locale, {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })

  if (!currency) return formatted
  if (KNOWN_CURRENCIES.has(currency)) return `${formatted} ${currency}`
  return `${formatted} ?`
}

export function isKnownCurrency(currency) {
  return KNOWN_CURRENCIES.has(currency)
}

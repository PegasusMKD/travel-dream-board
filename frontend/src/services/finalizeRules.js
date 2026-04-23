// Per-section rules for how many items can be marked as final (Selected).
// Enforced client-side before issuing PATCH requests; backend currently
// allows any count.
//
// accommodation: max 1 (one place to stay)
// transport:     max 2 (e.g. outbound + return)
// activities:    unlimited
const SECTION_FINAL_LIMIT = {
  accommodation: 1,
  transport: 2,
  activities: Infinity,
}

export function finalizeLimit(sectionType) {
  return SECTION_FINAL_LIMIT[sectionType] ?? Infinity
}

// Returns the items that should be unmarked (isFinal=false) before marking
// `targetItemId` as final, given the current items in the section. If no
// conflict exists, returns []. The strategy unmarks the oldest-finalized first.
export function itemsToUnselectFor(sectionType, items, targetItemId) {
  const limit = finalizeLimit(sectionType)
  if (!Number.isFinite(limit)) return []
  const others = items.filter((i) => i.id !== targetItemId && i.isFinal)
  // After marking target final, total = others.length + 1. We need <= limit.
  const overflow = others.length + 1 - limit
  if (overflow <= 0) return []
  // No reliable timestamp on items; just take the first N "others".
  return others.slice(0, overflow)
}

// Cycle order for the inline "advance status" affordance on item cards.
// Considering → Finalist → Booked → (stops; user uses sidebar for completed/rejected)
const ADVANCE_CYCLE = ['considering', 'finalist', 'booked']

export function nextStatus(current) {
  const idx = ADVANCE_CYCLE.indexOf(current)
  if (idx === -1 || idx >= ADVANCE_CYCLE.length - 1) return null
  return ADVANCE_CYCLE[idx + 1]
}

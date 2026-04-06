# Style Guide — Podróże Marzeń / Dream Travels

This document defines the visual language for the app. Reference it when building new components or modifying existing ones to keep the UI consistent.

---

## Typography

| Role | Font | Weight | Tailwind class |
|---|---|---|---|
| Primary font | Plus Jakarta Sans | — | `font-sans` (default) |
| Page headings | Plus Jakarta Sans | 800 | `text-3xl font-extrabold` or `text-2xl font-extrabold` |
| Section headings | Plus Jakarta Sans | 700 | `text-xl font-bold` |
| Card titles | Plus Jakarta Sans | 600 | `text-sm font-semibold` or `text-lg font-bold` |
| Body text | Plus Jakarta Sans | 400 | `text-sm` |
| Labels / metadata | Plus Jakarta Sans | 600 | `text-xs font-semibold uppercase tracking-wider` |
| Tiny metadata | Plus Jakarta Sans | 500 | `text-[11px] font-medium` |
| Monospace (refs) | System mono | — | `font-mono` |

**No serif fonts.** Playfair Display was removed for readability. Stick with Plus Jakarta Sans for everything — hierarchy comes from weight and size, not font family.

---

## Color System

Colors are defined as CSS custom properties in `src/index.css` under `@theme`. They are swapped via the `.dark` class on `<html>`.

### Accent (Rose)

The primary brand color. Used for buttons, links, selected states, and highlights.

| Token | Light | Dark | Usage |
|---|---|---|---|
| `accent-50` | `#fff1f2` | `#2a1520` | Subtle backgrounds (hover states, badges) |
| `accent-100` | `#ffe4e6` | `#3d1a2a` | Light accent backgrounds |
| `accent-200` | `#fecdd3` | — | Avatar backgrounds, rings |
| `accent-300` | `#fda4af` | — | Borders on selected items |
| `accent-400` | `#fb7185` | — | Icons, secondary accent |
| `accent-500` | `#f43f5e` | same | **Primary buttons, links, active tabs** |
| `accent-600` | `#e11d48` | same | Button hover states |
| `accent-700` | `#be123c` | same | Dark accent text |

### Surfaces

Background layers. Surface-0 is the "card" level, surface-50 is the page background.

| Token | Light | Dark | Usage |
|---|---|---|---|
| `surface-0` | `#ffffff` | `#18181b` | Cards, modals, sidebar, header |
| `surface-50` | `#f8f8fa` | `#0f0f12` | Page background |
| `surface-100` | `#f1f1f5` | `#1e1e24` | Input backgrounds, secondary fills |
| `surface-200` | `#e4e4eb` | `#2a2a33` | Borders, dividers |
| `surface-300` | `#c9c9d6` | `#3a3a48` | Scrollbar, disabled states |

### Text

| Token | Light | Dark | Usage |
|---|---|---|---|
| `text-primary` | `#1a1a2e` | `#f0f0f5` | Headings, titles, body text |
| `text-secondary` | `#555570` | `#b0b0c0` | Descriptions, notes, metadata |
| `text-tertiary` | `#8888a0` | `#808094` | Labels, placeholders, inactive tabs |
| `text-muted` | `#aaaabb` | `#606070` | Counters, timestamps, disabled text |

### Status Colors

Used via Tailwind utility classes, not custom tokens. Each status has a light bg + text pair and a dark mode variant.

| Status | Light bg | Light text | Dark bg | Dark text |
|---|---|---|---|---|
| Considering | `bg-surface-100` | `text-text-tertiary` | same | same |
| Finalist | `bg-amber-50` | `text-amber-600` | `bg-amber-950/40` | `text-amber-400` |
| Rejected | `bg-surface-100` | `text-text-muted` | same | same |
| Booked | `bg-emerald-50` | `text-emerald-600` | `bg-emerald-950/40` | `text-emerald-400` |
| Completed | `bg-blue-50` | `text-blue-600` | `bg-blue-950/40` | `text-blue-400` |

---

## Spacing & Layout

| Element | Value | Tailwind |
|---|---|---|
| Page max width | 80rem (1280px) | `max-w-7xl` |
| Page horizontal padding | 16/24/32px | `px-4 sm:px-6 lg:px-8` |
| Page bottom padding | 64px | `pb-16` |
| Card grid gap | 16px / 24px | `gap-4` (items) / `gap-6` (boards) |
| Section vertical gap | 40px | `space-y-10` |
| Card inner padding | 16px | `p-4` |
| Sidebar/modal inner padding | 24px | `p-6` |

### Grid Columns

| Context | Breakpoints |
|---|---|
| Board cards (home) | 1 col → 2 col (sm) → 3 col (lg) |
| Item cards (sections) | 1 col → 2 col (sm) → 3 col (lg) |
| Memory gallery | 2 columns → 3 columns (sm), masonry via CSS `columns` |

---

## Border Radius

| Element | Radius | Tailwind |
|---|---|---|
| Cards, modals, hero banner | 16px | `rounded-2xl` |
| Buttons, inputs, badges | 12px | `rounded-xl` |
| Small badges, status pills | 9999px (full) | `rounded-full` |
| Section icons, thumbnails | 8px | `rounded-lg` |
| Avatars | 9999px | `rounded-full` |

---

## Shadows

Keep shadows subtle. Cards use `shadow-sm` at rest and `shadow-md` or `shadow-lg` on hover. Modals use `shadow-2xl`.

| State | Tailwind |
|---|---|
| Card resting | `shadow-sm` |
| Card hover | `shadow-md` or `shadow-lg` |
| Modal / sidebar | `shadow-2xl` |
| Buttons | `shadow-sm` (optional) |

---

## Borders

- Default card border: `border border-surface-200`
- Hover border: `hover:border-accent-300`
- Selected/final items: `border-accent-400 ring-2 ring-accent-100`
- Rejected items: `border-surface-200 opacity-50`
- Dashed empty states: `border-2 border-dashed border-surface-200`
- Dividers: `border-t border-surface-200` or `border-b border-surface-200`

---

## Buttons

### Primary (accent)

```
bg-accent-500 hover:bg-accent-600 text-white font-semibold rounded-xl px-4 py-2
```

### Ghost / Secondary

```
bg-accent-50 hover:bg-accent-100 text-accent-500 font-semibold rounded-xl px-3 py-1.5
```

### Icon button

```
w-9 h-9 flex items-center justify-center rounded-lg hover:bg-surface-100 text-text-tertiary
```

### Destructive (use sparingly)

```
bg-red-500 text-white font-semibold rounded-lg px-3 py-2
```

All buttons must have `cursor-pointer` and `transition-colors`.

---

## Inputs

```
bg-surface-50 border border-surface-200 rounded-xl px-3 py-2.5 text-sm
text-text-primary placeholder:text-text-muted
focus:outline-none focus:border-accent-400 focus:ring-2 focus:ring-accent-100
```

---

## Animations & Transitions

| Element | Transition |
|---|---|
| Color/opacity changes | `transition-colors` (150ms default) |
| Card hover lift | `hover:-translate-y-1 transition-all duration-300` |
| Card image zoom | `group-hover:scale-105 transition-transform duration-500` |
| Sidebar slide-in | `translateX(100%) → 0`, 250ms ease-out |
| Lightbox fade-in | `opacity 0→1 + scale 0.97→1`, 200ms ease-out |
| Gallery image hover | `group-hover:scale-[1.03] transition-transform duration-300` |

Keep animations under 300ms for interactions. Use `duration-500` only for decorative transitions like image zoom.

---

## Component Patterns

### Cards

All cards follow this structure:
1. **Image area** — full-width, `object-cover`, with optional gradient overlay
2. **Content area** — `p-4`, contains title, badges, metadata
3. **Footer** — separated by `border-t border-surface-200 pt-3 mt-3`, contains actions

### Modals

- Backdrop: `bg-black/30 backdrop-blur-sm`
- Container: `bg-surface-0 rounded-2xl shadow-2xl border border-surface-200`
- Header: title (left) + close button (right), separated by `border-b`
- Body: `px-6 pb-6`

### Sidebar (detail panel)

- Fixed right, full height, `max-w-md`
- Same header/body pattern as modals
- Scrollable content area with `overflow-y-auto`
- Animated slide-in from right

### Empty States

- Dashed border container: `border-2 border-dashed border-surface-200 rounded-2xl py-10 text-center`
- Icon or emoji centered above text
- Muted description + accent-colored action button below

### Labels / Section Headers

Uppercase tracking labels for metadata fields:
```
text-xs font-semibold text-text-tertiary uppercase tracking-wider
```

---

## Dark Mode

Dark mode is toggled by adding/removing `.dark` on `<html>`. The CSS custom properties swap automatically.

Rules:
- Never hardcode light-mode hex values in components — always use the token classes (`bg-surface-0`, `text-text-primary`, etc.)
- For Tailwind's built-in colors used in status badges (amber, emerald, blue), use the `dark:` variant explicitly (e.g., `bg-amber-50 dark:bg-amber-950/40`)
- The `accent-50` and `accent-100` tokens have dark overrides — these are the only accent shades that change

---

## Internationalization

All user-facing strings live in `src/data/translations.js` with `pl` and `en` keys. Access via `useLang()` hook:

```jsx
const { t, lang } = useLang()
// t.addLink → "Dodaj link" or "Add link"
```

Rules:
- Never hardcode Polish or English strings in components
- Use `lang` for date formatting locale: `lang === 'pl' ? 'pl-PL' : 'en-US'`
- New components must add their strings to both `pl` and `en` in translations.js

---

## Iconography

All icons come from [Lucide React](https://lucide.dev). Default size: `w-4 h-4`. Use `strokeWidth={2}` (the default).

| Context | Size |
|---|---|
| Inline with text | `w-3.5 h-3.5` |
| Standard | `w-4 h-4` |
| Section headers | `w-4 h-4` inside a `w-8 h-8 rounded-lg bg-surface-100` container |
| Empty states | `w-8 h-8` or larger |
| Header logo | `w-5 h-5` inside a `w-9 h-9 rounded-xl bg-accent-500` container |

Do not mix icon libraries. Lucide only.

---

## File Structure

```
src/
├── components/       # Reusable UI components
├── context/          # React contexts (ThemeContext, LanguageContext)
├── data/             # Mock data, translations
├── pages/            # Route-level page components
├── index.css         # Tailwind imports, theme tokens, global styles
└── main.jsx          # App entry point with providers
```

New components go in `components/`. New pages go in `pages/`. Keep contexts in `context/`. All style tokens stay in `index.css`.

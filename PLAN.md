# Travel Dream Board — Product Specification

> A birthday gift built by a Senior Rust/Go/Java & Medior/Senior JS developer for his girlfriend.
> Hard deadline: ~1 month.

---

## 1. Context

### Why This Exists

She frequently researches trips and shares links — hotels, flights, activities — over chat with friends. These links get buried in conversation history, there's no way to compare options, and no structured record of what was considered. This app gives her a beautiful, organized home for travel planning that her friends can also collaborate on.

### About Her

| Interest | Relevance |
|---|---|
| Travel planning | Often researches trips, saves links in chats to friends — rarely has a structured place to store them |
| Shopping | Frequent online shopper; comfortable with link-based workflows |
| Pinterest / TikTok | Visual, scroll-driven aesthetic — prefers clean, beautiful interfaces over dense utilitarian ones |
| Language | Polish native; primarily uses Polish travel sites (e.g. eSky.pl for flights, Polish hotel aggregators) |

### Development Context

| Dimension | Detail |
|---|---|
| Timeline | ~1 month (hard birthday deadline) |
| Developer stack | Rust / Go / Java (backend), React 18 + Vite + TailwindCSS (frontend) |
| Existing frontend | React + Vite + TailwindCSS PoC already scaffolded |
| Infrastructure | GCP familiarity; can deploy simply for this project |
| Auth philosophy | Google OAuth only — no custom auth; friends use shareable links, no account required |
| AI API cost | Claude Haiku 4.5 — negligible at personal scale (~$0.000005/scrape, $5 free credits on signup) |

---

## 2. Core Concepts

### Trip Board
The top-level container. Each board represents a single planned (or dreamed-about) trip.

- Has a destination name, optional date range, and cover image (auto-fetched or manually set)
- Created by her; shareable with friends via a generated link
- Contains three fixed sections: Accommodation, Transport, Activities

### Sections
Every board always has exactly these three sections — no more, no less:

- 🏨 **Accommodation** — hotels, Airbnbs, apartments
- ✈️ **Transport** — flights, trains, car rentals, buses
- 📍 **Activities** — restaurants, attractions, experiences, day trips

### Items
A link saved into a section becomes an Item. Each item has:

- `url` — required
- `title` — auto-extracted via OG tags or AI fallback; always editable
- `image` — auto-extracted from OG image; falls back to a section icon
- `note` — freeform text she or friends can add (e.g. "only 2 rooms left!")
- `status` — one of: `considering` | `finalist` | `rejected` | `booked` | `completed`
- `is_final` — bool; marks the chosen option for the section
- `booking_ref` — optional confirmation number/reference (appears when status is `booked`)
- `votes` — up/down votes from collaborators (display name only, no account needed)
- `comments` — threaded comments from collaborators (display name + timestamp)

---

## 3. Finalization & Status Tracking

Each section supports marking item(s) as Final — the chosen option(s) for the trip.

| Section | Max Finals | Rationale |
|---|---|---|
| Accommodation | 1 | You book one place to stay |
| Transport | 1 | One outbound journey per trip |
| Activities | Many | Multiple places can be visited |

Items progress through a lightweight status track:

```
Considering → Finalist → Booked → Completed
                       ↘ Rejected (can happen at any stage)
```

No price tracking in v1 — explicitly out of scope.

---

## 4. Collaboration Model

She is the only user who needs an account (Google OAuth). Friends collaborate via shareable links — no registration required on their end.

| Role | How They Get Access | Permissions |
|---|---|---|
| Owner (her) | Google OAuth login | Full CRUD on boards, items, sections; manage share links |
| Collaborator (friends) | Shareable link with token | Add comments, cast votes; cannot delete or finalize items |

- Share links are board-scoped and role-scoped
- She generates them from board settings and can revoke them at any time
- Link format: `/boards/[uuid]?token=xxx`

---

## 5. Link Scraping Strategy

Scraping is **best-effort**. The goal is to reduce friction, not guarantee perfect extraction. The flow:

1. Fetch raw HTML server-side (with a real User-Agent)
2. Extract `og:title`, `og:image`, `og:description` — covers ~80% of major sites
3. Try `application/ld+json` Product/Hotel schema for structured data
4. If key fields still missing → send cleaned HTML to **Claude Haiku 4.5** for extraction
5. Return whatever was found; user can edit any field manually

### The eSky / Flight Aggregator Problem

eSky.pl and similar flight aggregators serve JavaScript-rendered **search result pages**, not individual product pages. These cannot be meaningfully scraped. The app handles this gracefully:

- The link is saved as a reference
- Title is auto-filled from OG tags (e.g. "Flights WAW → CDG | eSky")
- She adds a manual note with the relevant details (dates, price, times)

This is consistent with how she currently uses these links anyway — the link is context, not the data. The editable card makes the manual fallback frictionless.

---

## 6. Data Model

```
Board
├── id
├── name
├── destination
├── date_range (start, end) nullable
├── cover_image_url nullable
├── owner_id (Google sub)
├── created_at, updated_at
├── share_tokens[]
│   ├── token
│   ├── role: collaborator | viewer
│   ├── created_at
│   └── revoked_at nullable
└── sections: [accommodation, transport, activities]
    └── Item
        ├── id
        ├── board_id
        ├── section_type: accommodation | transport | activity
        ├── url
        ├── title
        ├── image_url nullable
        ├── note nullable
        ├── status: considering | finalist | rejected | booked | completed
        ├── is_final: bool
        ├── booking_ref nullable
        ├── created_at, updated_at
        ├── votes[]
        │   ├── display_name
        │   ├── value: up | down
        │   └── created_at
        └── comments[]
            ├── display_name
            ├── text
            └── created_at
```

---

## 7. Tech Stack

| Layer | Technology | Rationale |
|---|---|---|
| Frontend | React 18 + Vite + TailwindCSS | Already scaffolded; developer is Medior/Senior JS |
| Backend | Go (Gin or Chi) | Fast to write, good for HTTP APIs, developer's strength |
| Database | PostgreSQL | Relational; fits board/section/item/vote/comment hierarchy well |
| Migrations | golang-migrate | Already used by developer in other projects |
| Auth | Google OAuth 2.0 + JWT | One button; no custom auth complexity |
| Scraping | Go HTTP client + goquery | Lightweight HTML parsing for OG tags |
| AI Fallback | Claude Haiku 4.5 API | Extraction fallback for sites without clean OG tags |
| Deployment | GCP (Cloud Run or existing K8s) | Developer's existing infrastructure |

---

## 8. Feature Scope

### In Scope (v1 — Birthday)

- Trip board CRUD (create, rename, delete, set destination + date range)
- Fixed sections: Accommodation, Transport, Activities
- Item creation via URL paste with auto-scrape + manual edit fallback
- Item status progression: Considering → Finalist → Rejected → Booked → Completed
- Item finalization per section (1 final for Accommodation/Transport, many for Activities)
- Booking reference field on booked items
- Voting (up/down) per item — collaborators, display name only
- Commenting per item — collaborators, display name + timestamp
- Shareable links with Collaborator role (vote + comment, no delete/finalize)
- Google OAuth for owner account
- Polish language support in UI

### Explicitly Out of Scope (v1)

- Price tracking / re-scraping over time
- Email or push notifications
- Mobile app (PWA stretch goal only if time allows)
- Multiple owner accounts
- Full user registration for collaborators
- File or image uploads by users
- Map view of activities
- Viewer-only share link role (can add later, schema supports it)

---

## 9. Rough Milestone Plan

| Week | Focus | Deliverable |
|---|---|---|
| Week 1 | Backend foundations | DB schema + migrations, Go API scaffold, Google OAuth, basic board/item CRUD |
| Week 2 | Scraping + core UI | Scraper service (OG tags + AI fallback), React board view, item cards |
| Week 3 | Collaboration features | Voting, commenting, share link system, section finalization + status tracking |
| Week 4 | Polish + deployment | UI polish, error states, responsive layout, GCP deployment, final QA |

---

## 10. Scaffolding Notes for Claude Code

When scaffolding this project:

- **Monorepo structure** — `/frontend` (React + Vite) and `/backend` (Go) at the root
- **Start with the data model** — the board → section → item hierarchy is the core; get this right first before building anything else
- **golang-migrate** for DB migrations (developer already uses this pattern)
- **TailwindCSS only** — no component libraries unless explicitly chosen later
- **UI aesthetic** — clean, modern, slightly warm/feminine; think Pinterest-style card layout, not a developer dashboard. She scrolls Pinterest and TikTok daily — the visual bar is high
- **Makefile** for dev commands (developer's existing workflow)
- **No Windows-specific tooling** — developer runs Fedora + Hyprland/Wayland + NeoVim
- **Polish (pl) as the default UI language**

Start with: monorepo scaffold → DB schema → Go API routes (boards, items, votes, comments) → React board/item UI shell. Scraper and auth in week 2.

---

*Built with love. One month. Ship it. 🚀*

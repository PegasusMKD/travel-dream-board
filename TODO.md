# Travel Dream Board — TODO

## Week 1 — Backend Foundations

### Project Setup
- [X] Monorepo structure (`/frontend`, `/backend`)
- [X] Makefile with dev commands
- [X] Go module init
- [X] Railway project setup (backend service + Postgres database)
- [X] Environment variable config (`.env` + Railway vars)

### Database
- [X] golang-migrate setup
- [X] Migration: `boards`
- [X] Migration: `share_tokens`
- [X] Migration: `items`
- [X] Migration: `votes`
- [X] Migration: `comments`

### Auth
- [X] Google OAuth 2.0 integration
- [X] JWT issuance + validation middleware
- [X] Protected route middleware

### API — Boards
- [X] `POST /boards`
- [X] `GET /boards`
- [X] `GET /boards/:uuid`
- [X] `PATCH /boards/:uuid`
- [X] `DELETE /boards/:uuid`

### API — Accommodations
- [X] `POST /accomodations`
- [X] `GET /accomodations/:uuid`
- [X] `PATCH /accomodations/:uuid`
- [X] `DELETE /accomodations/:uuid`

### API — Activities
- [X] `POST /activities`
- [X] `GET /activities/:uuid`
- [X] `PATCH /activities/:uuid`
- [X] `DELETE /activities/:uuid`

### API — Transport
- [X] `POST /transport`
- [X] `GET /transport/:uuid`
- [X] `PATCH /transport/:uuid`
- [X] `DELETE /transport/:uuid`

---

## Week 2 — Scraping + Core UI

### Scraper Service
- [X] Go HTTP client with proper User-Agent
- [X] OG tag extraction (`og:title`, `og:image`, `og:description`)
- [X] `application/ld+json` schema extraction
- [X] Claude Haiku 4.5 AI fallback
- [X] Graceful partial-result handling (return what we have)

### Frontend Setup
- [X] Vite + React 18 + TailwindCSS scaffold
- [X] Routing setup
- [X] Google OAuth login page
- [X] API client / fetch wrapper

### Frontend — Boards
- [X] Board list view (dashboard)
- [X] Create board flow
- [X] Board detail view (3 sections)
- [X] Item card component
- [X] URL paste → scrape → editable card flow
- [X] Manual field editing on item card

---

## Week 3 — Collaboration Features

### API — Collaboration
- [X] `POST /boards/:uuid/share-tokens` (generate share link)
- [X] `DELETE /boards/:uuid/share-tokens/:token` (revoke)
- [X] Share token validation middleware
- [X] `POST /votes`
- [X] `PATCH /votes/:uuid` (change vote)
- [X] `DELETE /votes/:uuid`
- [X] `POST /comments`
- [X] `PATCH /comments/:uuid`
- [X] `DELETE /comments/:uuid`

### Frontend — Collaboration
- [ ] Share link generation UI
- [ ] Collaborator view (token-gated)
- [X] Voting UI on item cards
- [X] Comment thread on item cards
- [ ] Display name prompt for collaborators

### Frontend — Finalization & Status
- [ ] Status badge + progression UI
- [ ] Finalize item action (per section rules)
- [ ] Booking reference field (shown when status = booked)
- [ ] Rejected state styling

---

## Week 4 — Polish + Deployment

### UI Polish
- [ ] Empty states (no boards, no items, no votes)
- [ ] Loading states + skeletons
- [ ] Error states (scrape failed, network error)
- [ ] Responsive layout (mobile-friendly)
- [X] Polish UI language strings

### Deployment
- [X] Backend deployed to Railway
- [X] Postgres provisioned on Railway
- [X] Migrations run on Railway
- [X] Frontend deployed (Railway static or Vercel)
- [X] Environment variables set in Railway dashboard
- [X] Google OAuth redirect URIs updated for production domain
- [ ] Smoke test on production

### Final QA
- [ ] Happy path: create board → add items → share → friend votes → finalize → book
- [ ] Edge cases: scrape failure, revoked share link, duplicate votes
- [X] Cross-browser check

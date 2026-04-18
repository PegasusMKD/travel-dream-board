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
- [ ] Migration: `share_tokens`
- [X] Migration: `items`
- [X] Migration: `votes`
- [X] Migration: `comments`

### Auth
- [ ] Google OAuth 2.0 integration
- [ ] JWT issuance + validation middleware
- [ ] Protected route middleware

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
- [ ] Go HTTP client with proper User-Agent
- [ ] OG tag extraction (`og:title`, `og:image`, `og:description`)
- [ ] `application/ld+json` schema extraction
- [ ] Claude Haiku 4.5 AI fallback
- [ ] Graceful partial-result handling (return what we have)
- [ ] `POST /scrape` endpoint

### Frontend Setup
- [ ] Vite + React 18 + TailwindCSS scaffold
- [ ] Routing setup
- [ ] Google OAuth login page
- [ ] API client / fetch wrapper

### Frontend — Boards
- [ ] Board list view (dashboard)
- [ ] Create board flow
- [ ] Board detail view (3 sections)
- [ ] Item card component
- [ ] URL paste → scrape → editable card flow
- [ ] Manual field editing on item card

---

## Week 3 — Collaboration Features

### API — Collaboration
- [ ] `POST /boards/:uuid/share-tokens` (generate share link)
- [ ] `DELETE /boards/:uuid/share-tokens/:token` (revoke)
- [ ] Share token validation middleware
- [X] `POST /votes`
- [X] `PATCH /votes/:uuid` (change vote)
- [X] `DELETE /votes/:uuid`
- [X] `POST /comments`
- [X] `PATCH /comments/:uuid`
- [X] `DELETE /comments/:uuid`

### Frontend — Collaboration
- [ ] Share link generation UI
- [ ] Collaborator view (token-gated)
- [ ] Voting UI on item cards
- [ ] Comment thread on item cards
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
- [ ] Polish UI language strings

### Deployment
- [ ] Backend deployed to Railway
- [ ] Postgres provisioned on Railway
- [ ] Migrations run on Railway
- [ ] Frontend deployed (Railway static or Vercel)
- [ ] Environment variables set in Railway dashboard
- [ ] Google OAuth redirect URIs updated for production domain
- [ ] Smoke test on production

### Final QA
- [ ] Happy path: create board → add items → share → friend votes → finalize → book
- [ ] Edge cases: scrape failure, revoked share link, duplicate votes
- [ ] Cross-browser check

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
- [ ] `POST /boards`
- [ ] `GET /boards` (owner's boards)
- [ ] `GET /boards/:id`
- [ ] `PATCH /boards/:id`
- [ ] `DELETE /boards/:id`

### API — Items
- [ ] `POST /boards/:id/items`
- [ ] `PATCH /items/:id`
- [ ] `DELETE /items/:id`
- [ ] `PATCH /items/:id/finalize`
- [ ] `PATCH /items/:id/status`

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
- [ ] `POST /boards/:id/share-tokens` (generate share link)
- [ ] `DELETE /boards/:id/share-tokens/:token` (revoke)
- [ ] Share token validation middleware
- [ ] `POST /items/:id/votes`
- [ ] `DELETE /items/:id/votes` (change vote)
- [ ] `POST /items/:id/comments`
- [ ] `DELETE /comments/:id`

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

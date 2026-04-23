# Switch from like/dislike to 1–5 star ratings, with backend-enforced uniqueness

## Context

Today each item (accommodation, transport, activity) shows thumbs-up / thumbs-down counts. The frontend implements toggle-off and switch-vote logic locally (`frontend/src/pages/BoardDetail.jsx:97-118`), but the backend has no uniqueness guard — `POST /votes/` happily inserts duplicates. Counts can therefore be inflated by rapid clicks, direct API calls, or any stale-state race.

Rather than just hardening the binary vote, we're pivoting to a **1–5 star rating**. The entities are travel items — a review/rating paradigm fits them better than thumbs, and the schema already supports it: the column is literally `rank int`. Stars express preference intensity, which is more useful when ranking 5 candidate hotels than "7 likes vs 2 dislikes".

Backend must own uniqueness. Frontend checks stay as a UX optimization (immediate toggle-off/switch), but the DB enforces one-vote-per-(user, item).

## Scope

**In scope**
- Replace binary vote UI with 1–5 star ratings.
- Backend aggregation: switch from `likes`/`dislikes` counts to `avg_rating` + `rating_count`.
- Add UNIQUE constraint `(voted_by, voted_on_, voted_on_uuid)` with a migration that dedupes any pre-existing duplicates first.
- Keep the three-endpoint API (create / update / delete).
- Sidebar gets a dedicated "Ratings" section below Comments: star picker, "Clear rating" button, list of all raters with their star counts. Card becomes read-only (avg stars + count).
- Click-same-star-to-clear **AND** a "Clear rating" button both work.

**Out of scope (follow-ups)**
- Vote handler reads `user_uuid` from request body — user will verify separately.
- No half-stars. `rank` stays `int`, values 1–5.

## Execution checklist

### Backend
- [x] **1.** Create migration `000011_votes_unique_and_rank_range.up.sql` + `.down.sql` (dedupe, UNIQUE, CHECK 1..5, clamp existing data).
- [x] **2.** Rewrite aggregation queries in `accomodations.sql`, `activities.sql`, `transport.sql` — replace `likes`/`dislikes` with `avg_rating` + `rating_count`.
- [x] **3.** Regenerate sqlc.
- [x] **4.** Update `accomodations/model.go`, `activities/model.go`, `transport/model.go` — swap `Likes`/`Dislikes` for `AvgRating`/`RatingCount`.
- [x] **5.** Fix any callers of removed fields (none existed — grep clean, `go build` passes).

### Frontend
- [x] **6.** Update `mappers.js` — `mapVote` returns raw `rank`, no `value`. Items expose `avgRating`/`ratingCount`.
- [x] **7.** Update `BoardDetail.jsx` — `handleVote(item, rank)` with create/switch/toggle-off branches. Add `handleClearVote(item)`.
- [x] **8.** Update `ItemCard.jsx` — read-only stars + count. Remove `onVote`.
- [x] **9.** Update `ItemDetailSidebar.jsx` — remove old Votes block; add Ratings section BELOW Comments with interactive star picker, Clear button, flat rater list.
- [x] **10.** Update `translations.js` — add rating strings. `votes` key removed (no longer referenced).

### Verification
- [x] **11.** Backend compiles (`go build ./...` clean).
- [x] **12.** Frontend builds (`npm run build` clean, 510ms).
- [ ] **13.** Hand-test the UI — see manual steps below.

## Manual verification steps (for you to run)

1. **Apply the migration** against your dev DB:
   ```
   cd backend && migrate -path ./sql/migrations -database "<your-dsn>" up
   ```
   (Or your usual migrate command.) Watch for errors from duplicate-vote cleanup — expected to be none if you were the only tester.

2. **Start servers:**
   ```
   cd backend && go run ./cmd/...       # or however you normally run the API
   cd frontend && npm run dev
   ```

3. **UI checks:**
   - Open a board with items that used to have thumbs: cards should now show either "No ratings yet" (if no votes existed) or the migrated rating (old upvotes became 1 star, old downvotes also became 1 star per the clamp).
   - Click an item → sidebar opens. Scroll to bottom: "Ratings" section below Comments.
   - Click 4★ in the picker → card updates to "4.0 (1)", sidebar shows you at 4★.
   - Click 2★ → switches to 2★ (card shows "2.0 (1)").
   - Click 2★ again → vote cleared, card shows "No ratings yet".
   - Click 3★, then the "Clear rating" button → cleared.
   - Second browser/user rates 5★ → card shows avg, sidebar shows both raters sorted high-to-low.
   - Card clicks on stars should do nothing (card is read-only now).

4. **(Optional) Unique constraint smoke test** via curl/Postman: call `POST /api/v1/votes/` twice with the same `user_uuid`/`voted_on_uuid` — second call returns a 500 from the DB unique violation. This is expected; the frontend's branching means real users won't hit it, but the DB is authoritative.

## Files to modify

**Backend**
- `backend/sql/migrations/000006_votes_unique_per_user_per_item.up.sql` (new) + `.down.sql` (new)
- `backend/sql/queries/accomodations.sql`
- `backend/sql/queries/activities.sql`
- `backend/sql/queries/transport.sql`
- `backend/internal/accomodations/model.go`
- `backend/internal/activities/model.go`
- `backend/internal/transport/model.go`
- `backend/internal/db/*.sql.go` (regenerated)

**Frontend**
- `frontend/src/services/mappers.js`
- `frontend/src/pages/BoardDetail.jsx`
- `frontend/src/components/ItemCard.jsx`
- `frontend/src/components/ItemDetailSidebar.jsx`
- `frontend/src/data/translations.js`

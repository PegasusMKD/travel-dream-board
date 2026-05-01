# Travel Dream Board — POLISHING TODO

## Polishing Features

### Expand Transport - Legs (Times + Locations) ✅
- [x] Add fields on Transport for the exact departure and arrival times (outbound + inbound, each with from/to locations + datetimes — `outbound_departing_location`, `outbound_arriving_location`, `outbound_departing_at`, `outbound_arriving_at`, and `inbound_*` mirrors). Migration `000013_add_transport_legs`.
- [x] Expand the LLM tools to be able to parse that information from an image (or as a fallback) — shared `extractionSchema()` used by both text-fallback and image-extraction tool calls; prompts hint at round-trip vs one-way.
- [x] Expand the scraping logic to include those fields if possible? — OG/JSON-LD don't carry leg data in practice, so values come from Claude (text or image). `ScrapeResult` carries parsed `*time.Time` through to the transport service.
- [x] Should be manually editable — `ItemDetailSidebar` gets an Outbound + Return panel (4 inputs each: From, To, Departure, Arrival) when `sectionType === 'transport'`, plus a compact read-only summary with route arrow + formatted date/time.

### Expand Activities - Start and End Times ✅
- [x] Add fields on Activities for the exact start and end times — `start_at` / `end_at` timestamptz columns added in migration `000014_add_activity_times`. Plumbed through the Activity model, repository, and service via the shared `utility.TimestamptzFromTime` / `utility.TimePtrFromTimestamptz` helpers.
- [x] Expand the LLM tools to be able to parse that information from an image (or as a fallback) — `extractionSchema()` extended with `start_at` / `end_at` properties; both `fallbackToClaude` (text) and `ExtractFromImage` prompts now hint at activities/events.
- [x] Expand the scraping logic to include those fields if possible? — OG/JSON-LD don't carry event times in practice, so values come from Claude. `ScrapeResult` carries parsed `*time.Time` through to the activity service. Wall-clock parser promoted to `utility.ParseWallClockTime` and reused by transport too.
- [x] Should be manually enterable if you open the edit window — `ItemDetailSidebar` gets an `ActivityTimeEditor` (Start + End `datetime-local` inputs) when `sectionType === 'activities'`, plus an `ActivityTimeSummary` read-only row.

### Expand Transport - Total Journey Length ✅
- [x] Add fields on Transport for total journey length — `outbound_duration_minutes` / `inbound_duration_minutes` int columns added in migration `000015_add_transport_durations`. Stored as integer minutes (e.g. 225 = 3h 45m); covers layovers/multi-stop. Plumbed through transport model, repo, and service.
- [x] Expand the LLM tools to be able to parse that information from an image (or as a fallback) — `extractionSchema()` extended with two integer properties; both text and image prompts now ask for total leg duration in minutes alongside locations and datetimes.
- [x] Expand the scraping logic to include those fields if possible? — OG/JSON-LD don't carry leg durations, so values come from Claude. `ScrapeResult` carries `*int32` durations through to the transport service. Sidebar shows formatted "Xh Ym" in the leg summary and exposes a number input in the editor.


### Expand Activities - Location ✅
- [x] Add fields on Activities for exact location of the event — `location` text column added in migration `000016_add_activity_location`. Plumbed through Activity model, repo, and service.
- [x] Expand the LLM tools to be able to parse that information from an image (or as a fallback) — `extractionSchema()` extended with a `location` property; both `fallbackToClaude` (text) and `ExtractFromImage` prompts ask for venue + city or full address.
- [x] Expand the scraping logic to include those fields if possible? — OG/JSON-LD don't carry venue strings reliably, so values come from Claude. `ScrapeResult` carries `*string Location` through to the activity service.
- [x] Should be manually editable if you open the edit window — `ItemDetailSidebar`'s `ActivityTimeEditor` gains a Location text input; `ActivityTimeSummary` renders it on a separate row with a `MapPinned` icon when present.

### Expand All Items - Price & Description ✅
- [x] Add fields to denote an extracted price from the link — `price numeric(12,2)` + `description text` columns added to `accomodations`, `transport`, and `activities` in migration `000017_add_item_price_currency_description`. Plumbed through all three model/repo/service stacks via shared `utility.NumericFromString` / `utility.NumericToString` helpers.
- [x] Expand the LLM tools to be able to parse that information from an image (or as a fallback) — `extractionSchema()` extended with `price` (numeric string) and `currency` properties; both `fallbackToClaude` (text) and `ExtractFromImage` prompts now ask for total price + currency.
- [x] Expand the scraping logic to include those fields if possible? — OG/JSON-LD don't carry reliable price data, so values come from Claude. `ScrapeResult` now carries `*string Price` + `*db.CurrencyCode Currency`; services propagate them and the existing `Description` to each item.
- [x] Should be manually editable — `ItemDetailSidebar` gets a shared `PriceEditor` block (number input + currency `<select>` + description `<textarea>`) rendered for all section types, plus a compact `PriceSummary` read-only card with `Tag` icon.
- [x] The "price" field should support at least 3 currencies, PLN (zloty), EUR (euro), MKD (Macedonian Denar) — new `currency_code` Postgres enum with values `PLN`, `EUR`, `MKD`, `unknown`. `utility.ParseCurrencyCode` normalizes Claude output (handles `zł`, `€`, `ден`, full names) and falls back to `unknown` for unrecognized symbols.
- [x] It should also make it clear what currency it is displaying the information in — new `frontend/src/utils/formatPrice.js` formats as `'1 234,56 PLN'` (locale-aware grouping); `unknown` renders the amount with a `?` suffix and an italic "Unknown" hint in the sidebar summary. `ItemCard` shows an accent-tinted price chip when present.
- [x] Potentially also add a "currency" column? — done as part of the enum work; stored as `NullCurrencyCode`, exposed as `*db.CurrencyCode` on each item model and passed as a string through the JSON API.

### Implement "Memories" feature
- [ ] Create migration, queries and a module for Memories
- [ ] Create local folder for `/memories`
- [ ] Create a volume with path `/memories` on Railway
- [ ] Enable upload of `memories`
- [ ] Enforce access to memories, make sure whoever is trying to access them (either through URL or through the FE) has access to the Board

### Determine what code is not being used and can potentially be removed
- [ ] Create a Markdown file detailing the code that seems to not be used anymore
- [ ] Provide an explanation why it isn't needed anymore
- [ ] Provide an idea as to what to do with it

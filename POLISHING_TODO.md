# Travel Dream Board — POLISHING TODO

## Polishing Features

### Expand Transport - Legs (Times + Locations) ✅
- [x] Add fields on Transport for the exact departure and arrival times (outbound + inbound, each with from/to locations + datetimes — `outbound_departing_location`, `outbound_arriving_location`, `outbound_departing_at`, `outbound_arriving_at`, and `inbound_*` mirrors). Migration `000013_add_transport_legs`.
- [x] Expand the LLM tools to be able to parse that information from an image (or as a fallback) — shared `extractionSchema()` used by both text-fallback and image-extraction tool calls; prompts hint at round-trip vs one-way.
- [x] Expand the scraping logic to include those fields if possible? — OG/JSON-LD don't carry leg data in practice, so values come from Claude (text or image). `ScrapeResult` carries parsed `*time.Time` through to the transport service.
- [x] Should be manually editable — `ItemDetailSidebar` gets an Outbound + Return panel (4 inputs each: From, To, Departure, Arrival) when `sectionType === 'transport'`, plus a compact read-only summary with route arrow + formatted date/time.

### Expand Activities - Start and End Times
- [ ] Add fields on Activities for the exact start and end times
- [ ] Expand the LLM tools to be able to parse that information from an image (or as a fallback)
- [ ] Expand the scraping logic to include those fields if possible?
- [ ] Should be manually enterable if you open the edit window

### Expand Transport - Total Journey Length
- [ ] Add fields on Transport for total journey length
- [ ] Expand the LLM tools to be able to parse that information from an image (or as a fallback)
- [ ] Expand the scraping logic to include those fields if possible?


### Expand Activities - Location
- [ ] Add fields on Activities for exact location of the event
- [ ] Expand the LLM tools to be able to parse that information from an image (or as a fallback)
- [ ] Expand the scraping logic to include those fields if possible?
- [ ] Should be manually editable if you open the edit window

### Expand All Items - Price & Description
- [ ] Add fields to denote an extracted price from the link
- [ ] Expand the LLM tools to be able to parse that information from an image (or as a fallback)
- [ ] Expand the scraping logic to include those fields if possible?
- [ ] Should be manually editable
- [ ] The "price" field should support at least 3 currencies, PLN (zloty), EUR (euro), MKD (Macedonian Denar)
- [ ] It should also make it clear what currency it is displaying the information in
- [ ] Potentially also add a "currency" column?

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

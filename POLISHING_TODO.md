# Travel Dream Board — TODO

## Polishing Features

### Expand Transport & Activities - Start and End Times
- [ ] Add fields on Transport & Activities for exact start_time and end_time (first takeoff and last arrival back)
- [ ] Expand the LLM tools to be able to parse that information from an image (or as a fallback)
- [ ] Expand the scraping logic to include those fields if possible?


### Expand Transport - Total Journey Length
- [ ] Add fields on Transport for total journey length
- [ ] Expand the LLM tools to be able to parse that information from an image (or as a fallback)
- [ ] Expand the scraping logic to include those fields if possible?


### Expand Activities - Location
- [ ] Add fields on Activities for exact location of the event
- [ ] Expand the LLM tools to be able to parse that information from an image (or as a fallback)
- [ ] Expand the scraping logic to include those fields if possible?

### Expand All Items - Price & Description
- [ ] Add fields to denote an extracted price from the link
- [ ] Expand the LLM tools to be able to parse that information from an image (or as a fallback)
- [ ] Expand the scraping logic to include those fields if possible?

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

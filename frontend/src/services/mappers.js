export function mapBoardSummary(b) {
  return {
    id: b.uuid,
    name: b.name,
    destination: b.location_name,
    dateRange: b.starts_at && b.lasts_until
      ? { start: b.starts_at, end: b.lasts_until }
      : null,
    coverImage: b.thumbnail_url || '',
    memories: [],
    sections: {
      accommodation: [],
      transport: [],
      activities: [],
    },
  }
}

export function mapAggregatedBoard(b) {
  return {
    ...mapBoardSummary(b),
    memories: [],
    sections: {
      accommodation: (b.accomodations || []).map(mapItem),
      transport: (b.transport || []).map(mapItem),
      activities: (b.activities || []).map(mapItem),
    },
  }
}

function mapItem(item) {
  return {
    id: item.uuid,
    url: item.url,
    title: item.title,
    image: item.image_url || null,
    note: item.notes || '',
    status: item.status,
    isFinal: item.selected,
    bookingRef: item.booking_reference || null,
    votes: (item.votes || []).map(mapVote),
    comments: (item.comments || []).map(mapComment),
  }
}

function mapVote(v) {
  return {
    displayName: shortName(v.user_uuid),
    value: v.rank > 0 ? 'up' : 'down',
  }
}

function mapComment(c) {
  return {
    displayName: shortName(c.user_uuid),
    text: c.content,
    createdAt: null,
  }
}

function shortName(uuid) {
  if (!uuid) return 'User'
  return uuid.slice(0, 8)
}

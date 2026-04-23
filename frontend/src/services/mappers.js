const ITEM_SECTION_TO_BACKEND_KEY = {
  accommodation: 'accomodations',
  transport: 'transport',
  activities: 'activities',
}

const BACKEND_COMMENT_TARGET = {
  accommodation: 'accomodation',
  transport: 'transport',
  activities: 'activities',
}

const BACKEND_VOTE_TARGET = BACKEND_COMMENT_TARGET

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
    counts: {
      accommodation: b.accomodations_count || 0,
      transport: b.transport_count || 0,
      activities: b.activities_count || 0,
    },
  }
}

export function mapAggregatedBoard(b) {
  return {
    ...mapBoardSummary(b),
    sections: {
      accommodation: (b.accomodations || []).map(mapItem),
      transport: (b.transport || []).map(mapItem),
      activities: (b.activities || []).map(mapItem),
    },
  }
}

export function mapItem(item) {
  return {
    id: item.uuid,
    url: item.url,
    title: item.title,
    boardUuid: item.board_uuid,
    image: item.image_url || null,
    note: item.notes || '',
    status: item.status,
    isFinal: item.selected,
    bookingRef: item.booking_reference || null,
    avgRating: item.avg_rating || 0,
    ratingCount: item.rating_count || 0,
    votes: (item.votes || []).map(mapVote),
    comments: (item.comments || []).map(mapComment),
  }
}

export function mapVote(v) {
  return {
    id: v.uuid,
    userUuid: v.user_uuid,
    displayName: v.user_name || shortName(v.user_uuid),
    rank: v.rank,
  }
}

export function mapComment(c) {
  return {
    id: c.uuid,
    userUuid: c.user_uuid,
    displayName: c.user_name || shortName(c.user_uuid),
    text: c.content,
    createdAt: null,
  }
}

export function toBackendBoardPayload(form) {
  return {
    name: form.name,
    location_name: form.destination,
    starts_at: form.startDate ? new Date(form.startDate).toISOString() : null,
    lasts_until: form.endDate ? new Date(form.endDate).toISOString() : null,
    thumbnail_url: form.coverImage || null,
  }
}

export function toBackendItemPayload(item) {
  return {
    uuid: item.id,
    url: item.url,
    title: item.title,
    board_uuid: item.boardUuid,
    image_url: item.image || null,
    notes: item.note || null,
    status: item.status,
    selected: !!item.isFinal,
    booking_reference: item.bookingRef || null,
  }
}

export function backendCommentTarget(sectionType) {
  return BACKEND_COMMENT_TARGET[sectionType]
}

export function backendVoteTarget(sectionType) {
  return BACKEND_VOTE_TARGET[sectionType]
}

export function sectionToItemApi(api, sectionType) {
  const key = ITEM_SECTION_TO_BACKEND_KEY[sectionType]
  return api[key]
}

function shortName(uuid) {
  if (!uuid) return 'User'
  return uuid.slice(0, 8)
}

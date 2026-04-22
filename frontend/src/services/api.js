const API_BASE = '/api/v1'

let onUnauthorized = null

export function setUnauthorizedHandler(handler) {
  onUnauthorized = handler
}

export class AuthError extends Error {
  constructor(message) {
    super(message)
    this.name = 'AuthError'
  }
}

async function request(path, options = {}) {
  const res = await fetch(`${API_BASE}${path}`, {
    credentials: 'include',
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  })

  if (res.status === 401) {
    if (onUnauthorized) onUnauthorized()
    throw new AuthError('Unauthorized')
  }

  if (!res.ok) {
    const body = await res.json().catch(() => ({}))
    throw new Error(body.error || `Request failed: ${res.status}`)
  }

  if (res.status === 204 || res.headers.get('content-length') === '0') {
    return null
  }

  const text = await res.text()
  if (!text) return null
  return JSON.parse(text)
}

function itemEndpoints(base) {
  return {
    create: (url, boardUuid) => {
      const qs = new URLSearchParams({ url, boardUuid })
      return request(`/${base}/?${qs.toString()}`, { method: 'POST' })
    },
    update: (uuid, data) => request(`/${base}/${uuid}`, {
      method: 'PATCH',
      body: JSON.stringify(data),
    }),
    delete: (uuid) => request(`/${base}/${uuid}`, { method: 'DELETE' }),
  }
}

export const api = {
  auth: {
    me: () => request('/auth/me'),
    logout: () => request('/auth/logout', { method: 'POST' }),
    googleLoginUrl: `${API_BASE}/auth/google/login`,
  },
  boards: {
    getAll: () => request('/boards/'),
    getById: (uuid) => request(`/boards/${uuid}`),
    create: (data) => request('/boards/', { method: 'POST', body: JSON.stringify(data) }),
    update: (uuid, data) => request(`/boards/${uuid}`, { method: 'PATCH', body: JSON.stringify(data) }),
    delete: (uuid) => request(`/boards/${uuid}`, { method: 'DELETE' }),
    shareTokens: {
      list: (boardUuid) => request(`/boards/${boardUuid}/share-tokens/`),
      create: (boardUuid) => request(`/boards/${boardUuid}/share-tokens/`, { method: 'POST' }),
      delete: (boardUuid, token) => request(`/boards/${boardUuid}/share-tokens/${token}`, { method: 'DELETE' }),
    },
  },
  accomodations: itemEndpoints('accomodations'),
  activities: itemEndpoints('activities'),
  transport: itemEndpoints('transport'),
  comments: {
    create: (data) => request('/comments/', { method: 'POST', body: JSON.stringify(data) }),
    update: (uuid, content) => request(`/comments/${uuid}`, {
      method: 'PATCH',
      body: JSON.stringify({ content }),
    }),
    delete: (uuid) => request(`/comments/${uuid}`, { method: 'DELETE' }),
  },
  votes: {
    create: (data) => request('/votes/', { method: 'POST', body: JSON.stringify(data) }),
    update: (uuid, rank) => request(`/votes/${uuid}`, {
      method: 'PATCH',
      body: JSON.stringify({ rank }),
    }),
    delete: (uuid) => request(`/votes/${uuid}`, { method: 'DELETE' }),
  },
}

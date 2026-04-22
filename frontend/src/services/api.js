const API_BASE = '/api/v1'

let onUnauthorized = null

export function setUnauthorizedHandler(handler) {
  onUnauthorized = handler
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

  return res.json()
}

export class AuthError extends Error {
  constructor(message) {
    super(message)
    this.name = 'AuthError'
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
  },
}

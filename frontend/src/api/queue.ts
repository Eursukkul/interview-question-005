const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080'

type CurrentQueueResponse = {
  current_queue: string
  issued_at?: string
}

type QueueNumberResponse = {
  queue_number: string
  issued_at?: string
}

type ErrorResponse = {
  error: string
}

export type QueueSnapshot = {
  queueNumber: string
  issuedAt: string | null
}

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    headers: {
      'Content-Type': 'application/json'
    },
    ...options
  })

  const payload = (await response.json()) as T | ErrorResponse
  if (!response.ok) {
    const message =
      typeof payload === 'object' && payload !== null && 'error' in payload
        ? String(payload.error)
        : 'Request failed'
    throw new Error(message)
  }
  return payload as T
}

export async function getCurrentQueue(): Promise<QueueSnapshot> {
  const payload = await request<CurrentQueueResponse>('/api/queue/current')
  return {
    queueNumber: payload.current_queue ?? '00',
    issuedAt: payload.issued_at ?? null
  }
}

export async function getNextQueue(): Promise<QueueSnapshot> {
  const payload = await request<QueueNumberResponse>('/api/queue/next', { method: 'POST' })
  return {
    queueNumber: payload.queue_number,
    issuedAt: payload.issued_at ?? null
  }
}

export async function resetQueue(): Promise<QueueSnapshot> {
  const payload = await request<QueueNumberResponse>('/api/queue/reset', { method: 'POST' })
  return {
    queueNumber: payload.queue_number ?? '00',
    issuedAt: payload.issued_at ?? null
  }
}

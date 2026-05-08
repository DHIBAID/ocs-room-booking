const apiBase = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'

type ApiOptions = RequestInit & { token?: string | null }

export const apiFetch = async (path: string, options: ApiOptions = {}) => {
  const { token, ...rest } = options
  const res = await fetch(`${apiBase}${path}`, {
    ...rest,
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...(rest.headers ?? {}),
    },
  })

  const data = await res.json().catch(() => ({}))
  if (!res.ok) {
    throw new Error(data.error || 'Request failed')
  }
  return data
}

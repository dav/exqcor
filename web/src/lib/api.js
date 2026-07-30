// Thin fetch wrapper for the Exqcor API. All endpoints are same-origin.

export class ApiError extends Error {
  constructor(status, message) {
    super(message)
    this.status = status
  }
}

async function request(method, path, body) {
  const opts = { method, headers: {} }
  if (body !== undefined) {
    opts.headers['Content-Type'] = 'application/json'
    opts.body = JSON.stringify(body)
  }
  const res = await fetch(path, opts)
  const isJSON = res.headers.get('Content-Type')?.includes('application/json')
  const data = isJSON ? await res.json() : null
  if (!res.ok) {
    throw new ApiError(res.status, data?.error || res.statusText)
  }
  return data
}

export const api = {
  get: (path) => request('GET', path),
  post: (path, body) => request('POST', path, body ?? {}),
  put: (path, body) => request('PUT', path, body),
  del: (path) => request('DELETE', path),
}

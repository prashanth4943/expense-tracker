// api.js — thin wrapper around fetch for all backend calls.
// All amounts are sent as strings to preserve precision.

const BASE = import.meta.env.DEV ? '' : ''

async function request(method, path, body) {
  const opts = {
    method,
    headers: { 'Content-Type': 'application/json' },
  }
  if (body !== undefined) {
    opts.body = JSON.stringify(body)
  }
  const res = await fetch(BASE + path, opts)
  const data = await res.json()
  if (!res.ok) {
    throw new Error(data.error || `HTTP ${res.status}`)
  }
  return data
}

export function createExpense(payload) {
  return request('POST', '/expenses', payload)
}

export function listExpenses(params = {}) {
  const qs = new URLSearchParams()
  if (params.category) qs.set('category', params.category)
  if (params.sort)     qs.set('sort', params.sort)
  const query = qs.toString() ? '?' + qs.toString() : ''
  return request('GET', `/expenses${query}`)
}

export function listCategories() {
  return request('GET', '/categories')
}
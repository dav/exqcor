// EventSource wrapper. The browser handles reconnection and Last-Event-ID
// replay; we add named-event dispatch and a connection-state callback so
// pages can show a "reconnecting" notice on flaky venue Wi-Fi.

export function subscribe(handlers, onState) {
  const es = new EventSource('/api/events')
  for (const [type, fn] of Object.entries(handlers)) {
    es.addEventListener(type, (e) => {
      let data = null
      try {
        data = e.data ? JSON.parse(e.data) : null
      } catch {
        // ignore malformed frames
      }
      fn(data)
    })
  }
  es.onopen = () => onState?.('connected')
  es.onerror = () => onState?.('reconnecting')
  return () => es.close()
}

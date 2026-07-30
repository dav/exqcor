// Countdown helpers. The server sends ends_at plus its own clock
// (server_now); we trust the server clock and only use the local clock to
// advance between polls, so phone clock skew never matters.

export function clockOffset(serverNowISO) {
  return new Date(serverNowISO).getTime() - Date.now()
}

export function secondsLeft(endsAtISO, offset, nowMs) {
  return Math.ceil((new Date(endsAtISO).getTime() - (nowMs + offset)) / 1000)
}

export function formatClock(totalSeconds) {
  const s = Math.max(0, totalSeconds)
  return `${Math.floor(s / 60)}:${String(s % 60).padStart(2, '0')}`
}

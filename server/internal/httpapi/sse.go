package httpapi

import (
	"net/http"
	"strconv"
	"time"
)

// handleEvents streams real-time events over SSE. Events carry only ids —
// clients follow up through role-guarded endpoints — so the stream itself is
// safe for every role. Targeted events (your_turn) are filtered per device.
func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeErr(w, http.StatusInternalServerError, "streaming not supported")
		return
	}

	// Identify the audience member (if any) for targeted delivery.
	var audienceID int64
	if c, err := r.Cookie(deviceCookie); err == nil {
		if scriptID := s.activeScriptID(); scriptID != 0 {
			if m, err := s.Store.AudienceByDevice(scriptID, c.Value); err == nil && m != nil {
				audienceID = m.ID
			}
		}
	}

	sub, unsubscribe := s.Runtime.Hub.Subscribe(audienceID)
	defer unsubscribe()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	// Replay anything the client missed while reconnecting.
	if last, err := strconv.ParseInt(r.Header.Get("Last-Event-ID"), 10, 64); err == nil {
		for _, frame := range s.Runtime.Hub.Replay(last, sub) {
			w.Write(frame)
		}
	}
	w.Write([]byte(": connected\n\n"))
	flusher.Flush()

	heartbeat := time.NewTicker(25 * time.Second)
	defer heartbeat.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case frame := <-sub.Ch:
			w.Write(frame)
			flusher.Flush()
		case <-heartbeat.C:
			w.Write([]byte(": ping\n\n"))
			flusher.Flush()
		}
	}
}

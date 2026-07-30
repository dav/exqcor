package show

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Event is a real-time notification. Events deliberately carry only ids and
// types — never script content — so any role may receive any event; clients
// re-fetch details through role-guarded REST endpoints.
type Event struct {
	Seq  int64          `json:"seq"`
	Type string         `json:"type"`
	Data map[string]any `json:"data,omitempty"`

	// TargetAudienceID restricts delivery to one audience member (your_turn).
	TargetAudienceID int64 `json:"-"`
}

const replayBuffer = 256

type Subscriber struct {
	AudienceID int64
	Ch         chan []byte
}

// Hub fans events out to SSE subscribers, keeping a small replay ring so
// reconnecting phones (Last-Event-ID) don't miss a call.
type Hub struct {
	mu     sync.Mutex
	subs   map[*Subscriber]struct{}
	seq    int64
	recent []Event
}

func NewHub() *Hub {
	return &Hub{subs: map[*Subscriber]struct{}{}}
}

func (h *Hub) Subscribe(audienceID int64) (*Subscriber, func()) {
	sub := &Subscriber{AudienceID: audienceID, Ch: make(chan []byte, 16)}
	h.mu.Lock()
	h.subs[sub] = struct{}{}
	h.mu.Unlock()
	return sub, func() {
		h.mu.Lock()
		delete(h.subs, sub)
		h.mu.Unlock()
	}
}

// Broadcast queues an event for every (matching) subscriber. Slow clients
// are skipped rather than blocked on; EventSource reconnect + replay covers
// them.
func (h *Hub) Broadcast(typ string, data map[string]any, targetAudienceID int64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.seq++
	ev := Event{Seq: h.seq, Type: typ, Data: data, TargetAudienceID: targetAudienceID}
	h.recent = append(h.recent, ev)
	if len(h.recent) > replayBuffer {
		h.recent = h.recent[len(h.recent)-replayBuffer:]
	}
	frame := sseFrame(ev)
	for sub := range h.subs {
		if !deliverable(ev, sub) {
			continue
		}
		select {
		case sub.Ch <- frame:
		default:
		}
	}
}

// Replay returns frames for events after seq that the subscriber may see.
func (h *Hub) Replay(afterSeq int64, sub *Subscriber) [][]byte {
	h.mu.Lock()
	defer h.mu.Unlock()
	var out [][]byte
	for _, ev := range h.recent {
		if ev.Seq > afterSeq && deliverable(ev, sub) {
			out = append(out, sseFrame(ev))
		}
	}
	return out
}

func deliverable(ev Event, sub *Subscriber) bool {
	return ev.TargetAudienceID == 0 || ev.TargetAudienceID == sub.AudienceID
}

func sseFrame(ev Event) []byte {
	body, _ := json.Marshal(ev.Data)
	return fmt.Appendf(nil, "id: %d\nevent: %s\ndata: %s\n\n", ev.Seq, ev.Type, body)
}

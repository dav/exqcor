// Package show manages the live-performance state machine: active writer
// turns, their enforcement timers, and (via callbacks) real-time broadcast.
package show

import (
	"log"
	"sync"
	"time"

	"github.com/dav/exqcor/server/internal/store"
)

type Runtime struct {
	mu     sync.Mutex
	st     *store.Store
	timers map[int64]*time.Timer
	Hub    *Hub
}

func New(st *store.Store) *Runtime {
	return &Runtime{st: st, timers: map[int64]*time.Timer{}, Hub: NewHub()}
}

// StartTurn begins the next writer turn in a section and arms its
// enforcement timer.
func (r *Runtime) StartTurn(sectionID int64, writerName string, audienceMemberID *int64, durationSeconds int) (store.SubSection, error) {
	turn, err := r.st.StartTurn(sectionID, writerName, audienceMemberID, durationSeconds)
	if err != nil {
		return turn, err
	}
	r.armTimer(turn)
	r.Hub.Broadcast("turn_started", map[string]any{
		"section_id": turn.SectionID, "turn_id": turn.ID, "ends_at": turn.EndsAt,
	}, 0)
	return turn, nil
}

// EndTurn completes a turn now (writer done, admin cut, or timer expiry).
func (r *Runtime) EndTurn(turnID int64) error {
	turn, err := r.st.GetTurn(turnID)
	if err != nil {
		return err
	}
	if turn.CompletedAt != nil {
		return nil
	}
	if err := r.st.CompleteTurn(turnID); err != nil {
		return err
	}
	r.mu.Lock()
	if t, ok := r.timers[turnID]; ok {
		t.Stop()
		delete(r.timers, turnID)
	}
	r.mu.Unlock()
	// A writer who came from the audience queue is done with their visit.
	if memberID, err := r.st.TurnAudienceMemberID(turnID); err == nil && memberID != nil {
		if err := r.st.SetAudienceStatus(*memberID, "done"); err != nil {
			log.Printf("mark audience member %d done: %v", *memberID, err)
		}
		r.Hub.Broadcast("queue_changed", nil, 0)
	}
	r.Hub.Broadcast("turn_ended", map[string]any{
		"section_id": turn.SectionID, "turn_id": turn.ID,
	}, 0)
	return nil
}

// armTimer schedules server-side enforcement: the turn auto-completes at
// ends_at plus the grace period, whether or not the writer's device is alive.
func (r *Runtime) armTimer(turn store.SubSection) {
	if turn.EndsAt == nil {
		return
	}
	ends, err := time.Parse(time.RFC3339, *turn.EndsAt)
	if err != nil {
		log.Printf("turn %d: bad ends_at %q", turn.ID, *turn.EndsAt)
		return
	}
	d := time.Until(ends.Add(store.GraceSeconds * time.Second))
	if d < 0 {
		d = 0
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.timers[turn.ID] = time.AfterFunc(d, func() {
		if err := r.EndTurn(turn.ID); err != nil {
			log.Printf("auto-end turn %d: %v", turn.ID, err)
		}
	})
}

// Recover re-arms in-flight turns after a restart: expired ones complete
// immediately, running ones get their timer back. Yanking the power cord
// mid-show costs nothing but the in-flight keystroke.
func (r *Runtime) Recover() error {
	turns, err := r.st.ActiveTurns()
	if err != nil {
		return err
	}
	for _, turn := range turns {
		r.armTimer(turn)
		log.Printf("recovered active turn %d in section %d", turn.ID, turn.SectionID)
	}
	return nil
}

package httpapi

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/dav/exqcor/server/internal/store"
)

const keyActiveScript = "active_script_id"

func (s *Server) activeScriptID() int64 {
	v, _ := s.Store.Setting(keyActiveScript)
	id, _ := strconv.ParseInt(v, 10, 64)
	return id
}

// handleOpenShow arms a script as tonight's show.
func (s *Server) handleOpenShow(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ScriptID int64 `json:"script_id"`
	}
	if !readJSON(w, r, &req) {
		return
	}
	if _, err := s.Store.GetScript(req.ScriptID); err != nil {
		storeErr(w, err)
		return
	}
	if err := s.Store.SetSetting(keyActiveScript, strconv.FormatInt(req.ScriptID, 10)); err != nil {
		storeErr(w, err)
		return
	}
	s.broadcastShowState()
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

type sectionState struct {
	store.Section
	TurnCount  int   `json:"turn_count"`
	ActiveTurn *int64 `json:"active_turn_id"`
}

// showState assembles the live view of tonight's show; the fields are safe
// for every role (no line content beyond counts).
func (s *Server) showState() (map[string]any, error) {
	id := s.activeScriptID()
	if id == 0 {
		return map[string]any{"open": false}, nil
	}
	sc, err := s.Store.GetScript(id)
	if err != nil {
		if store.IsNotFound(err) {
			return map[string]any{"open": false}, nil
		}
		return nil, err
	}
	secs, err := s.Store.ListSections(id)
	if err != nil {
		return nil, err
	}
	states := []sectionState{}
	for _, sec := range secs {
		n, err := s.Store.TurnCount(sec.ID)
		if err != nil {
			return nil, err
		}
		st := sectionState{Section: sec, TurnCount: n}
		if turn, err := s.Store.ActiveTurn(sec.ID); err != nil {
			return nil, err
		} else if turn != nil {
			st.ActiveTurn = &turn.ID
		}
		states = append(states, st)
	}
	return map[string]any{
		"open":       true,
		"script":     sc,
		"sections":   states,
		"server_now": time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func (s *Server) handleShow(w http.ResponseWriter, r *http.Request) {
	state, err := s.showState()
	if err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, state)
}

// writeAccess splits writers into privileged devices (admin, station) and
// audience members writing from their own phone.
func (s *Server) writeAccess(r *http.Request) (privileged bool, m *store.AudienceView) {
	switch s.role(r) {
	case RoleAdmin, RoleStation:
		return true, nil
	case RoleAudience:
		return false, s.member(r)
	}
	return false, nil
}

// calledTo reports whether the member has been called to this section.
func calledTo(m *store.AudienceView, sectionID int64) bool {
	return m != nil && m.Status == "called" && m.CalledSectionID != nil && *m.CalledSectionID == sectionID
}

// isTurnWriter reports whether the member is the one writing this turn.
func (s *Server) isTurnWriter(m *store.AudienceView, turnID int64) bool {
	if m == nil {
		return false
	}
	id, err := s.Store.TurnAudienceMemberID(turnID)
	return err == nil && id != nil && *id == m.ID
}

// handleStartTurn begins the next writer turn — from a station, or from a
// called audience member's own phone.
func (s *Server) handleStartTurn(w http.ResponseWriter, r *http.Request) {
	sectionID, _ := pathID(r, "id")
	var req struct {
		WriterName string `json:"writer_name"`
	}
	if !readJSON(w, r, &req) {
		return
	}
	privileged, m := s.writeAccess(r)
	var audienceID *int64
	if !privileged {
		if !calledTo(m, sectionID) {
			writeErr(w, http.StatusForbidden, "it isn't your turn yet — watch for your number to be called")
			return
		}
		audienceID = &m.ID
		if req.WriterName == "" {
			req.WriterName = m.Name
		}
	}
	script, err := s.Store.SectionScript(sectionID)
	if err != nil {
		storeErr(w, err)
		return
	}
	turn, err := s.Runtime.StartTurn(sectionID, req.WriterName, audienceID, script.WritingSeconds)
	if errors.Is(err, store.ErrTurnActive) {
		writeErr(w, http.StatusConflict, err.Error())
		return
	}
	if err != nil {
		storeErr(w, err)
		return
	}
	if audienceID != nil {
		if err := s.Store.SetAudienceStatus(*audienceID, "writing"); err != nil {
			storeErr(w, err)
			return
		}
		s.Runtime.Hub.Broadcast("queue_changed", nil, 0)
	}
	writeJSON(w, http.StatusCreated, turn)
}

// handleTurnContext returns everything a writing device needs to render the
// active turn: the previous writer's last line ONLY, the pickable
// characters, the deadline, and the writer's own lines so far.
func (s *Server) handleTurnContext(w http.ResponseWriter, r *http.Request) {
	sectionID, _ := pathID(r, "id")
	sec, err := s.Store.GetSection(sectionID)
	if err != nil {
		storeErr(w, err)
		return
	}
	turn, err := s.Store.ActiveTurn(sectionID)
	if err != nil {
		storeErr(w, err)
		return
	}
	if privileged, m := s.writeAccess(r); !privileged {
		ownTurn := turn != nil && s.isTurnWriter(m, turn.ID)
		if !calledTo(m, sectionID) && !ownTurn {
			writeErr(w, http.StatusForbidden, "this act's writing view isn't open to you right now")
			return
		}
	}
	resp := map[string]any{
		"section":    sec,
		"turn":       nil,
		"server_now": time.Now().UTC().Format(time.RFC3339),
	}
	if turn != nil {
		lastLine, err := s.Store.LastVisibleLine(sectionID, turn.Ordering)
		if err != nil {
			storeErr(w, err)
			return
		}
		chars, err := s.Store.OnStageCharacters(sectionID)
		if err != nil {
			storeErr(w, err)
			return
		}
		myLines, err := s.Store.TurnLines(turn.ID)
		if err != nil {
			storeErr(w, err)
			return
		}
		resp["turn"] = turn
		resp["last_line"] = lastLine
		resp["characters"] = chars
		resp["my_lines"] = myLines
		resp["grace_seconds"] = store.GraceSeconds
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleAddLine(w http.ResponseWriter, r *http.Request) {
	turnID, _ := pathID(r, "id")
	var req struct {
		CharacterID int64  `json:"character_id"`
		Text        string `json:"text"`
	}
	if !readJSON(w, r, &req) {
		return
	}
	if len(req.Text) == 0 || len(req.Text) > 2000 {
		writeErr(w, http.StatusBadRequest, "line must be between 1 and 2000 characters")
		return
	}
	if privileged, m := s.writeAccess(r); !privileged && !s.isTurnWriter(m, turnID) {
		writeErr(w, http.StatusForbidden, "this isn't your turn")
		return
	}
	line, err := s.Store.AddLine(turnID, req.CharacterID, req.Text)
	if errors.Is(err, store.ErrTurnOver) {
		writeErr(w, http.StatusConflict, err.Error())
		return
	}
	if err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, line)
}

// handleEndTurn completes a turn: the writer tapping "I'm done", or an admin
// cutting it off.
func (s *Server) handleEndTurn(w http.ResponseWriter, r *http.Request) {
	turnID, _ := pathID(r, "id")
	if privileged, m := s.writeAccess(r); !privileged && !s.isTurnWriter(m, turnID) {
		writeErr(w, http.StatusForbidden, "this isn't your turn")
		return
	}
	if err := s.Runtime.EndTurn(turnID); err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// handleProgram returns tonight's playbill: show info and cast, safe for
// every role — no script content.
func (s *Server) handleProgram(w http.ResponseWriter, r *http.Request) {
	scriptID := s.activeScriptID()
	if scriptID == 0 {
		writeJSON(w, http.StatusOK, map[string]any{"open": false})
		return
	}
	sc, err := s.Store.GetScript(scriptID)
	if err != nil {
		storeErr(w, err)
		return
	}
	chars, err := s.Store.ListCharacters(scriptID)
	if err != nil {
		storeErr(w, err)
		return
	}
	actors, err := s.Store.ListActors()
	if err != nil {
		storeErr(w, err)
		return
	}
	secs, err := s.Store.ListSections(scriptID)
	if err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"open":       true,
		"script":     sc,
		"characters": chars,
		"actors":     actors,
		"sections":   secs,
	})
}

// handleFullScript returns the assembled script for reading and printing.
// Admin, stations, and actors only — the audience hears it on stage instead.
func (s *Server) handleFullScript(w http.ResponseWriter, r *http.Request) {
	id, _ := pathID(r, "id")
	fs, err := s.Store.FullScriptView(id)
	if err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, fs)
}

func (s *Server) handleScriptStats(w http.ResponseWriter, r *http.Request) {
	id, _ := pathID(r, "id")
	stats, err := s.Store.ScriptStats(id)
	if err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

func (s *Server) broadcastShowState() {
	s.Runtime.Hub.Broadcast("show_state", nil, 0)
}

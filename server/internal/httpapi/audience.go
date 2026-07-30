package httpapi

import (
	"errors"
	"net/http"

	"github.com/dav/exqcor/server/internal/store"
)

// member resolves the request's audience member in the active show, if any.
func (s *Server) member(r *http.Request) *store.AudienceView {
	c, err := r.Cookie(deviceCookie)
	if err != nil || c.Value == "" {
		return nil
	}
	scriptID := s.activeScriptID()
	if scriptID == 0 {
		return nil
	}
	m, err := s.Store.AudienceByDevice(scriptID, c.Value)
	if err != nil {
		return nil
	}
	return m
}

// handleAudienceJoin gives this device a number in tonight's queue
// (idempotent per device).
func (s *Server) handleAudienceJoin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if !readJSON(w, r, &req) {
		return
	}
	scriptID := s.activeScriptID()
	if scriptID == 0 {
		writeErr(w, http.StatusConflict, "no show is open yet")
		return
	}
	device := s.ensureDeviceCookie(w, r)
	m, err := s.Store.JoinAudience(scriptID, device, req.Name)
	if err != nil {
		storeErr(w, err)
		return
	}
	s.Runtime.Hub.Broadcast("queue_changed", nil, 0)
	writeJSON(w, http.StatusOK, s.meResponse(m))
}

func (s *Server) handleAudienceMe(w http.ResponseWriter, r *http.Request) {
	m := s.member(r)
	if m == nil {
		writeJSON(w, http.StatusOK, map[string]any{"joined": false})
		return
	}
	writeJSON(w, http.StatusOK, s.meResponse(*m))
}

func (s *Server) meResponse(m store.AudienceView) map[string]any {
	resp := map[string]any{"joined": true, "member": m}
	if sc, err := s.Store.GetScript(m.ScriptID); err == nil {
		resp["station_mode"] = sc.StationMode
		resp["script_title"] = sc.Title
	}
	if m.CalledSectionID != nil {
		if sec, err := s.Store.GetSection(*m.CalledSectionID); err == nil {
			resp["called_section"] = sec
		}
	}
	return resp
}

// --- admin queue management ---

func (s *Server) handleListAudience(w http.ResponseWriter, r *http.Request) {
	scriptID := s.activeScriptID()
	if scriptID == 0 {
		writeJSON(w, http.StatusOK, []store.AudienceView{})
		return
	}
	ms, err := s.Store.ListAudience(scriptID)
	if err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ms)
}

func (s *Server) handleCallNext(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SectionID int64 `json:"section_id"`
	}
	if !readJSON(w, r, &req) {
		return
	}
	scriptID := s.activeScriptID()
	if scriptID == 0 {
		writeErr(w, http.StatusConflict, "no show is open")
		return
	}
	m, err := s.Store.CallNext(scriptID, req.SectionID)
	if errors.Is(err, store.ErrNoOneWaiting) {
		writeErr(w, http.StatusConflict, err.Error())
		return
	}
	if err != nil {
		storeErr(w, err)
		return
	}
	sec, _ := s.Store.GetSection(req.SectionID)
	s.Runtime.Hub.Broadcast("queue_changed", nil, 0)
	s.Runtime.Hub.Broadcast("your_turn", map[string]any{
		"section_id": req.SectionID, "section_name": sec.Name,
	}, m.ID)
	writeJSON(w, http.StatusOK, m)
}

func (s *Server) handleAudienceStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := pathID(r, "id")
	status := r.PathValue("status")
	switch status {
	case "skip":
		status = "skipped"
	case "requeue":
		status = "waiting"
	default:
		writeErr(w, http.StatusBadRequest, "unknown action")
		return
	}
	if err := s.Store.SetAudienceStatus(id, status); err != nil {
		storeErr(w, err)
		return
	}
	m, err := s.Store.GetAudienceMember(id)
	if err != nil {
		storeErr(w, err)
		return
	}
	s.Runtime.Hub.Broadcast("queue_changed", nil, 0)
	writeJSON(w, http.StatusOK, m)
}

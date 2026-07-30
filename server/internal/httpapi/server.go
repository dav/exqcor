package httpapi

import (
	"encoding/json"
	"io/fs"
	"log"
	"net/http"

	"github.com/dav/exqcor/server/internal/show"
	"github.com/dav/exqcor/server/internal/store"
)

type Server struct {
	Store   *store.Store
	Runtime *show.Runtime
	Mux     *http.ServeMux
	BaseURL func() string // http://<lan-ip>:<port>, recomputed on demand
}

func New(st *store.Store, rt *show.Runtime, webFS fs.FS, baseURL func() string) *Server {
	s := &Server{Store: st, Runtime: rt, Mux: http.NewServeMux(), BaseURL: baseURL}
	s.routes(webFS)
	return s
}

func (s *Server) routes(webFS fs.FS) {
	s.Mux.Handle("GET /", http.FileServerFS(webFS))
	s.Mux.HandleFunc("GET /api/health", s.handleHealth)
	s.Mux.HandleFunc("GET /api/netinfo", s.handleNetInfo)
	s.Mux.HandleFunc("POST /api/netinfo", s.requireAdmin(s.handleSetNetInfo))
	s.Mux.HandleFunc("GET /api/qr.png", s.handleQRPNG)
	s.Mux.HandleFunc("GET /j/{token}", s.handleJoin)
	s.Mux.HandleFunc("POST /api/session/admin", s.handleAdminLogin)
	s.Mux.HandleFunc("GET /api/session", s.handleSession)

	admin := func(pattern string, h http.HandlerFunc) {
		s.Mux.HandleFunc(pattern, s.requireAdmin(h))
	}
	admin("GET /api/scripts", s.handleListScripts)
	admin("POST /api/scripts", s.handleCreateScript)
	admin("GET /api/scripts/{id}", s.handleGetScript)
	admin("PUT /api/scripts/{id}", s.handleUpdateScript)
	admin("DELETE /api/scripts/{id}", s.handleDeleteScript)
	admin("POST /api/scripts/{id}/duplicate", s.handleDuplicateScript)

	admin("GET /api/scripts/{id}/sections", s.handleListSections)
	admin("POST /api/scripts/{id}/sections", s.handleCreateSection)
	admin("PUT /api/sections/{id}", s.handleUpdateSection)
	admin("DELETE /api/sections/{id}", s.handleDeleteSection)
	admin("GET /api/sections/{id}/priming-line", s.handleGetPrimingLine)
	admin("POST /api/sections/{id}/priming-line", s.handleSetPrimingLine)

	admin("GET /api/scripts/{id}/characters", s.handleListCharacters)
	admin("POST /api/scripts/{id}/characters", s.handleCreateCharacter)
	admin("PUT /api/characters/{id}", s.handleUpdateCharacter)
	admin("DELETE /api/characters/{id}", s.handleDeleteCharacter)

	admin("GET /api/sections/{id}/character-sections", s.handleListCharacterSections)
	admin("POST /api/sections/{id}/character-sections", s.handleSetCharacterSection)

	admin("GET /api/sections/{id}/props", s.handleListProps)
	admin("POST /api/sections/{id}/props", s.handleCreateProp)
	admin("DELETE /api/props/{id}", s.handleDeleteProp)

	// Run-of-show + writing flow. Writing endpoints allow station devices
	// always, and audience members only for their own called turn (checked
	// inside the handlers).
	anyRole := func(pattern string, h http.HandlerFunc) {
		s.Mux.HandleFunc(pattern, s.requireRole(h, RoleStation, RoleAudience))
	}
	admin("POST /api/show/open", s.handleOpenShow)
	s.Mux.HandleFunc("GET /api/show", s.handleShow)
	s.Mux.HandleFunc("GET /api/program", s.handleProgram)
	anyRole("POST /api/sections/{id}/turns", s.handleStartTurn)
	anyRole("GET /api/sections/{id}/turn", s.handleTurnContext)
	anyRole("POST /api/turns/{id}/lines", s.handleAddLine)
	anyRole("POST /api/turns/{id}/done", s.handleEndTurn)

	// Script output (admin, stations, actors — not the audience: spoilers).
	s.Mux.HandleFunc("GET /api/scripts/{id}/full", s.requireRole(s.handleFullScript, RoleStation, RoleActor))
	admin("GET /api/scripts/{id}/stats", s.handleScriptStats)

	// Real-time + audience queue.
	s.Mux.HandleFunc("GET /api/events", s.handleEvents)
	s.Mux.HandleFunc("POST /api/audience/join", s.requireRole(s.handleAudienceJoin, RoleAudience))
	s.Mux.HandleFunc("GET /api/audience/me", s.handleAudienceMe)
	admin("GET /api/audience", s.handleListAudience)
	admin("POST /api/queue/call-next", s.handleCallNext)
	admin("POST /api/queue/{id}/{status}", s.handleAudienceStatus)

	admin("GET /api/tokens", s.handleTokens)
	admin("POST /api/tokens/regenerate", s.handleRegenerateToken)

	admin("GET /api/actors", s.handleListActors)
	admin("POST /api/actors", s.handleCreateActor)
	admin("PUT /api/actors/{id}", s.handleUpdateActor)
	admin("DELETE /api/actors/{id}", s.handleDeleteActor)
}

func (s *Server) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Mux.ServeHTTP(w, r)
	})
}

// --- JSON helpers ---

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writeJSON: %v", err)
	}
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func readJSON(w http.ResponseWriter, r *http.Request, v any) bool {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return false
	}
	return true
}

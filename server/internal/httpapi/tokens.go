package httpapi

import (
	"net/http"

	"github.com/dav/exqcor/server/internal/store"
)

// handleTokens returns the join URLs behind each role's QR code.
func (s *Server) handleTokens(w http.ResponseWriter, r *http.Request) {
	base := s.BaseURL()
	out := map[string]string{}
	for role, key := range map[string]string{
		RoleAudience: keyTokenAudience,
		RoleStation:  keyTokenStation,
		RoleActor:    keyTokenActor,
	} {
		v, err := s.Store.Setting(key)
		if err != nil {
			storeErr(w, err)
			return
		}
		out[role] = base + "/j/" + v
	}
	writeJSON(w, http.StatusOK, out)
}

// handleRegenerateToken rotates one role's join token. Old printed QR codes
// for that role stop working, and cookies issued from them die with it.
func (s *Server) handleRegenerateToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Role string `json:"role"`
	}
	if !readJSON(w, r, &req) {
		return
	}
	key := roleBasisKey(req.Role)
	if key == "" || req.Role == RoleAdmin {
		writeErr(w, http.StatusBadRequest, "role must be audience, station, or actor")
		return
	}
	if err := s.Store.SetSetting(key, store.RandomToken()); err != nil {
		storeErr(w, err)
		return
	}
	s.handleTokens(w, r)
}

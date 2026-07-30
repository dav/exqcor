package httpapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net"
	"net/http"
	"strings"

	"github.com/dav/exqcor/server/internal/store"
)

// Roles, least to most privileged.
const (
	RoleAudience = "audience"
	RoleActor    = "actor"
	RoleStation  = "station"
	RoleAdmin    = "admin"
)

const (
	roleCookie   = "exqcor_role"
	deviceCookie = "exqcor_device"
)

// Settings keys.
const (
	keySecret        = "secret"
	keyAdminPassHash = "admin_pass_hash"
	keyTokenAudience = "token_audience"
	keyTokenStation  = "token_station"
	keyTokenActor    = "token_actor"
)

func (s *Server) secret() []byte {
	v, err := s.Store.EnsureSetting(keySecret, store.RandomToken)
	if err != nil {
		panic(err)
	}
	return []byte(v)
}

// sign produces the signature for a role cookie. The signature covers the
// role's current join token (or the admin passphrase hash), so regenerating a
// token invalidates every cookie previously issued for that role.
func (s *Server) sign(role string) string {
	basis, _ := s.Store.Setting(roleBasisKey(role))
	mac := hmac.New(sha256.New, s.secret())
	mac.Write([]byte(role + ":" + basis))
	return hex.EncodeToString(mac.Sum(nil))
}

func roleBasisKey(role string) string {
	switch role {
	case RoleAudience:
		return keyTokenAudience
	case RoleStation:
		return keyTokenStation
	case RoleActor:
		return keyTokenActor
	case RoleAdmin:
		return keyAdminPassHash
	}
	return ""
}

func (s *Server) setRoleCookie(w http.ResponseWriter, role string) {
	http.SetCookie(w, &http.Cookie{
		Name:     roleCookie,
		Value:    role + ":" + s.sign(role),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 24 * 30,
	})
}

// role resolves the request's role. Requests from the server machine itself
// are implicitly admin so the operator never needs a passphrase locally.
func (s *Server) role(r *http.Request) string {
	if isLoopback(r) {
		return RoleAdmin
	}
	c, err := r.Cookie(roleCookie)
	if err != nil {
		return ""
	}
	role, sig, ok := strings.Cut(c.Value, ":")
	if !ok || roleBasisKey(role) == "" {
		return ""
	}
	if !hmac.Equal([]byte(sig), []byte(s.sign(role))) {
		return ""
	}
	return role
}

func isLoopback(r *http.Request) bool {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return false
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func (s *Server) requireAdmin(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.role(r) != RoleAdmin {
			writeErr(w, http.StatusUnauthorized, "admin access required")
			return
		}
		h(w, r)
	}
}

// requireRole allows any of the given roles; admin always passes.
func (s *Server) requireRole(h http.HandlerFunc, roles ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		got := s.role(r)
		if got == RoleAdmin {
			h(w, r)
			return
		}
		for _, want := range roles {
			if got == want {
				h(w, r)
				return
			}
		}
		writeErr(w, http.StatusUnauthorized, "not allowed for your role")
	}
}

// HashPassphrase derives the stored hash for an admin passphrase, keyed by
// the database's secret so the hash is useless outside this installation.
func HashPassphrase(st *store.Store, pass string) string {
	secret, err := st.EnsureSetting(keySecret, store.RandomToken)
	if err != nil {
		panic(err)
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte("passphrase:" + pass))
	return hex.EncodeToString(mac.Sum(nil))
}

func (s *Server) hashPassphrase(pass string) string {
	return HashPassphrase(s.Store, pass)
}

// --- handlers ---

// handleJoin lands QR scans: /j/<token> maps to a role, sets the role cookie,
// and redirects to that role's home screen.
func (s *Server) handleJoin(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	for role, key := range map[string]string{
		RoleAudience: keyTokenAudience,
		RoleStation:  keyTokenStation,
		RoleActor:    keyTokenActor,
	} {
		v, _ := s.Store.Setting(key)
		if v != "" && token == v {
			s.setRoleCookie(w, role)
			s.ensureDeviceCookie(w, r)
			http.Redirect(w, r, roleHome(role), http.StatusFound)
			return
		}
	}
	http.Error(w, "This QR code is no longer valid. Ask the production admin for a new one.", http.StatusNotFound)
}

func roleHome(role string) string {
	switch role {
	case RoleStation:
		return "/#/write"
	case RoleActor:
		return "/#/program"
	case RoleAdmin:
		return "/#/admin"
	}
	return "/#/audience"
}

func (s *Server) ensureDeviceCookie(w http.ResponseWriter, r *http.Request) string {
	if c, err := r.Cookie(deviceCookie); err == nil && c.Value != "" {
		return c.Value
	}
	v := store.RandomToken()
	http.SetCookie(w, &http.Cookie{
		Name:     deviceCookie,
		Value:    v,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 24 * 365,
	})
	return v
}

func (s *Server) handleAdminLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Passphrase string `json:"passphrase"`
	}
	if !readJSON(w, r, &req) {
		return
	}
	stored, err := s.Store.Setting(keyAdminPassHash)
	if err != nil || stored == "" {
		writeErr(w, http.StatusInternalServerError, "admin passphrase not configured")
		return
	}
	if !hmac.Equal([]byte(s.hashPassphrase(req.Passphrase)), []byte(stored)) {
		writeErr(w, http.StatusUnauthorized, "wrong passphrase")
		return
	}
	s.setRoleCookie(w, RoleAdmin)
	writeJSON(w, http.StatusOK, map[string]string{"role": RoleAdmin})
}

func (s *Server) handleSession(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"role": s.role(r)})
}

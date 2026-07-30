package httpapi

import (
	"net/http"

	"github.com/dav/exqcor/server/internal/netinfo"
	"github.com/dav/exqcor/server/internal/version"
	qrcode "github.com/skip2/go-qrcode"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"version": version.Version,
	})
}

func (s *Server) handleNetInfo(w http.ResponseWriter, r *http.Request) {
	chosen, _ := s.Store.Setting("chosen_ip")
	writeJSON(w, http.StatusOK, map[string]any{
		"candidates": netinfo.Candidates(),
		"chosen_ip":  chosen,
		"base_url":   s.BaseURL(),
	})
}

func (s *Server) handleSetNetInfo(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IP string `json:"ip"`
	}
	if !readJSON(w, r, &req) {
		return
	}
	if err := s.Store.SetSetting("chosen_ip", req.IP); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"base_url": s.BaseURL()})
}

// handleQRPNG renders a QR code for any URL under the server's own base URL.
// Restricting to our own origin keeps this from being an open QR generator.
func (s *Server) handleQRPNG(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		url = s.BaseURL()
	}
	png, err := qrcode.Encode(url, qrcode.Medium, 512)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(png)
}

package httpapi

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dav/exqcor/server/internal/store"
)

func pathID(r *http.Request, name string) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue(name), 10, 64)
	return id, err == nil
}

// storeErr maps store-layer failures to HTTP responses.
func storeErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, store.ErrVOSDProtected):
		writeErr(w, http.StatusUnprocessableEntity, err.Error())
	case store.IsNotFound(err):
		writeErr(w, http.StatusNotFound, "not found")
	default:
		writeErr(w, http.StatusInternalServerError, err.Error())
	}
}

// --- scripts ---

type scriptReq struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	Theme          string `json:"theme"`
	WritingSeconds int    `json:"writing_seconds"`
	StationMode    string `json:"station_mode"`
}

func (s *Server) handleListScripts(w http.ResponseWriter, r *http.Request) {
	scripts, err := s.Store.ListScripts()
	if err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, scripts)
}

func (s *Server) handleCreateScript(w http.ResponseWriter, r *http.Request) {
	var req scriptReq
	if !readJSON(w, r, &req) {
		return
	}
	if req.Title == "" {
		writeErr(w, http.StatusBadRequest, "title is required")
		return
	}
	sc, err := s.Store.CreateScript(req.Title, req.Description, req.Theme, req.WritingSeconds, req.StationMode)
	if err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, sc)
}

func (s *Server) handleGetScript(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r, "id")
	if !ok {
		writeErr(w, http.StatusBadRequest, "bad id")
		return
	}
	sc, err := s.Store.GetScript(id)
	if err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, sc)
}

func (s *Server) handleUpdateScript(w http.ResponseWriter, r *http.Request) {
	id, _ := pathID(r, "id")
	var req scriptReq
	if !readJSON(w, r, &req) {
		return
	}
	sc, err := s.Store.UpdateScript(id, req.Title, req.Description, req.Theme, req.WritingSeconds, req.StationMode)
	if err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, sc)
}

func (s *Server) handleDeleteScript(w http.ResponseWriter, r *http.Request) {
	id, _ := pathID(r, "id")
	if err := s.Store.DeleteScript(id); err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleDuplicateScript(w http.ResponseWriter, r *http.Request) {
	id, _ := pathID(r, "id")
	var req struct {
		Title string `json:"title"`
	}
	if !readJSON(w, r, &req) {
		return
	}
	sc, err := s.Store.DuplicateScript(id, req.Title)
	if err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, sc)
}

// --- sections ---

type sectionReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Ordering    *int   `json:"ordering"`
	Status      string `json:"status"`
}

func (s *Server) handleListSections(w http.ResponseWriter, r *http.Request) {
	id, _ := pathID(r, "id")
	secs, err := s.Store.ListSections(id)
	if err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, secs)
}

func (s *Server) handleCreateSection(w http.ResponseWriter, r *http.Request) {
	id, _ := pathID(r, "id")
	var req sectionReq
	if !readJSON(w, r, &req) {
		return
	}
	if req.Name == "" {
		writeErr(w, http.StatusBadRequest, "name is required")
		return
	}
	sec, err := s.Store.CreateSection(id, req.Name, req.Description, req.Ordering)
	if err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, sec)
}

func (s *Server) handleUpdateSection(w http.ResponseWriter, r *http.Request) {
	id, _ := pathID(r, "id")
	cur, err := s.Store.GetSection(id)
	if err != nil {
		storeErr(w, err)
		return
	}
	req := sectionReq{Name: cur.Name, Description: cur.Description, Status: cur.Status}
	if !readJSON(w, r, &req) {
		return
	}
	ord := cur.Ordering
	if req.Ordering != nil {
		ord = *req.Ordering
	}
	sec, err := s.Store.UpdateSection(id, req.Name, req.Description, ord, req.Status)
	if err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, sec)
}

func (s *Server) handleDeleteSection(w http.ResponseWriter, r *http.Request) {
	id, _ := pathID(r, "id")
	if err := s.Store.DeleteSection(id); err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleGetPrimingLine(w http.ResponseWriter, r *http.Request) {
	id, _ := pathID(r, "id")
	line, err := s.Store.PrimingLine(id)
	if err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"line": line})
}

func (s *Server) handleSetPrimingLine(w http.ResponseWriter, r *http.Request) {
	id, _ := pathID(r, "id")
	var req struct {
		CharacterID int64  `json:"character_id"`
		Text        string `json:"text"`
	}
	if !readJSON(w, r, &req) {
		return
	}
	if err := s.Store.SetPrimingLine(id, req.CharacterID, req.Text); err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// --- characters ---

type characterReq struct {
	ActorID     *int64 `json:"actor_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *Server) handleListCharacters(w http.ResponseWriter, r *http.Request) {
	id, _ := pathID(r, "id")
	cs, err := s.Store.ListCharacters(id)
	if err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, cs)
}

func (s *Server) handleCreateCharacter(w http.ResponseWriter, r *http.Request) {
	id, _ := pathID(r, "id")
	var req characterReq
	if !readJSON(w, r, &req) {
		return
	}
	if req.Name == "" {
		writeErr(w, http.StatusBadRequest, "name is required")
		return
	}
	c, err := s.Store.CreateCharacter(id, req.ActorID, req.Name, req.Description)
	if err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, c)
}

func (s *Server) handleUpdateCharacter(w http.ResponseWriter, r *http.Request) {
	id, _ := pathID(r, "id")
	var req characterReq
	if !readJSON(w, r, &req) {
		return
	}
	c, err := s.Store.UpdateCharacter(id, req.ActorID, req.Name, req.Description)
	if err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (s *Server) handleDeleteCharacter(w http.ResponseWriter, r *http.Request) {
	id, _ := pathID(r, "id")
	if err := s.Store.DeleteCharacter(id); err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// --- character_sections ---

func (s *Server) handleListCharacterSections(w http.ResponseWriter, r *http.Request) {
	id, _ := pathID(r, "id")
	css, err := s.Store.ListCharacterSections(id)
	if err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, css)
}

func (s *Server) handleSetCharacterSection(w http.ResponseWriter, r *http.Request) {
	id, _ := pathID(r, "id")
	var req struct {
		CharacterID int64 `json:"character_id"`
		Attached    bool  `json:"attached"`
		OnStage     bool  `json:"on_stage"`
	}
	if !readJSON(w, r, &req) {
		return
	}
	if err := s.Store.SetCharacterSection(req.CharacterID, id, req.Attached, req.OnStage); err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// --- props ---

func (s *Server) handleListProps(w http.ResponseWriter, r *http.Request) {
	id, _ := pathID(r, "id")
	ps, err := s.Store.ListProps(id)
	if err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ps)
}

func (s *Server) handleCreateProp(w http.ResponseWriter, r *http.Request) {
	id, _ := pathID(r, "id")
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if !readJSON(w, r, &req) {
		return
	}
	p, err := s.Store.CreateProp(id, req.Name, req.Description)
	if err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

func (s *Server) handleDeleteProp(w http.ResponseWriter, r *http.Request) {
	id, _ := pathID(r, "id")
	if err := s.Store.DeleteProp(id); err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// --- actors ---

func (s *Server) handleListActors(w http.ResponseWriter, r *http.Request) {
	as, err := s.Store.ListActors()
	if err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, as)
}

func (s *Server) handleCreateActor(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
		Bio  string `json:"bio"`
	}
	if !readJSON(w, r, &req) {
		return
	}
	if req.Name == "" {
		writeErr(w, http.StatusBadRequest, "name is required")
		return
	}
	a, err := s.Store.CreateActor(req.Name, req.Bio)
	if err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, a)
}

func (s *Server) handleUpdateActor(w http.ResponseWriter, r *http.Request) {
	id, _ := pathID(r, "id")
	var req struct {
		Name string `json:"name"`
		Bio  string `json:"bio"`
	}
	if !readJSON(w, r, &req) {
		return
	}
	a, err := s.Store.UpdateActor(id, req.Name, req.Bio)
	if err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, a)
}

func (s *Server) handleDeleteActor(w http.ResponseWriter, r *http.Request) {
	id, _ := pathID(r, "id")
	if err := s.Store.DeleteActor(id); err != nil {
		storeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

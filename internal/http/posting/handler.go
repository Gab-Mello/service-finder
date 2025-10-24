package posting

import (
	"encoding/json"
	"net/http"
	"strings"

	domain "github.com/Gab-Mello/service-finder/internal/posting"
)

type Handler struct{ svc *domain.Service }

func NewHandler(s *domain.Service) *Handler { return &Handler{svc: s} }

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	pid := strings.TrimSpace(r.URL.Query().Get("userId"))
	if pid == "" {
		writeErr(w, 400, "missing query parameter: userId")
		return
	}

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, 400, "invalid json")
		return
	}

	p, err := h.svc.Create(pid, req.Title, req.Description, req.Price, req.Category, req.City, req.District)
	if err != nil {
		writeErr(w, 400, err.Error())
		return
	}
	writeJSON(w, 201, p)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	pid := strings.TrimSpace(r.URL.Query().Get("userId"))
	if pid == "" {
		writeErr(w, 400, "missing query parameter: userId")
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/postings/")

	var patch map[string]any
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		writeErr(w, 400, "invalid json")
		return
	}

	p, err := h.svc.Update(pid, id, patch)
	if err != nil {
		writeErr(w, statusFor(err), err.Error())
		return
	}
	writeJSON(w, 200, p)
}

func (h *Handler) Archive(w http.ResponseWriter, r *http.Request) {
	pid := strings.TrimSpace(r.URL.Query().Get("userId"))
	if pid == "" {
		writeErr(w, 400, "missing query parameter: userId")
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/postings/")
	id = strings.TrimSuffix(id, "/archive")

	if err := h.svc.Archive(pid, id); err != nil {
		writeErr(w, statusFor(err), err.Error())
		return
	}
	writeJSON(w, 200, map[string]string{"status": "ok"})
}

func (h *Handler) ListMine(w http.ResponseWriter, r *http.Request) {
	pid := strings.TrimSpace(r.URL.Query().Get("userId"))
	if pid == "" {
		writeErr(w, 400, "missing query parameter: userId")
		return
	}
	list, _ := h.svc.ListMine(pid)
	writeJSON(w, 200, list)
}

func (h *Handler) GetPublic(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/postings/")
	p, err := h.svc.GetPublic(id)
	if err != nil {
		writeErr(w, 404, "not found")
		return
	}
	writeJSON(w, 200, p)
}

func (h *Handler) ListPublic(w http.ResponseWriter, r *http.Request) {
	list, _ := h.svc.ListPublic()
	writeJSON(w, 200, list)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
func statusFor(err error) int {
	switch err {
	case domain.ErrForbidden:
		return 403
	case domain.ErrNotFound:
		return 404
	default:
		return 400
	}
}
